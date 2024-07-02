
# openssl下载网址
```azure
http://slproweb.com/products/Win32OpenSSL.html
```

# 生成根证书/私钥、服务端证书/私钥、客户端证书/私钥

## 方法一

#### 创建根证书
- 创建一个CA私钥（根证书）
```
openssl genrsa -out ca.key 4096
```

- 创建一个conf 用来生成csr（请求签名证书文件）---ca.conf
```azure
cat > ca.conf << EOF
[req]
default_bits = 4096
distinguished_name = req_distinguished_name

[req_distinguished_name]
countryName = Country Name (2 letter code)
countryName_default=CN
stateOrProvinceName=State or Province Name (2 letter code)
stateOrProvinceName_default=BeiJing
localityName = Locality Name (eg, city)
localityName_default=BeiJing
organizationName=Organization Name (eg, company)
organizationName_default=915zb.com
commonName = Common Name (e.g. server FQDN or YOUR name)
commonName_max=64
commonName_default=localhost
EOF
```

- 生成csr文件
```
openssl req -new -sha256 -out ca.csr -key ca.key -config ca.conf
```
>该命令含义如下：
>req——执行证书签发命令
>-new——新证书签发请求
>-key——指定私钥路径
>-out——输出的csr文件的路径


- 生成证书 crt文件
```
openssl x509 -req -days 365 -in ca.csr -signkey ca.key -out ca.crt
```
>该命令的含义如下：
>x509——生成x509格式证书
>-req——输入csr文件
>-days——证书的有效期（天）
>-signkey——签发证书的私钥
>-in——要输入的csr文件
>-out——输出的cer证书文件


- crt转pem

```azure
openssl x509 -in ca.crt -out ca.pem -outform PEM
```

#### 生成服务器证书

- 生成server证书 - 新建server.conf
> commonName_default 客户端要根据这个字段做匹配
```
cat > server.conf << EOF
[req]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[req_distinguished_name]
countryName = Country Name (2 letter code)
countryName_default = CN
stateOrProvinceName = State or Province Name (2 letter code)
stateOrProvinceName_default = BeiJing
localityName = Locality Name (eg, city)
localityName_default = BeiJing
organizationName = Organization Name (eg, company)
organizationName_default = 915zb.com
commonName = Common Name (e.g. server FQDN or YOUR name)
commonName_max = 64
commonName_default = localhost
[ req_ext ]
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
DNS.2 = 915zb.com
DNS.3 = auth.915zb.com
DNS.4 = www.915zb.com
DNS.5 = design.915zb.com
DNS.6 = nianxi.915zb.com
IP  = 127.0.0.1
EOF
```
>[ req ]
>default_bits       = 2048
>distinguished_name = req_distinguished_name
>req_extensions     = req_ext

- 生成server.key
```
openssl genrsa -out server.key 2048
```

- 生成server.csr
```
openssl req -new -sha256 -out server.csr -key server.key -config server.conf
```

- 生成server.crt/pem
```
openssl x509 -req -days 365 -CA ca.pem -CAkey ca.key -CAcreateserial -in server.csr -out server.pem -extensions req_ext -extfile server.conf
```
>-CA——指定CA证书的路径
>-CAkey——指定CA证书的私钥路径
>-CAserial——指定证书序列号文件的路径
>-CAcreateserial——表示创建证书序列号文件(即上方提到的serial文件)，创建的序列号文件默认名称为-CA，指定的证书名称后加上.srl后缀


- 提取公钥
```azure
openssl rsa -in server.key   -pubout -out public.pem
```

- 转crt
```azure
openssl x509 -outform der -in server.pem -out server.crt
```

#### 生成客端证书

- 生成client.key
```
openssl ecparam -genkey -name secp384r1 -out client.key
```

- 生成client.csr
```
openssl req -new -key client.key -out client.csr -config server.conf
```

- 生成client.pem
```
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem -extensions req_ext -extfile server.conf
```
或
```
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem -extensions req_ext -extfile server.conf
```

# 生成http证书
 