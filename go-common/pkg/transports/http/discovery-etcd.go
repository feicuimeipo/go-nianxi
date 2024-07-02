package http

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"

	"time"
)

const schema = "ns"

type EtcdServer struct {
	cli *clientv3.Client
	ttl int64
}

func newEtcdServer(ttl int64, client *clientv3.Client) IDiscovery {
	r := EtcdServer{}
	r.cli = client
	r.ttl = ttl
	return &r
}

func (l *EtcdServer) client() interface{} {
	return l.cli
}

// 将服务地址注册到etcd中
func (l *EtcdServer) register(addr string, app string) error {

	//var err error
	//与etcd建立长连接，并保证连接不断(心跳检测)
	ticker := time.NewTicker(time.Second * time.Duration(l.ttl))
	go func() {
		key := "/" + schema + "/" + app + "/" + addr
		for {
			resp, err := l.cli.Get(context.Background(), key)
			//fmt.Printf("resp:%+v\n", resp)
			if err != nil {
				fmt.Printf("获取服务地址失败：%s", err)
			} else if resp.Count == 0 { //尚未注册
				err = keepAlive(app, addr, l.ttl, l.cli)
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
func keepAlive(serviceName, serverAddr string, ttl int64, cli *clientv3.Client) error {
	//创建租约
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		fmt.Printf("创建租期失败：%s\n", err)
		return err
	}

	//将服务地址注册到etcd中
	key := "/" + schema + "/" + serviceName + "/" + serverAddr
	_, err = cli.Put(context.Background(), key, serverAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		fmt.Printf("注册服务失败：%s", err)
		return err
	}

	//建立长连接
	ch, err := cli.KeepAlive(context.Background(), leaseResp.ID)
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
func (l *EtcdServer) deRegister(addr string, app string) error {
	if l.cli != nil {
		key := "/" + schema + "/" + app + "/" + addr
		l.cli.Delete(context.Background(), key)
	}
	return nil
}
