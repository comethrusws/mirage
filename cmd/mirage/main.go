package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"mirage/internal/config"
	"mirage/internal/proxy"
	"mirage/internal/recorder"
	"mirage/internal/ui"

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
			
			p := proxy.NewProxy(cfg, nil)
			
			dashboard := ui.NewUI(p)
			uiHandler := dashboard.Handler()
			
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/__mirage/") {
					uiHandler.ServeHTTP(w, r)
				} else {
					p.ServeHTTP(w, r)
				}
			})
			
			if err := http.ListenAndServe(addr, handler); err != nil {
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
			p := proxy.NewProxy(nil, rec)
			
			if err := http.ListenAndServe(addr, p); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		},
	}
	
	recordCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the proxy on")
	recordCmd.Flags().StringVarP(&outputFile, "output", "o", "traffic.json", "Output file for recorded traffic")

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the proxy on")
	startCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to scenarios config file")
	
	var scenariosCmd = &cobra.Command{
		Use:   "scenarios",
		Short: "Manage scenarios",
	}

	var listCmd = &cobra.Command{
		Use:   "list [config]",
		Short: "List scenarios in a config file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig(args[0])
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}
			fmt.Printf("Scenarios in %s:\n", args[0])
			for _, s := range cfg.Scenarios {
				fmt.Printf("- %s (Matches: %s %s)\n", s.Name, s.Match.Method, s.Match.Path)
			}
		},
	}
	scenariosCmd.AddCommand(listCmd)
	
	var replayCmd = &cobra.Command{
		Use:   "replay [traffic.json]",
		Short: "Replay recorded traffic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data, err := os.ReadFile(args[0])
			if err != nil {
				log.Fatalf("Failed to read file: %v", err)
			}
			
			var interactions []recorder.Interaction
			if err := json.Unmarshal(data, &interactions); err != nil {
				log.Fatalf("Failed to parse JSON: %v", err)
			}
			
			client := &http.Client{}
			fmt.Printf("Replaying %d interactions...\n", len(interactions))
			
			for i, interaction := range interactions {
				reqData := interaction.Request
				fmt.Printf("[%d] %s %s... ", i+1, reqData.Method, reqData.URL)
				
				req, err := http.NewRequest(reqData.Method, reqData.URL, strings.NewReader(reqData.Body))
				if err != nil {
					fmt.Printf("Failed to create request: %v\n", err)
					continue
				}
				
				for k, vv := range reqData.Headers {
					for _, v := range vv {
						req.Header.Add(k, v)
					}
				}
				
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}
				resp.Body.Close()
				fmt.Printf("Status: %d\n", resp.StatusCode)
			}
		},
	}

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(recordCmd)
	rootCmd.AddCommand(scenariosCmd)
	rootCmd.AddCommand(replayCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
