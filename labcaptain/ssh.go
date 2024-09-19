package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Function to run a remote command via SSH using the system's SSH agent
func runCommandOnServer(host, command string) error {
	return runCommandOnServerWithBuffer(host, command, nil, nil)
}
func runCommandOnServerWithBuffer(host, command string, stdoutBuffer, stderrBuffer *bytes.Buffer) error {
	// Connect to the local SSH agent
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return fmt.Errorf("could not connect to SSH agent: %v", err)
	}
	defer sshAgent.Close()

	// Create an SSH client configuration using the agent
	agentClient := agent.NewClient(sshAgent)
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Ignore host key verification for simplicity
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	if stdoutBuffer != nil {
		session.Stdout = stdoutBuffer
	} else {
		session.Stdout = os.Stdout
	}
	if stderrBuffer != nil {
		session.Stderr = stderrBuffer
	} else {
		session.Stderr = os.Stderr
	}

	// Run the command
	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	return nil
}
