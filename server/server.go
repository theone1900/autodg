/*
Copyright © 2020 Marvin

*/
package server

import (
	"autodg/pkg/localexec"
	"autodg/pkg/oracle"
	"autodg/pkg/sshexec"
	"autodg/service"
	"database/sql"
	"fmt"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

// Run 程序运行
func Run(cfg *service.CfgFile, mode string) error {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "prepare":
		// dg 初始化 - only prepare 阶段
		oraengine, err := NewOracleDBEngine(cfg.SourceConfig)
		if err != nil {
			return err
		}

		// 读取主库oracle_sid
		sid,err := oracle.GetOracleDBSid(oraengine)
		cfg.SourceConfig.OracleSid=sid
		fmt.Println("[oracle sid] : ",sid)

		// 读取主库db_name
		dbname,err := oracle.GetOracleDBname(oraengine)
		cfg.SourceConfig.OracleDBname=dbname
		fmt.Println("[oracle db-name] : ",dbname)

		// 读取主库db_unique_name
		db_unique_name,err := oracle.GetOracleUniquename(oraengine)
		cfg.SourceConfig.OracleUniqname=db_unique_name
		fmt.Println("[oracle db_unique_name] : ",db_unique_name)

		// 读取主库cluster_database 参数：
		israc,err := oracle.GetOracleClusterstat(oraengine)
		cfg.SourceConfig.IsRAC=israc
		if israc == "FALSE" {
			fmt.Println("[oracle cluster_database] : ",israc," ##This is a Single node DataBase")
		}else {
			fmt.Println("[oracle cluster_database] : ",israc," ##This is a RAC Cluster DataBase")
		}


		// 检查主库归档模式和force log状态
		arcmode,err :=oracle.GetOracleArcMode(oraengine)
		fmt.Println("[arch log mode] : ",arcmode)
		if arcmode != "ARCHIVELOG" {
			fmt.Println("[important]: oracle db is not run in ARCHIVELOG mode")
			os.Exit(1)
		}

		logmode,err :=oracle.GetOracleForcelog(oraengine)
		fmt.Println("[Force log mode] ：",logmode)
		if logmode != "YES"  {
			fmt.Println("[important]: oracle db is not run in  FORCE_LOGGING mode")
			os.Exit(1)
			//todo:
			//configForcelog(),是否自动配置force logging
		}

		//  主库数据库版本检查 <= 11201 的版本不支持
		dbversion,err :=oracle.GetOracleDbversion(oraengine)
		fmt.Println("[dbversion] ：",dbversion)
		//fmt.Println("dbversion[0:2]",dbversion[0:2])
		//fmt.Println(strings.Split(string(dbversion),".")[0])
		var lowdbveriosn int
		lowdbveriosn=11

		dv,err:=strconv.Atoi((strings.Split(string(dbversion),".")[0]))
		if  dv < lowdbveriosn  {
			fmt.Println("[important]: oracle dbversion is  <= 11201")
			os.Exit(1)
		}





		// tnsnames.ora orapwd 初始化
		// 检查主库密码文件是否存在
		isorapwexist,err := sshexec.CheckPriOrapwd(cfg.SourceConfig,sid)
		fmt.Println("[isorapwexist]",isorapwexist)
		a := strings.Index(isorapwexist,"1")
		if a > 1 {
			fmt.Println("[important] ：there is no orapw"+sid )
			//todo
			//创建oracle 密码文件：orapw()
		}

		// 更新主库tnsnames 文件
		sshexec.UpdatePriTns(cfg.SourceConfig)

		// 下载主库tnsnames.ora 文件
		sshexec.DownloadPriTns(cfg.SourceConfig)

		// 下载主库orapw 文件
		sshexec.DownloadPriOrapw(cfg.SourceConfig,sid)

		// 拷贝tnsnames，orapw 到备库
		service.Logger.Info("prepare DG", zap.String("Init standby DB tnsnames.ora file ",  cfg.SourceConfig.PrimaryOracleHome+"/network/admin/tnsnames.ora"))
		localexec.Copytns(cfg.SourceConfig)
		service.Logger.Info("prepare DG", zap.String("Init standby DB orapw file ",  cfg.SourceConfig.StandbyOracleHome+"/dbs/"+"orapw"+sid))
		localexec.CopyOrapw(cfg.SourceConfig,sid)


		// 初始化备库listener（single:oracle  RAC:GRID）
		// single: oracle user oracle_home
	    // RAC   : grid   user oracle_home
		service.Logger.Info("prepare DG", zap.String("Init standby DB listener.ora file ",  cfg.SourceConfig.StandbyGridHome+"/network/admin/tnsnames.ora"))
		localexec.InitStdListener(cfg.SourceConfig,dbname,sid)


		// 启动备库监听 lnsrctl start listener
		// single： oracle
		// RAC   ： grid
		service.Logger.Info("prepare DG", zap.String("Startup standby DB listener ","Listener starting........."))
		localexec.StartListener(cfg.SourceConfig)
		

		// 初始化standby instance pfile
		service.Logger.Info("prepare DG", zap.String("Init standby DB pfile ","pfile initing........."))
		localexec.InitStdInstancePfile(cfg.SourceConfig)



		// 启动备库实例到nomount
		service.Logger.Info("prepare DG", zap.String("Startup standby Instance ",  cfg.SourceConfig.StandbyOracleHome+"/dbs/pfile"+cfg.SourceConfig.OracleSid))
		localexec.StartStdInstance(cfg.SourceConfig)


		// tnsping (主库，备库)连接校验
		//todo:tnsping primary && tnsping standby
		//$ sqlplus sys/oracle@orabak as sysdba
		//$ sqlplus sys/oracle@orcl as sysdba



		// DG duplicate from active primary database
		service.Logger.Info("prepare DG", zap.String("Startup RMAN Duplicate Database From Active database ",cfg.SourceConfig.OracleSid))
		localexec.StartRmanDuplicate(cfg.SourceConfig)








		if err != nil {
			return err
		}
		//if err := prepare.TransferDBEnvPrepare(engine); err != nil {
		//	return err
		//}


	case "check":
		// 主库连通性检查，归档状态检查，forceinglog 检查，spfile 检查；
		engine, err :=  NewOracleDBEngine(cfg.SourceConfig)
		engine.Ping()
		if err != nil {
			return err
		}
		//log.Info(`[oracle DB PING] :`, zap.String("lv", "info"), zap.String("Stat", ))


		// 读取主库oracle version
		dbv,err := oracle.GetOracleDbversion(engine)
		if err != nil {
			return err
		}
		//fmt.Println("[oracle Dbversion] : ",dbv)
		//log.Info("[oracle Dbversion] : ",dbv)

		var lowdbveriosn int
		lowdbveriosn=11

		dv,err:=strconv.Atoi((strings.Split(string(dbv),".")[0]))
		if  dv < lowdbveriosn  {
			log.Warn(`oracle DBVersion`, zap.String("Dbversion", dbv),zap.String("extra","dbversion <= 11201  is not support!"))
			os.Exit(1)
		} else {
			log.Info(`oracle DBVersion`,zap.String("Dbversion", dbv))
		}




		// 读取主库oracle_sid
		sid,err := oracle.GetOracleDBSid(engine)
		if err != nil {
			return err
		}
		//fmt.Println("[oracle sid] : ",sid)
		//log.Info("[oracle sid] : ",sid)
		//log.Debug(`This is a debug message.`, zap.String("lv", "debug"), zap.Int("no", 1))
		log.Info(`oracle SID`,  zap.String("SID", sid))
		//log.Warn(`This is a warning message.`, zap.String("lv", "warning"), zap.String("extra", "some extra msg"))

		// 读取主库归档模式
		archemode ,err := oracle.GetOracleArcMode(engine)
		if err != nil {
			return err
		}
		//fmt.Println("[oracle archived_mode] : ",archemode)
		//log.Info("[oracle archived_mode] : ",archemode)
		if archemode != "ARCHIVELOG" {
			log.Warn(`oracle Archived_Mode`,  zap.String("Archived_Mode", archemode),zap.String("extra","oracle db is NO_ARCHIVELOG mode"))
		}else {
			log.Info(`oracle Archived_Mode`,  zap.String("Archived_Mode", archemode))
		}




		//读取主库forceing logging 模式
		forcelog ,err := oracle.GetOracleForcelog(engine)
		if err != nil {
			return err
		}
		//fmt.Println("[oracle force_logging] : ",forcelog)
		//log.Info("[oracle force_logging] : ",forcelog)
		log.Info(`oracle Force_Logging`,  zap.String("Force_Logging", forcelog))



		// 读取主库是否使用spfile
		spfile,err := oracle.CheckOracleSpfile(engine)
		if err != nil {
			return err
		}
		//fmt.Println("[oracle spfile status] : ",spfile)
		//log.Info("[oracle spfile status] : ",spfile)
		log.Info(`oracle SPFILE status`,  zap.String("Spfile", spfile))


		// 主库密码文件检查
		// todo: chakorapwd()   主库密码文件检查

		default:
		return fmt.Errorf("flag [mode] can not null or value configure error")
	}
	return nil
}

