package vsphere

import (
	"github.com/caarlos0/env/v6"
	"github.com/infraboard/cmdb/provider/vsphere/connectivity"
	"github.com/infraboard/cmdb/provider/vsphere/vm"
)

var (
	operator *Operator
)

func O() *Operator {
	if operator == nil {
		panic("please load config first")
	}
	return operator
}

func LoadOperatorFromEnv() error {
	client := &connectivity.VsphereClient{}
	if err := env.Parse(client); err != nil {
		return err
	}
	operator = NewOperator(client)
	return nil
}

func NewOperator(client *connectivity.VsphereClient) *Operator {
	return &Operator{
		client: client,
	}
}

type Operator struct {
	client *connectivity.VsphereClient
}

func (o *Operator) VmOperator() *vm.VMOperator {
	c, err := o.client.VimClient()
	if err != nil {
		panic(err)
	}
	return vm.NewVMOperator(c)
}
