/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package service

import (
	"encoding/pem"
	"errors"

	log "github.com/sirupsen/logrus"

	"crypto/x509"

	tjx509 "github.com/tjfoc/gmsm/x509"
	"github.com/xuperchain/xuper-ca/dao"
)

// 添加一个网络和网络管理员
func AddNetAdmin(net, address string, crypto string) error {
	if net == "" || address == "" {
		return ErrParam
	}
	log.Infof("GetNetAdmin: %+v", address)
	netAdminDao := dao.NetAdminDao{}
	// 校验网络是否已经存在
	netAdmin := netAdminDao.GetNetAdmin(net, address)
	if netAdmin != nil {
		log.Warnf("net admin is existed, net:%v, address:%v", net, address)
		return nil
		//return ErrNetExisted
	}
	nodeCert := &Cert{}
	var netHdPrikey string
	// 如果生成 gm 网络
	if crypto == "gm" {
		// 生成管理员证书
		gmRootCert, err := GetRootGMCert()
		if err != nil {
			return err
		}
		nodeCert, err = GenerateGMCert(gmRootCert, net, true, address)
		if err != nil {
			return err
		}
		// 生成网络根私钥
		netHdPrikey, err := GenerateGMNetHdPriKey()
		if err != nil {
			return err
		}
		log.Infof("GenerateGMNetHdPriKey:   %+v", netHdPrikey)
	} else {
		// 生成网络管理员证书
		rootCert, err := GetRootCert()
		if err != nil {
			return err
		}
		nodeCert, err = GenerateCert(rootCert, net, true, address)
		if err != nil {
			return err
		}
		// 生成网络根私钥
		netHdPrikey, err := GenerateNetHdPriKey()
		if err != nil {
			return err
		}
		log.Infof("GenerateNetHdPriKey:   %+v", netHdPrikey)
	}

	// 保存网络管理员
	_, err := netAdminDao.Insert(&dao.NetAdmin{
		Net:          net,
		Address:      address,
		SerialNum:    nodeCert.SerialNum,
		Cert:         nodeCert.Cert,
		PrivateKey:   nodeCert.PrivateKey,
		IsValid:      true,
		ValidTime:    nodeCert.ValidTime,
		HdPrivateKey: netHdPrikey,
		Crypto:       crypto,
	})
	if err != nil {
		return ErrDB
	}
	log.Infof("AddNetAdmin success")
	return nil
}

// 添加一个节点
func AddNode(net, adminAddress, address string) error {
	log.Infof("AddNode, net: %+v, address: %+v", net, address)
	// 校验参数
	if net == "" || address == "" {
		return ErrParam
	}
	netAdminDao := dao.NetAdminDao{}
	nodeDao := dao.NodeDao{}

	node := nodeDao.QueryValidNodeByNetAndAddress(net, address)
	if node != nil {
		log.Warnf("node is existed, net: %v, address: %v", net, address)
		return nil
		//return errors.New("node is existed")
	}

	// 校验网络是否已经存在
	netAdmin := netAdminDao.GetNetAdmin(net, adminAddress)
	if netAdmin == nil {
		return ErrCACert
	}

	// 获取nodeHdKeyStart
	nodeTotal, err := nodeDao.QueryTotalNode(net, adminAddress)
	if err != nil {
		log.Warnf("node total query error, error: %v", err)
		return nil
	}
	nodeCert := &Cert{}
	var nodeHdPriKey string
	// 根据网络的加密方式生成节点证书
	log.Info("netadmin.....", netAdmin)
	if netAdmin.Crypto == "gm" {
		// 获取网络管理员的GM证书
		rootGMCert, _, err := GetAdminGMCert(net, adminAddress)
		if err != nil {
			return err
		}
		// 生成节点GM证书
		nodeCert, err = GenerateGMCert(rootGMCert, net, false, address)
		if err != nil {
			return err
		}
		// 依据网络根私钥生成节点的GM根私钥
		nodeHdPriKey, err = GenerateNodeGMHdPriKey(nodeTotal, netAdmin.HdPrivateKey)
		if err != nil {
			return err
		}
	} else {
		// 获取网络管理员证书
		rootCert, _, err := GetAdminCert(net, adminAddress)
		if err != nil {
			return err
		}
		nodeCert, err = GenerateCert(rootCert, net, false, address)
		if err != nil {
			return err
		}
		// 依据网络根私钥生成节点根私钥
		nodeHdPriKey, err = GenerateNodeHdPriKey(nodeTotal, netAdmin.HdPrivateKey)
		if err != nil {
			return err
		}
	}

	// 保存节点
	_, err = nodeDao.Insert(&dao.Node{
		Net:          net,
		AdminAddress: adminAddress,
		Address:      address,
		SerialNum:    nodeCert.SerialNum,
		Cert:         nodeCert.Cert,
		PrivateKey:   nodeCert.PrivateKey,
		IsValid:      true,
		ValidTime:    nodeCert.ValidTime,
		HdPrivateKey: nodeHdPriKey,
	})

	if err != nil {
		return ErrDB
	}
	log.Infof("AddNetAdmin success")
	return nil
}

