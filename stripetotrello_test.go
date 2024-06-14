package stripetotrello

import (
	stripe "github.com/stripe/stripe-go/v76"
	"testing"
)

func TestHandler(t *testing.T) {

	type testCase struct {
		event  string
		err    error
		lenght int
	}

	client := NewClient()
	client.AppendHandler("cusomer.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	})
	client.AppendHandler("cusomer.updated", func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	})

	tcs := []testCase{
		{"customer.created", nil, 3},
		{"customer.updated", nil, 1},
		{"checkout.session.completed", StripeUnsupportedEventError{}, 0},
		{"subscription.created", StripeUnsupportedEventError{}, 0},
	}

	for _, tc := range tcs {
		res, err := client.Handler(tc.event)

		if err == nil && len(res) != tc.lenght {
			t.Errorf("Expected number of results %d, got %d", tc.lenght, len(res))
		}

		if err != nil && res != nil {
			t.Errorf("Expected number of results %d, got %d", tc.lenght, len(res))
		}
	}
}

func TestHandle(t *testing.T) {

	type testCase struct {
		event  string
		err    error
		lenght int
	}

	client := NewClient()
	client.AppendHandler("cusomer.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	})
	client.AppendHandler("cusomer.updated", func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	})

	tcs := []testCase{
		{"customer.created", nil, 3},
		{"customer.updated", nil, 1},
		{"checkout.session.completed", StripeUnsupportedEventError{}, 0},
		{"subscription.created", StripeUnsupportedEventError{}, 0},
	}

	for _, tc := range tcs {
		res, err := client.Handler(tc.event)

		if err == nil && len(res) != tc.lenght {
			t.Errorf("Expected number of results %d, got %d", tc.lenght, len(res))
		}

		if err != nil && res != nil {
			t.Errorf("Expected number of results %d, got %d", tc.lenght, len(res))
		}
	}
}
