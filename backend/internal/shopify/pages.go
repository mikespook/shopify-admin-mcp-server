package shopify

import (
	"context"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func ListPages(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	opts := goshopify.ListOptions{}
	if v, ok := args["limit"]; ok {
		if l, ok := toInt(v); ok {
			opts.Limit = l
		}
	}
	pages, err := client.Page.List(ctx, opts)
	if err != nil {
		return "", err
	}
	return toJSON(pages)
}

func GetPage(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	id, err := requireUint64(args, "id")
	if err != nil {
		return "", err
	}
	page, err := client.Page.Get(ctx, id, nil)
	if err != nil {
		return "", err
	}
	return toJSON(page)
}

func CreatePage(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	p := goshopify.Page{}
	if v, ok := args["title"].(string); ok {
		p.Title = v
	}
	if v, ok := args["body_html"].(string); ok {
		p.BodyHTML = v
	}
	if v, ok := args["handle"].(string); ok {
		p.Handle = v
	}
	if v, ok := args["published"].(bool); ok {
		p.Published = &v
	}
	page, err := client.Page.Create(ctx, p)
	if err != nil {
		return "", err
	}
	return toJSON(page)
}

func UpdatePage(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	id, err := requireUint64(args, "id")
	if err != nil {
		return "", err
	}
	p := goshopify.Page{Id: id}
	if v, ok := args["title"].(string); ok {
		p.Title = v
	}
	if v, ok := args["body_html"].(string); ok {
		p.BodyHTML = v
	}
	if v, ok := args["handle"].(string); ok {
		p.Handle = v
	}
	if v, ok := args["published"].(bool); ok {
		p.Published = &v
	}
	page, err := client.Page.Update(ctx, p)
	if err != nil {
		return "", err
	}
	return toJSON(page)
}
