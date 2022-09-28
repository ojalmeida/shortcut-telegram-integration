package config

import "github.com/spf13/viper"

var (
	Config Configuration
)

func init() {
	viper.SetDefault("database.addr", "localhost:5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.database-name", "postgres")
	viper.SetDefault("telegram.token", "")
	viper.SetDefault("telegram.authorization-token", "21f4257178c7cde44926b16f1482e2bf")
	viper.SetDefault("telegram.notification-rating", "7200")
	viper.SetDefault("shortcut.token", "")

	err := viper.BindEnv("database.addr", "DB_ADDR")
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
	err = viper.BindEnv("telegram.token", "TELEGRAM_TOKEN")
	if err != nil {
		return
	}
	err = viper.BindEnv("telegram.authorization-token", "TELEGRAM_AUTHORIZATION_TOKEN")
	if err != nil {
		return
	}
	err = viper.BindEnv("telegram.notification-rating", "TELEGRAM_NOTIFICATION_RATING")
	if err != nil {
		return
	}
	err = viper.BindEnv("shortcut.token", "SHORTCUT_TOKEN")
	if err != nil {
		return
	}

	Config.Database.Addr = viper.GetString("database.addr")
	Config.Database.User = viper.GetString("database.user")
	Config.Database.Password = viper.GetString("database.password")
	Config.Database.DatabaseName = viper.GetString("database.database-name")
	Config.Telegram.Token = viper.GetString("telegram.token")
	Config.Telegram.AuthorizationToken = viper.GetString("telegram.authorization-token")
	Config.Telegram.NotificationRating = viper.GetInt("telegram.notification-rating")
	Config.Shortcut.Token = viper.GetString("shortcut.token")

}

type Configuration struct {
	Database struct {
		Addr         string
		User         string
		Password     string
		DatabaseName string
	}
	Telegram struct {
		Token              string
		AuthorizationToken string
		NotificationRating int
	}
	Shortcut struct {
		Token string
	}
}
