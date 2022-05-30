
package service

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
	"time"
)

func RunSsh(cmds string,sshconfig SourceConfig) (string, error) {
	//var res string
	sshHost := sshconfig.Host
	sshUser := "root"
	sshPassword := sshconfig.RootPwd
	sshType := "password"
	sshPort := 22
	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         5 * time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	fmt.Println("[addr:]",addr)
	sshClient, err := ssh.Dial("tcp", addr, config)

	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()

	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}

	defer session.Close()
	//执行远程命令
	fmt.Println("[print current cmds]:",cmds)
	combo, err := session.CombinedOutput(cmds)
	if err != nil {
		log.Fatal("远程执行cmd 失败", err)
	}
	log.Println("命令输出:", string(combo))

	return string(combo),err

}



/*
	Oracle
*/
func (e *Engine) GetOracleDBVersion() (string, error) {
	querySQL := fmt.Sprintf(`select VALUE from NLS_DATABASE_PARAMETERS WHERE PARAMETER='NLS_RDBMS_VERSION'`)
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res[0]["VALUE"], err
	}
	return res[0]["VALUE"], nil
}


func (e *Engine) GetOracleDBSid() (string, error) {
	querySQL := fmt.Sprintf(`select value from v$parameter where NAME='instance_name'`)
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res[0]["VALUE"], err
	}
	return res[0]["VALUE"], nil
}

func (e *Engine) GetOracleDBCharacterSet() (string, error) {
	querySQL := fmt.Sprintf(`select userenv('language') AS LANG from dual`)
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res[0]["LANG"], err
	}
	return res[0]["LANG"], nil
}

func (e *Engine) GetOracleDBCharacterNLSCompCollation() (string, error) {
	querySQL := fmt.Sprintf(`select VALUE from NLS_DATABASE_PARAMETERS WHERE PARAMETER = 'NLS_COMP'`)
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return "", err
	}
	return res[0]["VALUE"], nil
}

func (e *Engine) GetOracleDBCharacterNLSSortCollation() (string, error) {
	querySQL := fmt.Sprintf(`select VALUE from NLS_DATABASE_PARAMETERS WHERE PARAMETER = 'NLS_SORT'`)
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return "", err
	}
	return res[0]["VALUE"], nil
}

func (e *Engine) GetOracleSchemaCollation(schemaName string) (string, error) {
	querySQL := fmt.Sprintf(`SELECT DECODE(DEFAULT_COLLATION,
'USING_NLS_COMP',(SELECT VALUE from NLS_DATABASE_PARAMETERS WHERE PARAMETER = 'NLS_COMP'),DEFAULT_COLLATION) DEFAULT_COLLATION FROM DBA_USERS WHERE USERNAME = '%s'`, strings.ToUpper(schemaName))
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return "", err
	}
	return res[0]["DEFAULT_COLLATION"], nil
}



func (e *Engine) GetOraclePartitionTableINFO(schemaName, tableName string) ([]map[string]string, error) {
	querySQL := fmt.Sprintf(`SELECT L.PARTITIONING_TYPE,
       L.SUBPARTITIONING_TYPE,
       L.PARTITION_EXPRESS,
       LISTAGG(skc.COLUMN_NAME, ',') WITHIN GROUP (ORDER BY skc.COLUMN_POSITION) AS SUBPARTITION_EXPRESS
FROM (SELECT pt.OWNER,
             pt.TABLE_NAME,
             pt.PARTITIONING_TYPE,
             pt.SUBPARTITIONING_TYPE,
             LISTAGG(ptc.COLUMN_NAME, ',') WITHIN GROUP (ORDER BY ptc.COLUMN_POSITION) AS PARTITION_EXPRESS
      FROM DBA_PART_TABLES pt,
           DBA_PART_KEY_COLUMNS ptc
      WHERE pt.OWNER = ptc.OWNER
        AND pt.TABLE_NAME = ptc.NAME
        AND ptc.OBJECT_TYPE = 'TABLE'
        AND UPPER(pt.OWNER) = UPPER('%s')
        AND UPPER(pt.TABLE_NAME) = UPPER('%s')
      GROUP BY pt.OWNER, pt.TABLE_NAME, pt.PARTITIONING_TYPE,
               pt.SUBPARTITIONING_TYPE) L,
     DBA_SUBPART_KEY_COLUMNS skc
WHERE L.OWNER = skc.OWNER
  AND L.TABLE_NAME = skc.NAME
GROUP BY  L.PARTITIONING_TYPE,
       L.SUBPARTITIONING_TYPE,
       L.PARTITION_EXPRESS`, schemaName, tableName)
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res, err
	}
	return res, nil
}

