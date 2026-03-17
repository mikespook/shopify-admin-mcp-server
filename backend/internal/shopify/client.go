package shopify

import (
	goshopify "github.com/bold-commerce/go-shopify/v4"
)

// NewClient creates a Shopify Admin API client using an access token.
func NewClient(store, token string) (*goshopify.Client, error) {
	app := goshopify.App{}
	return goshopify.NewClient(app, store, token, goshopify.WithVersion("2024-01"))
}
