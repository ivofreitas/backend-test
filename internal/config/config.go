package config

import (
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Env values
type Env struct {
	Server   Server
	Log      Log
	Doc      Doc
	MySQL    MySQL
	Security Security
}

// Server config
type Server struct {
	Host     string
	BasePath string
	Port     string
}

type Authorization struct {
	Secret string
}

// Log config
type Log struct {
	Enabled bool
	Level   string
}

// Doc - swagger information
type Doc struct {
	Title       string
	Description string
	Enabled     bool
	Version     string
}

// MySQL - db conn information
type MySQL struct {
	Username     string
	Password     string
	Host         string
	Database     string
	PoolConn     int
	Timeout      time.Duration
	ConnLifetime time.Duration
}

type Security struct {
	SecretKey string
}

var (
	env  *Env
	once sync.Once
)

// GetEnv returns env values
func GetEnv() *Env {

	once.Do(func() {

		viper.AutomaticEnv()
		if err := godotenv.Load("internal/config/.env"); err != nil {
			log.Println(err)
		}

		env = new(Env)
		env.Server.Host = viper.GetString("SERVER_HOST")
		env.Server.BasePath = viper.GetString("SERVER_BASE_PATH")
		env.Server.Port = viper.GetString("SERVER_PORT")

		env.Log.Enabled = viper.GetBool("LOG_ENABLED")
		env.Log.Level = viper.GetString("LOG_LEVEL")

		env.Doc.Title = viper.GetString("DOC_TITLE")
		env.Doc.Description = viper.GetString("DOC_DESCRIPTION")
		env.Doc.Enabled = viper.GetBool("DOC_ENABLED")
		env.Doc.Version = viper.GetString("DOC_VERSION")

		env.MySQL.Username = viper.GetString("MYSQL_USERNAME")
		env.MySQL.Password = viper.GetString("MYSQL_PASSWORD")
		env.MySQL.Host = viper.GetString("MYSQL_HOST")
		env.MySQL.Database = viper.GetString("MYSQL_DATABASE")
		env.MySQL.PoolConn = viper.GetInt("MYSQL_POOL_CONN")
		env.MySQL.Timeout = viper.GetDuration("MYSQL_QUERY_TIMEOUT")
		env.MySQL.ConnLifetime = viper.GetDuration("MYSQL_CONN_LIFETIME")

		env.Security.SecretKey = viper.GetString("SECURITY_SECRET_KEY")
	})
	return env
}
