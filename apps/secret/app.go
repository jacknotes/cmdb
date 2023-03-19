package secret

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/infraboard/mcube/crypto/cbc"
	"github.com/jacknotes/cmdb/conf"

	"github.com/go-playground/validator/v10"
	request "github.com/infraboard/mcube/http/request"
	"github.com/rs/xid"
)

const (
	AppName = "secret"
)

var (
	validate = validator.New()
)

func NewDefaultSecret() *Secret {
	return &Secret{
		Data: &CreateSecretRequest{
			RequestRate: 5,
		},
	}
}

func NewSecretSet() *SecretSet {
	return &SecretSet{
		Items: []*Secret{},
	}
}

func (s *SecretSet) Add(item *Secret) {
	s.Items = append(s.Items, item)
}

func NewSecret(req *CreateSecretRequest) (*Secret, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &Secret{
		Id:       xid.New().String(),
		CreateAt: time.Now().UnixMilli(),
		Data:     req,
	}, nil
}

func NewCreateSecretRequest() *CreateSecretRequest {
	return &CreateSecretRequest{
		RequestRate: 5,
	}
}

func (req *CreateSecretRequest) Validate() error {
	if len(req.AllowRegions) == 0 {
		return fmt.Errorf("required less one allow_regions")
	}
	return validate.Struct(req)
}

func NewQuerySecretRequestFromHTTP(r *http.Request) *QuerySecretRequest {
	qs := r.URL.Query()

	return &QuerySecretRequest{
		Page:     request.NewPageRequestFromHTTP(r),
		Keywords: qs.Get("keywords"),
	}
}

func NewQuerySecretRequest() *QuerySecretRequest {
	return &QuerySecretRequest{
		Page:     request.NewDefaultPageRequest(),
		Keywords: "",
	}
}

func NewDescribeSecretRequest(id string) *DescribeSecretRequest {
	return &DescribeSecretRequest{
		Id: id,
	}
}

func NewDeleteSecretRequestWithID(id string) *DeleteSecretRequest {
	return &DeleteSecretRequest{
		Id: id,
	}
}

func (s *CreateSecretRequest) AllowRegionString() string {
	return strings.Join(s.AllowRegions, ",")
}

func (s *CreateSecretRequest) LoadAllowRegionFromString(regions string) {
	if regions != "" {
		s.AllowRegions = strings.Split(regions, ",")
	}
}

// @ciphered@xxxxxx
// @plan@123456
// 123456
func (s *CreateSecretRequest) EncryptAPISecret(key string) error {
	// 判断文本是否已经加密
	if strings.HasPrefix(s.ApiSecret, conf.CIPHER_TEXT_PREFIX) {
		return fmt.Errorf("text has ciphered")
	}
	// 由于需要解密，所以采用对称加密算法进行加密
	cipherText, err := cbc.Encrypt([]byte(s.ApiSecret), []byte(key))
	if err != nil {
		return err
	}
	// TODO 为什么需要使用base64编码
	// Base64是常见传输8bit数据的编码方式
	// Base64编码的作用：由于某些系统中只能使用ASCII字符。Base64就是用来将非ASCII字符的数据转换成ASCII字符的一种方法。
	base64Str := base64.StdEncoding.EncodeToString(cipherText)
	s.ApiSecret = fmt.Sprintf("%s%s", conf.CIPHER_TEXT_PREFIX, base64Str)
	return nil
}

func (s *CreateSecretRequest) DecryptAPISecret(key string) error {
	// 如果文本已经是明文
	if !strings.HasPrefix(s.ApiSecret, conf.CIPHER_TEXT_PREFIX) {
		return fmt.Errorf("text is plan text")
	}
	// 获取base64文本
	base64CipherText := strings.TrimPrefix(s.ApiSecret, conf.CIPHER_TEXT_PREFIX)
	// 将base64文本利用cbc算法解密
	planText, err := cbc.Decrypt([]byte(base64CipherText), []byte(key))
	if err != nil {
		return err
	}
	// 将[]byte转为string
	s.ApiSecret = string(planText)
	fmt.Printf("%s, abc", s.ApiSecret)
	return nil
}

// 敏感数据脱敏
func (s *CreateSecretRequest) Desense() {
	if s.ApiSecret != "" {
		s.ApiSecret = "******"
	}
}
