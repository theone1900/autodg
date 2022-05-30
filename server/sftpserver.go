package server

import (
	"autodg/service"
	"fmt"
	"golang.org/x/crypto/ssh"
    "github.com/pkg/sftp"
	)
//获取ftp连接

func NewSftpConnectEngine( sshCfg service.SourceConfig) (*sftp.Client) {

	var (
		client       *ssh.Client
		err          error
	)

	client,err =  NewSSHConnectEngine(sshCfg)

	// 此时获取了sshClient，下面使用sshClient构建sftpClient
	ftpclient, err := sftp.NewClient(client)
	if err != nil {
		fmt.Println("创建ftp客户端失败", err)
		panic(err)
	}
	return ftpclient
	}


