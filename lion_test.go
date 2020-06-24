package lion_test

import (
	"fmt"
	"testing"

	"code.amh-group.com/devops/lion"
)

func TestLion(t *testing.T) {
	l, err := lion.Login("http://dev-lion.amh-group.com/lion/open", "<YourAppName>", "<YourPassword>")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(*l)
	config, err := l.QueryConfigByProject("qa", "phantom-service-rel")
	fmt.Println(config)
	config, err = l.QueryConfigByKey("qa", []string{"phantom-service-rel.database.addr", "phantom-service-rel.database.port"})
	fmt.Println(config)
	config, err = l.QueryConfigByKey("qa", "phantom-service-rel.database.addr")
	fmt.Println(config)
}
