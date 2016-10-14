package sign

// Signer is something which can sign data for a Gengo API call
type Signer interface {
	Sign(data string) string
}
