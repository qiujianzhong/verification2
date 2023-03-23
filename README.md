# 短信验证码获取
> 可以自动获取多个环境的验证码，并且取最新的返回到页面。

# 使用
修改verification2.go中的数据库连接

编译
```code
#编译Linux版本
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build verification2.go
#编译mac版本
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build verification2.go
#编译windoes版本
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build verification2.go
```
# 运行
#Linux后台运行
nohub ./verification2 &

打开环境链接 http://ip:8889/


# 界面
![](./截图.png "页面")