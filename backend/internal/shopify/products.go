package shopify

import (
	"context"
	"encoding/json"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func ListProducts(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	opts := goshopify.ListOptions{}
	if v, ok := args["limit"]; ok {
		if l, ok := toInt(v); ok {
			opts.Limit = l
		}
	}
	products, err := client.Product.List(ctx, opts)
	if err != nil {
		return "", err
	}
	return toJSON(products)
}

func GetProduct(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	id, err := requireUint64(args, "id")
	if err != nil {
		return "", err
	}
	product, err := client.Product.Get(ctx, id, nil)
	if err != nil {
		return "", err
	}
	return toJSON(product)
}

func CreateProduct(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	p := goshopify.Product{}
	if v, ok := args["title"].(string); ok {
		p.Title = v
	}
	if v, ok := args["body_html"].(string); ok {
		p.BodyHTML = v
	}
	if v, ok := args["vendor"].(string); ok {
		p.Vendor = v
	}
	if v, ok := args["product_type"].(string); ok {
		p.ProductType = v
	}
	if v, ok := args["tags"].(string); ok {
		p.Tags = v
	}
	if v, ok := args["status"].(string); ok {
		p.Status = goshopify.ProductStatus(v)
	}
	product, err := client.Product.Create(ctx, p)
	if err != nil {
		return "", err
	}
	return toJSON(product)
}

func UpdateProduct(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	id, err := requireUint64(args, "id")
	if err != nil {
		return "", err
	}
	p := goshopify.Product{Id: id}
	if v, ok := args["title"].(string); ok {
		p.Title = v
	}
	if v, ok := args["body_html"].(string); ok {
		p.BodyHTML = v
	}
	if v, ok := args["vendor"].(string); ok {
		p.Vendor = v
	}
	if v, ok := args["product_type"].(string); ok {
		p.ProductType = v
	}
	if v, ok := args["tags"].(string); ok {
		p.Tags = v
	}
	if v, ok := args["status"].(string); ok {
		p.Status = goshopify.ProductStatus(v)
	}
	product, err := client.Product.Update(ctx, p)
	if err != nil {
		return "", err
	}
	return toJSON(product)
}

// helpers

func toJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func requireUint64(args map[string]any, key string) (uint64, error) {
	v, ok := args[key]
	if !ok {
		return 0, fmt.Errorf("missing required argument: %s", key)
	}
	n, ok := toInt(v)
	if !ok || n <= 0 {
		return 0, fmt.Errorf("argument %s must be a positive integer", key)
	}
	return uint64(n), nil
}

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	}
	return 0, false
}
