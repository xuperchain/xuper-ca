package service_test

import (
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"testing"
	"time"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"github.com/xuperchain/xuper-ca/crypto"
	"github.com/xuperchain/xuper-ca/service"
)

// 测试 GM 生成证书
func Test_GenertRootCert_GM(t *testing.T) {

	address := "TeyyPLpp9L7QAcxHangtcHTu7HUZ6iydY"
	net := "test"

	validTime := time.Now().AddDate(10, 0, 0)
	rd.Seed(time.Now().UnixNano())
	serialNum := big.NewInt(rd.Int63())

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
		NotBefore:             time.Now(), //证书有效期开始时间
		NotAfter:              validTime,  //证书有效期结束时间
		BasicConstraintsValid: true,       //基本的有效性约束
		IsCA:                  true,       //是否是根证书
		//ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},            //证书用途(客户端认证，数据加密)
		//KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageCertSign, //产生密钥对的作用
		//EmailAddresses:        []string{"xchain-help@baidu.com"},
		//IPAddresses:    []net.IP{net.ParseIP("127.0.0.1")},
		//DNSNames: []string{"ca.server.com"},
	}

	// 通过随机数种子生成sm2私钥
	pk, err := sm2.GenerateKey(nil)
	if err != nil {
		fmt.Println("生成sm2私钥失败", err.Error())
		return
	}

	var caCert *service.GMOriginalCert

	// 生成证书
	var certByte []byte
	if caCert == nil {
		// 自签发证书
		certByte, err = x509.CreateCertificate(cert, &x509.Certificate{}, &pk.PublicKey, pk)
	} else {
		// 有root证书 生成中间证书 或者节点证书
		certByte, err = x509.CreateCertificate(cert, caCert.Cert, &pk.PublicKey, caCert.PrivateKey)
	}
	if err != nil {
		fmt.Println("生成证书失败", err.Error())
	}
	//编码证书文件和私钥文件
	caPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certByte,
	}
	cacert := pem.EncodeToMemory(caPem)

	buf, err := x509.MarshalSm2PrivateKey(pk, nil)
	if err != nil {
		fmt.Println("MarshalECPrivateKey failed", err.Error())
	}
	keyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: buf,
	}
	key := pem.EncodeToMemory(keyPem)
	fmt.Println(&service.Cert{
		SerialNum:  cert.SerialNumber.String(),
		Cert:       string(cacert),
		PrivateKey: string(key),
		ValidTime:  int(validTime.Unix()),
	})
	// return &service.Cert{
	// 	SerialNum:  cert.SerialNumber.String(),
	// 	Cert:       string(cacert),
	// 	PrivateKey: string(key),
	// 	ValidTime:  int(validTime.Unix()),
	// }, nil

	err = service.WriteCert("./cert/", &service.Cert{
		SerialNum:  cert.SerialNumber.String(),
		Cert:       string(cacert),
		PrivateKey: string(key),
		ValidTime:  int(validTime.Unix()),
	})
	if err != nil {
		fmt.Println("write root cert failed, ", err)
		return
	}
	fmt.Println("init root cert success")

}

// 测试有root 证书 生成中间证书或者节点证书
func Test_GenertCert_GM(t *testing.T) {

	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	net := "test"
	validTime := time.Now().AddDate(10, 0, 0)
	rd.Seed(time.Now().UnixNano())
	serialNum := big.NewInt(rd.Int63())

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
		NotBefore:             time.Now(), //证书有效期开始时间
		NotAfter:              validTime,  //证书有效期结束时间
		BasicConstraintsValid: true,       //基本的有效性约束
		IsCA:                  true,       //是否是根证书
		//ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},            //证书用途(客户端认证，数据加密)
		//KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageCertSign, //产生密钥对的作用
		//EmailAddresses:        []string{"xchain-help@baidu.com"},
		//IPAddresses:    []net.IP{net.ParseIP("127.0.0.1")},
		//DNSNames: []string{"ca.server.com"},
	}

	// 通过随机数种子生成sm2私钥
	pk, err := sm2.GenerateKey(nil)
	if err != nil {
		fmt.Println("生成sm2私钥失败", err.Error())
		return
	}

	caCert, err := GetRootCert()
	if err != nil {
		fmt.Println("获取根证书失败", err.Error())
		return
	}

	// 生成证书
	var certByte []byte
	if caCert == nil {
		// 自签发证书
		certByte, err = x509.CreateCertificate(cert, &x509.Certificate{}, &pk.PublicKey, pk)
	} else {
		// 有root证书 生成中间证书 或者节点证书
		certByte, err = x509.CreateCertificate(cert, caCert.Cert, &pk.PublicKey, caCert.PrivateKey)
	}
	if err != nil {
		fmt.Println("生成证书失败", err.Error())
	}
	//编码证书文件和私钥文件
	caPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certByte,
	}
	cacert := pem.EncodeToMemory(caPem)

	buf, err := x509.MarshalSm2PrivateKey(pk, nil)
	if err != nil {
		fmt.Println("MarshalECPrivateKey failed", err.Error())
	}
	keyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: buf,
	}
	key := pem.EncodeToMemory(keyPem)
	fmt.Println(&service.Cert{
		SerialNum:  cert.SerialNumber.String(),
		Cert:       string(cacert),
		PrivateKey: string(key),
		ValidTime:  int(validTime.Unix()),
	})
	// return &service.Cert{
	// 	SerialNum:  cert.SerialNumber.String(),
	// 	Cert:       string(cacert),
	// 	PrivateKey: string(key),
	// 	ValidTime:  int(validTime.Unix()),
	// }, nil

	err = service.WriteCert("./middle/", &service.Cert{
		SerialNum:  cert.SerialNumber.String(),
		Cert:       string(cacert),
		PrivateKey: string(key),
		ValidTime:  int(validTime.Unix()),
	})
	if err != nil {
		fmt.Println("write root cert failed, ", err)
		return
	}
	fmt.Println("init root cert success")

}

