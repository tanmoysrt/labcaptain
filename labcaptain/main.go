package main

import (
	"log"
	"os"
)

func main() {
	// if user is not root, exit
	if os.Geteuid() != 0 {
		log.Fatal("This program must be run as root")
	}

	// check for SSH_AUTH_SOCK env variable
	if _, ok := os.LookupEnv("SSH_AUTH_SOCK"); !ok {
		log.Fatal("SSH_AUTH_SOCK environment variable not set")
	}

	initiateDB()
	// Run the root command
	rootCmd.Execute()
}
