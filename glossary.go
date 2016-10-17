package gengo

// Glossary is implemented in Go, so we can code nicely here thank god

import (
	"encoding/json"
	"fmt"

	"github.com/trinchan/gengo/language"
)

const (
	glossaryNamespace = "/translate/glossary"
)

// ListGlossariesResponse defines the response for the ListGlossaries() endpoint.
type ListGlossariesResponse struct {
	Glossaries []Glossary
}

// UnmarshalJSON implements the Unmarshaler interface so that we can keep our response types consistent.
func (g *ListGlossariesResponse) UnmarshalJSON(b []byte) error {
	gs := new([]Glossary)
	err := json.Unmarshal(b, gs)
	if err != nil {
		return err
	}
	g.Glossaries = *gs
	return nil
}

// Glossary defines a Glossary for ListGlossariesResponse.
type Glossary struct {
	ID         int                `json:"id"`
	UserID     int                `json:"customer_user_id"`
	SourceID   int                `json:"source_language_id"`
	SourceCode language.Code      `json:"source_language_code"`
	Targets    []GlossaryLanguage `json:"target_languages"`
	Public     bool               `json:"is_public"`
	UnitCount  int                `json:"unit_count"`
	Ctime      Time               `json:"ctime"`
	Title      string             `json:"title"`
	Status     int                `json:"status"`
}

// GlossaryLanguage defines the ID and code pair for glossary responses
type GlossaryLanguage struct {
	ID   int           `json:"0"`
	Code language.Code `json:"1"`
}

// ListGlossaries retrieves a list of glossaries that belongs to the authenticated user.
func (c *Client) ListGlossaries() (*ListGlossariesResponse, error) {
	glr := new(ListGlossariesResponse)
	err := c.get(glossaryNamespace, nil, glr)
	return glr, err
}

// GetGlossaryRequest defines the request parameters for the GetGlossaryByID() endpoint.
type GetGlossaryRequest struct {
	ID int `json:"id"`
}

// NewGetGlossaryRequest creates a new GetGlossaryRequest with the given id and options.
func NewGetGlossaryRequest(id int) *GetGlossaryRequest {
	g := &GetGlossaryRequest{
		ID: id,
	}
	return g
}

// GlossaryResponse defines the response from the GetGlossaryByID() endpoint.
type GlossaryResponse struct {
	Glossary Glossary
}

// UnmarshalJSON implements the Unmarshaler interface so we caan keep our response types consistent.
func (g *GlossaryResponse) UnmarshalJSON(b []byte) error {
	gl := new(Glossary)
	err := json.Unmarshal(b, gl)
	if err != nil {
		return err
	}
	g.Glossary = *gl
	return nil
}

// GetGlossaryByID retrieves a glossary by ID.
func (c *Client) GetGlossaryByID(req *GetGlossaryRequest) (*GlossaryResponse, error) {
	gr := new(GlossaryResponse)
	err := c.get(glossaryNamespace+fmt.Sprintf("/%d", req.ID), nil, gr)
	return gr, err
}
