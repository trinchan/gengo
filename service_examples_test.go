package gengo

import (
	"fmt"

	"github.com/trinchan/gengo/lang"
)

// Get all language pairs
func ExampleClient_LanguagePairs_all() {
	g := NewFromEnv()
	req := NewLanguagePairsRequest()
	r, err := g.LanguagePairs(req)
	if err != nil {
		fmt.Printf("Error retrieving language pairs: %v\n", err)
	}
	for _, lp := range r.LanguagePairs {
		fmt.Printf("%s -> %s (%s) - %0.2f %s\n", lp.Source, lp.Target, lp.Tier, lp.UnitPrice, lp.Currency)
	}
}

// Get all language pairs with a specific source language
func ExampleClient_LanguagePairs_source() {
	g := NewFromEnv()
	req := NewLanguagePairsRequest(WithSource(language.English))
	r, err := g.LanguagePairs(req)
	if err != nil {
		fmt.Printf("Error retrieving language pairs: %v\n", err)
	}
	for _, lp := range r.LanguagePairs {
		fmt.Printf("%s -> %s (%s) - %0.2f %s\n", lp.Source, lp.Target, lp.Tier, lp.UnitPrice, lp.Currency)
	}
}
