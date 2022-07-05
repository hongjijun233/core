package cpumem

import (
	"context"
	"strconv"

	enginetypes "github.com/projecteru2/core/engine/types"
	"github.com/projecteru2/core/resources"
	"github.com/projecteru2/core/resources/cpumem/models"
	"github.com/projecteru2/core/resources/cpumem/types"
	coretypes "github.com/projecteru2/core/types"
)

// Plugin wrapper of CPUMem
type Plugin struct {
	c *models.CPUMem
}

// NewPlugin creates a new Plugin
func NewPlugin(config coretypes.Config) (*Plugin, error) {
	c, err := models.NewCPUMem(config)
	if err != nil {
		return nil, err
	}
	return &Plugin{c: c}, nil
}

// GetDeployArgs .
func (c *Plugin) GetDeployArgs(ctx context.Context, nodeName string, deployCount int, resourceOpts coretypes.WorkloadResourceOpts) (*resources.GetDeployArgsResponse, error) {
	workloadResourceOpts := &types.WorkloadResourceOpts{}
	if err := workloadResourceOpts.ParseFromRawParams(coretypes.RawParams(resourceOpts)); err != nil {
		return nil, err
	}
	engineArgs, resourceArgs, err := c.c.GetDeployArgs(ctx, nodeName, deployCount, workloadResourceOpts)
	if err != nil {
		return nil, err
	}

	resp := &resources.GetDeployArgsResponse{}
	err = resources.ToResp(map[string]interface{}{
		"engine_args":   engineArgs,
		"resource_args": resourceArgs,
	}, resp)
	return resp, err
}

// GetReallocArgs .
func (c *Plugin) GetReallocArgs(ctx context.Context, nodeName string, originResourceArgs coretypes.WorkloadResourceArgs, resourceOpts coretypes.WorkloadResourceOpts) (*resources.GetReallocArgsResponse, error) {
	workloadResourceOpts := &types.WorkloadResourceOpts{}
	if err := workloadResourceOpts.ParseFromRawParams(coretypes.RawParams(resourceOpts)); err != nil {
		return nil, err
	}
	originWorkloadResourceArgs := &types.WorkloadResourceArgs{}
	if err := originWorkloadResourceArgs.ParseFromRawParams(coretypes.RawParams(originResourceArgs)); err != nil {
		return nil, err
	}

	engineArgs, delta, resourceArgs, err := c.c.GetReallocArgs(ctx, nodeName, originWorkloadResourceArgs, workloadResourceOpts)
	if err != nil {
		return nil, err
	}

	resp := &resources.GetReallocArgsResponse{}
	err = resources.ToResp(map[string]interface{}{
		"engine_args":   engineArgs,
		"delta":         delta,
		"resource_args": resourceArgs,
	}, resp)
	return resp, err
}

// GetRemapArgs .
func (c *Plugin) GetRemapArgs(ctx context.Context, nodeName string, workloadMap map[string]*coretypes.Workload) (*resources.GetRemapArgsResponse, error) {
	workloadResourceArgsMap, err := c.workloadMapToWorkloadResourceArgsMap(workloadMap)
	if err != nil {
		return nil, err
	}

	engineArgs, err := c.c.GetRemapArgs(ctx, nodeName, workloadResourceArgsMap)
	if err != nil {
		return nil, err
	}

	resp := &resources.GetRemapArgsResponse{}
	err = resources.ToResp(map[string]interface{}{
		"engine_args": engineArgs,
	}, resp)
	return resp, err
}

// GetNodesDeployCapacity .
func (c *Plugin) GetNodesDeployCapacity(ctx context.Context, nodeNames []string, resourceOpts coretypes.WorkloadResourceOpts) (*resources.GetNodesDeployCapacityResponse, error) {
	workloadResourceOpts := &types.WorkloadResourceOpts{}
	if err := workloadResourceOpts.ParseFromRawParams(coretypes.RawParams(resourceOpts)); err != nil {
		return nil, err
	}

	nodesDeployCapacity, total, err := c.c.GetNodesDeployCapacity(ctx, nodeNames, workloadResourceOpts)
	if err != nil {
		return nil, err
	}

	resp := &resources.GetNodesDeployCapacityResponse{}
	err = resources.ToResp(map[string]interface{}{
		"nodes": nodesDeployCapacity,
		"total": total,
	}, resp)
	return resp, err
}

// GetMostIdleNode .
func (c *Plugin) GetMostIdleNode(ctx context.Context, nodeNames []string) (*resources.GetMostIdleNodeResponse, error) {
	nodeName, priority, err := c.c.GetMostIdleNode(ctx, nodeNames)
	if err != nil {
		return nil, err
	}

	resp := &resources.GetMostIdleNodeResponse{}
	err = resources.ToResp(map[string]interface{}{
		"node":     nodeName,
		"priority": priority,
	}, resp)
	return resp, err
}

