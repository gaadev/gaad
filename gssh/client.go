package gssh

import (
	"golang.org/x/crypto/ssh"
)

type Client struct {
	Network string

	// the host  to connect to
	HostPort string

	// the client config to use
	ClientConfig *ssh.ClientConfig
}

func NewClient(hostPort string, clientConfig *ssh.ClientConfig) *Client {
	return &Client{Network: "tcp", HostPort: hostPort, ClientConfig: clientConfig}
}
