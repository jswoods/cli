package api

import (
	"bytes"
	"context"
	"net/http"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type Client struct {
	config Config
}

type Config struct {
	Address string
	Token   string
}

func DefaultConfig() Config {
	config := Config{
		Address: "",
		Token:   "",
	}

	return config
}

func (c *Client) Address() string {
	return c.config.Address
}

func (c *Client) Token() string {
	return c.config.Token
}

func NewClient(config Config) (*Client, error) {
	return &Client{
		config: config,
	}, nil
}

func (c *Client) newGraphQLClient() *graphql.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.config.Token},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	return graphql.NewClient(c.Address()+"graphql", httpClient)
}

func (c *Client) Query(query interface{}, variables map[string]interface{}) error {
	client := c.newGraphQLClient()
	graphqlErr := client.Query(context.Background(), query, variables)
	return graphqlErr
}

func (c *Client) Mutate(mutation interface{}, variables map[string]interface{}) error {
	client := c.newGraphQLClient()
	graphqlErr := client.Mutate(context.Background(), mutation, variables)
	return graphqlErr
}

func (c *Client) httpGET(path string) (*http.Response, error) {
	url := c.Address() + path
	req, reqErr := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Authorization", "Bearer "+c.Token())
	req.Header.Set("Accept", "application/json")
	var client = &http.Client{}

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func (c *Client) HttpPOST(path string, jsonStr *bytes.Buffer) (*http.Response, error) {
	url := c.Address() + path
	req, reqErr := http.NewRequest("POST", url, jsonStr)
	req.Header.Set("Authorization", "Bearer "+c.config.Token)
	req.Header.Set("Content-Type", "application/json")
	var client = &http.Client{}

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func optBoolArg(v *bool) *graphql.Boolean {
	var argPtr *graphql.Boolean
	if v != nil {
		argPtr = graphql.NewBoolean(graphql.Boolean(*v))
	}
	return argPtr
}

func optStringArg(v *string) *graphql.String {
	var argPtr *graphql.String
	if v != nil {
		argPtr = graphql.NewString(graphql.String(*v))
	}
	return argPtr
}
