## Lion
> 用于满帮运维开发团队接入lion配置管理中心

## Example(Golang)
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
    if err != nil {
    	fmt.Println(err.Error())
    	return
    }  	
    fmt.Println(config)
    
    // 获取用户下 部分项目的部分配置
  	config, err = l.QueryConfigByKey("qa", []string{"phantom-service-rel.database.addr", "phantom-service-rel.database.port"})
    if err != nil {
        fmt.Println(err.Error())
        return
    }
  	fmt.Println(config)

    // 获取用户下 单个配置项
  	config, err = l.QueryConfigByKey("qa", "phantom-service-rel.database.addr")
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(config)
}
```

## Example(Python)
> beidou-api/comm/utils.py line 198-269
```python
def get_uniform_conf(key, parent='beidou', path=None, filename='cfg.toml', source='local', **kwargs):
    """
    get direct key value, ignore the config file type
    :param key:
    :param parent:
    :param path:
    :param filename:
    :param source:
    :return:
    """
    if source == 'local':
        if not path:
            path = '{0}/conf/{1}/{2}'.format(os.getcwd(), current_env, filename)
        else:
            path = '{0}/{1}'.format(path, filename)
        extend_name = filename.split('.')[-1]
        if extend_name == 'toml':
            return toml.load(path).get(parent, {}).get(key)
        elif extend_name == 'ini':
            cf = ConfigParser.ConfigParser()
            try:
                return cf.read(path).get(parent, key)
            except Exception:
                return
        else:
            pass
    elif source == 'lion':
        retry_count = 1
        # add retry logic
        while retry_count < 3:
            try:
                app = "devops"
                headers = {
                    'Content-Type': 'application/x-www-form-urlencoded',
                }
                url = get_uniform_conf('{0}_url'.format('prd' if current_env == 'prod' else 'dev'), parent='lion')
                data = urlencode(dict(
                    app=app,
                    pwd='j32f[f;fQFk5~dash' if current_env == 'prod' else '1jkf43FDr[tad',
                ))
                http_base = HttpBase(headers=headers, data=data, timeout=3)
                http_base.prepared_request('POST', url + u'/getToken.do')
                response = http_base.send_request()
                if response:
                    data = urlencode(dict(
                        env='product' if current_env == 'prod' else 'prelease' if current_env == 'qa' else 'qa',
                        token=response.json().get('result', {}).get('token', ''),
                        key=key if kwargs.get("all") else None,
                        app=app,
                        projectName=kwargs.get('project_name')
                    ))
                    http_base.data = data
                    http_base.prepared_request(
                        'POST', url + u'/configQuery.do'
                    )
                    result = http_base.send_request()
                    if result and result.json().get('code', -1) == 0:
                        logger.info(u'get lion config success')
                        return result.json()
                else:
                    logger.info(u'get lion config failure')
                retry_count += 1
            except Exception as e:
                print_exc()
                logger.info(u'get lion config occur exception, {0}'.format(e.message))
                retry_count += 1
        if retry_count >= 100:
            return True
        else:
            return False
    else:
        pass
```
