package cloud66

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// Config is used to configure the creation of a client
type Config struct {
	// Address is the address of the Nomad agent
	Address string

	// HttpClient is the client to use. Default will be used if not provided.
	//
	// If set, it expected to be configured for tls already, and TLSConfig is ignored.
	// You may use ConfigureTLS() function to aid with initialization.
	HttpClient *http.Client

	// HttpAuth is the auth info to use for http access.
	AccessToken string

	// WaitTime limits how long a Watch will block. If not provided,
	// the agent default values will be used.
	WaitTime time.Duration

	Headers http.Header
}

// DefaultConfig returns a default configuration for the client
func DefaultConfig() *Config {
	config := &Config{
		Address: "https://app.cloud66.com/api/3",
	}
	if token := os.Getenv("CLOUD66_ACCESS_TOKEN"); token != "" {
		config.AccessToken = token
	}

	return config
}

// Client provides a client to the Cloud66 API
type Client struct {
	httpClient *http.Client
	config     Config
}

func defaultHttpClient() *http.Client {
	httpClient := cleanhttp.DefaultClient()
	transport := httpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	return httpClient
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {
	// bootstrap the config
	defConfig := DefaultConfig()

	if config.Address == "" {
		config.Address = defConfig.Address
	} else if _, err := url.Parse(config.Address); err != nil {
		return nil, fmt.Errorf("invalid address '%s': %v", config.Address, err)
	}

	httpClient := config.HttpClient
	if httpClient == nil {
		httpClient = defaultHttpClient()
	}

	client := &Client{
		config:     *config,
		httpClient: httpClient,
	}
	return client, nil
}
