package gengo

import "github.com/trinchan/gengo/lang"

const (
	accountNamespace = "/account"
)

// AccountStatsResponse defines the response for the AccountStatus() endpoint.
type AccountStatsResponse struct {
	CreditsSpent Float64 `json:"credits_spent"`
	Processing   Float64 `json:"processing"`
	UserSince    Time    `json:"user_since"`
	Currency     string  `json:"currency"`
	BillingType  string  `json:"billing_type"`
	CustomerType string  `json:"customer_type"`
}

// AccountStats retrieves account stats, such as orders made.
func (c *Client) AccountStats() (*AccountStatsResponse, error) {
	asr := new(AccountStatsResponse)
	err := c.get(accountNamespace+"/stats", nil, asr)
	return asr, err
}

// MeResponse defines the response for the Me() endpoint.
type MeResponse struct {
	Email        string `json:"email"`
	Name         string `json:"full_name"`
	DisplayName  string `json:"display_name"`
	LanguageCode string `json:"language_code"`
}

// Me retrieves account information, such as email.
func (c *Client) Me() (*MeResponse, error) {
	mr := new(MeResponse)
	err := c.get(accountNamespace+"/me", nil, mr)
	return mr, err
}

// BalanceResponse defines the response for the Balance() endpoint.
type BalanceResponse struct {
	Credits  Float64 `json:"credits"`
	Currency string  `json:"currency"`
}

// Balance retrieves account balance in credits.
func (c *Client) Balance() (*BalanceResponse, error) {
	br := new(BalanceResponse)
	err := c.get(accountNamespace+"/balance", nil, br)
	return br, err
}

// PreferredTranslatorsResponse defines the response for the PreferredTranslators() endpoint.
type PreferredTranslatorsResponse struct {
	PreferredTranslators []PreferredTranslatorResponse
}

// PreferredTranslatorResponse defines the structure for a single PreferredTranslator in a PreferredTranslatorsResponse.
type PreferredTranslatorResponse struct {
	lang.Pair
	Tier
	Translators []PreferredTranslator `json:"translators"`
}

// PreferredTranslators retrieves preferred translators set by user.
func (c *Client) PreferredTranslators() (*PreferredTranslatorsResponse, error) {
	ptr := []PreferredTranslatorResponse{}
	err := c.get(accountNamespace+"/preferred_translators", nil, ptr)
	return &PreferredTranslatorsResponse{PreferredTranslators: ptr}, err
}
