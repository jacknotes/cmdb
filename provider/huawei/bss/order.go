package bss

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/model"
	"github.com/infraboard/mcube/pager"

	"github.com/infraboard/cmdb/apps/order"
	"github.com/infraboard/cmdb/apps/resource"
	"github.com/infraboard/cmdb/provider"
	"github.com/infraboard/cmdb/utils"
)

func (o *BssOperator) DescribeOrder(ctx context.Context, r *provider.DescribeRequest) (*order.Order, error) {
	return nil, nil
}

func (o *BssOperator) PageQueryOrder(req *provider.QueryOrderRequest) pager.Pager {
	p := newOrderPager(o, req)
	p.SetRate(req.Rate)
	return p
}

// 客户购买包年/包月资源后,可以查看待审核、处理中、已取消、已完成和待支付等状态的订单
// 参考文档: https://apiexplorer.developer.huaweicloud.com/apiexplorer/doc?product=BSS&api=ListCustomerOrders
func (o *BssOperator) doQueryOrder(req *model.ListCustomerOrdersRequest) (*order.OrderSet, error) {
	set := order.NewOrderSet()

	resp, err := o.client.ListCustomerOrders(req)
	if err != nil {
		return nil, err
	}

	set.Total = int64(*resp.TotalCount)
	set.Items = o.transferOrderSet(resp.OrderInfos).Items

	// 补充订单关联的资源
	o.fillResourceId(set)

	return set, nil
}

func (o *BssOperator) transferOrderSet(list *[]model.CustomerOrderV2) *order.OrderSet {
	set := order.NewOrderSet()
	items := *list
	for i := range items {
		set.Add(o.transferOrder(items[i]))
	}
	return set
}

func (o *BssOperator) transferOrder(ins model.CustomerOrderV2) *order.Order {
	b := order.NewDefaultOrder()
	b.Vendor = resource.VENDOR_ALIYUN
	b.Id = tea.StringValue(ins.OrderId)
	b.OrderType = praseOrderType(ins.OrderType)
	b.Status = praseOrderStatus(ins.Status)
	b.ProductCode = tea.StringValue(ins.ServiceTypeCode)
	b.ProductName = tea.StringValue(ins.ServiceTypeName)
	b.CreateAt = utils.ParseDefaultSecondTime(tea.StringValue(ins.CreateTime))
	b.PayAt = utils.ParseDefaultSecondTime(tea.StringValue(ins.PaymentTime))

	// 金额信息
	cost := b.Cost
	cost.SalePrice = utils.PtrFloat64(ins.OfficialAmount)
	cost.RealCost = utils.PtrFloat64(ins.AmountAfterDiscount)
	return b
}

func (o *BssOperator) fillResourceId(set *order.OrderSet) {
	wg := &sync.WaitGroup{}
	for i := range set.Items {
		wg.Add(1)
		go func(orderId string) {
			defer wg.Done()
			o.log.Debugf("query order: %s resources", orderId)
			req := &model.ListPayPerUseCustomerResourcesRequest{
				Body: &model.QueryResourcesReq{
					OrderId: tea.String(orderId),
				},
			}
			od := set.GetOrderById(orderId)
			if o != nil {
				return
			}

			err := o.doOrderResource(req, od)
			if err != nil {
				o.log.Errorf("query order resource error, %s", err)
			}

		}(set.Items[i].Id)
	}

	wg.Wait()
}

// 客户在伙伴销售平台查询某个或所有的包年/包月资源(ListPayPerUseCustomerResources)
// 参考文档: https://apiexplorer.developer.huaweicloud.com/apiexplorer/doc?product=BSS&api=ListPayPerUseCustomerResources
func (o *BssOperator) doOrderResource(req *model.ListPayPerUseCustomerResourcesRequest, od *order.Order) error {
	o.tb.Wait(1)

	resp, err := o.client.ListPayPerUseCustomerResources(req)
	if err != nil {
		return err
	}

	if resp.Data == nil {
		return nil
	}

	pi, err := json.Marshal(*resp.Data)
	if err == nil {
		od.ProductInfo = string(pi)
	}

	for _, d := range *resp.Data {
		if tea.Int32Value(d.IsMainResource) == 1 {
			od.ResourceId = append(od.ResourceId, tea.StringValue(d.ResourceId))
		}
	}

	return nil
}

// 客户可以在伙伴销售平台查看订单详情
// 参考文档: https://apiexplorer.developer.huaweicloud.com/apiexplorer/sdk?product=BSS&api=ShowCustomerOrderDetails
