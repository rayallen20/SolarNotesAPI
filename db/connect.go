package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"SolarNotesAPI/conf"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB 在服务启动并成功连接后供 DAO 使用。
var DB *gorm.DB

// Connect 在 HTTP 服务启动前打开连接池并检查 MySQL 连接。
func Connect(cfg conf.Database) error {
	connection := driver.NewConfig()
	connection.Net = "tcp"
	connection.Addr = net.JoinHostPort(cfg.Host, cfg.Port)
	connection.User = cfg.User
	connection.Passwd = cfg.Password
	connection.DBName = cfg.Name
	connection.Params = map[string]string{"charset": "utf8mb4"}
	connection.ParseTime = true
	connection.Loc = time.Local
	connection.Timeout = 5 * time.Second
	connection.ReadTimeout = 10 * time.Second
	connection.WriteTimeout = 10 * time.Second

	connector, err := driver.NewConnector(connection)
	if err != nil {
		return fmt.Errorf("初始化 MySQL 驱动失败: %w", err)
	}
	pool := sql.OpenDB(connector)
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(time.Hour)

	handle, err := gorm.Open(mysql.New(mysql.Config{Conn: pool}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		// 初始化错误由启动入口统一记录，避免驱动原文包含认证信息。
		Logger: logger.Discard,
	})
	if err != nil {
		_ = pool.Close()
		var serverErr *driver.MySQLError
		if errors.As(err, &serverErr) {
			return fmt.Errorf("连接数据库失败 (MySQL 错误码 %d)", serverErr.Number)
		}
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	handle.Logger = logger.New(log.New(os.Stdout, "[gorm] ", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
	})
	DB = handle
	return nil
}
