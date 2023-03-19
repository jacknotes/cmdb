package cvm_test

import (
	"fmt"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/jacknotes/cmdb/provider/txyun/connectivity"
	"github.com/jacknotes/cmdb/provider/txyun/cvm"
	"github.com/stretchr/testify/assert"

	txcvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var (
	operator *cvm.CVMOperator
)

func TestQueryCVMInstance(t *testing.T) {
	should := assert.New(t)

	set, err := operator.Query(txcvm.NewDescribeInstancesRequest())
	should.NoError(err)
	fmt.Println(set)
}

func TestPageQueryCVMInstance(t *testing.T) {
	should := assert.New(t)
	// 每秒5个请求速率限制, 1/5
	pg := operator.PageQuery(cvm.NewPageQueryRequest(5))
	HasNext := true
	for HasNext {
		ps := pg.Next()

		should.NoError(ps.Err)
		fmt.Println(ps.Data)
		HasNext = ps.HasNext
	}
}

func init() {
	err := connectivity.LoadClientFromEnv()
	if err != nil {
		panic("load client from env error")
	}
	c := connectivity.C()
	operator = cvm.NewCVMOperator(c.CvmClient())
	//开启开发者模式，打印日志
	zap.DevelopmentSetup()
}
