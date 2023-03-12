package redis

import (
	"context"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	cmdbRedis "github.com/infraboard/cmdb/apps/redis"
	"github.com/infraboard/cmdb/apps/resource"
	"github.com/infraboard/cmdb/provider"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/pager"
)

func (o *RedisOperator) DescribeRedis(ctx context.Context, r *provider.DescribeRequest) (
	*cmdbRedis.Redis, error) {
	if err := r.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	req := redis.NewDescribeInstancesRequest()
	req.InstanceId = tea.String(r.Id)
	req.Limit = tea.Uint64(1)

	set, err := o.Query(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "No instance found") {
			return nil, exception.NewNotFound(err.Error())
		}
		return nil, err
	}

	if set.Length() == 0 {
		return nil, exception.NewNotFound("redis %s not found", r.Id)
	}

	return set.Items[0], nil
}

func (o *RedisOperator) PageQueryRedis(req *provider.QueryRequest) pager.Pager {
	return newPager(20, o)
}

// 查询Redis实例列表
// 参考: https://console.cloud.tencent.com/api/explorer?Product=redis&Version=2018-04-12&Action=DescribeInstances&SignVersion=
func (o *RedisOperator) Query(ctx context.Context, req *redis.DescribeInstancesRequest) (*cmdbRedis.Set, error) {
	resp, err := o.client.DescribeInstancesWithContext(ctx, req)
	if err != nil {
		return nil, err
	}

	return o.transferSet(resp.Response), nil
}

func (o *RedisOperator) transferSet(items *redis.DescribeInstancesResponseParams) *cmdbRedis.Set {
	set := cmdbRedis.NewSet()
	for i := range items.InstanceSet {
		set.Add(o.transferOne(items.InstanceSet[i]))
	}
	return set
}

func (o *RedisOperator) transferOne(ins *redis.InstanceSet) *cmdbRedis.Redis {
	r := cmdbRedis.NewDefaultRedis()
	b := r.Resource.Meta

	b.CreateAt = o.parseTime(tea.StringValue(ins.Createtime))
	b.Id = tea.StringValue(ins.InstanceId)

	info := r.Resource.Spec
	info.Vendor = resource.VENDOR_TENCENT
	info.Region = tea.StringValue(ins.Region)
	info.ExpireAt = o.parseTime(tea.StringValue(ins.DeadlineTime))
	info.Category = tea.StringValue(ins.ProductType)
	info.Type = o.ParseType(ins.Type)
	info.Name = tea.StringValue(ins.InstanceName)

	r.Resource.Status.Phase = praseStatus(ins.Status)
	r.Resource.Status.PrivateAddress = []string{tea.StringValue(ins.WanIp)}
	r.Resource.Cost.PayMode = o.parsePAY_MODE(ins.BillingMode)

	info.Memory = int32(tea.Float64Value(ins.Size))
	info.BandWidth = int32(tea.Int64Value(ins.NetLimit))

	desc := r.Describe
	desc.MaxConnection = tea.Int64Value(ins.ClientLimitMax)
	desc.EngineType = tea.StringValue(ins.Engine)
	desc.EngineVersion = o.ParseType(ins.Type)
	desc.ConnectAddr = tea.StringValue(ins.WanIp)
	desc.ConnectPort = tea.Int64Value(ins.Port)
	return r
}
