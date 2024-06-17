package trello

import (
	"fmt"
	"strings"
)

type (
	NewCardReqs struct {
		Cards []NewCardReq `json:"cards"`
	}

	NewCardReq struct {
		Name        string `json:"name"`
		Description string `json:"desc"`
		IDList      string `json:"id_list"`
	}

	GetListReq struct {
		Cards      StatusEnum
		CardFields string
		Fields     string
		Filter     StatusEnum
		ID         string
	}

	GetListReqParams struct {
		Cards      string `json:"cards"`
		CardFields string `json:"card_fields"`
		Filter     string `json:"filter"`
		Fields     string `json:"fields"`
	}

	CreateBoardReq struct {
		Name           string
		Description    string
		IDOrganization string

		IDBoardSource  string
		KeepFromSource string
		EmailsToInvite []string
	}
)

func (c CreateBoardReq) Ready() bool {
	if c.Name == "" || c.Description == "" || c.IDOrganization == "" {
		return false
	}
	return true
}

func NewCreateBoardReq(cfgs ...func(*CreateBoardReq)) *CreateBoardReq {
	output := &CreateBoardReq{}
	for _, f := range cfgs {
		f(output)
	}

	return output
}

func CreateBoardWithEmailsToInvite(emails []string) func(*CreateBoardReq) {
	return func(c *CreateBoardReq) {
		c.EmailsToInvite = emails
	}
}

func CreateBoardWithName(name string) func(*CreateBoardReq) {
	return func(c *CreateBoardReq) {
		c.Name = name
	}
}

func CreateBoardWithDescription(desc string) func(*CreateBoardReq) {
	desc = strings.Replace(strings.Replace(desc, " ", "%20", -1), "\n", "%0A", -1)
	return func(c *CreateBoardReq) {
		c.Description = desc
	}
}

func CreateBoardWithSource(source string) func(*CreateBoardReq) {
	return func(c *CreateBoardReq) {
		c.IDBoardSource = source
	}
}

func CreateBoardWithOrganization(org string) func(*CreateBoardReq) {
	return func(c *CreateBoardReq) {
		c.IDOrganization = org
	}
}

func CreateBoardWithKeepFromSource(keep string) func(*CreateBoardReq) {
	if keep == "" {
		keep = "cards"
	}

	return func(c *CreateBoardReq) {
		c.KeepFromSource = keep
	}
}
func (c *CreateBoardReq) String() string {
	output := []string{
		fmt.Sprintf("name=%s", strings.ToLower(c.Name)),
		fmt.Sprintf("desc=%s", strings.ToLower(c.Description)),
	}

	if c.IDBoardSource != "" {
		output = append(output, fmt.Sprintf("idBoardSource=%s", c.IDBoardSource))
	}

	if c.KeepFromSource != "none" {
		output = append(output, fmt.Sprintf("keepFromSource=%s", c.KeepFromSource))
	}

	if c.IDOrganization != "" {
		output = append(output, fmt.Sprintf("idOrganization=%s", c.IDOrganization))
	}

	return strings.Join(output, "&")
}

func (t *GetListReqParams) String() string {
	var output []string

	if t.Cards != "NONE" {
		output = append(output, fmt.Sprintf("cards=%s", strings.ToLower(t.Cards)))
		output = append(output, fmt.Sprintf("card_fields=%s", t.CardFields))
	}

	if t.Filter != "NONE" {
		output = append(output, fmt.Sprintf("filter=%s", strings.ToLower(t.Filter)))
		output = append(output, fmt.Sprintf("fields=%s", t.Fields))
	}

	return strings.Join(output, "&")
}

func (t *GetListReq) GetBody() (*GetListReqParams, error) {
	cards, err := t.Cards.GetValue()
	if err != nil {
		return nil, err
	}
	filter, err := t.Filter.GetValue()
	if err != nil {
		return nil, err
	}

	return &GetListReqParams{
		Cards:      cards,
		Fields:     t.Fields,
		Filter:     filter,
		CardFields: t.CardFields,
	}, nil
}

func NewGetlistReq(id string) GetListReq {
	return GetListReq{
		Cards:      NONE,
		CardFields: "all",
		Filter:     NONE,
		Fields:     "all",
		ID:         id,
	}
}