// 获取网络管理员的证书
func GetAdminCert(net, adminAddress string) (*OriginalCert, string, error) {
	log.Infof("GetAdminCert, net: %+v", net)

	var cacert string
	var key string
	var address string

	// 获取网络管理员的证书
	// 网络名为其他时指各个网络的管理员
	netAdminDao := &dao.NetAdminDao{}
	admin := netAdminDao.GetNetAdmin(net, adminAddress)
	if admin == nil {
		return nil, "", errors.New("can't get net root ca, net is " + net)
	}
	cacert = admin.Cert
	key = admin.PrivateKey
	address = admin.Address

	// 解析证书
	caBlock, _ := pem.Decode([]byte(cacert))
	certificate, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, "", err
	}

	//解析私钥
	keyBlock, _ := pem.Decode([]byte(key))
	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, "", err
	}

	log.Infof("GetAdminCert success")
	return &OriginalCert{
		Cert:       certificate,
		PrivateKey: privateKey,
	}, address, nil
}

// 获取网络管理员的证书（GM）
func GetAdminGMCert(net, adminAddress string) (*GMOriginalCert, string, error) {
	log.Infof("GetAdminCert, net: %+v", net)

	var cacert string
	var key string
	var address string

	// 获取网络管理员的证书
	// 网络名为其他时指各个网络的管理员
	netAdminDao := &dao.NetAdminDao{}
	admin := netAdminDao.GetNetAdmin(net, adminAddress)
	if admin == nil {
		return nil, "", errors.New("can't get net root ca, net is " + net)
	}
	cacert = admin.Cert
	key = admin.PrivateKey
	address = admin.Address

	// 解析证书
	caBlock, _ := pem.Decode([]byte(cacert))
	certificate, err := tjx509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, "", err
	}
	//解析私钥

	privateKey, err := tjx509.ReadPrivateKeyFromPem([]byte(key), nil)
	if err != nil {
		return nil, "", err
	}

	log.Infof("GetAdminCert success")
	return &GMOriginalCert{
		Cert:       certificate,
		PrivateKey: privateKey,
	}, address, nil
}

// 获取节点的证书
func GetNode(net, address string) (*Cert, string, error) {
	log.Infof("GetNode, net: %+v", net)
	nodeDao := dao.NodeDao{}
	ret := nodeDao.QueryValidNodeByNetAndAddress(net, address)
	if ret == nil {
		return nil, "", ErrCertNoExisted
	}
	adminDao := dao.NetAdminDao{}
	adminCert := adminDao.GetNetAdmin(net, ret.AdminAddress)
	if adminCert == nil {
		return nil, "", ErrCACert
	}

	log.Infof("GetNode success")
	return &Cert{
		Address:    address,
		SerialNum:  ret.SerialNum,
		Cert:       ret.Cert,
		PrivateKey: ret.PrivateKey,
		ValidTime:  ret.ValidTime,
		CaCert:     adminCert.Cert,
	}, ret.HdPrivateKey, nil
}

// 获取增量撤销列表
func GetRevokeList(net, latestSerialNum string) (*[]dao.Revoke, error) {
	log.Infof("GetRevokeList, cursorLeastId: %+v", latestSerialNum)
	revokeDao := dao.RevokeDao{}
	ret, err := revokeDao.GetList(net, latestSerialNum)
	if err != nil {
		return nil, err
	}
	log.Infof("GetRevokeList success")
	return ret, nil
}

// 撤销一个节点
func RevokeNode(net, address string) (bool, error) {
	log.Infof("RevokeNode, net: %+v, address: %+v", net, address)
	nodeDao := dao.NodeDao{}
	ret, err := nodeDao.RevokeNodeByNetAndAddress(net, address)
	if ret == true {
		log.Infof("RevokeNode success")
	} else {
		log.Warning("RevokeNode falied, err:", err)
	}
	return ret, err
}

// 解密一笔交易
func DecryptByHdKey(net, adminAddress, hdPubKey, cypherText string) (string, error) {
	log.Infof("EcryptByHdKey, net: %+v, address: %+v, hdPubKey: %+v, cypherText: %+v", net, adminAddress, hdPubKey, []byte(cypherText))
	netAdminDao := &dao.NetAdminDao{}
	admin := netAdminDao.GetNetAdmin(net, adminAddress)
	if admin == nil || admin.HdPrivateKey == "" {
		return "", ErrParam
	}
	ret, err := DecryptByNetHdPriKey(admin.HdPrivateKey, hdPubKey, cypherText)
	if err != nil {
		log.Warning("EcryptByNetHdPriKey, err:", err)
		return "", err
	}
	return ret, err
}
