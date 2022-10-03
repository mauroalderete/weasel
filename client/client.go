package client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client *ethclient.Client
}

func (c *Client) Connect(server string) error {

	client, err := ethclient.Dial(server)
	if err != nil {
		return fmt.Errorf("failed to connect to ethclient: %v", err)
	}

	c.Close()
	c.client = client
	return nil
}

func (c *Client) Close() {
	if c.client == nil {
		return
	}

	c.client.Close()
}

func (c *Client) Client() *ethclient.Client {
	return c.client
}
