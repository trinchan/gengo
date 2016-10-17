package gengo

import "fmt"

func ExampleNew() {
	publicKey := "{PUBLIC_KEY}"
	privateKey := "{PRIVATE_KEY}"
	baseURL := SandboxBaseURL
	g := New(publicKey, privateKey, baseURL)
	r, err := g.AccountStats()
	if err != nil {
		fmt.Printf("Error retrieving account stats: %v\n", err)
	}
	fmt.Printf("User since: %s\n", r.UserSince)
	fmt.Printf("Credits spent: %0.2f %s\n", r.CreditsSpent, r.Currency)
}

func ExampleNewFromEnv() {
	g := NewFromEnv() // Client options are loaded from the environment
	r, err := g.AccountStats()
	if err != nil {
		fmt.Printf("Error retrieving account stats: %v\n", err)
	}
	fmt.Printf("User since: %s\n", r.UserSince)
	fmt.Printf("Credits spent: %0.2f %s\n", r.CreditsSpent, r.Currency)
}
