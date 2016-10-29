Golang client for the [Gengo](https://gengo.com) API

Documentation and testing WIP

```go
package main

import (
	"fmt"

	"github.com/trinchan/gengo"
	"github.com/trinchan/gengo/lang"
)

func main() {
	// Set GENGO_PUBLIC_KEY, GENGO_PRIVATE_KEY environment variables
	g := gengo.NewFromEnv()
	// or use gengo.New
	// g := gengo.New(publicKey, privateKey, gengo.SandboxURL)
	req := gengo.NewLanguagePairsRequest(gengo.WithSource(lang.Japanese))
	r, err := g.LanguagePairs(req)
	if err != nil {
		panic(fmt.Errorf("Error retrieving language pairs: %v", err))
	}
	for _, lp := range r.LanguagePairs {
		fmt.Printf("%s -> %s (%s) - %0.2f %s\n", lp.Source, lp.Target, lp.Tier, lp.UnitPrice, lp.Currency)
	}
}
```
