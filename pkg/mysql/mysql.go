package mysql

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	db         *sqlx.DB // 主库：users / contact / apply 等（不含私聊 messages 表）
	messagesDB *sqlx.DB // 私聊消息库 messages 表；未配置 mysql_messages.dbname 时回落到主库 db
	docsDB     *sqlx.DB // 文件登记表 docs；未配置 mysql_docs.dbname 时为 nil
)

func InitMysql() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	return
}

func DB() *sqlx.DB {
	return db
}

// connectSecondary 连接独立库；prefix 为 viper 键前缀如 mysql_messages、mysql_docs。
func connectSecondary(prefix string) (*sqlx.DB, error) {
	dbname := strings.TrimSpace(viper.GetString(prefix + ".dbname"))
	if dbname == "" {
		return nil, nil
	}
	host := viper.GetString(prefix + ".host")
	if host == "" {
		host = viper.GetString("mysql.host")
	}
	port := viper.GetInt(prefix + ".port")
	if port == 0 {
		port = viper.GetInt("mysql.port")
	}
	user := viper.GetString(prefix + ".user")
	if user == "" {
		user = viper.GetString("mysql.user")
	}
	pass := viper.GetString(prefix + ".password")
	if pass == "" {
		pass = viper.GetString("mysql.password")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		user, pass, host, port, dbname)
	conn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// InitMessagesMysql 私聊消息库（表 messages）。未配置 dbname 时回落到主库 DB()。
func InitMessagesMysql() error {
	conn, err := connectSecondary("mysql_messages")
	if err != nil {
		zap.L().Error("connect messages DB failed", zap.Error(err))
		return err
	}
	if conn == nil {
		zap.L().Info("mysql_messages.dbname empty, table messages uses main mysql.dbname")
		messagesDB = nil
		return nil
	}
	messagesDB = conn
	zap.L().Info("messages mysql connected", zap.String("dbname", viper.GetString("mysql_messages.dbname")))
	return nil
}

// MessagesDB 私聊 messages 表所在连接；未配独立库时与主库相同。
func MessagesDB() *sqlx.DB {
	if messagesDB != nil {
		return messagesDB
	}
	return db
}

// InitDocsMysql 文件登记表 docs 所在库；未配置 dbname 时不连接、不写 docs。
func InitDocsMysql() error {
	conn, err := connectSecondary("mysql_docs")
	if err != nil {
		zap.L().Error("connect docs DB failed", zap.Error(err))
		return err
	}
	if conn == nil {
		zap.L().Info("mysql_docs.dbname empty, skip docs metadata DB")
		docsDB = nil
		return nil
	}
	docsDB = conn
	zap.L().Info("docs mysql connected", zap.String("dbname", viper.GetString("mysql_docs.dbname")))
	return nil
}

// DocsDB 可能为 nil。
func DocsDB() *sqlx.DB {
	return docsDB
}
