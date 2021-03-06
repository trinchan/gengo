Golang client for the [Gengo](https://gengo.com) API

[Documentation](https://godoc.org/github.com/trinchan/gengo) and testing WIP

```go
package main

import (
	"fmt"

	. "github.com/trinchan/gengo"
	"github.com/trinchan/gengo/lang"
)

func main() {
	// Set GENGO_PUBLIC_KEY, GENGO_PRIVATE_KEY environment variables
	g := NewFromEnv()
	// or use New
	// g := New(publicKey, privateKey, SandboxURL)
	req := NewLanguagePairsRequest(WithSource(lang.Japanese))
	r, err := g.LanguagePairs(req)
	if err != nil {
		panic(fmt.Errorf("Error retrieving language pairs: %v", err))
	}
	for _, lp := range r.LanguagePairs {
		fmt.Printf("%s -> %s (%s) - %0.2f %s\n", lp.Source, lp.Target, lp.Tier, lp.UnitPrice, lp.Currency)
	}
}
```
