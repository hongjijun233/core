package models

import (
	"fmt"
	"strings"

	"github.com/projecteru2/core/resources/volume/types"
)

// GetMetricsDescription .
func (v *Volume) GetMetricsDescription() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":   "storage_used",
			"help":   "node used storage.",
			"type":   "gauge",
			"labels": []string{"podname", "nodename"},
		},
		{
			"name":   "storage_capacity",
			"help":   "node available storage.",
			"type":   "gauge",
			"labels": []string{"podname", "nodename"},
		},
	}
}

func (v *Volume) ConvertNodeResourceInfoToMetrics(podname string, nodename string, nodeResourceCapacity *types.NodeResourceArgs, nodeResourceUsage *types.NodeResourceArgs) []map[string]interface{} {
	cleanedNodeName := strings.ReplaceAll(nodename, ".", "_")
	metrics := []map[string]interface{}{
		{
			"name":   "storage_used",
			"labels": []string{podname, nodename},
			"value":  fmt.Sprintf("%v", nodeResourceUsage.Storage),
			"key":    fmt.Sprintf("core.node.%s.storage.used", cleanedNodeName),
		},
		{
			"name":   "storage_capacity",
			"labels": []string{podname, nodename},
			"value":  fmt.Sprintf("%v", nodeResourceCapacity.Storage),
			"key":    fmt.Sprintf("core.node.%s.storage.used", cleanedNodeName),
		},
	}

	return metrics
}
