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

type RequestType int

const (
	RequestTypeGet RequestType = iota
	RequestTypePost
	RequestTypePut
	RequestTypeDelete
)

var requestTypeName = map[RequestType]string{
	RequestTypeGet:    "GET",
	RequestTypePost:   "POST",
	RequestTypePut:    "PUT",
	RequestTypeDelete: "DELETE",
}

func (rt RequestType) String() string {
	return requestTypeName[rt]
}

func (rt RequestType) Headers() map[string]string {
	switch rt {
	case RequestTypeGet:
		return nil
	case RequestTypePost:
		return map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	case RequestTypePut:
		return map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	case RequestTypeDelete:
		return map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	default:
		return nil
	}
}

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

func (c *MAASClient) Do(ctx context.Context, requestType RequestType, path string, body io.Reader) (string, error) {
	fullURL := fmt.Sprintf("%s%s", c.baseURL, path)

	timeoutContext, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutContext, requestType.String(), fullURL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	nonce, err := generateNonce()
	if err != nil {
		return "", err
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

	if requestType.Headers() != nil {
		for key, value := range requestType.Headers() {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("MAAS API error: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("MAAS API returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	return string(responseBody), nil
}
