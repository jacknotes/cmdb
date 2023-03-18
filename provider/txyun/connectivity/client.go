package connectivity

import (
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/jacknotes/cmdb/utils"
	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	cloudaudit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	client *TencentCloudClient
)

func C() *TencentCloudClient {
	if client == nil {
		panic("please load config first")
	}
	return client
}

func LoadClientFromEnv() error {
	client = &TencentCloudClient{}
	if err := env.Parse(client); err != nil {
		return err
	}

	return nil
}

// NewTencentCloudClient client
func NewTencentCloudClient(credentialID, credentialKey, region string) *TencentCloudClient {
	return &TencentCloudClient{
		Region:    region,
		SecretID:  credentialID,
		SecretKey: credentialKey,
	}
}

// TencentCloudClient client for all TencentCloud service
type TencentCloudClient struct {
	Region    string `env:"TX_CLOUD_REGION"`
	SecretID  string `env:"TX_CLOUD_SECRET_ID"`
	SecretKey string `env:"TX_CLOUD_SECRET_KEY"`

	cvmConn       *cvm.Client
	cdbConn       *cdb.Client
	sqlserverConn *sqlserver.Client
	vpcConn       *vpc.Client
	redisConn     *redis.Client
	cosConn       *cos.Client
	cbsConn       *cbs.Client
	clbConn       *clb.Client
	mongoConn     *mongodb.Client
	dnsConn       *dnspod.Client
	billConn      *billing.Client
	auditConn     *cloudaudit.Client
}

// UseCvmClient cvm
func (me *TencentCloudClient) CvmClient() *cvm.Client {
	if me.cvmConn != nil {
		return me.cvmConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	cvmConn, _ := cvm.NewClient(credential, me.Region, cpf)
	me.cvmConn = cvmConn
	return me.cvmConn
}

// UseCvmClient cvm
func (me *TencentCloudClient) DnsClient() *dnspod.Client {
	if me.dnsConn != nil {
		return me.dnsConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	dnsConn, _ := dnspod.NewClient(credential, me.Region, cpf)
	me.dnsConn = dnsConn
	return me.dnsConn
}

// UseCvmClient cvm
func (me *TencentCloudClient) AuditClient() *cloudaudit.Client {
	if me.auditConn != nil {
		return me.auditConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	auditConn, _ := cloudaudit.NewClient(credential, me.Region, cpf)
	me.auditConn = auditConn
	return me.auditConn
}

// UseCvmClient cvm
func (me *TencentCloudClient) ClbClient() *clb.Client {
	if me.clbConn != nil {
		return me.clbConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	clbConn, _ := clb.NewClient(credential, me.Region, cpf)
	me.clbConn = clbConn
	return me.clbConn
}

// UseCvmClient cvm
func (me *TencentCloudClient) CBSClient() *cbs.Client {
	if me.cbsConn != nil {
		return me.cbsConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	cbsConn, _ := cbs.NewClient(credential, me.Region, cpf)
	me.cbsConn = cbsConn
	return me.cbsConn
}

// UseCvmClient cvm
func (me *TencentCloudClient) VpcClient() *vpc.Client {
	if me.vpcConn != nil {
		return me.vpcConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	vpcConn, _ := vpc.NewClient(credential, me.Region, cpf)
	me.vpcConn = vpcConn
	return me.vpcConn
}

// UseBillingClient billing客户端
func (me *TencentCloudClient) BillingClient() *billing.Client {
	if me.billConn != nil {
		return me.billConn
	}
	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	billConn, _ := billing.NewClient(credential, "", cpf)
	me.billConn = billConn

	return me.billConn
}

// CDBClient cdb
func (me *TencentCloudClient) CDBClient() *cdb.Client {
	if me.cdbConn != nil {
		return me.cdbConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	cdbConn, _ := cdb.NewClient(credential, me.Region, cpf)
	me.cdbConn = cdbConn
	return me.cdbConn
}

// CDBClient cdb
func (me *TencentCloudClient) SQLServerClient() *sqlserver.Client {
	if me.sqlserverConn != nil {
		return me.sqlserverConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	sqlserverConn, _ := sqlserver.NewClient(credential, me.Region, cpf)
	me.sqlserverConn = sqlserverConn
	return me.sqlserverConn
}

// RedisClient cdb
func (me *TencentCloudClient) RedisClient() *redis.Client {
	if me.redisConn != nil {
		return me.redisConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	conn, _ := redis.NewClient(credential, me.Region, cpf)
	me.redisConn = conn
	return me.redisConn
}

// RedisClient cdb
func (me *TencentCloudClient) MongoClient() *mongodb.Client {
	if me.mongoConn != nil {
		return me.mongoConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	conn, _ := mongodb.NewClient(credential, me.Region, cpf)
	me.mongoConn = conn
	return me.mongoConn
}

// CDBClient cdb
func (me *TencentCloudClient) CosClient() *cos.Client {
	if me.cosConn != nil {
		return me.cosConn
	}

	me.cosConn = cos.NewClient(nil, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  me.SecretID,
			SecretKey: me.SecretKey,
		},
	})

	return me.cosConn
}

// 获取客户端账号ID
func (me *TencentCloudClient) Account() (string, error) {
	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	stsConn, _ := sts.NewClient(credential, me.Region, cpf)

	req := sts.NewGetCallerIdentityRequest()

	resp, err := stsConn.GetCallerIdentity(req)
	if err != nil {
		return "", fmt.Errorf("unable to initialize the STS client: %#v", err)
	}

	return utils.PtrStrV(resp.Response.AccountId), nil
}
