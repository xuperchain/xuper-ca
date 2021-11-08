package service

import (
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"strings"
	"time"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"github.com/xuperchain/xuper-ca/config"
)

type GMOriginalCert struct {
	SerialNum  string
	CaCert     *x509.Certificate
	Cert       *x509.Certificate
	PrivateKey *sm2.PrivateKey
	ValidTime  int
}

func GenerateGMCert(caCert *GMOriginalCert, net string, root bool, address string) (*Cert, error) {
	validTime := time.Now().AddDate(10, 0, 0)
	rd.Seed(time.Now().UnixNano())
	serialNum := big.NewInt(rd.Int63())
	//isCA := false
	if root {
		net = net + ".root"
		//isCA = true
	}

	//[ v3_req ]
	//basicConstraints = CA:FALSE
	//keyUsage = nonRepudiation, digitalSignature

	//[ v3enc_req ]
	//basicConstraints = CA:FALSE
	//keyUsage = keyAgreement, keyEncipherment, dataEncipherment

	// cert := &x509.Certificate{
	// 	SerialNumber: serialNum, //证书序列号
	// 	Subject: pkix.Name{ // 证书的当前身份
	// 		Country:      []string{"XCHAIN"},
	// 		SerialNumber: address,
	// 		//	//Organization:       []string{"BD"},
	// 		//	//OrganizationalUnit: []string{"BD"},
	// 		//	//Province:           []string{"BJ"},
	// 		CommonName: net,

	// 		//	//Locality:           []string{"BJ"},
	// 	},
	// 	SignatureAlgorithm:    x509.SM2WithSM3,
	// 	NotBefore:             time.Now(), //证书有效期开始时间
	// 	NotAfter:              validTime,  //证书有效期结束时间
	// 	BasicConstraintsValid: true,       //基本的有效性约束
	// 	IsCA:                  true,       //是否是根证书
	// 	// ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	// 	//ExtKeyUsage:           []sm2.ExtKeyUsage{sm2.ExtKeyUsageClientAuth, sm2.ExtKeyUsageServerAuth},            //证书用途(客户端认证，数据加密)
	// 	//KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageCertSign, //产生密钥对的作用

	// 	//KeyUsageDigitalSignature/KeyUsageContentCommitment
	// 	//KeyUsageDataEncipherment/KeyUsageKeyEncipherment/KeyUsageKeyAgreement

	// 	KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement, //产生密钥对的作用(加密证书)
	// 	//KeyUsage: sm2.KeyUsageCertSign | sm2.KeyUsageDigitalSignature , //产生密钥对的作用(签名证书)
	// 	//EmailAddresses:        []string{"xchain-help@baidu.com"},
	// 	//IPAddresses:    []net.IP{net.ParseIP("127.0.0.1")},
	// 	//DNSNames: []string{"ca.server.com"},
	// }
	template := &x509.Certificate{
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
		NotBefore:             time.Now(), //证书有效期开始时间
		NotAfter:              validTime,  //证书有效期结束时间
		BasicConstraintsValid: true,       //基本的有效性约束
		IsCA:                  true,       //是否是根证书
		//ExtKeyUsage:           []sm2.ExtKeyUsage{sm2.ExtKeyUsageClientAuth, sm2.ExtKeyUsageServerAuth},
		//ExtKeyUsage:           []sm2.ExtKeyUsage{sm2.ExtKeyUsageClientAuth, sm2.ExtKeyUsageServerAuth},            //证书用途(客户端认证，数据加密)
		//KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageCertSign, //产生密钥对的作用

		//KeyUsageDigitalSignature/KeyUsageContentCommitment
		//KeyUsageDataEncipherment/KeyUsageKeyEncipherment/KeyUsageKeyAgreement

		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment | x509.KeyUsageCertSign | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement, //产生密钥对的作用(加密证书)
		//KeyUsage: sm2.KeyUsageCertSign | sm2.KeyUsageDigitalSignature , //产生密钥对的作用(签名证书)
		//EmailAddresses:        []string{"xchain-help@baidu.com"},
		//IPAddresses:    []net.IP{net.ParseIP("127.0.0.1")},
		//DNSNames: []string{"ca.server.com"},
		SignatureAlgorithm: x509.SM2WithSM3,
	}
	//生成公钥私钥对
	priKey, err := sm2.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	pubKey, _ := priKey.Public().(*sm2.PublicKey)

	var certpem []byte
	if caCert == nil {
		// 自签发证书
		certpem, err = x509.CreateCertificateToPem(template, &x509.Certificate{}, pubKey, priKey)
	} else {
		certpem, err = x509.CreateCertificateToPem(template, caCert.Cert, pubKey, caCert.PrivateKey)
	}
	if err != nil {
		return nil, err
	}

	privPem, err := x509.WritePrivateKeyToPem(priKey, nil)
	if err != nil {
		return nil, err
	}

	return &Cert{
		SerialNum:  template.SerialNumber.String(),
		Cert:       string(certpem),
		PrivateKey: string(privPem),
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

func GetRootGMCert() (*GMOriginalCert, error) {
	path := config.GetCertPath()
	if strings.LastIndex(path, "/") != len([]rune(path))-1 {
		path = path + "/gm/"
	} else {
		path = path + "gm/"
	}

	// 网络名为root时是指从文件中加载caserver的根管理员
	caFile, err := ioutil.ReadFile(path + CERT_NAME)
	if err != nil {
		return nil, err
	}
	cacert := string(caFile)

	// 解析证书
	caBlock, _ := pem.Decode([]byte(cacert))
	certificate, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, err
	}
	//解析私钥
	keyFile, err := ioutil.ReadFile(path + PRIVATEKEY_NAME)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ReadPrivateKeyFromPem(keyFile, nil)
	if err != nil {
		return nil, err
	}

	return &GMOriginalCert{
		Cert:       certificate,
		PrivateKey: privateKey,
	}, nil
}
