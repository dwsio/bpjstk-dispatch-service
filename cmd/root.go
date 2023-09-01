package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/config"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/server"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"

	"github.com/spf13/cobra"
)

var priority string
var channel string

func init() {
	runCmd.Flags().StringVarP(&channel, "channel", "c", "", "Channel name to be handled (required)")
	runCmd.PersistentFlags().StringVarP(&priority, "priority", "p", "normal", "Set the priority")
}

var rootCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Notification dispatcher service.",
	Long:  "Notification dispatcher service.",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dispatch.",
	Long:  "Print the version number of Dispatch.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dispatch Service v1.0")
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Dispatch service.",
	Long:  "Run the Dispatch service.",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.LoadConfigFromOS()

		cfg.Project.Priority = priority

		ip, err := getCurrentIPv4()
		if err != nil {
			log.Fatalf("Failed to get IP: %v", err)
		}
		cfg.Project.ServerIP = ip

		s := server.NewServer(cfg)

		channels := []string{
			constants.CHANNEL_JMO,
			constants.CHANNEL_SMILE,
			constants.CHANNEL_SIPP,
			constants.CHANNEL_SIDIA,
			constants.CHANNEL_PERISAI,
		}

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

		if !found {
			log.Fatalf("Invalid channel: %s", channel)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd, runCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}

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
