package tools

import (
	"context"

	"github.com/awcullen/opcua/client"
)

type ToolsManager struct {
	client *client.Client
}

func (t *ToolsManager) RegisterGetAllVariables(ctx context.Context)
