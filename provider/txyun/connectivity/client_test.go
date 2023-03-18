package connectivity_test

import (
	"fmt"
	"testing"

	"github.com/jacknotes/cmdb/provider/txyun/connectivity"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	should := assert.New(t)

	err := connectivity.LoadClientFromEnv()
	if should.NoError(err) {
		c := connectivity.C()
		fmt.Println(c.Account())
	}
}
