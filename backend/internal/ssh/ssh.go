package ssh

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client represents an SSH client connection
type Client struct {
	config  *ssh.ClientConfig
	host    string
	port    int
	mu      sync.Mutex
	session *ssh.Session
}

// NewClient creates a new SSH client
func NewClient(host string, port int, username string, authMethods ...ssh.AuthMethod) (*Client, error) {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	return &Client{
		config: config,
		host:   host,
		port:   port,
	}, nil
}

// Connect establishes an SSH connection
func (c *Client) Connect() (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := ssh.Dial("tcp", addr, c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}
	return conn, nil
}

// Execute runs a command and returns the output
func (c *Client) Execute(ctx context.Context, command string) (string, error) {
	conn, err := c.Connect()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up output buffer
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Execute command with context
	done := make(chan error, 1)
	go func() {
		done <- session.Run(command)
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("command execution failed: %w, stderr: %s", err, stderr.String())
		}
		return stdout.String(), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// ExecuteStreaming runs a command and streams the output
func (c *Client) ExecuteStreaming(ctx context.Context, command string, outputChan chan<- string) error {
	conn, err := c.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up output pipe
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				outputChan <- string(buf[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	// Start the session
	if err := session.Start(command); err != nil {
		return err
	}

	// Wait for completion with context
	done := make(chan error, 1)
	go func() {
		done <- session.Wait()
	}()

	select {
	case err := <-done:
		close(outputChan)
		return err
	case <-ctx.Done():
		session.Signal(ssh.SIGKILL)
		return ctx.Err()
	}
}

// Dial establishes a direct TCP connection through the SSH tunnel
func (c *Client) Dial(network, addr string) (net.Conn, error) {
	conn, err := c.Connect()
	if err != nil {
		return nil, err
	}

	return conn.Dial(network, addr)
}

// Pool manages a pool of SSH connections
type Pool struct {
	clients    map[string]*Client
	maxSize    int
	mu         sync.RWMutex
	ttl        time.Duration
}

// NewPool creates a new SSH connection pool
func NewPool(maxSize int, ttl time.Duration) *Pool {
	return &Pool{
		clients: make(map[string]*Client),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// GetClient retrieves or creates an SSH client
func (p *Pool) GetClient(host string, port int, username string, authMethods ...ssh.AuthMethod) (*Client, error) {
	key := fmt.Sprintf("%s:%d:%s", host, port, username)

	p.mu.Lock()
	defer p.mu.Unlock()

	if client, exists := p.clients[key]; exists {
		return client, nil
	}

	if len(p.clients) >= p.maxSize {
		return nil, fmt.Errorf("connection pool is full")
	}

	client, err := NewClient(host, port, username, authMethods...)
	if err != nil {
		return nil, err
	}

	p.clients[key] = client
	return client, nil
}

// Release returns a client to the pool
func (p *Pool) Release(client *Client) {
	// In a real implementation, we'd track connection state
	// and recycle or close connections as needed
}

// Close closes all connections in the pool
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, client := range p.clients {
		// Close underlying connections
	}
	p.clients = make(map[string]*Client)
}

// PasswordAuth creates a password authentication method
func PasswordAuth(password string) ssh.AuthMethod {
	return ssh.Password(password)
}

// KeyAuth creates a public key authentication method
func KeyAuth(keyPath string) (ssh.AuthMethod, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}