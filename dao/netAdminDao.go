/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package dao

import (
	"log"
	"time"
)

type NetAdmin struct {
	Id           int    `db:"id"`
	Net          string `db:"net"`
	Address      string `db:"address"`
	SerialNum    string `db:"serial_num"`
	Cert         string `db:"cert"`
	PrivateKey   string `db:"private_key"`
	CreateTime   int    `db:"create_time"`
	UpdateTime   int    `db:"update_time"`
	IsValid      bool   `db:"is_valid"`
	ValidTime    int    `db:"valid_time"`
	HdPrivateKey string `db:"hd_private_key"`
}

// ca内的网络管理数据层
type NetAdminDao struct {
}

func (netAdminDao *NetAdminDao) Insert(netAdmin *NetAdmin) (int64, error) {
	// 需要校验吗?

	netAdmin.CreateTime = int(time.Now().Unix())
	netAdmin.UpdateTime = int(time.Now().Unix())

	caDb := GetDbInstance()
	result, err := caDb.db.Exec(
		"INSERT INTO net_admin(`net`, `address`, `serial_num`, `cert`, `private_key`, `create_time`, `update_time`, `is_valid`, "+
			"`valid_time`, `hd_private_key`) VALUES (?,?,?,?,?,?,?,?,?,?)",
		netAdmin.Net,
		netAdmin.Address,
		netAdmin.SerialNum,
		netAdmin.Cert,
		netAdmin.PrivateKey,
		netAdmin.CreateTime,
		netAdmin.UpdateTime,
		netAdmin.IsValid,
		netAdmin.ValidTime,
		netAdmin.HdPrivateKey)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
	}
	return id, err
}

func (netAdminDao *NetAdminDao) GetNetAdmin(net, adminAddress string) *NetAdmin {
	var netAdmin NetAdmin
	caDb := GetDbInstance()
	err := caDb.db.Get(&netAdmin,
		"SELECT * FROM net_admin WHERE net=? and address=? and is_valid=?", net, adminAddress, true)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &netAdmin
}
