# Kate is a lightweight HTTP API framework.(For Go1.6+)

## GETTING STARTED
Please see [helloworld.go](https://github.com/k81/kate/blob/master/examples/helloworld.go)

## Features
- 请求的handler支持context
- 日志支持context，比如某一次http请求可以用session来区分
- 支持优雅重启，重启或关闭服务时，保证所有正在处理的请求都处理完
- 中间件可扩展,优雅重启、panic自动捕获、日志、签名校验、超时处理等都是在中间件中实现
- 待续...

## TODO
- ORM and Cache support
- Circuit Breaker and RateLimit support
- Others...
