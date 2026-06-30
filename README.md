# go-endoflife-api

A Go client for the [endoflife.date](https://endoflife.date) [API v1](https://endoflife.date/docs/api/v1/).

endoflife.date documents EOL dates and support lifecycles for hundreds of products (operating systems, languages, frameworks, hardware, …). This SDK lets you discover and query that data from Go.

## Install

```sh
go get github.com/shyim/go-endoflife-api
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	endoflife "github.com/shyim/go-endoflife-api"
)

func main() {
	client := endoflife.NewClient()

	resp, err := client.Product(context.Background(), "ubuntu")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Result.Label)
	for _, rel := range resp.Result.Releases {
		fmt.Printf("  %-8s eol=%-5v maintained=%v\n", rel.Name, rel.IsEol, rel.IsMaintained)
	}
}
```

## Endpoints

Every endpoint takes a `context.Context` as its first argument.

| Method | API |
| --- | --- |
| `Index(ctx)` | `GET /` |
| `Products(ctx)` | `GET /products` |
| `ProductsFull(ctx)` | `GET /products/full` |
| `Product(ctx, product)` | `GET /products/{product}` |
| `ProductRelease(ctx, product, release)` | `GET /products/{product}/releases/{release}` |
| `ProductLatestRelease(ctx, product)` | `GET /products/{product}/releases/latest` |
| `Categories(ctx)` | `GET /categories` |
| `ProductsByCategory(ctx, category)` | `GET /categories/{category}` |
| `Tags(ctx)` | `GET /tags` |
| `ProductsByTag(ctx, tag)` | `GET /tags/{tag}` |
| `IdentifierTypes(ctx)` | `GET /identifiers` |
| `IdentifiersByType(ctx, type)` | `GET /identifiers/{identifier_type}` |

## Configuration

`NewClient` accepts options:

```go
client := endoflife.NewClient(
	endoflife.WithHTTPClient(myHTTPClient), // any type with Do(*http.Request)
	endoflife.WithUserAgent("my-app/1.0"),
	endoflife.WithBaseURL("https://endoflife.date/api/v1"),
)
```

The default HTTP client has a 30s timeout and follows the `301` redirects that the
API issues when products, categories or tags are renamed.

## Error handling

Non-`200` responses are returned as `*APIError`. Two helpers cover the common cases:

```go
_, err := client.Product(ctx, "unknown")
switch {
case endoflife.IsNotFound(err):
	// 404 – product does not exist
case endoflife.IsTooManyRequests(err):
	// 429 – inspect APIError.RetryAfter
	var apiErr *endoflife.APIError
	errors.As(err, &apiErr)
	time.Sleep(apiErr.RetryAfter)
}
```

## Dates and null values

- Date-only fields use the `Date` type, which (un)marshals `"YYYY-MM-DD"`.
- Fields the API documents as nullable are modeled as pointers (`*Date`, `*string`,
  `*bool`), where `nil` means the value is JSON `null` / unknown.
- Fields the API omits entirely for inapplicable products (`isEoas`, `isEoes`,
  `isDiscontinued`, and their `*From` dates, plus `custom`) use `omitempty` pointers.

## License

MIT
