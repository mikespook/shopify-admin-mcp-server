package shopify

import (
	"context"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func ListCollections(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	opts := goshopify.ListOptions{}
	if v, ok := args["limit"]; ok {
		if l, ok := toInt(v); ok {
			opts.Limit = l
		}
	}
	collections, err := client.CustomCollection.List(ctx, opts)
	if err != nil {
		return "", err
	}
	return toJSON(collections)
}

func GetCollection(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	id, err := requireUint64(args, "id")
	if err != nil {
		return "", err
	}
	collection, err := client.CustomCollection.Get(ctx, id, nil)
	if err != nil {
		return "", err
	}
	return toJSON(collection)
}

func CreateCollection(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	c := goshopify.CustomCollection{}
	if v, ok := args["title"].(string); ok {
		c.Title = v
	}
	if v, ok := args["body_html"].(string); ok {
		c.BodyHTML = v
	}
	if v, ok := args["handle"].(string); ok {
		c.Handle = v
	}
	collection, err := client.CustomCollection.Create(ctx, c)
	if err != nil {
		return "", err
	}
	return toJSON(collection)
}

func UpdateCollection(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	id, err := requireUint64(args, "id")
	if err != nil {
		return "", err
	}
	c := goshopify.CustomCollection{Id: id}
	if v, ok := args["title"].(string); ok {
		c.Title = v
	}
	if v, ok := args["body_html"].(string); ok {
		c.BodyHTML = v
	}
	if v, ok := args["handle"].(string); ok {
		c.Handle = v
	}
	collection, err := client.CustomCollection.Update(ctx, c)
	if err != nil {
		return "", err
	}
	return toJSON(collection)
}