// NewEngineDB 数据库引擎初始化
func NewEngineDB(sourceCfg service.SourceConfig) (*service.Engine, error) {
	var (
		engine1 *service.Engine
		oraDB  *sql.DB
		err    error
	)
	oraDB, err = NewOracleDBEngine(sourceCfg)

	if err != nil {
		fmt.Println("[oradb engine ]","engine oradb err")
		panic(err)
	}
	//engine, err = NewMySQLEngineGeneralDB()
	//if err != nil {
	//	return engine, err
	//}
    fmt.Printf("[oradb11]",oraDB )
	fmt.Println("#####################")
	fmt.Printf("[engine.OracleDB]",&engine1.OracleDB)
	engine1.OracleDB = oraDB
	return engine1, nil

}


// NewEngineSSH SSH引擎初始化
//func NewEngineSSH(sourceCfg service.SourceConfig) (*ssh.Client, error) {
//	var (
//		engine *service.Engine
//		sshserver *ssh.Client
//		err    error
//	)
//	sshserver, err = NewSSHConnectEngine(sourceCfg)
//	if err != nil {
//		return engine, err
//	}
//	engine.SshExc = sshserver
//
//	sftpserver = NewSftpConnectEngine(sourceCfg)
//	if err != nil {
//		return engine, err
//	}
//	engine.SftpExc = sftpserver
//	return engine, nil
//}