// GetNodeResourceInfo .
func (c *Plugin) GetNodeResourceInfo(ctx context.Context, nodeName string, workloads []*coretypes.Workload) (*resources.GetNodeResourceInfoResponse, error) {
	return c.getNodeResourceInfo(ctx, nodeName, workloads, false)
}

// FixNodeResource .
func (c *Plugin) FixNodeResource(ctx context.Context, nodeName string, workloads []*coretypes.Workload) (*resources.GetNodeResourceInfoResponse, error) {
	return c.getNodeResourceInfo(ctx, nodeName, workloads, true)
}

// SetNodeResourceUsage .
func (c *Plugin) SetNodeResourceUsage(ctx context.Context, nodeName string, resourceOpts coretypes.NodeResourceOpts, resourceArgs coretypes.NodeResourceArgs, workloadResourceArgs []coretypes.WorkloadResourceArgs, delta bool, incr bool) (*resources.SetNodeResourceUsageResponse, error) {
	var nodeResourceOpts *types.NodeResourceOpts
	var nodeResourceArgs *types.NodeResourceArgs
	var workloadResourceArgsList []*types.WorkloadResourceArgs

	if resourceOpts != nil {
		nodeResourceOpts = &types.NodeResourceOpts{}
		if err := nodeResourceOpts.ParseFromRawParams(coretypes.RawParams(resourceOpts)); err != nil {
			return nil, err
		}
	}

	if resourceArgs != nil {
		nodeResourceArgs = &types.NodeResourceArgs{}
		if err := nodeResourceArgs.ParseFromRawParams(coretypes.RawParams(resourceArgs)); err != nil {
			return nil, err
		}
	}

	if workloadResourceArgs != nil {
		workloadResourceArgsList = make([]*types.WorkloadResourceArgs, len(workloadResourceArgs))
		for i, workloadResourceArg := range workloadResourceArgs {
			workloadResourceArgsList[i] = &types.WorkloadResourceArgs{}
			if err := workloadResourceArgsList[i].ParseFromRawParams(coretypes.RawParams(workloadResourceArg)); err != nil {
				return nil, err
			}
		}
	}

	before, after, err := c.c.SetNodeResourceUsage(ctx, nodeName, nodeResourceOpts, nodeResourceArgs, workloadResourceArgsList, delta, incr)
	if err != nil {
		return nil, err
	}

	resp := &resources.SetNodeResourceUsageResponse{}
	err = resources.ToResp(map[string]interface{}{
		"before": before,
		"after":  after,
	}, resp)
	return resp, err
}

// SetNodeResourceCapacity .
func (c *Plugin) SetNodeResourceCapacity(ctx context.Context, nodeName string, resourceOpts coretypes.NodeResourceOpts, resourceArgs coretypes.NodeResourceArgs, delta bool, incr bool) (*resources.SetNodeResourceCapacityResponse, error) {
	var nodeResourceOpts *types.NodeResourceOpts
	var nodeResourceArgs *types.NodeResourceArgs

	if resourceOpts != nil {
		nodeResourceOpts = &types.NodeResourceOpts{}
		if err := nodeResourceOpts.ParseFromRawParams(coretypes.RawParams(resourceOpts)); err != nil {
			return nil, err
		}
	}
	if resourceArgs != nil {
		nodeResourceArgs = &types.NodeResourceArgs{}
		if err := nodeResourceArgs.ParseFromRawParams(coretypes.RawParams(resourceArgs)); err != nil {
			return nil, err
		}
	}

	before, after, err := c.c.SetNodeResourceCapacity(ctx, nodeName, nodeResourceOpts, nodeResourceArgs, delta, incr)
	if err != nil {
		return nil, err
	}

	resp := &resources.SetNodeResourceCapacityResponse{}
	err = resources.ToResp(map[string]interface{}{
		"before": before,
		"after":  after,
	}, resp)
	return resp, err
}

// SetNodeResourceInfo .
func (c *Plugin) SetNodeResourceInfo(ctx context.Context, nodeName string, resourceCapacity coretypes.NodeResourceArgs, resourceUsage coretypes.NodeResourceArgs) (*resources.SetNodeResourceInfoResponse, error) {
	capacity := &types.NodeResourceArgs{}
	if err := capacity.ParseFromRawParams(coretypes.RawParams(resourceCapacity)); err != nil {
		return nil, err
	}

	usage := &types.NodeResourceArgs{}
	if err := usage.ParseFromRawParams(coretypes.RawParams(resourceUsage)); err != nil {
		return nil, err
	}

	if err := c.c.SetNodeResourceInfo(ctx, nodeName, capacity, usage); err != nil {
		return nil, err
	}
	return &resources.SetNodeResourceInfoResponse{}, nil
}

