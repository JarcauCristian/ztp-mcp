package resources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
)

func CreateAvailableHosts() mcp.Resource {
	return mcp.NewResource(
		"vm-hosts://available",
		"Available VMs Hosts",
		mcp.WithResourceDescription("Returns the available VM hosts from the ZTP agent conected."),
		mcp.WithMIMEType("application/json"),
	)
}

func HandleAvailableHosts(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	apiUrl := fmt.Sprintf("%s/MAAS/api/2.0/vm-hosts/", os.Getenv("MAAS_BASE_URL"))

	response, err := http.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("MAAS API error: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %w", err)
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(body),
		},
	}, nil
}
