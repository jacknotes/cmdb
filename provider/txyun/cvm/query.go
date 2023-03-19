package cvm

import (
	"fmt"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/jacknotes/cmdb/apps/host"
	"github.com/jacknotes/cmdb/utils"
)

// 查看实例列表
// 查看实例列表: https://console.cloud.tencent.com/api/explorer?Product=cvm&Version=2017-03-12&Action=DescribeInstances&SignVersion=
func (o *CVMOperator) Query(req *cvm.DescribeInstancesRequest) (*host.HostSet, error) {
	resp, err := o.client.DescribeInstances(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.ToJsonString())
	// 需要把腾讯云CVM的数据结构转化为我们定义的Host
	set := o.transferSet(resp.Response.InstanceSet)
	set.Total = utils.PtrInt64(resp.Response.TotalCount)
	return set, nil
}

func NewPageQueryRequest(reqPerSecond int) *PageQueryRequest {
	return &PageQueryRequest{
		reqPerSecond: reqPerSecond,
	}
}

type PageQueryRequest struct {
	reqPerSecond int
}

func (o *CVMOperator) PageQuery(req *PageQueryRequest) host.Pager {
	return newPager(20, o, req.reqPerSecond)
}
