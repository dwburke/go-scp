package scp

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	PrivateKey string
	Host       string
	Port       int
	Username   string
	signer     *ssh.Signer
	connection *ssh.Client
	sftp       *sftp.Client
}

func New(private_key string, host string, port int, username string) (*Config, error) {
	if port <= 0 {
		port = 22
	}

	config := Config{PrivateKey: private_key, Host: host, Port: port, Username: username}

	return &config, nil
}

func LoadPrivateKey(keyfile string) (string, error) {
	if keyfile == "" {
		keyfile = os.Getenv("HOME") + "/.ssh/id_rsa"
	}

	b, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return "", err
	}

	str := string(b)

	return str, nil
}

func (config *Config) Close() {
	if config.sftp != nil {
		config.sftp.Close()
	}

	if config.connection != nil {
		config.connection.Close()
	}
}

func (config *Config) Connect() error {

	if config.signer != nil && config.connection != nil {
		// TODO test connection?
		return nil
	}

	if config.signer == nil {
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return err
		}

		config.signer = &signer
	}

	if config.connection == nil {
		ssh_config := &ssh.ClientConfig{
			User:            config.Username,
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(*config.signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		conn, err := ssh.Dial("tcp", config.Host+":"+strconv.Itoa(config.Port), ssh_config)
		if err != nil {
			return err
		}

		config.connection = conn
	}

	return nil
}

func (config *Config) ScpSession() error {
	if config.sftp != nil {
		return nil
	}

	if config.connection == nil {
		config.Connect()
	}

	sftp_client, err := sftp.NewClient(config.connection)
	if err != nil {
		return err
	}

	config.sftp = sftp_client

	return nil
}

func (config *Config) Get(source string, destination string) error {
	err := config.Connect()
	if err != nil {
		return err
	}

	err = config.ScpSession()
	if err != nil {
		return err
	}

	srcFile, err := config.sftp.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile.WriteTo(dstFile)

	return nil
}
