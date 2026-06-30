package endoflife

import (
	"context"
	"net/url"
)

// Index lists the main endoflife.date API endpoints.
//
// GET /
func (c *Client) Index(ctx context.Context) (*UriListResponse, error) {
	var out UriListResponse
	if err := c.get(ctx, "/", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Products lists all products with a summary of each.
//
// GET /products
func (c *Client) Products(ctx context.Context) (*ProductListResponse, error) {
	var out ProductListResponse
	if err := c.get(ctx, "/products", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ProductsFull lists all products with their complete details. This returns a
// large payload; prefer Products when only summaries are needed.
//
// GET /products/full
func (c *Client) ProductsFull(ctx context.Context) (*FullProductListResponse, error) {
	var out FullProductListResponse
	if err := c.get(ctx, "/products/full", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Product retrieves the full details of a single product, including its
// release cycles.
//
// GET /products/{product}
func (c *Client) Product(ctx context.Context, product string) (*ProductResponse, error) {
	var out ProductResponse
	if err := c.get(ctx, "/products/"+url.PathEscape(product), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ProductRelease retrieves a specific release cycle of a product.
//
// GET /products/{product}/releases/{release}
func (c *Client) ProductRelease(ctx context.Context, product, release string) (*ProductReleaseResponse, error) {
	path := "/products/" + url.PathEscape(product) + "/releases/" + url.PathEscape(release)
	var out ProductReleaseResponse
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ProductLatestRelease retrieves the latest release cycle of a product.
//
// GET /products/{product}/releases/latest
func (c *Client) ProductLatestRelease(ctx context.Context, product string) (*ProductReleaseResponse, error) {
	path := "/products/" + url.PathEscape(product) + "/releases/latest"
	var out ProductReleaseResponse
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Categories lists all categories.
//
// GET /categories
func (c *Client) Categories(ctx context.Context) (*UriListResponse, error) {
	var out UriListResponse
	if err := c.get(ctx, "/categories", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ProductsByCategory lists all product summaries within a category.
//
// GET /categories/{category}
func (c *Client) ProductsByCategory(ctx context.Context, category string) (*ProductListResponse, error) {
	var out ProductListResponse
	if err := c.get(ctx, "/categories/"+url.PathEscape(category), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Tags lists all tags.
//
// GET /tags
func (c *Client) Tags(ctx context.Context) (*UriListResponse, error) {
	var out UriListResponse
	if err := c.get(ctx, "/tags", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ProductsByTag lists all product summaries having the given tag.
//
// GET /tags/{tag}
func (c *Client) ProductsByTag(ctx context.Context, tag string) (*ProductListResponse, error) {
	var out ProductListResponse
	if err := c.get(ctx, "/tags/"+url.PathEscape(tag), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IdentifierTypes lists all identifier types known to endoflife.date (e.g. purl).
//
// GET /identifiers
func (c *Client) IdentifierTypes(ctx context.Context) (*UriListResponse, error) {
	var out UriListResponse
	if err := c.get(ctx, "/identifiers", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IdentifiersByType lists all identifiers for the given type, each referencing
// its related product.
//
// GET /identifiers/{identifier_type}
func (c *Client) IdentifiersByType(ctx context.Context, identifierType string) (*IdentifierListResponse, error) {
	var out IdentifierListResponse
	if err := c.get(ctx, "/identifiers/"+url.PathEscape(identifierType), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
