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
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(startCmd)
	serverCmd.AddCommand(serverAddCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverRemoveCmd)
	serverCmd.AddCommand(serverEnableCmd)
	serverCmd.AddCommand(serverDisableCmd)
	serverCmd.AddCommand(serverSetupPodmanCmd)
	serverCmd.AddCommand(serverSetupPrometheusCmd)
	generateCmd.AddCommand(generatePrometheusConfigCmd)
	generatePrometheusConfigCmd.Flags().BoolP("save", "s", false, "Save the config to /etc/prometheus/prometheus.yml")
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
		server, serverErr := GetServerByIP(args[0])
		if serverErr != nil {
			fmt.Println("Failed to get server with IP " + args[0])
			return
		}
		err := SetupPodman(args[0] + ":22")
		if err != nil {
			fmt.Println("Failed to setup podman on server with IP " + cmd.Flag("ip").Value.String())
			fmt.Println(err.Error())
			return
		} else {
			err := server.SetPodmanInstalled(true)
			if err != nil {
				fmt.Println("Failed to save podman installed status of server")
				fmt.Println(err.Error())
			} else {
				fmt.Println("Podman setup successfully")
			}
		}
	},
}

var serverSetupPrometheusCmd = &cobra.Command{
	Use:   "setup-prometheus",
	Short: "Setup prometheus exporter on the server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("IP address is required")
			return
		}
		server, serverErr := GetServerByIP(args[0])
		if serverErr != nil {
			fmt.Println("Failed to get server with IP " + args[0])
			return
		}
		err := SetupPrometheusExporter(args[0] + ":22")
		if err != nil {
			fmt.Println("Failed to setup prometheus exporter on server with IP " + args[0])
			fmt.Println(err.Error())
			return
		} else {
			err := server.SetPrometheusExportedEnabled(true)
			if err != nil {
				fmt.Println("Failed to save prometheus exporter enabled status of server")
				fmt.Println(err.Error())
			} else {
				fmt.Println("Prometheus exporter setup successfully")
			}
		}
	},
}

var serverEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("IP address is required")
			return
		}
		server, serverErr := GetServerByIP(args[0])
		if serverErr != nil {
			fmt.Println("Failed to get server with IP " + args[0])
			return
		}
		err := server.Enable()
		if err != nil {
			fmt.Println("Failed to enable server with IP " + args[0])
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("Server enabled successfully")
		}
	},
}

var serverDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("IP address is required")
			return
		}
		server, serverErr := GetServerByIP(args[0])
		if serverErr != nil {
			fmt.Println("Failed to get server with IP " + args[0])
			return
		}
		err := server.Disable()
		if err != nil {
			fmt.Println("Failed to disable server with IP " + args[0])
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("Server disabled successfully")
		}
	},
}

// Generate
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate config",
	Run: func(cmd *cobra.Command, args []string) {
		// print help
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

var generatePrometheusConfigCmd = &cobra.Command{
	Use:   "prometheus-config",
	Short: "Generate prometheus config",
	Run: func(cmd *cobra.Command, args []string) {
		configContent, err := generatePrometheusConfig()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		isSave := cmd.Flag("save").Value.String() == "true"
		if isSave {
			err := os.WriteFile("/etc/prometheus/prometheus.yml", []byte(configContent), 0644)
			if err != nil {
				fmt.Println("Failed to save prometheus config to /etc/prometheus/prometheus.yml")
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Prometheus config saved to /etc/prometheus/prometheus.yml")
		} else {
			fmt.Println(configContent)
		}
	},
}

// Start API Server
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start API server and worker",
	Run: func(cmd *cobra.Command, args []string) {
		go startAPIServer()
		go processPendingLabsDeployment()
		go processExpiredLabsDeletion()
		go processAutoExpiratonOfLabs()
		<-make(chan bool)
	},
}
