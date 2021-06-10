/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package service

import (
	"github.com/xuperchain/xuper-ca/config"
	"github.com/xuperchain/xuper-ca/dao"
)

// 判断address是否有节点权限
func CheckNode(address, net string) bool {
	if config.GetCaAdmin() == address {
		return true
	}

	// 校验是不是节点自己
	nodeDao := &dao.NodeDao{}
	node := nodeDao.QueryValidNodeByNetAndAddress(net, address)
	if node != nil {
		return true
	}

	// 校验是不是节点的网络管理员
	node = nodeDao.QueryValidNodeByNetAndAdmin(net, address)
	if node != nil {
		return true
	}

	return false
}

// 判断address是否有网络管理员权限
func CheckNetAdmin(address, net string) bool {
	if config.GetCaAdmin() == address {
		return true
	}

	netDao := dao.NetAdminDao{}
	admin := netDao.GetNetAdmin(net, address)
	if admin != nil {
		return true
	}
	return false
}

// 判断address是否有ca根权限
func CheckCaAdmin(address string) bool {
	if config.GetCaAdmin() != address {
		return false
	}
	return true
}
