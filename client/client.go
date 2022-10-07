package client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client *ethclient.Client
	url    string
}

func (c *Client) Connect(url string) error {

	c.url = url
	client, err := ethclient.Dial(c.url)
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

func (c *Client) Url() string {
	return c.url
}
