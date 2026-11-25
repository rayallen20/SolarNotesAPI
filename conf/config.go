package conf

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Service  Service  `yaml:"service"`
	Database Database `yaml:"database"`
	Asset    Asset    `yaml:"asset"`
	Storage  Storage  `yaml:"storage"`
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

type Asset struct {
	BaseURL string `yaml:"base_url"`
}

type Storage struct {
	RootDir string `yaml:"root_dir"`
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
	baseURL, err := url.Parse(c.Asset.BaseURL)
	if err != nil || (baseURL.Scheme != "http" && baseURL.Scheme != "https") ||
		baseURL.Hostname() == "" || baseURL.User != nil || strings.ContainsAny(c.Asset.BaseURL, "?#") {
		return errors.New("asset.base_url 必须为包含主机的 HTTP(S) 地址, 且不能含用户信息, 查询参数或片段")
	}
	if port := baseURL.Port(); port != "" {
		n, err := strconv.ParseUint(port, 10, 16)
		if err != nil || n == 0 {
			return errors.New("asset.base_url 的端口必须为 1 到 65535 之间的整数")
		}
	}
	if !filepath.IsAbs(c.Storage.RootDir) {
		return errors.New("storage.root_dir 必须为绝对路径")
	}
	return nil
}
