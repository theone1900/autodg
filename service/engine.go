
package service

import (
	"database/sql"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	//"gorm.io/gorm"
)

//var (
//	// Oracle/Mysql 对于 'NULL' 统一字符 NULL 处理，查询出来转成 NULL,所以需要判断处理
//	// 查询字段值 NULL
//	// 如果字段值 = NULLABLE 则表示值是 NULL
//	// 如果字段值 = "" 则表示值是空字符串
//	// 如果字段值 = 'NULL' 则表示值是 NULL 字符串
//	// 如果字段值 = 'null' 则表示值是 null 字符串
//	IsNull = "NULLABLE"
//)

// 定义数据库引擎
type Engine struct {
	OracleDB *sql.DB
	//MysqlDB  *sql.DB
	//GormDB   *gorm.DB
	SshExc   *ssh.Client
	SftpExc  *sftp.Client
}




// 查询返回表字段列和对应的字段行数据
func Query(db *sql.DB, querySQL string) ([]string, []map[string]string, error) {
	var (
		cols []string
		res  []map[string]string
	)
	rows, err := db.Query(querySQL)
	if err != nil {
		return cols, res, fmt.Errorf("general sql [%v] query failed: [%v]", querySQL, err.Error())
	}
	defer rows.Close()

	//不确定字段通用查询，自动获取字段名称
	cols, err = rows.Columns()
	if err != nil {
		return cols, res, fmt.Errorf("general sql [%v] query rows.Columns failed: [%v]", querySQL, err.Error())
	}

	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return cols, res, fmt.Errorf("general sql [%v] query rows.Scan failed: [%v]", querySQL, err.Error())
		}

		row := make(map[string]string)
		for k, v := range values {
			key := cols[k]
			// 数据库类型 MySQL NULL 是 NULL，空字符串是空字符串
			// 数据库类型 Oracle NULL、空字符串归于一类 NULL
			// Oracle/Mysql 对于 'NULL' 统一字符 NULL 处理，查询出来转成 NULL,所以需要判断处理
			if v == nil { // 处理 NULL 情况，当数据库类型 MySQL 等于 nil
				//row[key] = IsNull
			} else {
				// 处理空字符串以及其他值情况
				// 数据统一 string 格式显示
				row[key] = string(v)
			}
		}
		res = append(res, row)
	}

	if err = rows.Err(); err != nil {
		return cols, res, fmt.Errorf("general sql [%v] query rows.Next failed: [%v]", querySQL, err.Error())

	}
	return cols, res, nil
}


