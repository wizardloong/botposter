package handler

import (
	"bytes"
	"os"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

type ArticleHandler struct {
	Handler
}

func (h *ArticleHandler) article(m *tb.Message) {
	url := m.Text

	post, image, err := h.services.RewriteArticle(url)
	if err != nil {
		h.bot.Send(m.Sender, err.Error())
		return
	}

	h.bot.Send(m.Sender, post, &tb.SendOptions{ParseMode: tb.ModeHTML})
	h.bot.Send(m.Sender, "Would you like to add some words?", &tb.SendOptions{ParseMode: tb.ModeHTML})

	h.bot.Handle(tb.OnText, func(m *tb.Message) {
		var finalPost string
		if len(m.Text) != 0 {
			finalPost = post + "\n--------\n" + m.Text
		} else {
			finalPost = post
		}

		h.bot.Send(m.Sender, "Here is completed post:\n"+finalPost, &tb.SendOptions{ParseMode: tb.ModeHTML})

		publishButton := tb.InlineButton{
			Unique: "publish",
			Text:   "Publish",
		}
		editButton := tb.InlineButton{
			Unique: "edit",
			Text:   "Edit",
		}
		h.bot.Handle(&publishButton, func(c *tb.Callback) {
			h.publish(m.Sender, image, finalPost)
		})
		h.bot.Handle(&editButton, func(c *tb.Callback) {
			h.edit(m.Sender, post)
		})

		// h.bot.Send(m.Sender, "Опубликовать или Исправить?", &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{publishButton, editButton}}})
	})
}

func (h *ArticleHandler) publish(sender *tb.User, image []byte, finalPost string) {
	channelID := os.Getenv("CHANNEL_ID")

	h.bot.Send(sender, "Successfully published.")
	chatID, _ := strconv.ParseInt(channelID, 10, 64)
	photo := &tb.Photo{
		File:    tb.FromReader(bytes.NewReader(image)),
		Caption: finalPost,
	}
	h.bot.Send(tb.ChatID(chatID), photo, &tb.SendOptions{ParseMode: tb.ModeHTML})
}

func (h *ArticleHandler) edit(sender *tb.User, post string) {
	h.bot.Send(sender, "What to edit?")
	h.bot.Handle(tb.OnText, func(m *tb.Message) {
		finalPost := post + "\n" + m.Text
		h.bot.Send(m.Sender, "Here is completed post:\n"+finalPost, &tb.SendOptions{ParseMode: tb.ModeHTML})
	})
}
