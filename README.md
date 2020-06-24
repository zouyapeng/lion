## Lion
> 用于满帮运维开发团队接入lion配置管理中心

## Example
```go
package main

import (
  "fmt"

  "code.amh-group.com/devops/lion"
)

func main() {
    // 初始化lion配置
    l, err := lion.Init("http://dev-lion.amh-group.com/lion/open", "<YourAppName>", "<YourPassword>")
  	if err != nil {
  		fmt.Println(err.Error())
  		return
  	}

    // 获取整个项目的所有配置
  	config, err := l.QueryConfigByProject("qa", "phantom-service-rel")
  	fmt.Println(config)
    
    // 获取用户下 部分项目的部分配置
  	config, err = l.QueryConfigByKey("qa", []string{"phantom-service-rel.database.addr", "phantom-service-rel.database.port"})
  	fmt.Println(config)

    // 获取用户下 单个配置项
  	config, err = l.QueryConfigByKey("qa", "phantom-service-rel.database.addr")
  	fmt.Println(config)
}
```
