package stripetotrello

import (
	"fmt"
	"strings"
	"sync"

	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

type (
	StripeEventHandler       func(event *stripe.Event) error
	StripeFailedEventHandler func(event *stripe.Event, err error) error
	Client                   struct {
		stripeWebhookSecret string

		handlers       map[string][]StripeEventHandler
		successHandler map[string]StripeEventHandler
		failureHandler map[string]StripeFailedEventHandler
	}

	StripeEventError struct {
		fn   string
		args []interface{}
		err  error
	}

	StripeEventErrors []StripeEventError
)

func NewClient(cfgs ...func(*Client)) *Client {
	c := &Client{}
	for _, f := range cfgs {
		f(c)
	}
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
		return nil, newError("Client.Handler", []interface{}{eventType}, fmt.Errorf(fmt.Sprintf("No %s found in available handlers", eventType)))
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
	h, ok := st.handlers[eventType]
	if !ok {
		st.handlers[eventType] = handlers
	}

	h = append(h, handlers...)
	st.handlers[eventType] = h
}

func (st *Client) AddSuccessHandler(eventType string, handler StripeEventHandler) {
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

	for i, h := range handlers {
		if err := h(event); err != nil {
			fh, ok := st.failureHandler[string(event.Type)]
			if !ok {
				return newError(fmt.Sprintf("Client.Handle.handlers[%d]", i), []interface{}{event}, err)
			}
			return fh(event, err)
		}
	}

	h, ok := st.successHandler[string(event.Type)]
	if !ok {
		return nil
	}

	return h(event)
}

func (st *Client) HandleParallel(event *stripe.Event) error {
	handlers, err := st.Handler(string(event.Type))
	if err != nil {
		return newError("Client.HandleParallel", []interface{}{event}, err)
	}
	var wg sync.WaitGroup

	errors := make(chan StripeEventError, len(handlers))

	for i, h := range handlers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h(event); err != nil {
				errors <- newError(fmt.Sprintf("Client.Handle.handlers[%d]", i), []interface{}{event}, err)
			}
		}()
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		errs := StripeEventErrors{}
		for err := range errors {
			errs = append(errs, err)
		}
		fh, ok := st.failureHandler[string(event.Type)]
		nErr := newError("Client.Handle", []interface{}{event}, errs)
		if !ok {
			return nErr
		}
		return fh(event, nErr)
	}

	h, ok := st.successHandler[string(event.Type)]
	if !ok {
		return nil
	}

	return h(event)
}
