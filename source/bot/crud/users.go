package crud

import (
	"errors"
)

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrMentionNameAlreadyExists = errors.New("mention username already in use")
)

type User struct {
	ID                  string `db:"id"`
	TelegramChatID      int64  `db:"telegram_chat_id"`
	ShortcutMentionName string `db:"shortcut_mention_name"`
	NotificationRate    int    `db:"notification_rate"`
}

func CreateUser(u User) (created User, err error) {

	query := `INSERT INTO "users" VALUES ($1, $2, $3, $4) RETURNING *`
	err = db.Get(&created, query, u.ID, u.TelegramChatID, u.ShortcutMentionName, u.NotificationRate)

	return
}

func GetUser(id string) (retrieved User, err error) {

	query := `SELECT * FROM "users" WHERE id = $1 LIMIT 1`
	err = db.Get(&retrieved, query, id)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			err = ErrUserNotFound

		}

	}

	return
}

func GetUserByChatID(chatID string) (retrieved User, err error) {
	query := `SELECT * FROM "users" WHERE telegram_chat_id = $1 LIMIT 1`
	err = db.Get(&retrieved, query, chatID)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			err = ErrUserNotFound

		}

	}

	return
}

func GetUserByMentionName(name string) (retrieved User, err error) {

	query := `SELECT * FROM "users" WHERE shortcut_mention_name = $1 LIMIT 1`
	err = db.Get(&retrieved, query, name)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			err = ErrUserNotFound

		}

	}

	return
}

func UpdateUser(id string, u User) (updated User, err error) {

	query := `UPDATE "users" SET 
                   telegram_chat_id = $2, 
                   shortcut_mention_name = $3,
                   notification_rate = $4
               WHERE id = $1 RETURNING *`
	err = db.Get(&updated, query, id, u.TelegramChatID, u.ShortcutMentionName, u.NotificationRate)

	return
}

func DeleteUser(id string) (err error) {

	query := `DELETE FROM "users" WHERE id = $1`
	_, err = db.Exec(query, id)

	return
}
