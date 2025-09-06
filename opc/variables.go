package opc

import (
	"context"
	"fmt"
	"strings"

	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
)

type NodeInfo struct {
	NodeID     string
	BrowseName string
}

func GetNodeIDsAndBrowseNames(ctx context.Context, client *client.Client) (string, error) {
	objectsID := ua.ParseNodeID("ns=0;i=85") // Objects folder- todo: make root configurable
	results, err := recursiveBrowse(ctx, client, objectsID)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "No nodes found in the OPC UA server.", nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d nodes:\n", len(results))
	for _, n := range results {
		fmt.Fprintf(&sb, "- NodeID: %s, BrowseName: %s\n", n.NodeID, n.BrowseName)
	}
	return sb.String(), nil
}

func recursiveBrowse(ctx context.Context, client *client.Client, nodeID ua.NodeID) ([]NodeInfo, error) {
	varNodeInfo := []NodeInfo{}

	browseReq := &ua.BrowseRequest{
		View: ua.ViewDescription{},
		NodesToBrowse: []ua.BrowseDescription{
			{
				NodeID:          nodeID,
				BrowseDirection: ua.BrowseDirectionForward,
				ReferenceTypeID: ua.ReferenceTypeIDHierarchicalReferences,
				IncludeSubtypes: true,
			},
		},
	}

	browseResp, err := client.Browse(ctx, browseReq)
	if err != nil || len(browseResp.Results) == 0 {
		return nil, nil
	}

	for _, ref := range browseResp.Results[0].References {
		targetID := ref.NodeID.NodeID
		if targetID == nil {
			continue
		}

		if ref.NodeClass == ua.NodeClassVariable {
			varNodeInfo = append(varNodeInfo, NodeInfo{
				NodeID:     fmt.Sprint(ref.NodeID.NodeID),
				BrowseName: ref.BrowseName.String(),
			})
		}

		// Recurse if Object (keep exploring deeper)
		if ref.NodeClass == ua.NodeClassObject {
			childVarNodes, err := recursiveBrowse(ctx, client, ref.NodeID.NodeID)
			if err != nil {
				return nil, err
			}

			varNodeInfo = append(varNodeInfo, childVarNodes...)
		}
	}

	// TODO: add support for handling continuation points
	// Handle continuation points
	// cp := browseResp.Results[0].ContinuationPoint
	// for cp != nil && len(cp) > 0 {
	// 	nextReq := &ua.BrowseNextRequest{
	// 		ReleaseContinuationPoints: false,
	// 		ContinuationPoints:        [][]byte{cp},
	// 	}
	// 	nextResp, err := client.BrowseNext(ctx, nextReq)
	// 	if err != nil || len(nextResp.Results) == 0 {
	// 		break
	// 	}
	// 	for _, ref := range nextResp.Results[0].References {
	// 		if ref.TargetNodeID.NodeID == nil {
	// 			continue
	// 		}
	// 		_ = collectNodes(ctx, client, *ref.TargetNodeID.NodeID, out)
	// 	}
	// 	cp = nextResp.Results[0].ContinuationPoint
	// }

	return varNodeInfo, nil
}
