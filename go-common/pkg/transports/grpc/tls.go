package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"os"
	"path/filepath"
)

var (
	//port               = ":50051"
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata!")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token, 调用认证失败!")
)

func tlsClient(crtFile string, keyFile string, caFile string, serverName string) (credentials.TransportCredentials, error) {
	if !filepath.IsAbs(crtFile) {
		crtFile = filepath.Join("configs", crtFile)
	}
	if !filepath.IsAbs(keyFile) {
		keyFile = filepath.Join("configs", keyFile)
	}
	if !filepath.IsAbs(caFile) {
		caFile = filepath.Join("configs", caFile)
	}
	crtFile = filepath.Clean(crtFile)
	keyFile = filepath.Clean(keyFile)
	caFile = filepath.Clean(caFile)

	// TLS认证
	//从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	//key, _ := tls.LoadX509KeyPair("certs/server.pem","certs/server.key")
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool() //初始化一个CertPool
	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPool.AppendCertsFromPEM(ca) //解析传入的证书，解析成功会将其加到池子中
	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, //客户端证书
		ServerName:   serverName,              //"www.p-pp.cn",注意这里的参数为配置文件中所允许的ServerName，也就是其中配置的DNS...
		RootCAs:      certPool,
	})
	return cred, nil
}
func tlsServer(crtFile string, keyFile string, caFile string) (credentials.TransportCredentials, error) {
	if !filepath.IsAbs(crtFile) {
		crtFile = filepath.Join("configs", crtFile)
	}
	if !filepath.IsAbs(keyFile) {
		keyFile = filepath.Join("configs", keyFile)
	}
	if !filepath.IsAbs(caFile) {
		caFile = filepath.Join("configs", caFile)
	}
	crtFile = filepath.Clean(crtFile)
	keyFile = filepath.Clean(keyFile)
	caFile = filepath.Clean(caFile)

	// TLS认证
	//从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	//key, _ := tls.LoadX509KeyPair("certs/server.pem","certs/server.key")
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool() //初始化一个CertPool
	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPool.AppendCertsFromPEM(ca)         //解析传入的证书，解析成功会将其加到池子中
	cred := credentials.NewTLS(&tls.Config{ //构建基于TLS的TransportCredentials选项
		Certificates: []tls.Certificate{cert},        //服务端证书链，可以有多个
		ClientAuth:   tls.RequireAndVerifyClientCert, //要求必须验证客户端证书
		ClientCAs:    certPool,                       //设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})
	return cred, nil
}
