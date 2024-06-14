package stripetotrello

import (
	"fmt"
	"testing"

	stripe "github.com/stripe/stripe-go/v76"
)

func TestHandler(t *testing.T) {

	type testCase struct {
		event  string
		err    error
		lenght int
	}

	client := NewClient()
	client.AppendHandler("customer.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		return nil, nil
	})
	client.AppendHandler("customer.updated", func(_ *stripe.Event) (interface{}, error) {
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
		event      stripe.Event
		shouldFail bool
	}

	type resStruct struct {
		data string
	}

	client := NewClient()
	client.AppendHandler("customer.created", func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 1")
		return "testing 1", nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 3")
		return resStruct{"testing"}, nil
	})

	client.AppendHandler("customer.updated", func(_ *stripe.Event) (interface{}, error) {
		return 0.0, nil
	})

	client.AppendHandler("customer.deleted", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, nil)
	})

	client.AppendHandler("subscription.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, nil)
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	})

	tcs := []testCase{
		{stripe.Event{Type: "customer.created"}, false},
		{stripe.Event{Type: "customer.updated"}, false},
		{stripe.Event{Type: "customer.deleted"}, true},
		{stripe.Event{Type: "subscription.created"}, true},
	}

	for _, tc := range tcs {
		err := client.Handle(&tc.event)
		if err != nil && !tc.shouldFail {
			t.Errorf("Event should have NOT failed event type = %s", tc.event.Type)
		}

		if err == nil && tc.shouldFail {
			t.Errorf("Event should have failed event type = %s", tc.event.Type)
		}
	}
}

func TestHandleParallel(t *testing.T) {
	type testCase struct {
		event      stripe.Event
		shouldFail bool
	}

	type resStruct struct {
		data string
	}

	client := NewClient()
	client.AppendHandler("customer.created", func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 1")
		return "testing 1", nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 3")
		return resStruct{"testing"}, nil
	})

	client.AppendHandler("customer.updated", func(_ *stripe.Event) (interface{}, error) {
		return 0.0, nil
	})

	client.AppendHandler("customer.deleted", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, nil)
	})

	client.AppendHandler("subscription.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, nil)
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	})

	tcs := []testCase{
		{stripe.Event{Type: "customer.created"}, false},
		{stripe.Event{Type: "customer.updated"}, false},
		{stripe.Event{Type: "customer.deleted"}, true},
		{stripe.Event{Type: "subscription.created"}, true},
	}

	for _, tc := range tcs {
		err := client.HandleParallel(&tc.event)
		if err != nil && !tc.shouldFail {
			t.Errorf("Event should have NOT failed event type = %s", tc.event.Type)
		}

		if err == nil && tc.shouldFail {
			t.Errorf("Event should have failed event type = %s", tc.event.Type)
		}
	}
}

func TestHandleWithErrorHandler(t *testing.T) {
	type testCase struct {
		event      stripe.Event
		shouldFail bool
	}

	type resStruct struct {
		data string
	}

	client := NewClient()
	client.AppendHandler("customer.created", func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 1")
		return "testing 1", nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 3")
		return resStruct{"testing"}, nil
	})

	client.AppendHandler("customer.updated", func(_ *stripe.Event) (interface{}, error) {
		return 0.0, nil
	})

	client.AppendHandler("customer.deleted", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, nil)
	})

	client.AppendHandler("subscription.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, fmt.Errorf("test"))
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	})

	client.AddFailureHandler("subscription.created", func(_ *stripe.Event, err error) error {
		expectedErr := newError("It fails", []interface{}{"error"}, fmt.Errorf("test"))
		if err.Error() == expectedErr.Error() {
			fmt.Println("Matches")
			return nil
		}
		return err
	})

	tcs := []testCase{
		{stripe.Event{Type: "customer.created"}, false},
		{stripe.Event{Type: "customer.updated"}, false},
		{stripe.Event{Type: "customer.deleted"}, true},
		{stripe.Event{Type: "subscription.created"}, false},
	}

	for _, tc := range tcs {
		err := client.Handle(&tc.event)
		if err != nil && !tc.shouldFail {
			t.Errorf("Event should have NOT failed event type = %s", tc.event.Type)
		}

		if err == nil && tc.shouldFail {
			t.Errorf("Event should have failed event type = %s", tc.event.Type)
		}
	}
}

func TestHandleWithErrorHandlerParallel(t *testing.T) {
	type testCase struct {
		event      stripe.Event
		shouldFail bool
	}

	type resStruct struct {
		data string
	}

	client := NewClient()
	client.AppendHandler("customer.created", func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 1")
		return "testing 1", nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 3")
		return resStruct{"testing"}, nil
	})

	client.AppendHandler("customer.updated", func(_ *stripe.Event) (interface{}, error) {
		return 0.0, nil
	})

	client.AppendHandler("customer.deleted", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, nil)
	})

	client.AppendHandler("subscription.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, fmt.Errorf("test"))
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	})

	client.AddFailureHandler("subscription.created", func(_ *stripe.Event, err error) error {
		return nil
	})

	tcs := []testCase{
		{stripe.Event{Type: "customer.created"}, false},
		{stripe.Event{Type: "customer.updated"}, false},
		{stripe.Event{Type: "customer.deleted"}, true},
		{stripe.Event{Type: "subscription.created"}, false},
	}

	for _, tc := range tcs {
		err := client.HandleParallel(&tc.event)
		if err != nil && !tc.shouldFail {
			t.Errorf("Event should have NOT failed event type = %s", tc.event.Type)
		}

		if err == nil && tc.shouldFail {
			t.Errorf("Event should have failed event type = %s", tc.event.Type)
		}
	}
}