// AddNode .
func (c *Plugin) AddNode(ctx context.Context, nodeName string, resourceOpts coretypes.NodeResourceOpts, nodeInfo *enginetypes.Info) (*resources.AddNodeResponse, error) {
	nodeResourceOpts := &types.NodeResourceOpts{}
	if err := nodeResourceOpts.ParseFromRawParams(coretypes.RawParams(resourceOpts)); err != nil {
		return nil, err
	}

	// set default value
	if nodeInfo != nil {
		if len(nodeResourceOpts.CPUMap) == 0 {
			nodeResourceOpts.CPUMap = types.CPUMap{}
			for i := 0; i < nodeInfo.NCPU; i++ {
				nodeResourceOpts.CPUMap[strconv.Itoa(i)] = c.c.Config.Scheduler.ShareBase
			}
		}

		if nodeResourceOpts.Memory == 0 {
			nodeResourceOpts.Memory = nodeInfo.MemTotal * 8 / 10 // use 80% of real memory
		}
	}

	nodeResourceInfo, err := c.c.AddNode(ctx, nodeName, nodeResourceOpts)
	if err != nil {
		return nil, err
	}

	resp := &resources.AddNodeResponse{}
	err = resources.ToResp(map[string]interface{}{
		"capacity": nodeResourceInfo.Capacity,
		"usage":    nodeResourceInfo.Usage,
	}, resp)
	return resp, err
}

// RemoveNode .
func (c *Plugin) RemoveNode(ctx context.Context, nodeName string) (*resources.RemoveNodeResponse, error) {
	if err := c.c.RemoveNode(ctx, nodeName); err != nil {
		return nil, err
	}
	return &resources.RemoveNodeResponse{}, nil
}

// Name .
func (c *Plugin) Name() string {
	return "cpumem"
}

func (c *Plugin) workloadMapToWorkloadResourceArgsMap(workloadMap map[string]*coretypes.Workload) (*types.WorkloadResourceArgsMap, error) {
	workloadResourceArgsMap := types.WorkloadResourceArgsMap{}
	for workloadID, workload := range workloadMap {
		workloadResourceArgs := &types.WorkloadResourceArgs{}
		if err := workloadResourceArgs.ParseFromRawParams(coretypes.RawParams(workload.ResourceArgs[c.Name()])); err != nil {
			return nil, err
		}
		workloadResourceArgsMap[workloadID] = workloadResourceArgs
	}

	return &workloadResourceArgsMap, nil
}

func (c *Plugin) workloadListToWorkloadResourceArgsMap(workloads []*coretypes.Workload) (*types.WorkloadResourceArgsMap, error) {
	workloadMap := map[string]*coretypes.Workload{}
	for _, workload := range workloads {
		workloadMap[workload.ID] = workload
	}

	return c.workloadMapToWorkloadResourceArgsMap(workloadMap)
}

func (c *Plugin) getNodeResourceInfo(ctx context.Context, nodeName string, workloads []*coretypes.Workload, fix bool) (*resources.GetNodeResourceInfoResponse, error) {
	workloadResourceArgsMap, err := c.workloadListToWorkloadResourceArgsMap(workloads)
	if err != nil {
		return nil, err
	}

	nodeResourceInfo, diffs, err := c.c.GetNodeResourceInfo(ctx, nodeName, workloadResourceArgsMap, fix)
	if err != nil {
		return nil, err
	}

	resp := &resources.GetNodeResourceInfoResponse{}
	err = resources.ToResp(map[string]interface{}{
		"resource_info": nodeResourceInfo,
		"diffs":         diffs,
	}, resp)
	return resp, err
}

// GetMetricsDescription .
func (c *Plugin) GetMetricsDescription(ctx context.Context) (*resources.GetMetricsDescriptionResponse, error) {
	resp := &resources.GetMetricsDescriptionResponse{}
	err := resources.ToResp(c.c.GetMetricsDescription(), resp)
	return resp, err
}

// ResolveNodeResourceInfoToMetrics .
func (c *Plugin) ResolveNodeResourceInfoToMetrics(ctx context.Context, podName string, nodeName string, info *resources.NodeResourceInfo) (*resources.ResolveNodeResourceInfoToMetricsResponse, error) {
	capacity, usage := &types.NodeResourceArgs{}, &types.NodeResourceArgs{}
	if err := capacity.ParseFromRawParams(coretypes.RawParams(info.Capacity)); err != nil {
		return nil, err
	}
	if err := usage.ParseFromRawParams(coretypes.RawParams(info.Usage)); err != nil {
		return nil, err
	}

	metrics := c.c.ResolveNodeResourceInfoToMetrics(podName, nodeName, capacity, usage)
	resp := &resources.ResolveNodeResourceInfoToMetricsResponse{}
	err := resources.ToResp(metrics, resp)
	return resp, err
}
