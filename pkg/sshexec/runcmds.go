package sshexec

import (
	"autodg/service"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"time"

	//"golang.org/x/crypto/ssh"
	"log"
	"os"
)

/*
	sshserver
*/
// 更新主库tnsnames
// single: oraclehome
// RAC   : gridhome
func UpdatePriTns(sshconfig service.SourceConfig) (string, error) {
	var res string
	var err error
	if sshconfig.IsRAC == "FALSE" {
		//service.Logger.Info("prepare DG", zap.String("Get Primary DB (Sinage & RAC)","##This is a Single node DataBase"))
		tnsnames := fmt.Sprintf(`
pri1900 =
		(DESCRIPTION =
			(ADDRESS_LIST =
			(ADDRESS = (PROTOCOL = TCP)(HOST = %s)(PORT = 1521))
		)
		(CONNECT_DATA =
			(SERVICE_NAME = %s)
))

std1900 =
		(DESCRIPTION =
			(ADDRESS_LIST =
			(ADDRESS = (PROTOCOL = TCP)(HOST = %s)(PORT = 1521))
		)
		(CONNECT_DATA =
			(SERVICE_NAME = %s)
))`, sshconfig.Host, sshconfig.ServiceName, sshconfig.StandbyHostIps[0], sshconfig.ServiceName)
		cmds := fmt.Sprintf(`echo "%s" >> %s/network/admin/tnsnames.ora`, tnsnames, sshconfig.PrimaryOracleHome)
		fmt.Printf("[UpdatePriTns cmds:]", cmds)

		res, err := service.RunSsh(cmds, sshconfig)
		if err != nil {
			return res, err
		}
		return res, err

	} else {
		//fmt.Println("[oracle cluster_database] : ",israc," ##This is a RAC Cluster DataBase")
		//service.Logger.Info("prepare DG", zap.String("Get Primary DB (Sinage & RAC)","##This is a RAC Cluster DataBase"))
		tnsnames := fmt.Sprintf(`
pri1900 =
		(DESCRIPTION =
			(ADDRESS_LIST =
			(ADDRESS = (PROTOCOL = TCP)(HOST = %s)(PORT = 1521))
		)
		(CONNECT_DATA =
			(SERVICE_NAME = %s)
))

std1900 =
		(DESCRIPTION =
			(ADDRESS_LIST =
			(ADDRESS = (PROTOCOL = TCP)(HOST = %s)(PORT = 1521))
		)
		(CONNECT_DATA =
			(SERVICE_NAME = %s)
))`, sshconfig.Host, sshconfig.ServiceName, sshconfig.StandbyHostIps[0], sshconfig.ServiceName)
		var cmds = fmt.Sprintf(`echo "%s" >> %s/network/admin/tnsnames.ora`, tnsnames, sshconfig.PrimaryGridHome)
		fmt.Println("[UpdatePriTns cmds:]", cmds)

		res, err := service.RunSsh(cmds, sshconfig)
		if err != nil {
			return res, err
		}
		return res, err
	}
	return res, err
}

// DownloadPriTns 下载主库tnsnames
func DownloadPriTns(sftpconfig service.SourceConfig) (string, error) {

	//var res string
	var sftpClient *sftp.Client

	sshHost := sftpconfig.Host
	sshUser := "root"
	sshPassword := sftpconfig.RootPwd
	sshType := "password"
	sshPort := 22
	//创建ssh登陆配置
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
	fmt.Println("[addr:]", addr)
	sshClient, err := ssh.Dial("tcp", addr, config)

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		log.Fatalln("error occurred:", err)
	}
	srcFile, err := sftpClient.Open(fmt.Sprintf(`%s/network/admin/tnsnames.ora`, sftpconfig.PrimaryOracleHome)) //远程

	dstPath := "tnsnames.ora"
	dstFile, err := os.Create(dstPath) //本地

	//defer func() {
	//	_ = srcFile.Close()
	//	_ = dstFile.Close()
	//}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		log.Fatalln("error occurred", err)
	}
	fmt.Println("tnsnames文件下载完毕")

	return "tnsnames 文件下载完毕", err
}

// CheckPriOrapwd 检查主库密码文件orapw 是否存在
func CheckPriOrapwd(sshconfig service.SourceConfig, oracle_sid string) (string, error) {

	var cmds = fmt.Sprintf(`ls -l %s/dbs/orapw%s |wc -l 2>/dev/null`, sshconfig.PrimaryOracleHome, oracle_sid)
	fmt.Println("[check orapwd start]")
	res, err := service.RunSsh(cmds, sshconfig)
	fmt.Println("[check orapwd end]")
	if err != nil {
		return res, err
	}
	return res, nil
}

// DownloadPriOrapw 下载主库orapw 密码文件
func DownloadPriOrapw(sftpconfig service.SourceConfig, oracle_sid string) (string, error) {

	//var res string
	var sftpClient *sftp.Client

	sshHost := sftpconfig.Host
	sshUser := "root"
	sshPassword := sftpconfig.RootPwd
	sshType := "password"
	sshPort := 22
	//创建ssh登陆配置
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
	fmt.Println("[addr:]", addr)
	sshClient, err := ssh.Dial("tcp", addr, config)

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		log.Fatalln("error occurred:", err)
	}
	srcFile, err := sftpClient.Open(fmt.Sprintf(`%s/dbs/orapw%s`, sftpconfig.PrimaryOracleHome, oracle_sid)) //远程
	fmt.Println(fmt.Sprintf(`%s/dbs/orapw%s`, sftpconfig.PrimaryOracleHome, oracle_sid))
	dstPath := fmt.Sprintf(`orapw%s`, oracle_sid)
	dstFile, err := os.Create(dstPath) //本地

	//defer func() {
	//	_ = srcFile.Close()
	//	_ = dstFile.Close()
	//}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		log.Fatalln("error occurred", err)
	}
	fmt.Println("orapw文件下载完毕")

	return "orapw文件下载完毕", err
}
