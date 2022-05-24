package ssh

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
)

// SSHConfig 本地 ssh 配置
type SSHConfig struct {
	Port           int
	Address        string
	Username       string
	PrivateKeyPath string
}

type ViaSSHDialer struct {
	client *ssh.Client
}

func (v *ViaSSHDialer) DialContext(ctx context.Context, addr string) (net.Conn, error) {
	return v.client.Dial("tcp", addr)
}

// PublicKeyFile 读取私钥转公钥
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// AgentSSH 也可以通过 ssh-add /path/to/your/private/certificate/file 方式
func AgentSSH(file string) ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return PublicKeyFile(file)
}

// GetSSHConnection 获得 ssh 连接
func GetSSHConnection(name string, sshConfig SSHConfig) *ssh.Client {
	sshClientConfig := &ssh.ClientConfig{
		User: sshConfig.Username,
		Auth: []ssh.AuthMethod{AgentSSH(sshConfig.PrivateKeyPath)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshCon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshConfig.Address, sshConfig.Port), sshClientConfig)
	if err != nil {
		log.Error(err)
		return nil
	}
	mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{sshCon}).DialContext)
	return sshCon
}
