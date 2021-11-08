/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package dao

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"github.com/xuperchain/xuper-ca/config"
)

const defaultSQLiteSchema = `
create table node (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    net varchar(100) NOT NULL,
	adminAddress varchar(50) NOT NULL,
    address varchar(50) NOT NULL,
    serial_num varchar(100) NOT NULL,
    cert TEXT NOT NULL,
    private_key TEXT NOT NULL,
    create_time int(10) NOT NULL,
    update_time int(10) NOT NULL,
    is_valid BOOLEAN DEFAULT TRUE,
    valid_time int(10) NOT NULL,
    hd_private_key TEXT NOT NULL DEFAULT ''
);
CREATE INDEX idx_node_net_addr ON node(net, address);
CREATE UNIQUE INDEX uidx_node_serial ON node(serial_num);


create table net_admin (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    net varchar(100) NOT NULL,
    address varchar(50) NOT NULL,
    serial_num varchar(100) NOT NULL,
    cert TEXT NOT NULL,
    private_key TEXT NOT NULL,
    create_time int(10) NOT NULL,
    update_time int(10) NOT NULL,
    is_valid BOOLEAN DEFAULT TRUE,
    valid_time int(10) NOT NULL,
    hd_private_key TEXT NOT NULL DEFAULT '',
	is_gm BOOLEAN DEFAULT FALSE
);
CREATE UNIQUE INDEX uidx_net_serial ON net_admin(serial_num);
CREATE INDEX idx_net ON net_admin(net);

create table revoke (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	net varchar(100) NOT NULL,
    serial_num varchar(100) NOT NULL,
    create_time int(10) NOT NULL
);
CREATE UNIQUE INDEX uidx_revoke_serial ON revoke(serial_num);
`

type CaDb struct {
	db *sqlx.DB
}

var caDb *CaDb

func GetDbInstance() *CaDb {
	if caDb == nil {
		caDb = NewCaDb()
	}
	return caDb
}

func NewCaDb() (connection *CaDb) {
	var err error
	var db *sqlx.DB
	db, err = sqlx.Connect(config.GetDBConfig().DbType, config.GetDBConfig().DbPath)

	if err != nil {
		log.Print(err)
	}

	return &CaDb{
		db: db,
	}
}

//DB目前仅支持sqlite3, ca启动时会校验是否存在sqlite3 db文件, 如果db不存在则自动创建
func InitTables() {
	if config.GetDBConfig().DbType != "sqlite3" {
		log.Println("not using sqlite, please check tables by self")
		return
	}

	_, err := os.Stat(config.GetDBConfig().DbPath) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return
		}
		// db 文件不存在
		file, err := os.Create(config.GetDBConfig().DbPath)
		if err != nil {
			panic("create db failed")
		}
		file.Close()

		dbConn := NewCaDb()

		// execute a query on the server
		_, err = dbConn.db.Exec(defaultSQLiteSchema)
		if err != nil {
			log.Println(err)
			panic("create table failed")
		}
		log.Println("create tables success")

	}
	return
}
