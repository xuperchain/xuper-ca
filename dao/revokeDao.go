/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package dao

import (
	log "github.com/sirupsen/logrus"
)

type Revoke struct {
	Id         int    `db:"id"`
	Net        string `db:"net"`
	SerialNum  string `db:"serial_num"`
	CreateTime int    `db:"create_time"`
}

// ca网络下的撤销列表数据层
type RevokeDao struct {
}

// cursor
func (revokeDao *RevokeDao) GetList(net, latestSerialNum string) (*[]Revoke, error) {
	ret := []Revoke{}
	caDb := GetDbInstance()

	var id int
	if latestSerialNum == "" {
		id = 0
	} else {
		var cursor Revoke
		err := caDb.db.Get(&cursor, "SELECT * FROM revoke WHERE net=? AND serial_num=?", net, latestSerialNum)
		if err != nil {
			log.Warning("GetRevokeList failed, err:", err)
		}
		id = cursor.Id
	}

	rows, err := caDb.db.Queryx("SELECT * FROM revoke where net=? AND id>?", net, id)
	if err != nil {
		return nil, err
	}
	var row Revoke
	rows.StructScan(row)
	// iterate over each row
	for rows.Next() {
		row := Revoke{}
		err = rows.Scan(&row.Id, &row.Net, &row.SerialNum, &row.CreateTime)
		if err != nil {
			return nil, err
		}
		ret = append(ret, row)
	}
	return &ret, nil
}
