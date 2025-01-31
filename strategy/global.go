package strategy

import (
	"context"
	"fmt"
	"sort"

	"github.com/projecteru2/core/log"
	"github.com/projecteru2/core/types"
	"github.com/projecteru2/core/utils"

	"github.com/pkg/errors"
)

// GlobalPlan 基于全局资源配额
// 尽量使得资源消耗平均
func GlobalPlan(ctx context.Context, infos []Info, need, total, _ int) (map[string]int, error) {
	if total < need {
		return nil, errors.WithStack(types.NewDetailedErr(types.ErrInsufficientRes,
			fmt.Sprintf("need: %d, available: %d", need, total)))
	}
	strategyInfos := make([]Info, len(infos))
	copy(strategyInfos, infos)
	sort.Slice(infos, func(i, j int) bool { return infos[i].Capacity > infos[j].Capacity })
	length := len(strategyInfos)
	i := 0

	deployMap := make(map[string]int)
	for need > 0 {
		p := i
		deploy := 0
		delta := 0.0
		if i < length-1 {
			delta = utils.Round(strategyInfos[i+1].Usage - strategyInfos[i].Usage)
			i++
		}
		for j := 0; j <= p && need > 0 && delta >= 0; j++ {
			// 减枝
			if strategyInfos[j].Capacity == 0 {
				continue
			}
			cost := utils.Round(strategyInfos[j].Rate)
			deploy = int(delta / cost)
			if deploy == 0 {
				deploy = 1
			}
			if deploy > strategyInfos[j].Capacity {
				deploy = strategyInfos[j].Capacity
			}
			if deploy > need {
				deploy = need
			}
			strategyInfos[j].Capacity -= deploy
			deployMap[strategyInfos[j].Nodename] += deploy
			need -= deploy
		}
	}
	// 这里 need 一定会为 0 出来，因为 volTotal 保证了一定大于 need
	// 这里并不需要再次排序了，理论上的排序是基于资源使用率得到的 Deploy 最终方案
	log.Debugf(ctx, "[GlobalPlan] strategyInfos: %v", strategyInfos)
	return deployMap, nil
}
