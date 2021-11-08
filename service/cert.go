/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"os"
	"strings"
	"time"

	"github.com/xuperchain/xuper-ca/config"
	"github.com/xuperchain/xuper-ca/util"
)

const CERT_NAME = "root.crt"
const PRIVATEKEY_NAME = "root.key"

type Cert struct {
	Address    string
	SerialNum  string
	Cert       string
	PrivateKey string
	CaCert     string
	ValidTime  int
}

type OriginalCert struct {
	SerialNum  string
	CaCert     *x509.Certificate
	Cert       *x509.Certificate
	PrivateKey *rsa.PrivateKey
	ValidTime  int
}

// 生成证书
func GenerateCert(caCert *OriginalCert, net string, root bool, address string) (*Cert, error) {
	validTime := time.Now().AddDate(10, 0, 0)
	rd.Seed(time.Now().UnixNano())
	serialNum := big.NewInt(rd.Int63())

	if root {
		net = net + ".root"
	}

	cert := &x509.Certificate{
		SerialNumber: serialNum, //证书序列号
		Subject: pkix.Name{ // 证书的当前身份
			Country:      []string{"XCHAIN"},
			SerialNumber: address,
			//	//Organization:       []string{"BD"},
			//	//OrganizationalUnit: []string{"BD"},
			//	//Province:           []string{"BJ"},
			CommonName: net,
			//	//Locality:           []string{"BJ"},
		},
		NotBefore:             time.Now(),                                                                 //证书有效期开始时间
		NotAfter:              validTime,                                                                  //证书有效期结束时间
		BasicConstraintsValid: true,                                                                       //基本的有效性约束
		IsCA:                  true,                                                                       //是否是根证书
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途(客户端认证，数据加密)
		//KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageCertSign, //产生密钥对的作用
		//EmailAddresses:        []string{"xchain-help@baidu.com"},
		//IPAddresses:    []net.IP{net.ParseIP("127.0.0.1")},
		//DNSNames: []string{"ca.server.com"},
	}

	//生成公钥私钥对
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	var certByte []byte
	if caCert == nil {
		// 自签发证书
		certByte, err = x509.CreateCertificate(rand.Reader, cert, &x509.Certificate{}, &priKey.PublicKey, priKey)
	} else {
		certByte, err = x509.CreateCertificate(rand.Reader, cert, caCert.Cert, &priKey.PublicKey, caCert.PrivateKey)
	}
	if err != nil {
		return nil, err
	}
	//编码证书文件和私钥文件
	caPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certByte,
	}
	cacert := pem.EncodeToMemory(caPem)

	buf := x509.MarshalPKCS1PrivateKey(priKey)
	keyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: buf,
	}
	key := pem.EncodeToMemory(keyPem)

	return &Cert{
		SerialNum:  cert.SerialNumber.String(),
		Cert:       string(cacert),
		PrivateKey: string(key),
		ValidTime:  int(validTime.Unix()),
	}, nil

	//// 解析证书
	//caBlock, _ := pem.Decode([]byte(cacert))
	//certificate, err := x509.ParseCertificate(caBlock.Bytes)
	//if err != nil {
	//	return nil, err
	//}
	////解析私钥
	//keyBlock, _ := pem.Decode([]byte(key))
	//privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	//if err != nil {
	//	return nil, err
	//}

	//return &OriginalCert{
	//	SerialNum:"",
	//	Cert:certificate,
	//	PrivateKey:privateKey,
	//	ValidTime: int(validTime.Unix()),
	//}, nil
}

// 获取ca根证书
func GetRootCert() (*OriginalCert, error) {
	path := config.GetCertPath()
	if strings.LastIndex(path, "/") != len([]rune(path))-1 {
		path = path + "/default/"
	} else {
		path = path + "default/"
	}

	// 网络名为root时是指从文件中加载caserver的根管理员
	caFile, err := ioutil.ReadFile(path + CERT_NAME)
	if err != nil {
		return nil, err
	}
	cacert := string(caFile)

	keyFile, err := ioutil.ReadFile(path + PRIVATEKEY_NAME)
	if err != nil {
		return nil, err
	}
	key := string(keyFile)
	// 解析证书
	caBlock, _ := pem.Decode([]byte(cacert))
	certificate, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, err
	}

	//解析私钥
	keyBlock, _ := pem.Decode([]byte(key))
	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &OriginalCert{
		Cert:       certificate,
		PrivateKey: privateKey,
	}, nil
}

// 将证书写入证书文件
func WriteCert(path string, cert *Cert) error {
	if strings.LastIndex(path, "/") != len([]rune(path))-1 {
		path = path + "/"
	}
	// 判断文件夹是否存在 不存在新建文件夹
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	err = util.WriteFileUsingFilename(path+PRIVATEKEY_NAME, []byte(cert.PrivateKey))
	if err != nil {
		return err
	}
	err = util.WriteFileUsingFilename(path+CERT_NAME, []byte(cert.Cert))
	if err != nil {
		return err
	}
	return nil
}
