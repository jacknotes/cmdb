package cvm

import (
	"encoding/json"
	"fmt"

	"github.com/jacknotes/cmdb/apps/host"
	"github.com/jacknotes/cmdb/utils"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// 创建实例: https://cloud.tencent.com/document/api/213/15730
func (o *CVMOperator) Create(req *cvm.DescribeInstancesRequest) (*host.HostSet, error) {
	resp, err := o.client.DescribeInstances(req)
	if err != nil {
		return nil, err
	}

	set := o.transferSet(resp.Response.InstanceSet)
	set.Total = utils.PtrInt64(resp.Response.TotalCount)

	return set, nil
}

// 查询可用区列表: https://cloud.tencent.com/document/product/213/15707
func (o *CVMOperator) DescribeZones() error {
	req := cvm.NewDescribeZonesRequest()
	resp, err := o.client.DescribeZones(req)
	if err != nil {
		return err
	}

	for i := range resp.Response.ZoneSet {
		zone := resp.Response.ZoneSet[i]
		fmt.Println(*zone.Zone, *zone.ZoneName, "id", *zone.ZoneId, *zone.ZoneState)
	}
	return nil
}

// 查询实例机型列表: https://cloud.tencent.com/document/api/213/15749
// 实例规格说明文档: https://cloud.tencent.com/document/product/213/11518
func (o *CVMOperator) DescribeInstanceType() error {
	req := cvm.NewDescribeInstanceTypeConfigsRequest()
	resp, err := o.client.DescribeInstanceTypeConfigs(req)
	if err != nil {
		return err
	}
	v, _ := json.Marshal(resp)
	fmt.Println(string(v))
	return nil
}
