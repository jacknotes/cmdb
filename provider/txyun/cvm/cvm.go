package cvm

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
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
