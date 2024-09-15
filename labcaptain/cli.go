package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "labcaptain",
	Short: "LabCaptain is a daemon + CLI tool to deploy any kind of lab environment in cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// print help
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

func Execute() {
	cobra.EnableCommandSorting = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverAddCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverRemoveCmd)
	serverCmd.AddCommand(serverSetupPodmanCmd)
}

// Server
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage your deployment servers",
	Run: func(cmd *cobra.Command, args []string) {
		// print help
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

var serverAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("IP address is required")
			return
		}
		server := Server{
			IP: args[0],
		}
		err := server.Create()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Server created successfully")
		fmt.Println("Server IP: " + server.IP)
	},
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all servers",
	Run: func(cmd *cobra.Command, args []string) {
		servers, err := GetAllServers()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
		fmt.Fprintln(w, "IP ADDRESS\tPODMAN INSTALLED\tPROMETHEUS EXPORTED\tENABLED")
		for _, server := range servers {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%t\t%t\t%t", server.IP, server.PodmanInstalled, server.PrometheusExportedEnabled, server.Enabled))
		}
		w.Flush()
	},
}

var serverRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("IP address is required")
			return
		}
		err := DeleteServer(args[0])
		if err != nil {
			fmt.Println("Failed to drop server with IP " + cmd.Flag("ip").Value.String())
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Server removed successfully")
	},
}

var serverSetupPodmanCmd = &cobra.Command{
	Use:   "setup-podman",
	Short: "Setup podman on the server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("IP address is required")
			return
		}
		err := SetupPodman(args[0] + ":22")
		if err != nil {
			fmt.Println("Failed to setup podman on server with IP " + cmd.Flag("ip").Value.String())
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("Podman setup successfully")
		}
	},
}
