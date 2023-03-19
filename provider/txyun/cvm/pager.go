package cvm

import (
	"github.com/jacknotes/cmdb/apps/host"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

func newPager(pageSize int, operator *CVMOperator) *pager {
	req := cvm.NewDescribeInstancesRequest()
	req.Limit = common.Int64Ptr(int64(pageSize))

	return &pager{
		size:     pageSize,
		number:   1,
		operator: operator,
		req:      req,
		log:      zap.L().Named("Pager"),
	}
}

type pager struct {
	size     int
	number   int
	total    int64
	operator *CVMOperator
	req      *cvm.DescribeInstancesRequest
	log      logger.Logger
}

func (p *pager) Next() *host.PagerResult {
	result := host.NewPageResult()

	// 每调用一次，就需要构造一个翻页的请求
	resp, err := p.operator.Query(p.nextReq())
	if err != nil {
		result.Err = err
		return result
	}

	// 完成一页请求过后，需要修改total
	p.total = resp.Total
	p.log.Debugf("get %d hosts", len(resp.Items))

	result.Data = resp
	result.HasNext = p.hasNext()

	p.number++
	return result
}

// func (p *pager) Scan(ctx context.Context, set pager.Set) error {
// 	resp, err := p.operator.Query(p.nextReq())
// 	if err != nil {
// 		return err
// 	}
// 	p.CheckHasNext(resp)
// 	p.log.Debugf("get %d hosts", len(resp.Items))
// 	set.Add(resp.ToAny()...)

// 	return nil
// }

func (p *pager) nextReq() *cvm.DescribeInstancesRequest {
	p.log.Debugf("请求第%d页数据", p.number)
	p.req.Offset = common.Int64Ptr(p.offset())
	p.req.Limit = common.Int64Ptr(int64(p.size))
	return p.req
}

func (p *pager) hasNext() bool {
	return int64(p.number*p.size) < p.total
}

func (p *pager) offset() int64 {
	return int64(p.size * (p.number - 1))
}
