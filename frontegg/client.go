package frontegg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"net/http"
	"os"
)

type Client struct {
	cfg   *fronteggConfig
	token string
	conn  http.Client
}

func NewClientFromEnv(ctx context.Context, d *plugin.QueryData) (*Client, error) {
	clientIdFromEnv := os.Getenv("FRONTEGG_CLIENT_ID")
	secretKeyFromEnv := os.Getenv("FRONTEGG_SECRET_KEY")
	parsedConfig := &fronteggConfig{}
	if clientIdFromEnv != "" {
		parsedConfig.ClientID = &clientIdFromEnv
	}
	if secretKeyFromEnv != "" {
		parsedConfig.SecretKey = &secretKeyFromEnv
	}

	if parsedConfig.ClientID == nil {
		return nil, errors.New("'clientId' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	if parsedConfig.SecretKey == nil {
		return nil, errors.New("'secretKey' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	client := Client{
		cfg: parsedConfig,
		// TODO: set a meaningful timeout
		conn: http.Client{},
	}
	// TODO: cache the token somehow
	err := client.GetToken()
	if err != nil {
		return nil, err
	}
	return &client, nil
}

type AuthRequest struct {
	ClientID  string `json:"clientId"`
	SecretKey string `json:"secret"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
}

func (c *Client) GetToken() error {
	input := AuthRequest{
		ClientID:  *c.cfg.ClientID,
		SecretKey: *c.cfg.SecretKey,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.frontegg.com/auth/vendor/", buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := c.conn.Do(req)
	if err != nil {
		return err
	}

	var output AuthResponse
	err = json.NewDecoder(res.Body).Decode(&output)
	if err != nil {
		return err
	}

	c.token = output.Token
	return nil
}

func (c Client) RawRequest(method string, path string) (*http.Response, error) {
	url := "https://api.frontegg.com/" + path
	fmt.Println("URL!", url)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.token)
	res, err := c.conn.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) RequestAsStruct(method string, path string, res *interface{}) error {
	raw, err := c.RawRequest(method, path)
	if err != nil {
		return err
	}
	err = json.NewDecoder(raw.Body).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) RequestAsInterface(method string, path string) (interface{}, error) {
	raw, err := c.RawRequest(method, path)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.NewDecoder(raw.Body).Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
