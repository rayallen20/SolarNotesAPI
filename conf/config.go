package conf

import (
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Service  Service  `yaml:"service"`
	Database Database `yaml:"database"`
}

type Service struct {
	Port string `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

func (c Config) validate() error {
	for _, port := range []struct{ field, value string }{
		{"service.port", c.Service.Port},
		{"database.port", c.Database.Port},
	} {
		n, err := strconv.ParseUint(port.value, 10, 16)
		if err != nil || n == 0 {
			return fmt.Errorf("%s 必须为 1 到 65535 之间的整数", port.field)
		}
	}

	for _, item := range []struct{ field, value string }{
		{"database.host", c.Database.Host},
		{"database.user", c.Database.User},
		{"database.password", c.Database.Password},
		{"database.name", c.Database.Name},
	} {
		if strings.TrimSpace(item.value) == "" {
			return fmt.Errorf("%s 不能为空", item.field)
		}
	}
	return nil
}
