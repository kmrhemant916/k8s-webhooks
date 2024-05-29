package helpers

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type RolePermission struct {
	Role       string   `yaml:"role"`
	Permissions []string `yaml:"permissions"`
}

type GroupRole struct {
	Group       string   `yaml:"group"`
	Roles []string `yaml:"roles"`
}

type Config struct {
    Service struct {
        Port string `yaml:"port"`
    } `yaml:"service"`
    Database struct {
		Name string `yaml:"name"`
		Password string `yaml:"password"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Username string `yaml:"username"`
	} `yaml:"database"`
    Rabbitmq struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"rabbitmq"`
	Roles []string `yaml:"roles"`
	Permissions []string `yaml:"permissions"`
	RolePermissions []RolePermission `yaml:"rolePermissions"`
	Groups []string `yaml:"groups"`
	GroupRoles []GroupRole `yaml:"groupRoles"`
	JWTKey string `yaml:"jwt_key"`
}

func (c *Config) ReadConf(f string) (*Config, error) {
    buf, err := ioutil.ReadFile(f)
    if err != nil {
        return nil, err
    }
    err = yaml.Unmarshal(buf, c)
    if err != nil {
        return nil, fmt.Errorf("in file %q: %w", f, err)
    }
    return c, err
}