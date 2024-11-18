package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

var (
	sshMapMutex    *sync.Mutex      = &sync.Mutex{}
	sshHostPoolMap map[string]*Pool = make(map[string]*Pool)
)

// Function to run a remote command via SSH using the system's SSH agent
func runCommandOnServer(host, command string) error {
	return runCommandOnServerWithBuffer(host, command, nil, nil)
}
func runCommandOnServerWithBuffer(host, command string, stdoutBuffer, stderrBuffer *bytes.Buffer) error {
	var p *Pool

	sshMapMutex.Lock()
	p = sshHostPoolMap[host]
	if p == nil {
		p, err := NewPool(host)
		if err != nil {
			return err
		}
		sshHostPoolMap[host] = p
	}
	sshMapMutex.Unlock()
	c, err := p.borrowClient()
	if err != nil {
		return err
	}
	defer p.returnClient(c)

	return runCommand(c, command, stdoutBuffer, stderrBuffer)
}

func runCommand(sc *sshClient, command string, stdoutBuffer, stderrBuffer *bytes.Buffer) error {
	// Create a session
	session, err := sc.client.NewSession()
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
