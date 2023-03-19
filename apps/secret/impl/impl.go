package impl

import (
	"database/sql"

	"github.com/jacknotes/cmdb/conf"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"github.com/jacknotes/cmdb/apps/host"
	"github.com/jacknotes/cmdb/apps/secret"
)

var (
	// Service 服务实例
	svr = &impl{}
)

type impl struct {
	db  *sql.DB
	log logger.Logger

	host host.ServiceServer
	secret.UnimplementedServiceServer
}

func (s *impl) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)
	return nil
}

func (s *impl) Name() string {
	return secret.AppName
}

// 需要提前注册到grpc server中
func (s *impl) Registry(server *grpc.Server) {
	secret.RegisterServiceServer(server, svr)
}

// 将该rpc应用注册到ioc容器中
func init() {
	app.RegistryGrpcApp(svr)
}
