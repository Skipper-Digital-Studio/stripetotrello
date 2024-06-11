package stripetrello

import (
	"fmt"
	"strings"
	"sync"

	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

type (
	StripeEventHandler func(event *stripe.Event) error
	StripeTrello       struct {
		stripeWebhookSecret string

		handlers map[string][]StripeEventHandler
	}

	StripeEventError struct {
		fn   string
		args []interface{}
		err  error
	}

	StripeEventErrors []StripeEventError
)

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

func (st StripeTrello) Handler(eventType string) ([]StripeEventHandler, error) {
	handler, ok := st.handlers[eventType]
	if !ok {
		return nil, newError("StripeTrello.Handler", []interface{}{eventType}, fmt.Errorf(fmt.Sprintf("No %s found in available handlers", eventType)))
	}
	return handler, nil
}

func (st StripeTrello) Event(raw []byte, signature string) (*stripe.Event, error) {
	event, err := webhook.ConstructEvent(raw, signature, st.stripeWebhookSecret)
	if err != nil {
		return nil, newError("StripeTrello.Event", []interface{}{raw, signature}, err)
	}

	return &event, nil
}

func (st *StripeTrello) AppendHandler(eventType string, handlers ...StripeEventHandler) {
	h, ok := st.handlers[eventType]
	if !ok {
		st.handlers[eventType] = handlers
	}

	h = append(h, handlers...)
	st.handlers[eventType] = h
}

func (st *StripeTrello) Handle(event *stripe.Event) error {
	handlers, err := st.Handler(string(event.Type))
	if err != nil {
		return newError("StripeTrello.Handle", []interface{}{event}, err)
	}

	for i, h := range handlers {
		if err := h(event); err != nil {
			return newError(fmt.Sprintf("StripeTrello.Handle.handlers[%d]", i), []interface{}{event}, err)
		}
	}

	return nil
}

func (st *StripeTrello) HandleParallel(event *stripe.Event) error {
	handlers, err := st.Handler(string(event.Type))
	if err != nil {
		return newError("StripeTrello.HandleParallel", []interface{}{event}, err)
	}
	var wg sync.WaitGroup

	errors := make(chan StripeEventError, len(handlers))

	for i, h := range handlers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h(event); err != nil {
				errors <- newError(fmt.Sprintf("StripeTrello.Handle.handlers[%d]", i), []interface{}{event}, err)
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
		return newError("StripeTrello.Handle", []interface{}{event}, errs)
	}

	return nil
}
