package crud

import (
	"errors"
)

var (
	ErrChatNotFound = errors.New("chat not found")
)

type Chat struct {
	ID         int64 `db:"id"`
	Authorized bool  `db:"authorized"`
}

func CreateChat(c Chat) (created Chat, err error) {

	query := `INSERT INTO "chats" VALUES ($1, $2) RETURNING *`
	err = db.Get(&created, query, c.ID, c.Authorized)

	return
}

func GetChat(id int64) (retrieved Chat, err error) {

	query := `SELECT * FROM "chats" WHERE id = $1 LIMIT 1`
	err = db.Get(&retrieved, query, id)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			err = ErrChatNotFound

		}

	}

	return
}

func UpdateChat(id int64, c Chat) (updated Chat, err error) {

	query := `UPDATE "chats" SET 
                   authorized = $2
               WHERE id = $1 RETURNING *`
	err = db.Get(&updated, query, id, c.Authorized)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			err = ErrChatNotFound

		}

	}

	return
}

func DeleteChat(id int64) (err error) {

	query := `DELETE FROM "chats" WHERE id = $1`
	_, err = db.Exec(query, id)

	return
}
