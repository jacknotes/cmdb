package api_test

import (
	"testing"

	"github.com/jacknotes/cmdb/apps/secret"
)

var (
	encryptKey = "sdfsdfsfdsfd"
)

func TestSecretEncrypt(t *testing.T) {
	ins := secret.NewDefaultSecret()
	ins.Data.ApiSecret = "123456"
	ins.Data.EncryptAPISecret(encryptKey)
	t.Log(ins.Data.ApiSecret)
	ins.Data.DecryptAPISecret(encryptKey)
	t.Log(ins.Data.ApiSecret)
}
