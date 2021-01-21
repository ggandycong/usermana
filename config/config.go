package config

import "time"

const (
	// RedisAddr redis地址.
	RedisAddr string = "localhost:6379"
	// RedisPoolSize 连接redis最多的连接(Maximum number of socket connections.).
	RedisPoolSize int = 30
	// TokenMaxExTime token生存时间.
	TokenMaxExTime int = 3600

	// MysqlDB 连接数据库地址.
	MysqlDB string = "root:11111111@(127.0.0.1:3306)/test_db?charset=utf8"
	// ConnMaxLifetime 数据库一个连接的最大生命周期.
	ConnMaxLifetime time.Duration = 2 * time.Second
	// MaxIdleConns 连接池中最大空闲连接数.
	MaxIdleConns int = 1000
	// MaxOpenConns 同时连接数据库中最多连接数.
	MaxOpenConns int = 2000

	// TCPServerLogPath TCP服务日志.
	TCPServerLogPath string = "./log/tcp_server.log"
	// TCPServerAddr tcp server ip:port.
	TCPServerAddr string = ":3194"
	// TCPClientPoolSize 客户端tcp连接池大小.
	TCPClientPoolSize int = 200

	// HTTPServerLogPath HTTP服务日志.
	HTTPServerLogPath string = "./log/http_server.log"
	// HTTPServerAddr HTTP服务地址.
	HTTPServerAddr string = ":1088"

	// StaticFilePath http静态文件服务地址.
	StaticFilePath string = "../static/"

	// DefaultImagePath 默认头像.
	DefaultImagePath string = "andy.jpeg"
)