// 生成网络根私钥 测试
func Test_GenerateNetHdPriKey(t *testing.T) {
	gmCryptoClient := crypto.GetGMHdCryptoClient()
	ecdsaAccount, err := gmCryptoClient.CreateNewAccountWithMnemonic(service.Language, service.StrengthHard)
	if err != nil {
		fmt.Println("CreateNewAccountWithMnemonic failed and err is:", err.Error())
		return
	}
	fmt.Printf("mnemonic is %v, Address: %v , jsonPrivateKey: %v, jsonPublicKey: %v", ecdsaAccount.Mnemonic, ecdsaAccount.Address, ecdsaAccount.JsonPrivateKey, ecdsaAccount.JsonPublicKey)
	jsonMasterKey, err := gmCryptoClient.GenerateMasterKeyByMnemonic(ecdsaAccount.Mnemonic, service.Language)
	if err != nil {
		fmt.Printf("GenerateMasterKeyByMnemonic failed and err is: %v", err)
		return
	}
	fmt.Println(jsonMasterKey)
}

// 测试生成节点私钥（根据网络根私钥）
func Test_GenerateNodeHdPriKey(t *testing.T) {

	gmCryptoClient := crypto.GetGMHdCryptoClient()
	netHdPriKey := "{\"Version\":\"BCC5AA==\",\"Depth\":0,\"ParentFP\":\"AAAAAA==\",\"ChildNum\":0,\"ChainCode\":\"bd2Cp1f51nDLc6F07IYgR+x8wsnX411LdvV4VRPG1QE=\",\"Key\":\"gXzRX5O6zUmftVnZdvnTSH3LHhlHpxi2zZcuSwAMrJQ=\",\"PubKey\":null,\"AccountNum\":{\"0\":0},\"Cryptography\":2,\"IsPrivate\":true}"
	// 兼容旧网络
	if netHdPriKey == "" {
		return
	}
	nodeHdKeyStart := service.HardenedKeyStart + uint32(1)

	childHdKey, err := gmCryptoClient.GenerateChildKey(netHdPriKey, nodeHdKeyStart)
	if err != nil {
		fmt.Printf("GenerateChildKey failed and err is: %v", err)
		return
	}
	fmt.Println(childHdKey)

}

// 获取根证书
func GetRootCert() (*service.GMOriginalCert, error) {
	// path := config.GetCertPath()
	// if strings.LastIndex(path, "/") != len([]rune(path))-1 {
	// 	path = path + "/"
	// }

	// 网络名为root时是指从文件中加载caserver的根管理员
	caFile, err := ioutil.ReadFile("./cert/" + service.CERT_NAME)
	if err != nil {
		fmt.Println("get caserver cert failed", err.Error())
		return nil, err
	}
	cacert := string(caFile)

	// 解析证书
	caBlock, _ := pem.Decode([]byte(cacert))
	certificate, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		fmt.Println("解析证书失败", err.Error())
		return nil, err
	}

	//解析私钥
	keyFile, err := ioutil.ReadFile("./cert/" + service.PRIVATEKEY_NAME)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ReadPrivateKeyFromPem(keyFile, nil)
	if err != nil {
		fmt.Println("解析私钥失败", err.Error())
		return nil, err
	}

	return &service.GMOriginalCert{
		Cert:       certificate,
		PrivateKey: privateKey,
	}, nil
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
