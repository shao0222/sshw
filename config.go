package sshw

import (
	"io/ioutil"
	"os/user"
	"path"

	"time"

	"fmt"

	"crypto/md5"

	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	MOtp []*MOtp `json:"motp"`
	Host []*Node `json:"host"`
}
type Node struct {
	Name       string  `json:"name"`
	Host       string  `json:"host"`
	User       string  `json:"user"`
	Port       int     `json:"port"`
	KeyPath    string  `json:"keypath"`
	Passphrase string  `json:"passphrase"`
	Password   string  `json:"password"`
	Children   []*Node `json:"children"`
	MOtp       string  `json:"motp"`
}

func (n *Node) String() string {
	return n.Name
}

func (n *Node) user() string {
	if n.User == "" {
		return "root"
	}
	return n.User
}

func (n *Node) port() int {
	if n.Port <= 0 {
		return 22
	}
	return n.Port
}

func (n *Node) password() ssh.AuthMethod {
	if n.Password == "" {
		if n.MOtp != "" {
			for _, m := range config.MOtp {
				if m.Name == n.MOtp {
					return ssh.Password(m.String())
				}
			}
		}
		return nil
	}
	return ssh.Password(n.Password)
}

type MOtp struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
	Pin    string `json:"pin"`
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
}

func (m *MOtp) String() string {
	str := fmt.Sprintf("%d%s%s", time.Now().Unix()/10, m.Secret, m.Pin)
	return m.Prefix + fmt.Sprintf("%x", md5.Sum([]byte(str)))[:6] + m.Suffix
}

var (
	config Config
)

func GetNodeConfig() []*Node {
	return config.Host
}

func getMotpConfig() []*MOtp {
	return config.MOtp
}

func LoadConfig() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(path.Join(u.HomeDir, ".sshw"))
	if err != nil {
		return err
	}

	var c Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	config = c

	return nil
}
