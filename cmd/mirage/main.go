package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mirage/internal/config"
	"mirage/internal/proxy"
	"mirage/internal/recorder"

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
			p := proxy.NewProxy(cfg, nil)
			
			// Start server
			if err := http.ListenAndServe(addr, p); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		},
	}

	var outputFile string
	var recordCmd = &cobra.Command{
		Use:   "record",
		Short: "Start proxy in recording mode",
		Run: func(cmd *cobra.Command, args []string) {
			addr := fmt.Sprintf(":%d", port)
			fmt.Printf("Starting Mirage recorder on %s, saving to %s...\n", addr, outputFile)
			
			rec := recorder.NewRecorder(outputFile)
			p := proxy.NewProxy(nil, rec) // No config in record mode? Or maybe allow both? For now clean slate.
			
			if err := http.ListenAndServe(addr, p); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		},
	}
	
	recordCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the proxy on")
	recordCmd.Flags().StringVarP(&outputFile, "output", "o", "traffic.json", "Output file for recorded traffic")

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the proxy on")
	startCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to scenarios config file")
	
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(recordCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
