package gengo

import (
	"fmt"

	"github.com/trinchan/gengo/language"
)

func ExampleClient_PostJobs() {
	g := NewFromEnv()
	jobs := []*JobRequest{
		NewJobRequest("Text to translate", language.NewPair(language.English, language.Japanese), TierStandard),
		NewJobRequest("翻訳するテキスト", language.NewPair(language.Japanese, language.English), TierStandard),
	}
	req := NewPostJobsRequest(jobs)
	r, err := g.PostJobs(req)
	if err != nil {
		fmt.Printf("Error posting jobs: %v\n", err)
	}
	fmt.Printf("Order ID: %d\n", r.OrderID)
	fmt.Printf("New jobs posted: %d\n", r.Count)
	fmt.Printf("Credits used: %0.2f %s\n", r.CreditsUsed, r.Currency)
	for _, job := range r.Jobs {
		fmt.Printf("Duplicate job: %d\n", job.ID)
	}
}
