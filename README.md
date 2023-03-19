# go-std-http-example

演示go标准库中的[net/http](https://pkg.go.dev/net/http@go1.20.2)怎样使用，示例的接口包括

## HTTP Methods

> 测试基础的HTTP方法

| API   | 请求方法 | 功能     |
|:----- |:---- |:------ |
| /get  | GET  | 获取请求信息 |
| /post | POST | 获取请求信息 |

## Request inspection

> 获取请求头信息

| API         | 请求方法 | 功能                  |
| ----------- | ---- | ------------------- |
| /headers    | GET  | 获取请求的请求头信息          |
| /user-agent | GET  | 获取请求头中的user-agent的值 |

## Dynamic data

> 处理动态数据

| API             | 请求方法 | 功能                     | 示例                                   |
| --------------- | ---- | ---------------------- | ------------------------------------ |
| /base64/{value} | GET  | 解码一个base64编码的字符串并且返回结果 | /base64/Z28tc3RkLWh0dHAtZXhhbXBsZQ== |

## Cookies

> cookie相关功能

| API             | 请求方法 | 功能                 | 示例                            |
| --------------- | ---- | ------------------ | ----------------------------- |
| /cookies        | GET  | 查看请求中带有的所有cookie信息 | /cookies                      |
| /cookies/set    | GET  | 请求服务器设置cookie      | /cookies/set?name=Mike&age=18 |
| /cookies/delete | GET  | 请求服务器删除某个cookie    | /cookies/delete?name=         |

## Images

> 请求图片数据接口

| API         | 请求方法 | 功能            |
| ----------- | ---- | ------------- |
| /image/jpeg | GET  | 返回一张jpeg格式的图片 |
| /image/png  | GET  | 返回一张png格式的图片  |

## Redirect

> 重定向相关功能接口

| API                | 请求方法 | 功能    | 示例                     |
| ------------------ | ---- | ----- | ---------------------- |
| /absolute-redirect | GET  | 重定向n次 | /absolute-redirect?n=3 |

## Echo

> 回声接口

| API             | 请求方法 | 功能     |
| --------------- | ---- | ------ |
| /echo/{content} | GET  | echo功能 |

## Response formats

> 返回不同格式的响应格式内容

| API      | 请求方法 | 功能                              |
| -------- | ---- | ------------------------------- |
| /deflate | GET  | 返回content-encoding为deflate的响应结果 |
| /gzip    | GET  | 返回content-encoding为gzip的响应结果    |
| /html    | GET  | 返回html格式的响应                     |
| /xml     | GET  | 返回xml格式的响应                      |

## Caches

> 模拟缓存的一些操作

| API          | 请求方法 | 功能                                                              | 示例                |
| ------------ | ---- | --------------------------------------------------------------- | ----------------- |
| /cache       | GET  | 如果请求头中带有If-Modified-Since或If-None-Match，则返回304；否则就是和/get接口一样的结果 | /cache            |
| /cache/{n}   | GET  | 请求服务器在响应头中添加cache-control，其max-age=n                            | /cache/60         |
| /etag/{etag} | GET  | 请求服务器在响应头中添加ETag，值为设置的{etag}                                    | /etag/coffeemaker |

## Authorization

> HTTP基础的认证功能

| API                         | 请求方法 | 功能                  | 示例                      |
| --------------------------- | ---- | ------------------- | ----------------------- |
| /basic-auth/username/passwd | GET  | 提示用户要进行HTTP Basic认证 | /basic-auth/ryan/123456 |