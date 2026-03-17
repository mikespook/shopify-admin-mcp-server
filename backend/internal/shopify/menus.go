package shopify

import (
	"context"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

// Menu represents a Shopify navigation menu.
// go-shopify does not implement the Navigation API, so we use raw HTTP calls.
type Menu struct {
	ID     uint64     `json:"id"`
	Handle string     `json:"handle"`
	Title  string     `json:"title"`
	Items  []MenuItem `json:"items"`
}

type MenuItem struct {
	ID       uint64     `json:"id"`
	Title    string     `json:"title"`
	URL      string     `json:"url"`
	Type     string     `json:"type"`
	Items    []MenuItem `json:"items,omitempty"`
}

type menusResource struct {
	Menus []Menu `json:"menus"`
}

type menuResource struct {
	Menu *Menu `json:"menu"`
}

func ListMenus(ctx context.Context, client *goshopify.Client, _ map[string]any) (string, error) {
	resource := new(menusResource)
	if err := client.Get(ctx, "menus.json", resource, nil); err != nil {
		return "", err
	}
	return toJSON(resource.Menus)
}

func GetMenu(ctx context.Context, client *goshopify.Client, args map[string]any) (string, error) {
	handle, ok := args["handle"].(string)
	if !ok || handle == "" {
		return "", fmt.Errorf("missing required argument: handle")
	}
	resource := new(menuResource)
	if err := client.Get(ctx, fmt.Sprintf("menus/%s.json", handle), resource, nil); err != nil {
		return "", err
	}
	return toJSON(resource.Menu)
}
