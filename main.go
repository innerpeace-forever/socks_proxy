package main

import (
	"github.com/cihub/seelog"
	"github.com/innerpeace-forever/apf"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	app := apf.New().Configure(
		apf.WithConfiguration(apf.TOML("./conf/conf.toml")),
		apf.WithLogger(),
		apf.WithCli(apf.NewCli("Test")))

	err := app.Run(func(app *apf.Application) error {
		seelog.Info("Start Socks Proxy!")
		port := int(((app.Config().Other["Service"].(map[string]interface{}))["Port"]).(int64))
		address := ":" + strconv.Itoa(port)

		l, err := net.Listen("tcp", address)
		if err != nil {
			seelog.Infof("Listen[%v] Failed! ERR:%v", address, err)
			return err
		}

		for {
			client, err := l.Accept()
			if err != nil {
				seelog.Errorf("Accept Failed! Error:%v", err)
				continue
			}
			go handleClientRequest(client)

			if waitStopSignalNonBlock() {
				break
			}
		}

		return nil
	})

	if err != nil {
		seelog.Info("Run Failed! %v", err)
	}

	seelog.Info("Stopped!")
	app.Flush()
	os.Exit(0)
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		seelog.Error("handleClientRequest client is nil")
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		seelog.Errorf("handleClientRequest client read head failed! Error:%v", err)
		return
	}

	if b[0] == 0x05 {
		// Just process Socks5 protocol
		// Without authentication
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		if err != nil || n < 4 {
			seelog.Errorf("handleClientRequest Read Failed! Error%:%v", err)
			return
		}
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //Domain
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			seelog.Errorf("handleClientRequest Dial(%s, %s) Failed! Error:%v", host, port, err)
			return
		}
		seelog.Infof("New Connection, Dial(%s:%s).", host, port)
		defer server.Close()
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功

		// Transport datagram
		go io.Copy(server, client)
		io.Copy(client, server)
	}
}

func waitStopSignalNonBlock() bool {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-stop:
		seelog.Infof("get stop signal[%v]", s)
		return true
	default:
		return false
	}
}
