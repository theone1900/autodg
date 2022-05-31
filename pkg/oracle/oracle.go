package oracle

import (
	"autodg/service"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// GetOracleDBSid 获取oracle_sid
func GetOracleDBSid(oracledb *sql.DB) (string, error) {
	startTime := time.Now()
	service.Logger.Info("get oracle oracle_sid start")

	querySQL := fmt.Sprintf(`select value from v$parameter where NAME='instance_name'`)
	_, res, err := service.Query(oracledb, querySQL)
	if err != nil {
		return res[0]["VALUE"], err
	}

	endTime := time.Now()
	service.Logger.Info("get oracle oracle_sid finished",
		zap.String("CMDS", querySQL),
		zap.String("cost", endTime.Sub(startTime).String()))

	return res[0]["VALUE"], nil
}

// GetOracleDBname 获取db-name
func GetOracleDBname(oracledb *sql.DB) (string, error) {
	querySQL := fmt.Sprintf(`select value from v$parameter where NAME='db_name'`)
	_, res, err := service.Query(oracledb, querySQL)
	if err != nil {
		return res[0]["VALUE"], err
	}
	return res[0]["VALUE"], nil
}

// GetOracleUniquename 获取db_unique_name
func GetOracleUniquename(oracledb *sql.DB) (string, error) {
	querySQL := fmt.Sprintf(`select value from v$parameter where NAME='db_unique_name'`)
	_, res, err := service.Query(oracledb, querySQL)
	if err != nil {
		return res[0]["VALUE"], err
	}
	return res[0]["VALUE"], nil
}

func GetOracleArcMode(oracledb *sql.DB) (string, error) {

	querySQL := fmt.Sprintf(`select LOG_MODE from v$database`)
	_, res, err := service.Query(oracledb, querySQL)
	if err != nil {
		return res[0]["LOG_MODE"], err
	}
	return res[0]["LOG_MODE"], nil
}

func GetOracleForcelog(oracledb *sql.DB) (string, error) {

	querySQL := fmt.Sprintf(`select force_logging from v$database`)
	_, res, err := service.Query(oracledb, querySQL)
	if err != nil {
		return res[0]["FORCE_LOGGING"], err
	}
	return res[0]["FORCE_LOGGING"], nil
}

// GetOracleDbversion 获取oracle 数据库版本
func GetOracleDbversion(oracledb *sql.DB) (string, error) {
	querySQL := fmt.Sprintf(`select VALUE from NLS_DATABASE_PARAMETERS WHERE PARAMETER='NLS_RDBMS_VERSION'`)
	_, res, err := service.Query(oracledb, querySQL)
	if err != nil {
		return res[0]["VALUE"], err
	}
	return res[0]["VALUE"], nil
}

// GetOracleClusterstat get oracle is cluster or not
func GetOracleClusterstat(oracledb *sql.DB) (string, error) {
	querySQL := fmt.Sprintf(`select VALUE from v$parameter  where name='cluster_database'`)
	_, res, err := service.Query(oracledb, querySQL)
	if err != nil {
		return res[0]["VALUE"], err
	}
	return res[0]["VALUE"], nil
}

// CheckOracleSpfile check primary database  using spfile or not?
func CheckOracleSpfile(oracledb *sql.DB) (string, error) {
	querySQL := fmt.Sprintf(` select count(*) VALUE from v$parameter  where name='spfile'`)
	_, res, err := service.Query(oracledb, querySQL)
	if err != nil {
		return res[0]["VALUE"], err
	}
	return res[0]["VALUE"], nil
}
