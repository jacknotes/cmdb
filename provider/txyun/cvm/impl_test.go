package cvm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/infraboard/cmdb/apps/disk"
	"github.com/infraboard/cmdb/apps/eip"
	"github.com/infraboard/cmdb/apps/host"
	"github.com/infraboard/cmdb/provider"
	"github.com/infraboard/cmdb/provider/txyun"
	"github.com/infraboard/cmdb/utils"
	"github.com/infraboard/mcube/logger/zap"

	op "github.com/infraboard/cmdb/provider/txyun/cvm"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var (
	operator *op.CVMOperator
	ctx      = context.Background()
)

func TestPageQueryHost(t *testing.T) {
	pager := operator.PageQueryHost(provider.NewQueryRequest())

	for pager.Next() {
		set := host.NewHostSet()
		if err := pager.Scan(ctx, set); err != nil {
			panic(err)
		}
		for i := range set.Items {
			fmt.Println(set.Items[i])
		}

	}
}

func TestPageQueryDisk(t *testing.T) {
	pager := operator.PageQueryDisk(provider.NewQueryRequest())

	for pager.Next() {
		set := disk.NewDiskSet()
		if err := pager.Scan(ctx, set); err != nil {
			panic(err)
		}
		fmt.Println(set)
	}
}

func TestDescribeDisk(t *testing.T) {
	req := provider.NewDescribeRequest("xxx")
	ins, err := operator.DescribeDisk(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)

}

func TestPageQueryEip(t *testing.T) {
	pager := operator.PageQueryEip(provider.NewQueryRequest())

	for pager.Next() {
		set := eip.NewEIPSet()
		if err := pager.Scan(ctx, set); err != nil {
			panic(err)
		}
		fmt.Println(set)
	}
}

func TestDescribeEcs(t *testing.T) {
	req := provider.NewDescribeRequest("xxxx")
	ins, err := operator.DescribeHost(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestInquiryPrice(t *testing.T) {
	req := cvm.NewInquiryPriceRunInstancesRequest()
	req.Placement = &cvm.Placement{
		Zone: utils.StringPtr("ap-shanghai-2"),
	}
	req.ImageId = utils.StringPtr("img-l5eqiljn")
	req.InstanceType = utils.StringPtr("S4.SMALL1")
	req.InstanceChargeType = utils.StringPtr("SPOTPAID")
	if err := operator.InquiryNewPrice(req); err != nil {

	}
}

func TestDescribeZones(t *testing.T) {
	operator.DescribeZones()
}

func TestDescribeInstanceType(t *testing.T) {
	operator.DescribeInstanceType()
}

func TestCreate(t *testing.T) {
}

func init() {
	zap.DevelopmentSetup()

	err := txyun.LoadOperatorFromEnv()
	if err != nil {
		panic(err)
	}

	c := txyun.O().Client()
	operator = op.NewCVMOperator(c.CvmClient(), c.CBSClient(), c.VpcClient())
}
