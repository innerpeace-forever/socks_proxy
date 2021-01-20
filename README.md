# socks_proxy

一个使用golang实现的socks5代理服务。不支持socks4，也不支持用户名和密码验证。

什么场景上可以使用？

* 学习如何编写一个socks5 proxy服务；
* 部署一个自己使用的socks5 proxy用来访问外网；

### Installation

在linux服务器上，将源代码clone到$GOPATH/src/目录下，并执行编译。

```shell
git clone https://github.com/innerpeace-forever/socks_proxy.git
mv socks_proxy $GOPATH/src/
cd $GOPATH/src/socks_proxy
sh -x build/build.sh
```

### Configuration

除了基本的日志配置外（可以不用动），最主要的是配置端口。在conf/conf.toml上

```toml
[Other.Service]
Port = 8080
```

