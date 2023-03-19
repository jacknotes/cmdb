package impl_test

import (
	"context"
	"os"
	"testing"

	"github.com/infraboard/mcube/app"
	"github.com/jacknotes/cmdb/apps/secret"
	"github.com/jacknotes/cmdb/conf"
	"go.uber.org/zap"
)

var (
	ins secret.ServiceServer
)

func TestImpl_QuerySecret(t *testing.T) {
	querySecret, err := ins.QuerySecret(context.Background(), secret.NewQuerySecretRequest())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(querySecret)
}

func TestImpl_DescribeSecret(t *testing.T) {
	describeSecret, err := ins.DescribeSecret(context.Background(), secret.NewDescribeSecretRequest("ca8sio13n7plv10h6fu0"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(describeSecret)
}

func TestImpl_CreateSecret(t *testing.T) {
	req := secret.NewCreateSecretRequest()
	req.Description = "测试用例"
	req.ApiKey = os.Getenv("TX_CLOUD_SECRET_ID")
	req.ApiSecret = os.Getenv("TX_CLOUD_SECRET_KEY")
	req.AllowRegions = []string{"*"}
	ss, err := ins.CreateSecret(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ss)
}

func init() {
	// 从环境变量中加载测试配置
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}
	// 全局日志对象初始化
	zap.Development()

	// 初始化所有实例
	if err := app.InitAllApp(); err != nil {
		panic(err)
	}

	// 实例化ins对象
	ins = app.GetGrpcApp(secret.AppName).(secret.ServiceServer)
}