func TestHandleWithErrorAndSuccessHandlerParallel(t *testing.T) {
	type testCase struct {
		event      stripe.Event
		shouldFail bool
	}

	type resStruct struct {
		data string
	}

	client := NewClient()
	client.AppendHandler("customer.created", func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 1")
		return "testing 1", nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 3")
		return resStruct{"testing"}, nil
	})
	client.AddSuccessHandler("customer.created", func(_ *stripe.Event, results []interface{}) error {
		for _, r := range results {
			switch rType := r.(type) {
			case int:
				fmt.Println("INT")
			case string:
				fmt.Println("STRING")
			case resStruct:
				fmt.Println("CUSTOM STRUCT")
			default:
				return newError("Error unexpected response type", []interface{}{rType}, fmt.Errorf("ERROR"))
			}
		}
		return nil
	})

	client.AppendHandler("customer.updated", func(_ *stripe.Event) (interface{}, error) {
		return 0.0, nil
	})
	client.AddSuccessHandler("customer.updated", func(_ *stripe.Event, results []interface{}) error {
		for _, r := range results {
			switch rType := r.(type) {
			case int:
				fmt.Println("INT")
			case string:
				fmt.Println("STRING")
			case resStruct:
				fmt.Println("CUSTOM STRUCT")
			default:
				return newError("Error unexpected response type", []interface{}{rType}, fmt.Errorf("ERROR"))
			}
		}
		return nil
	})

	client.AppendHandler("customer.deleted", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, nil)
	})

	client.AppendHandler("subscription.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, fmt.Errorf("test"))
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	})

	client.AddFailureHandler("subscription.created", func(_ *stripe.Event, err error) error {
		return nil
	})

	tcs := []testCase{
		{stripe.Event{Type: "customer.created"}, false},
		{stripe.Event{Type: "customer.updated"}, true},
		{stripe.Event{Type: "customer.deleted"}, true},
		{stripe.Event{Type: "subscription.created"}, false},
	}

	for _, tc := range tcs {
		err := client.HandleParallel(&tc.event)
		if err != nil && !tc.shouldFail {
			t.Errorf("Event should have NOT failed event type = %s", tc.event.Type)
		}

		if err == nil && tc.shouldFail {
			t.Errorf("Event should have failed event type = %s", tc.event.Type)
		}
	}
}

func TestHandleWithErrorAndSuccessHandler(t *testing.T) {
	type testCase struct {
		event      stripe.Event
		shouldFail bool
	}

	type resStruct struct {
		data string
	}

	client := NewClient()
	client.AppendHandler("customer.created", func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 1")
		return "testing 1", nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 3")
		return resStruct{"testing"}, nil
	})
	client.AddSuccessHandler("customer.created", func(_ *stripe.Event, results []interface{}) error {
		for _, r := range results {
			switch rType := r.(type) {
			case int:
				fmt.Println("INT")
			case string:
				fmt.Println("STRING")
			case resStruct:
				fmt.Println("CUSTOM STRUCT")
			default:
				return newError("Error unexpected response type", []interface{}{rType}, fmt.Errorf("ERROR"))
			}
		}
		return nil
	})

	client.AppendHandler("customer.updated", func(_ *stripe.Event) (interface{}, error) {
		return 0.0, nil
	})
	client.AddSuccessHandler("customer.updated", func(_ *stripe.Event, results []interface{}) error {
		for _, r := range results {
			switch rType := r.(type) {
			case int:
				fmt.Println("INT")
			case string:
				fmt.Println("STRING")
			case resStruct:
				fmt.Println("CUSTOM STRUCT")
			default:
				return newError("Error unexpected response type", []interface{}{rType}, fmt.Errorf("ERROR"))
			}
		}
		return nil
	})

	client.AppendHandler("customer.deleted", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, nil)
	})

	client.AppendHandler("subscription.created", func(_ *stripe.Event) (interface{}, error) {
		return nil, newError("It fails", []interface{}{"error"}, fmt.Errorf("test"))
	}, func(_ *stripe.Event) (interface{}, error) {
		fmt.Println("handler 2")
		return 2, nil
	})

	client.AddFailureHandler("subscription.created", func(_ *stripe.Event, err error) error {
		return nil
	})

	tcs := []testCase{
		{stripe.Event{Type: "customer.created"}, false},
		{stripe.Event{Type: "customer.updated"}, true},
		{stripe.Event{Type: "customer.deleted"}, true},
		{stripe.Event{Type: "subscription.created"}, false},
	}

	for _, tc := range tcs {
		err := client.HandleParallel(&tc.event)
		if err != nil && !tc.shouldFail {
			t.Errorf("Event should have NOT failed event type = %s", tc.event.Type)
		}

		if err == nil && tc.shouldFail {
			t.Errorf("Event should have failed event type = %s", tc.event.Type)
		}
	}
}
