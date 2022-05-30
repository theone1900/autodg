/*
Copyright © 2099 huanglj.
*/
package server

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/wentaojin/transferdb/utils"

	"autodg/service"

	"github.com/godror/godror"

	_ "github.com/godror/godror"
)

// 创建 oracle 数据库引擎
func NewOracleDBEngine(oraCfg service.SourceConfig) (*sql.DB, error) {
	// https://pkg.go.dev/github.com/godror/godror
	// https://github.com/godror/godror/blob/db9cd12d89cdc1c60758aa3f36ece36cf5a61814/doc/connection.md

	connString := fmt.Sprintf("oracle://%s:%s@%s/%s?%s",
		oraCfg.Username, oraCfg.Password, utils.StringsBuilder(oraCfg.Host, ":", strconv.Itoa(oraCfg.Port)),
		oraCfg.ServiceName, oraCfg.ConnectParams)
    fmt.Println(connString)

	oraDSN, err := godror.ParseDSN(connString)
	if err != nil {
		return nil, err
	}

	oraDSN.OnInitStmts = oraCfg.SessionParams

	sqlDB := sql.OpenDB(godror.NewConnector(oraDSN))
	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetMaxOpenConns(0)
	sqlDB.SetConnMaxLifetime(0)

	err = sqlDB.Ping()
	if err != nil {
		return sqlDB, fmt.Errorf("error on ping oracle database connection:%v", err)
	}
	fmt.Println("oracle db connect succes")
	return sqlDB, nil
}
