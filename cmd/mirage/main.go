package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mirage/internal/config"
	"mirage/internal/proxy"

	"github.com/spf13/cobra"
)

func main() {
	var port int
	var configPath string

	var rootCmd = &cobra.Command{
		Use:   "mirage",
		Short: "Mirage is an API mocking gateway",
		Long:  `Mirage intercepts HTTP requests and allows mocking responses, recording traffic, and simulating network conditions.`,
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the proxy server",
		Run: func(cmd *cobra.Command, args []string) {
			addr := fmt.Sprintf(":%d", port)
			fmt.Printf("Starting Mirage proxy on %s...\n", addr)
			
			// Load config if provided
			var cfg *config.Config
			if configPath != "" {
				var err error
				cfg, err = config.LoadConfig(configPath)
				if err != nil {
					log.Fatalf("Failed to load config: %v", err)
				}
				fmt.Printf("Loaded configuration from %s (%d scenarios)\n", configPath, len(cfg.Scenarios))
			} else {
				fmt.Println("No config file specified (-c), running in pure proxy mode")
			}
			
			// Initialize proxy handler
			p := proxy.NewProxy(cfg)
			
			// Start server
			if err := http.ListenAndServe(addr, p); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		},
	}

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the proxy on")
	startCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to scenarios config file")
	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
