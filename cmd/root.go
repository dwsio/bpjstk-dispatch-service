package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/cmd/server"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"

	"github.com/spf13/cobra"
)

var priority string
var channel string

// Initialize command-line flags
func init() {
	runCmd.Flags().StringVarP(&channel, "channel", "c", "", "Channel name to be handled (required)")
	runCmd.PersistentFlags().StringVarP(&priority, "priority", "p", "normal", "Set the priority")
}

// Define the root command
var rootCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Notification dispatcher service.",
	Long:  "Notification dispatcher service.",
}

// Define the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dispatch.",
	Long:  "Print the version number of Dispatch.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dispatch Service v1.0")
	},
}

// Define the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Dispatch service.",
	Long:  "Run the Dispatch service.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration from OS environment
		cfg := config.LoadConfigFromOS()

		// Set priority in configuration
		cfg.Project.Priority = priority

		// Get the IPv4 address of the server
		ip, err := getCurrentIPv4()
		if err != nil {
			log.Fatalf("Failed to get IP: %v", err)
		}
		cfg.Project.ServerIP = ip

		// Create a new server instance
		s := server.NewServer(cfg)

		// Define valid channels
		channels := []string{
			constants.CHANNEL_JMO,
			constants.CHANNEL_SMILE,
			constants.CHANNEL_SIPP,
			constants.CHANNEL_SIDIA,
			constants.CHANNEL_PERISAI,
		}

		// Find the selected channel and start the server
		found := false
		for _, c := range channels {
			if strings.EqualFold(channel, c) {
				if err := s.Run(c, priority); err != nil {
					log.Fatalf("Failed to run server: %v", err)
				}

				found = true
				break
			}
		}

		// If the selected channel is invalid
		if !found {
			log.Fatalf("Invalid channel: %s", channel)
		}
	},
}

// Initialize root command and its subcommands
func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCmd)
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}

// Function to retrieve the IPv4 address of the server
func getCurrentIPv4() (string, error) {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip = ipNet.IP.String()
		}
	}
	if ip == "" {
		return "", errors.New("no IP address found")
	}
	return ip, nil
}
