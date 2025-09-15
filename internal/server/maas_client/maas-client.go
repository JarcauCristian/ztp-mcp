package maas_client

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	defaultClient *MAASClient
	once          sync.Once
	initErr       error
)

func GetClient() (*MAASClient, error) {
	once.Do(func() {
		defaultClient, initErr = NewMAASClientFromEnv()
	})

	return defaultClient, initErr
}

func MustClient() *MAASClient {
	client, err := GetClient()
	if err != nil {
		panic(fmt.Errorf("failed to initialize MAAS client: %w", err))
	}
	return client
}

func generateNonce() (string, error) {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

type MAASResponse struct {
	Body       string
	StatusCode int
	Headers    http.Header
}

type MAASClient struct {
	baseURL     string
	consumerKey string
	token       string
	secret      string
}

func NewMAASClientFromEnv() (*MAASClient, error) {
	baseURL := os.Getenv("MAAS_BASE_URL")
	apiKey := os.Getenv("MAAS_API_KEY")
	if baseURL == "" {
		return nil, fmt.Errorf("MAAS_BASE_URL environment variable not set")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("MAAS_API_KEY environment variable not set")
	}
	parts := strings.Split(apiKey, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("MAAS_API_KEY must be in the format consumer_key:token:secret")
	}
	return &MAASClient{
		baseURL:     baseURL,
		consumerKey: parts[0],
		token:       parts[1],
		secret:      parts[2],
	}, nil
}

func (c *MAASClient) Get(ctx context.Context, path string) (MAASResponse, error) {
	fullURL := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return MAASResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	nonce, err := generateNonce()
	if err != nil {
		return MAASResponse{}, err
	}

	signature := "&" + url.QueryEscape(c.secret)

	authHeader := fmt.Sprintf(
		`OAuth oauth_version="1.0", oauth_signature_method="PLAINTEXT", oauth_consumer_key="%s", oauth_token="%s", oauth_signature="%s", oauth_nonce="%s", oauth_timestamp="%s"`,
		url.QueryEscape(c.consumerKey),
		url.QueryEscape(c.token),
		url.QueryEscape(signature),
		nonce,
		timestamp,
	)

	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return MAASResponse{}, fmt.Errorf("MAAS API error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MAASResponse{}, fmt.Errorf("failed to read the response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return MAASResponse{}, fmt.Errorf("MAAS API returned status %d: %s", resp.StatusCode, string(body))
	}

	return MAASResponse{
		Body:       string(body),
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}

func (c *MAASClient) Post(ctx context.Context, path string, body io.Reader) (MAASResponse, error) {
	fullURL := fmt.Sprintf("%s%s", c.baseURL, path)

	timeoutContext, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutContext, "POST", fullURL, body)
	if err != nil {
		return MAASResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	nonce, err := generateNonce()
	if err != nil {
		return MAASResponse{}, err
	}

	signature := "&" + url.QueryEscape(c.secret)

	authHeader := fmt.Sprintf(
		`OAuth oauth_version="1.0", oauth_signature_method="PLAINTEXT", oauth_consumer_key="%s", oauth_token="%s", oauth_signature="%s", oauth_nonce="%s", oauth_timestamp="%s"`,
		url.QueryEscape(c.consumerKey),
		url.QueryEscape(c.token),
		url.QueryEscape(signature),
		nonce,
		timestamp,
	)

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return MAASResponse{}, fmt.Errorf("MAAS API error: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return MAASResponse{}, fmt.Errorf("failed to read the response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return MAASResponse{}, fmt.Errorf("MAAS API returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	return MAASResponse{
		Body:       string(responseBody),
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}
