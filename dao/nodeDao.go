/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package dao

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	Id           int    `db:"id"`
	Net          string `db:"net"`
	AdminAddress string `db:"adminAddress"`
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

// ca网络下的节点管理数据层
type NodeDao struct {
}

func (nodeDao *NodeDao) Insert(node *Node) (int64, error) {
	node.CreateTime = int(time.Now().Unix())
	node.UpdateTime = int(time.Now().Unix())

	caDb := GetDbInstance()
	result, err := caDb.db.Exec(
		"INSERT INTO node(`net`, `adminAddress`, `address`, `serial_num`, `cert`, `private_key`, `create_time`, `update_time`, "+
			"`is_valid`, `valid_time`, `hd_private_key`) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		node.Net,
		node.AdminAddress,
		node.Address,
		node.SerialNum,
		node.Cert,
		node.PrivateKey,
		node.CreateTime,
		node.UpdateTime,
		node.IsValid,
		node.ValidTime,
		node.HdPrivateKey)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Warning(err)
	}
	return id, nil
}

func (nodeDao *NodeDao) QueryValidNodeByNetAndAddress(net, address string) *Node {
	caDb := GetDbInstance()
	ret := Node{}
	err := caDb.db.Get(&ret, "SELECT * FROM node WHERE is_valid=? AND net=? AND address=?", true, net, address)
	if err != nil {
		log.Warning(err)
		return nil
	}
	return &ret
}

func (nodeDao *NodeDao) QueryValidNodeByNetAndAdmin(net, adminAddress string) *Node {
	caDb := GetDbInstance()
	ret := Node{}
	err := caDb.db.Get(&ret, "SELECT * FROM node WHERE is_valid=? AND net=? AND adminAddress=?", true, net, adminAddress)
	if err != nil {
		log.Warning(err)
		return nil
	}
	return &ret
}

func (nodeDao *NodeDao) QueryTotalNode(net, adminAddress string) (uint32, error) {
	var total uint32 = 0
	caDb := GetDbInstance()
	err := caDb.db.QueryRow("SELECT COUNT(*) FROM node WHERE net=? AND adminAddress=? AND is_valid=?", net, adminAddress, true).Scan(&total)
	if err != nil {
		log.Warning(err)
		return 0, err
	}
	return total, nil
}

func (nodeDao *NodeDao) RevokeNodeByNetAndAddress(net, address string) (bool, error) {
	caDb := GetDbInstance()
	tx, _ := caDb.db.Begin()
	defer tx.Rollback()

	// 1.查出撤销节点的授权记录
	nodeRecords, err := tx.Query("SELECT * FROM node WHERE is_valid=? AND net=? AND address=?", true, net, address)
	if err != nil {
		log.Warning("RevokeNodeByNetAndAddress get authorization records failed")
		return false, err
	}

	// 2.更新revoke 表
	for nodeRecords.Next() {
		var node Node
		err := nodeRecords.Scan(&node.Id, &node.Net, &node.Address, &node.AdminAddress, &node.SerialNum, &node.Cert, &node.PrivateKey,
			&node.CreateTime, &node.UpdateTime, &node.IsValid, &node.ValidTime, &node.HdPrivateKey)
		if err != nil {
			log.Warning(err)
			return false, err
		}

		_, err = tx.Exec("INSERT INTO revoke(`net`, `serial_num`, `create_time`) VALUES (?,?,?)",
			node.Net,
			node.SerialNum,
			int(time.Now().Unix()))
		if err != nil {
			log.Warning(err)
			return false, err
		}
	}

	// 3.更新node表
	_, err = tx.Exec("UPDATE node SET is_valid=? WHERE net=? AND address=?", false, net, address)
	if err != nil {
		log.Warning("RevokeNodeByNetAndAddress update authorization records failed")
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}
