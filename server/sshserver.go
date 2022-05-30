package server

import (
	"autodg/service"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"time"
)

func NewSSHConnectEngine( sshCfg service.SourceConfig) ( *ssh.Client, error ) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		//session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(sshCfg.RootPwd))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:               "root",
		Auth:               auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback:    hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf( "%s:%d", sshCfg.Host, 22 )

	if client, err = ssh.Dial( "tcp", addr, clientConfig ); err != nil {
		return nil, err
	}

	//// create session
	//if session, err = client.NewSession(); err != nil {
	//	return nil, err
	//}

	return client, nil
}

func RunSsh(cmds string,sshconfig service.SourceConfig) {
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
	 err = session.Run(cmds)
	if err != nil {
		log.Fatal("远程执行cmd 失败", err)
	}
	log.Println("命令输出:", string("ok"))

}











