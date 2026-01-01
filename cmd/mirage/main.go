package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"mirage/internal/config"
	"mirage/internal/logger"
	"mirage/internal/proxy"
	"mirage/internal/recorder"
	"mirage/internal/ui"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

func main() {
	var port int
	var configPath string
	var noBrowser bool

	var rootCmd = &cobra.Command{
		Use:     "mirage",
		Short:   "Mirage is an API mocking gateway",
		Long:    `Mirage intercepts HTTP requests and allows mocking responses, recording traffic, and simulating network conditions.`,
		Version: version,
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the proxy server",
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintBanner(version)
			
			addr := fmt.Sprintf(":%d", port)
			dashboardURL := fmt.Sprintf("http://localhost:%d/__mirage/", port)
			
			var cfg *config.Config
			if configPath != "" {
				var err error
				cfg, err = config.LoadConfig(configPath)
				if err != nil {
					logger.LogError(fmt.Sprintf("Failed to load config: %v", err))
					os.Exit(1)
				}
				logger.LogSuccess(fmt.Sprintf("Loaded %d scenarios from %s", len(cfg.Scenarios), configPath))
			} else {
				logger.LogInfo("No config specified, running in pure proxy mode")
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
			
			logger.LogSuccess(fmt.Sprintf("Server started on %s", addr))
			logger.LogInfo(fmt.Sprintf("Dashboard: %s", dashboardURL))
			fmt.Println()
			
			if !noBrowser {
				go browser.OpenURL(dashboardURL)
			}
			
			if err := http.ListenAndServe(addr, handler); err != nil {
				logger.LogError(fmt.Sprintf("Server failed: %v", err))
				os.Exit(1)
			}
		},
	}

	var outputFile string
	var recordCmd = &cobra.Command{
		Use:   "record",
		Short: "Start proxy in recording mode",
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintBanner(version)
			
			addr := fmt.Sprintf(":%d", port)
			
			rec := recorder.NewRecorder(outputFile)
			p := proxy.NewProxy(nil, rec)
			
			logger.LogSuccess(fmt.Sprintf("Recording started on %s", addr))
			logger.LogInfo(fmt.Sprintf("Saving to %s", outputFile))
			fmt.Println()
			
			if err := http.ListenAndServe(addr, p); err != nil {
				logger.LogError(fmt.Sprintf("Server failed: %v", err))
				os.Exit(1)
			}
		},
	}
	
	recordCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the proxy on")
	recordCmd.Flags().StringVarP(&outputFile, "output", "o", "traffic.json", "Output file for recorded traffic")

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the proxy on")
	startCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to scenarios config file")
	startCmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Don't open browser automatically")
	
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
				logger.LogError(fmt.Sprintf("Failed to load config: %v", err))
				os.Exit(1)
			}
			logger.LogSuccess(fmt.Sprintf("Scenarios in %s:", args[0]))
			for _, s := range cfg.Scenarios {
				fmt.Printf("  • %s (%s %s)\n", s.Name, s.Match.Method, s.Match.Path)
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
				logger.LogError(fmt.Sprintf("Failed to read file: %v", err))
				os.Exit(1)
			}
			
			var interactions []recorder.Interaction
			if err := json.Unmarshal(data, &interactions); err != nil {
				logger.LogError(fmt.Sprintf("Failed to parse JSON: %v", err))
				os.Exit(1)
			}
			
			client := &http.Client{}
			logger.LogInfo(fmt.Sprintf("Replaying %d interactions...", len(interactions)))
			
			for i, interaction := range interactions {
				reqData := interaction.Request
				fmt.Printf("[%d] %s %s... ", i+1, reqData.Method, reqData.URL)
				
				req, err := http.NewRequest(reqData.Method, reqData.URL, strings.NewReader(reqData.Body))
				if err != nil {
					logger.LogError(fmt.Sprintf("Failed to create request: %v", err))
					continue
				}
				
				for k, vv := range reqData.Headers {
					for _, v := range vv {
						req.Header.Add(k, v)
					}
				}
				
				resp, err := client.Do(req)
				if err != nil {
					logger.LogError(err.Error())
					continue
				}
				resp.Body.Close()
				logger.LogSuccess(fmt.Sprintf("Status: %d", resp.StatusCode))
			}
		},
	}

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(recordCmd)
	rootCmd.AddCommand(scenariosCmd)
	rootCmd.AddCommand(replayCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.LogError(err.Error())
		os.Exit(1)
	}
}
