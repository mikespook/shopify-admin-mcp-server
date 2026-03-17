package mcp

func GetTools() []ToolDefinition {
	return []ToolDefinition{
		// Products
		{
			Name:        "shopify_list_products",
			Description: "List all products in the Shopify store",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"limit": {Type: "integer", Description: "Maximum number of products to return (default 50, max 250)"},
					"page_info": {Type: "string", Description: "Pagination cursor for the next page"},
				},
			},
		},
		{
			Name:        "shopify_get_product",
			Description: "Get a single product by ID",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"id": {Type: "integer", Description: "Product ID"},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "shopify_create_product",
			Description: "Create a new product in the Shopify store",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"title":        {Type: "string", Description: "Product title"},
					"body_html":    {Type: "string", Description: "Product description (HTML)"},
					"vendor":       {Type: "string", Description: "Product vendor"},
					"product_type": {Type: "string", Description: "Product type"},
					"tags":         {Type: "string", Description: "Comma-separated list of tags"},
					"status":       {Type: "string", Description: "Product status: active, draft, or archived"},
				},
				Required: []string{"title"},
			},
		},
		{
			Name:        "shopify_update_product",
			Description: "Update an existing product",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"id":           {Type: "integer", Description: "Product ID"},
					"title":        {Type: "string", Description: "Product title"},
					"body_html":    {Type: "string", Description: "Product description (HTML)"},
					"vendor":       {Type: "string", Description: "Product vendor"},
					"product_type": {Type: "string", Description: "Product type"},
					"tags":         {Type: "string", Description: "Comma-separated list of tags"},
					"status":       {Type: "string", Description: "Product status: active, draft, or archived"},
				},
				Required: []string{"id"},
			},
		},
		// Collections
		{
			Name:        "shopify_list_collections",
			Description: "List all custom collections in the Shopify store",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"limit": {Type: "integer", Description: "Maximum number of collections to return (default 50, max 250)"},
				},
			},
		},
		{
			Name:        "shopify_get_collection",
			Description: "Get a single custom collection by ID",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"id": {Type: "integer", Description: "Collection ID"},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "shopify_create_collection",
			Description: "Create a new custom collection",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"title":    {Type: "string", Description: "Collection title"},
					"body_html": {Type: "string", Description: "Collection description (HTML)"},
					"handle":   {Type: "string", Description: "URL-friendly handle"},
				},
				Required: []string{"title"},
			},
		},
		{
			Name:        "shopify_update_collection",
			Description: "Update an existing custom collection",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"id":       {Type: "integer", Description: "Collection ID"},
					"title":    {Type: "string", Description: "Collection title"},
					"body_html": {Type: "string", Description: "Collection description (HTML)"},
					"handle":   {Type: "string", Description: "URL-friendly handle"},
				},
				Required: []string{"id"},
			},
		},
		// Pages
		{
			Name:        "shopify_list_pages",
			Description: "List all pages in the Shopify store",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"limit": {Type: "integer", Description: "Maximum number of pages to return (default 50, max 250)"},
				},
			},
		},
		{
			Name:        "shopify_get_page",
			Description: "Get a single page by ID",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"id": {Type: "integer", Description: "Page ID"},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "shopify_create_page",
			Description: "Create a new page in the Shopify store",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"title":     {Type: "string", Description: "Page title"},
					"body_html": {Type: "string", Description: "Page content (HTML)"},
					"handle":    {Type: "string", Description: "URL-friendly handle"},
					"published": {Type: "boolean", Description: "Whether the page is published"},
				},
				Required: []string{"title"},
			},
		},
		{
			Name:        "shopify_update_page",
			Description: "Update an existing page",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"id":        {Type: "integer", Description: "Page ID"},
					"title":     {Type: "string", Description: "Page title"},
					"body_html": {Type: "string", Description: "Page content (HTML)"},
					"handle":    {Type: "string", Description: "URL-friendly handle"},
					"published": {Type: "boolean", Description: "Whether the page is published"},
				},
				Required: []string{"id"},
			},
		},
		// Menus
		{
			Name:        "shopify_list_menus",
			Description: "List all navigation menus in the Shopify store",
			InputSchema: InputSchema{
				Type: "object",
			},
		},
		{
			Name:        "shopify_get_menu",
			Description: "Get a single navigation menu by handle",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"handle": {Type: "string", Description: "Menu handle (e.g. main-menu, footer)"},
				},
				Required: []string{"handle"},
			},
		},
	}
}
