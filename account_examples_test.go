package gengo

import (
	"fmt"
	"log"
)

func ExampleClient_AccountStats() {
	g := NewFromEnv()
	r, err := g.AccountStats()
	if err != nil {
		fmt.Printf("Error retrieving account stats: %v\n", err)
	}
	fmt.Printf("User since: %s\n", r.UserSince)
	fmt.Printf("Credits spent: %0.2f %s\n", r.CreditsSpent, r.Currency)
}

func ExampleClient_Balance() {
	g := NewFromEnv()
	r, err := g.Balance()
	if err != nil {
		fmt.Printf("Error retrieving account balance: %v\n", err)
	}
	log.Printf("Balance: %0.2f %s", r.Credits, r.Currency)
}

func ExampleClient_PreferredTranslators() {
	g := NewFromEnv()
	r, err := g.PreferredTranslators()
	if err != nil {
		fmt.Printf("Error retrieving preferred translators: %v\n", err)
	}
	for _, lp := range r.PreferredTranslators {
		fmt.Printf("%s - %s (%s)\n", lp.Source, lp.Target, lp.Tier)
		for _, translator := range lp.Translators {
			fmt.Printf("\t%d - Last seen: %s\n", translator.ID, translator.LastLogin)
		}
	}
}
