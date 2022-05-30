package prepare

import (
	"autodg/pkg/sshexec"
	"autodg/service"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

// dg同步环境准备
func SyncDg(engine *service.Engine,cfg *service.CfgFile) error {
	startTime := time.Now()
	service.Logger.Info("prepare autodg env start")

	res, err := engine.UpdatePriTns(cfg.SourceConfig.Host,cfg.SourceConfig.StandbyHostIps[0],cfg.SourceConfig.ServiceName,cfg.SourceConfig.PrimaryOracleHome)

	if err !=nil {
		return err
	}


	//if err := engine.InitDefaultValueMap(); err != nil {
	//	return err
	//}
	endTime := time.Now()
	service.Logger.Info("prepare tansferdb env finished",
		zap.String("cost", endTime.Sub(startTime).String()))
	return nil
}


///*
//	sshserver
//*/
//// 更新主库tnsnames
//func  UpdatePriTns(priip string,stdip string,servicename string,priOraHome string,) (string, error) {
//
//	tnsnames := fmt.Sprintf(`
//	pri1900 =
//		(DESCRIPTION =
//			(ADDRESS_LIST =
//			(ADDRESS = (PROTOCOL = TCP)(HOST = '%s')(PORT = 1521))
//		)
//		(CONNECT_DATA =
//			(SERVICE_NAME = '%s')
//		))
//
//	std1900 =
//		(DESCRIPTION =
//			(ADDRESS_LIST =
//			(ADDRESS = (PROTOCOL = TCP)(HOST = '%s')(PORT = 1521))
//		)
//		(CONNECT_DATA =
//			(SERVICE_NAME = '%s')
//		))`,priip,servicename,servicename,stdip,servicename)
//	cmds := fmt.Sprintf(`echo '%s' >> '%s'/tnsnames.ora`,tnsnames,priOraHome)
//
//	res, err := sshexec.runcmds(e.SshExc, cmds)
//	if err != nil {
//		return res, err
//	}
//	return res , nil
//}
//// 下载主库tnsnames
//func  DownloadPriTns(priOraHome string) (string, error) {
//
//	srcFile ,err := e.SftpExc.Open(fmt.Sprintf(`'%s'/network/admin/tnsnames.ora`,priOraHome)) //远程
//	dstPath := "./"
//	dstFile, err := os.Create(dstPath) //本地
//
//	defer func() {
//		_ = srcFile.Close()
//		_ = dstFile.Close()
//	}()
//
//	if _, err := srcFile.WriteTo(dstFile); err != nil {
//		log.Fatalln("error occurred", err)
//	}
//	fmt.Println("文件下载完毕")
//
//	return "文件下载完毕",err
//}


//
//// 检查主库密码文件orapw 是否存在
//func  CheckPriOrapwd(priip string,priOraHome string,oracle_sid string) (string, error) {
//
//	cmds := fmt.Sprintf(`'%s'/dbs/orapwd'%s'`,priOraHome,oracle_sid)
//
//	res, err := sshexec.runcmds(e.SshExc, cmds)
//	if err != nil {
//		return res, err
//	}
//	return res , nil
//}




