package cvm

import (
	"time"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/jacknotes/cmdb/apps/host"
	"github.com/jacknotes/cmdb/apps/resource"
	"github.com/jacknotes/cmdb/utils"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func NewCVMOperator(client *cvm.Client) *CVMOperator {
	return &CVMOperator{
		client: client,
		log:    zap.L().Named("CVM"),
	}
}

type CVMOperator struct {
	client *cvm.Client
	log    logger.Logger
}

func (o *CVMOperator) transferSet(items []*cvm.Instance) *host.HostSet {
	set := host.NewHostSet()
	for i := range items {
		ins := o.transferOne(items[i])
		set.Add(ins)
	}
	return set
}

func (o *CVMOperator) transferOne(ins *cvm.Instance) *host.Host {
	h := host.NewDefaultHost()
	h.Base.Vendor = resource.Vendor_TENCENT
	h.Base.Region = o.client.GetRegion()
	h.Base.Zone = utils.PtrStrV(ins.Placement.Zone)
	h.Base.CreateAt = o.parseTime(utils.PtrStrV(ins.CreatedTime))
	h.Base.Id = utils.PtrStrV(ins.InstanceId)

	h.Information.ExpireAt = o.parseTime(utils.PtrStrV(ins.ExpiredTime))
	h.Information.Type = utils.PtrStrV(ins.InstanceType)
	h.Information.Name = utils.PtrStrV(ins.InstanceName)
	h.Information.Status = utils.PtrStrV(ins.InstanceState)
	h.Information.Tags = transferTags(ins.Tags)
	h.Information.PublicIp = utils.SlicePtrStrv(ins.PublicIpAddresses)
	h.Information.PrivateIp = utils.SlicePtrStrv(ins.PrivateIpAddresses)
	h.Information.PayType = utils.PtrStrV(ins.InstanceChargeType)

	h.Describe.Cpu = utils.PtrInt64(ins.CPU)
	h.Describe.Memory = utils.PtrInt64(ins.Memory)
	h.Describe.OsName = utils.PtrStrV(ins.OsName)
	h.Describe.SerialNumber = utils.PtrStrV(ins.Uuid)
	h.Describe.ImageId = utils.PtrStrV(ins.ImageId)

	if ins.InternetAccessible != nil {
		h.Describe.InternetMaxBandwidthOut = utils.PtrInt64(ins.InternetAccessible.InternetMaxBandwidthOut)
	}
	h.Describe.KeyPairName = utils.SlicePtrStrv(ins.LoginSettings.KeyIds)
	h.Describe.SecurityGroups = utils.SlicePtrStrv(ins.SecurityGroupIds)
	return h
}

func transferTags(tags []*cvm.Tag) (ret []*resource.Tag) {
	return nil
}

func (o *CVMOperator) parseTime(t string) int64 {
	ts, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err != nil {
		o.log.Errorf("parse time %s error, %s", t, err)
		return 0
	}
	return ts.UnixMilli()
}
