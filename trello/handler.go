package trello

import (
	"encoding/json"
	"fmt"

	stripe "github.com/stripe/stripe-go/v76"
)

type (
	TrelloError struct {
		fn   string
		args []interface{}
		err  error
	}
)

func (t TrelloError) Error() string {
	return fmt.Sprintf("Error calling %s - with args %v - result in error %s", t.fn, t.args, t.err.Error())
}

func NewTrelloError(fn string, args []interface{}, err error) TrelloError {
	return TrelloError{
		fn,
		args,
		err,
	}
}

func (c *Client) DefaultHandlerBuilder(opts ...func(*CreateBoardReq)) func(*stripe.Event) (interface{}, error) {
	return func(event *stripe.Event) (interface{}, error) {
		errorFN := func(err error) error {
			return NewTrelloError("trello.handler.DefaultHandlerBuilder", []interface{}{event}, err)

		}
		req := NewCreateBoardReq(opts...)

		switch event.Type {
		case "customer.subscription.created":
			if err := defaultCheckoutSessionCompleted(event, req); err != nil {
				return nil, errorFN(err)
			}
		case "customer.created":
			if err := defaultCustomerCreated(event, req); err != nil {
				return nil, errorFN(err)
			}
		case "checkout.session.completed":
			if err := defaultCheckoutSessionCompleted(event, req); err != nil {
				return nil, errorFN(err)
			}
		default:
			return nil, errorFN(fmt.Errorf("Unsupported event %s", string(event.Type)))
		}

		board, err := c.BoardByName(req.Name)
		if err == nil {
			if err := c.SendInvites(req.EmailsToInvite, board.Id); err != nil {
				return nil, errorFN(err)
			}
		}
		res, err := c.NewBoard(*req)
		if err != nil {
			return nil, errorFN(err)
		}

		return res, nil
	}
}

func (c *Client) SendInvites(emails []string, boardID string) error {
	for _, email := range emails {
		if err := c.Invite(email, boardID); err != nil {
			return NewTrelloError("trello.handler.SendInvites", []interface{}{emails, boardID}, err)
		}
	}
	return nil
}

func defaultSubscriptionCreated(event *stripe.Event, req *CreateBoardReq) error {
	var s stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &s); err != nil {
		return NewTrelloError("trello.handler.defaultSubscriptionCreated", []interface{}{event, req}, err)
	}

	if req.Name == "" {
		req.Name = s.Customer.Name
	}

	if req.Description == "" {
		req.Description = fmt.Sprintf("New board created for the customer with name = %s and ID = %s", s.Customer.Name, s.Customer.ID)
	}

	if req.EmailsToInvite == nil || len(req.EmailsToInvite) == 0 {
		req.EmailsToInvite = []string{s.Customer.Email}
	}

	return nil
}

func defaultCustomerCreated(event *stripe.Event, req *CreateBoardReq) error {
	var c stripe.Customer
	if err := json.Unmarshal(event.Data.Raw, &c); err != nil {
		return NewTrelloError("trello.handler.defaultCustomerCreated", []interface{}{event, req}, err)
	}

	if req.Name == "" {
		req.Name = c.Name
	}

	if req.Description == "" {
		req.Description = fmt.Sprintf("New board created for the customer with name = %s and ID = %s", c.Name, c.ID)
	}
	if req.EmailsToInvite == nil || len(req.EmailsToInvite) == 0 {
		req.EmailsToInvite = []string{c.Email}
	}
	return nil
}

func defaultCheckoutSessionCompleted(event *stripe.Event, req *CreateBoardReq) error {
	var s stripe.CheckoutSession
	err := json.Unmarshal(event.Data.Raw, &s)
	if err != nil {
		return NewTrelloError("trello.handler.defaultCheckoutSessionCompleted", []interface{}{event, req}, err)
	}
	if req.Name == "" {
		req.Name = fmt.Sprintf("%s-%s", s.Customer.Name, s.Customer.ID)
	}

	if req.Description == "" {
		req.Description = fmt.Sprintf("New board created for the customer with name = %s and ID = %s", s.Customer.Name, s.Customer.ID)
	}
	if req.EmailsToInvite == nil || len(req.EmailsToInvite) == 0 {
		req.EmailsToInvite = []string{s.Customer.Email}
	}
	return nil
}
