package conf

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Load 读取当前环境的配置，并在返回前校验必填字段和端口。
func Load() (*Config, error) {
	filename := "config.dev.yaml"
	switch os.Getenv("SOLAR_NOTES_ENV") {
	case "", "dev":
	case "prod":
		filename = "config.prod.yaml"
	default:
		return nil, errors.New("SOLAR_NOTES_ENV 仅支持 dev 或 prod")
	}

	executable, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("定位可执行文件失败: %w", err)
	}
	path := filepath.Join(filepath.Dir(executable), filename)
	_, err = os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		path = filename
	} else if err != nil {
		return nil, fmt.Errorf("定位配置文件失败: %w", err)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		// YAML 解析错误可能包含原始配置值，包括误填的凭据。
		return nil, fmt.Errorf("%s 无法解析，请检查 YAML 语法和字段名", filename)
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("%s 必须仅包含一个完整的 YAML 文档", filename)
	}
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("%s 配置无效: %w", filename, err)
	}
	config.Asset.BaseURL = strings.TrimRight(config.Asset.BaseURL, "/")
	return &config, nil
}
