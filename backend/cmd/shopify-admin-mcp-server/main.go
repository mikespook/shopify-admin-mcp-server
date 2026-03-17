package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mikespook/possum"
	"github.com/mikespook/shopify-admin-mcp-server/internal/mcp"
	"github.com/mikespook/shopify-admin-mcp-server/internal/shopify"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// version is set at build time via -ldflags "-X main.version=vX.Y.Z"
var version = "dev"

type Config struct {
	Store string `yaml:"store"`
	Token string `yaml:"token"`
	Addr  string `yaml:"addr"`
}

func main() {
	showVersion := flag.Bool("version", false, "Print version and exit")
	configPath := flag.String("config", "./config.yaml", "Path to YAML config file")
	addr := flag.String("addr", "localhost:9093", "Listen address")
	store := flag.String("store", "", "Shopify store domain (e.g. myshop.myshopify.com)")
	token := flag.String("token", "", "Shopify Admin API access token")
	stdio := flag.Bool("stdio", false, "Use stdio transport instead of HTTP")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	// Track which flags were explicitly set on the command line
	addrSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "addr" {
			addrSet = true
		}
	})

	// 2. Env vars
	if *store == "" {
		*store = os.Getenv("SHOPIFY_STORE")
	}
	if *token == "" {
		*token = os.Getenv("SHOPIFY_TOKEN")
	}

	// 3. Config file
	configExists := false
	cfg := Config{}
	if data, err := os.ReadFile(*configPath); err == nil {
		configExists = true
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("failed to parse config %s: %v", *configPath, err)
		}
		if *store == "" {
			*store = cfg.Store
		}
		if *token == "" {
			*token = cfg.Token
		}
		if !addrSet && cfg.Addr != "" {
			*addr = cfg.Addr
		}
	}

	// 4. Interactive prompt — only when stdin is a real TTY (skipped in stdio mode or Docker)
	if !*stdio && term.IsTerminal(int(os.Stdin.Fd())) {
		if *store == "" {
			fmt.Print("Shopify store domain: ")
			fmt.Scanln(store)
		}
		if *token == "" {
			fmt.Print("Shopify access token: ")
			b, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				log.Fatalf("failed to read token: %v", err)
			}
			*token = string(b)
		}
	}

	if *store == "" || *token == "" {
		log.Fatal("store and token are required")
	}

	// 5. Create config file if it didn't exist
	if !configExists {
		cfg = Config{Store: *store, Token: *token, Addr: *addr}
		data, err := yaml.Marshal(cfg)
		if err == nil {
			if err := os.WriteFile(*configPath, data, 0600); err != nil {
				log.Printf("warning: could not save config to %s: %v", *configPath, err)
			} else {
				log.Printf("config saved to %s", *configPath)
			}
		}
	}

	client, err := shopify.NewClient(*store, *token)
	if err != nil {
		log.Fatalf("failed to create Shopify client: %v", err)
	}

	server := mcp.NewServer(client, version)

	if *stdio {
		log.Printf("shopify-admin-mcp-server listening on stdio")
		server.ServeStdio()
		return
	}

	cors := &possum.CORSConfig{
		AllowOrigin:  "*",
		AllowMethods: "POST, OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}
	cors.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/mcp", possum.Chain(
		server.Handler,
		possum.Log,
		possum.Cors(cors),
	))

	log.Printf("shopify-admin-mcp-server listening on %s", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}
