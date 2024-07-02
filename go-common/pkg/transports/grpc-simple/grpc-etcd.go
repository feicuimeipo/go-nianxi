package grpc_simple

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type EtcdServer struct {
	cli         *clientv3.Client
	serviceName string
	addr        string
}

func InitRegisterServer(serverAddr string, serviceName string, cli *clientv3.Client) {
	fmt.Printf("grpcservice address: %s\n", serverAddr)
	l := EtcdServer{}
	l.cli = cli
	l.serviceName = serviceName
	l.addr = serverAddr

	l.register(serviceName, serverAddr, 5)

	//关闭信号处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		l.unRegister(serviceName, serverAddr)
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
}

// 将服务地址注册到etcd中
func (l *EtcdServer) register(serviceName string, serverAddr string, ttl int64) error {

	//与etcd建立长连接，并保证连接不断(心跳检测)
	ticker := time.NewTicker(time.Second * time.Duration(ttl))
	go func() {
		key := "/" + schema + "/" + serviceName + "/" + serverAddr
		for {
			resp, err := l.cli.Get(context.Background(), key)
			//fmt.Printf("resp:%+v\n", resp)
			if err != nil {
				fmt.Printf("获取服务地址失败：%s", err)
			} else if resp.Count == 0 { //尚未注册
				err = l.keepAlive(serviceName, serverAddr, ttl)
				if err != nil {
					fmt.Printf("保持连接失败：%s", err)
				}
			}
			<-ticker.C
		}
	}()

	return nil
}

// 保持服务器与etcd的长连接
func (l *EtcdServer) keepAlive(serviceName, serverAddr string, ttl int64) error {
	//创建租约
	leaseResp, err := l.cli.Grant(context.Background(), ttl)
	if err != nil {
		fmt.Printf("创建租期失败：%s\n", err)
		return err
	}

	//将服务地址注册到etcd中
	key := "/" + schema + "/" + serviceName + "/" + serverAddr
	_, err = l.cli.Put(context.Background(), key, serverAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		fmt.Printf("注册服务失败：%s", err)
		return err
	}

	//建立长连接
	ch, err := l.cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		fmt.Printf("建立长连接失败：%s\n", err)
		return err
	}

	//清空keepAlive返回的channel
	go func() {
		for {
			<-ch
		}
	}()
	return nil
}

// 取消注册
func (l *EtcdServer) unRegister(serviceName, serverAddr string) {
	if l.cli != nil {
		key := "/" + schema + "/" + serviceName + "/" + serverAddr
		l.cli.Delete(context.Background(), key)
	}
}
