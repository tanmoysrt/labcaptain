package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	MAX_CLIENTS       = 10
	IDLE_TIMEOUT      = time.Duration(1 * time.Second)
	POOL_WAIT_TIMEOUT = time.Duration(1 * time.Second)
)

// ErrPoolClosed is thrown when a closed Pool is used.
var ErrPoolClosed = errors.New("pool closed")

type sshClient struct {
	client  *ssh.Client
	lastUse time.Time
}

type Pool struct {
	clients         chan *sshClient
	stopBorrow      chan bool
	closed          atomic.Bool
	createdClients  atomic.Int32
	host            string
	config          *ssh.ClientConfig
	maxClients      int
	idleTimeout     time.Duration
	poolWaitTimeout time.Duration
}

func NewPool(host string) (*Pool, error) {
	// Connect to the local SSH agent
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, fmt.Errorf("could not connect to SSH agent: %v", err)
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

	p := &Pool{
		clients:         make(chan *sshClient, MAX_CLIENTS),
		stopBorrow:      make(chan bool),
		host:            host,
		config:          config,
		maxClients:      MAX_CLIENTS,
		idleTimeout:     IDLE_TIMEOUT,
		poolWaitTimeout: POOL_WAIT_TIMEOUT,
	}

	go p.sweepClients(time.Second * 2)

	return p, nil
}

// Close closes the pool.
func (p *Pool) Close() {
	p.closed.Store(true)
	close(p.stopBorrow)

	// If the sweeper isn't already running, run it.
	if p.idleTimeout.Seconds() <= 1 {
		p.sweepClients(time.Second * 1)
	}
}

func (p *Pool) newClient() (*sshClient, error) {
	// Connect to the SSH server
	client, err := ssh.Dial("tcp", p.host, p.config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	sshClient := &sshClient{
		client:  client,
		lastUse: time.Now(),
	}
	return sshClient, nil
}

func (p *Pool) borrowClient() (*sshClient, error) {
	switch {
	case p.closed.Load():
		return nil, ErrPoolClosed
	case (int(p.createdClients.Load()) < p.maxClients) && (len(p.clients) == 0):
		p.createdClients.Add(1)
		cl, err := p.newClient()
		if err != nil {
			// Decrement counter on failed connection creation.
			p.createdClients.Add(-1)
			return nil, err
		}
		return cl, nil
	}

	select {
	case c := <-p.clients:
		return c, nil
	case <-p.stopBorrow:
		return nil, ErrPoolClosed
	case <-time.After(p.poolWaitTimeout):
		return nil, errors.New("timed out waiting for free conn in pool")
	}

}

func (p *Pool) returnClient(c *sshClient) (err error) {

	defer func() {
		if err != nil {
			p.createdClients.Add(-1)
			c.client.Close()
		}
	}()

	select {
	case p.clients <- c:
		return nil
	case <-time.After(p.poolWaitTimeout):
		return errors.New("timed out returning connection to pool")
	case <-p.stopBorrow:
		return ErrPoolClosed
	}

}

func (p *Pool) sweepClients(interval time.Duration) {
	activeClients := make([]*sshClient, cap(p.clients))

	for {
		<-time.After(interval)
		var (
			num            = len(p.clients)
			createdClients = p.createdClients.Load()
			closed         = p.closed.Load()
		)
		if closed && createdClients == 0 {
			return
		}
		for i := 0; i < num; i++ {
			var c *sshClient
			// Pick a connection to check from the pool.
			select {
			case c = <-p.clients:
			default:
				continue
			}

			if closed || time.Since(c.lastUse) > p.idleTimeout {
				p.createdClients.Add(-1)
				c.client.Close()
				continue
			}

			activeClients = append(activeClients, c)
		}

		for _, c := range activeClients {
			select {
			case p.clients <- c:
			default:
				c.client.Close()
				p.createdClients.Add(-1)
			}
		}
	}
}
