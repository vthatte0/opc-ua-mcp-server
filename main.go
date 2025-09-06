package main

import (
	"context"
	"log"

	"github.com/awcullen/opcua/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/vthatte/opc-ua-mcp-server/opc"
)

func main() {
	ctx := context.Background()

	// Get available servers
	client, err := client.Dial(
		ctx,
		"opc.tcp://localhost:4840",
		client.WithInsecureSkipVerify(),
	)

	if err != nil {
		log.Panic(err)
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"Opc-UA",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Get all variables tool
	allVarsTool := mcp.NewTool("findAllOpcVars",
		mcp.WithDescription("Get all available variables from the OPC UA server"),
	)

	// Add the tool to get all variables
	s.AddTool(allVarsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resString, err := opc.GetNodeIDsAndBrowseNames(ctx, client)
		if err != nil {
			return mcp.NewToolResultError("failed to get opc browse names"), nil
		}

		return mcp.NewToolResultText(resString), nil
	})

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		log.Panic(err)
	}

}
