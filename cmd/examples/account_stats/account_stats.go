package main

import (
	"log"

	"github.com/trinchan/gengo"
)

func main() {
	// Set GENGO_PUBLIC_KEY, GENGO_PRIVATE_KEY, GENGO_PRODUCTION vars
	g := gengo.NewFromEnv()
	r, err := g.AccountStats()
	if err != nil {
		log.Printf("Error retrieving account stats: %v", err)
	}
	log.Printf("User since: %s", r.UserSince)
	log.Printf("Credits spent: %0.2f %s", r.CreditsSpent, r.Currency)
}
