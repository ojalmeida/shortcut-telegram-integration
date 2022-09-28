package telegram

import (
	"bot/crud"
	"errors"
)

func isAuthorized(chatID int64) (ok bool) {

	chat, err := crud.GetChat(chatID)
	if err != nil {
		return
	}

	ok = chat.Authorized

	return
}

func authorize(chatID int64) (err error) {

	chat := crud.Chat{ID: chatID, Authorized: true}

	_, err = crud.UpdateChat(chatID, chat)
	if err != nil {
		if errors.Is(err, crud.ErrChatNotFound) {

			_, err = crud.CreateChat(chat)

		}
	}

	return
}
