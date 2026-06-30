package endoflife_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	endoflife "github.com/shyim/go-endoflife-api"
)

func ExampleClient_Product() {
	client := endoflife.NewClient()

	resp, err := client.Product(context.Background(), "go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Result.Label)
	for i := range resp.Result.Releases {
		rel := &resp.Result.Releases[i]
		fmt.Printf("%s: EOL=%v\n", rel.Name, rel.IsEol)
	}
}

func ExampleIsNotFound() {
	client := endoflife.NewClient()

	_, err := client.Product(context.Background(), "totally-unknown-product")
	if endoflife.IsNotFound(err) {
		fmt.Println("product not found")
		return
	}

	var apiErr *endoflife.APIError
	if errors.As(err, &apiErr) && apiErr.RetryAfter > 0 {
		fmt.Printf("rate limited, retry after %s\n", apiErr.RetryAfter)
	}
}
