package stripetotrello

import (
	"fmt"
	"strings"
	"sync"

	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

type (
	StripeSuccessEventHandler func(event *stripe.Event, results []interface{}) error
	StripeFailedEventHandler  func(event *stripe.Event, err error) error
	StripeEventHandler        func(event *stripe.Event) (interface{}, error)

	Client struct {
		stripeWebhookSecret string

		handlers       map[string][]StripeEventHandler
		successHandler map[string]StripeSuccessEventHandler
		failureHandler map[string]StripeFailedEventHandler
	}

	StripeEventError struct {
		fn   string
		args []interface{}
		err  error
	}

	StripeUnsupportedEventError struct {
		event string
	}

	StripeEventErrors []StripeEventError
)

func NewUnsupportedError(event string) StripeUnsupportedEventError {
	return StripeUnsupportedEventError{
		event: event,
	}
}

func (s StripeUnsupportedEventError) Error() string {
	return fmt.Sprintf("Unsupported event detected: %s", s.event)
}

func NewClient(cfgs ...func(*Client)) *Client {
	c := &Client{}
	for _, f := range cfgs {
		f(c)
	}
	c.handlers = make(map[string][]StripeEventHandler)
	c.successHandler = make(map[string]StripeSuccessEventHandler)
	c.failureHandler = make(map[string]StripeFailedEventHandler)
	return c
}

func WithStripeWebhookSecret(secret string) func(*Client) {
	return func(c *Client) {
		c.stripeWebhookSecret = secret
	}
}

func (sees StripeEventErrors) Error() string {
	var output []string
	for _, err := range sees {
		output = append(output, err.Error())
	}

	return strings.Join(output, " - ")
}

func newError(fn string, args []interface{}, err error) StripeEventError {
	return StripeEventError{
		fn,
		args,
		err,
	}
}

func (see StripeEventError) Error() string {
	return fmt.Sprintf("Error calling %s - with args %v - result in error %s", see.fn, see.args, see.err.Error())
}

func (st Client) Handler(eventType string) ([]StripeEventHandler, error) {
	handler, ok := st.handlers[eventType]
	if !ok {
		return nil, NewUnsupportedError(fmt.Sprintf("No %s found in available handlers", eventType))
	}
	return handler, nil
}

func (st Client) Event(raw []byte, signature string) (*stripe.Event, error) {
	event, err := webhook.ConstructEvent(raw, signature, st.stripeWebhookSecret)
	if err != nil {
		return nil, newError("Client.Event", []interface{}{raw, signature}, err)
	}

	return &event, nil
}

func (st *Client) AppendHandler(eventType string, handlers ...StripeEventHandler) {
	if st.handlers == nil {
		st.handlers = make(map[string][]StripeEventHandler)
	}
	h, ok := st.handlers[eventType]
	if !ok {
		st.handlers[eventType] = handlers
	}

	h = append(h, handlers...)
	st.handlers[eventType] = h
}

func (st *Client) AddSuccessHandler(eventType string, handler StripeSuccessEventHandler) {
	st.successHandler[eventType] = handler
}

func (st *Client) AddFailureHandler(eventType string, handler StripeFailedEventHandler) {
	st.failureHandler[eventType] = handler
}

func (st *Client) Handle(event *stripe.Event) error {
	handlers, err := st.Handler(string(event.Type))
	if err != nil {
		return newError("Client.Handle", []interface{}{event}, err)
	}

	results := make([]interface{}, len(handlers))
	for i, h := range handlers {
		res, err := h(event)
		if err != nil {
			fh, ok := st.failureHandler[string(event.Type)]
			if !ok {
				return newError(fmt.Sprintf("Client.Handle.handlers[%d]", i), []interface{}{event}, err)
			}
			return fh(event, err)
		}
		results[i] = res
	}

	h, ok := st.successHandler[string(event.Type)]
	if !ok {
		return nil
	}

	if err = h(event, results); err != nil {
		return err
	}
	return nil
}

func (st *Client) HandleParallel(event *stripe.Event) error {
	handlers, err := st.Handler(string(event.Type))
	switch err.(type) {
	case StripeEventError:
		return newError("Client.HandleParallel", []interface{}{event}, err)
	case StripeUnsupportedEventError:
		return err
	}
	if err != nil {
		return newError("Client.HandleParallel", []interface{}{event}, err)
	}
	var wg sync.WaitGroup

	errors := make(chan StripeEventError, len(handlers))
	results := make(chan interface{}, len(handlers))

	for i, h := range handlers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := h(event)
			if err != nil {
				errors <- newError(fmt.Sprintf("Client.Handle.handlers[%d]", i), []interface{}{event}, err)
			}
			results <- res
		}()
	}

	wg.Wait()
	close(errors)
	close(results)

	if len(errors) > 0 {
		errs := StripeEventErrors{}
		for err := range errors {
			errs = append(errs, err)
		}
		nErr := newError("Client.Handle", []interface{}{event}, errs)
		fh, ok := st.failureHandler[string(event.Type)]
		if !ok {
			return nErr
		}
		tt := fh(event, nErr)
		return tt
	}

	if len(results) != len(handlers) {
		nErr := newError("Client.HandleParallel", []interface{}{event}, fmt.Errorf("Not all the handlers return a valid response"))
		fh, ok := st.failureHandler[string(event.Type)]
		if !ok {
			return nErr
		}
		return fh(event, nErr)
	}

	rs := []interface{}{}
	for r := range results {
		rs = append(rs, r)
	}

	sh, ok := st.successHandler[string(event.Type)]
	if !ok {
		return nil
	}

	return sh(event, rs)
}
