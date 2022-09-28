package config

import "github.com/spf13/viper"

var (
	Config Configuration
)

func init() {
	viper.SetDefault("server.addr", "0.0.0.0:8080")
	viper.SetDefault("database.addr", "localhost:5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.database-name", "postgres")
	viper.SetDefault("shortcut.token", "")

	err := viper.BindEnv("server.addr", "SERVER_ADDR")
	if err != nil {
		return
	}
	err = viper.BindEnv("database.addr", "DB_ADDR")
	if err != nil {
		return
	}
	err = viper.BindEnv("database.user", "DB_USER")
	if err != nil {
		return
	}
	err = viper.BindEnv("database.password", "DB_PASSWORD")
	if err != nil {
		return
	}
	err = viper.BindEnv("database.database-name", "DB_DBNAME")
	if err != nil {
		return
	}
	err = viper.BindEnv("shortcut.token", "SHORTCUT_TOKEN")
	if err != nil {
		return
	}

	Config.Server.Addr = viper.GetString("server.addr")
	Config.Database.Addr = viper.GetString("database.addr")
	Config.Database.User = viper.GetString("database.user")
	Config.Database.Password = viper.GetString("database.password")
	Config.Database.DatabaseName = viper.GetString("database.database-name")
	Config.Shortcut.Token = viper.GetString("shortcut.token")

}

type Configuration struct {
	Server struct {
		Addr string
	}
	Database struct {
		Addr         string
		User         string
		Password     string
		DatabaseName string
	}
	Shortcut struct {
		Token string
	}
}
