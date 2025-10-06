package tags

import (
	"context"

	"github.com/JarcauCristian/ztp-mcp/internal/server/tools"
	"github.com/mark3labs/mcp-go/mcp"
)

type Tag struct{}

func (Tag) Register() {}

type DeleteTag struct{}

func (DeleteTag) Create() mcp.Tool {
	return mcp.NewTool(
		"delete_tag",
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The name of the tag that will be created."),
		),
		mcp.WithToolAnnotation(tools.CreateToolAnnotation("Delete Tag", false, true, false, true)),
		mcp.WithDescription("Tool used to delete a specific tag from the running instance."),
	)
}

func (DeleteTag) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return nil, nil
}
