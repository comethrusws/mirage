package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mirage/internal/proxy"

	"github.com/spf13/cobra"
)

func main() {
	var port int

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
			
			// Initialize proxy handler
			p := proxy.NewProxy()
			
			// Start server
			if err := http.ListenAndServe(addr, p); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		},
	}

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the proxy on")
	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
