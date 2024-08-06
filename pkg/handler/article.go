package handler

import (
	"bytes"
	"os"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (h *Handler) article(m *tb.Message) {
	url := m.Text

	post, image, err := h.services.RewriteArticle(url)
	if err != nil {
		h.bot.Send(m.Sender, err.Error())
		return
	}

	h.preview(m.Sender, image, post)

	choiceYes := tb.InlineButton{
		Unique: "choice_yes",
		Text:   "Yes",
	}
	choiceNo := tb.InlineButton{
		Unique: "choice_publish",
		Text:   "No",
	}

	h.bot.Handle(&choiceYes, func(c *tb.Callback) {
		h.bot.Handle(tb.OnText, func(m *tb.Message) {
			post = post + "\n\n" + m.Text
			h.prepareAndPublish(c.Sender, post, image)
		})
	})
	h.bot.Handle(&choiceNo, func(c *tb.Callback) {
		h.prepareAndPublish(c.Sender, post, image)
	})

	h.bot.Send(m.Sender, "Would you like to add some words?", &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{choiceYes, choiceNo}}})
}

func (h *Handler) prepareAndPublish(sender *tb.User, post string, image []byte) {

	h.bot.Send(sender, "Here is completed post:", &tb.SendOptions{ParseMode: tb.ModeHTML})
	h.preview(sender, image, post)

	publishButton := tb.InlineButton{
		Unique: "publish",
		Text:   "Publish",
	}
	cancelButton := tb.InlineButton{
		Unique: "cancel",
		Text:   "Cancel",
	}
	h.bot.Handle(&publishButton, func(c *tb.Callback) {
		h.publish(sender, image, post)
	})
	h.bot.Handle(&cancelButton, func(c *tb.Callback) {
		h.bot.Send(sender, "Ok", &tb.SendOptions{ParseMode: tb.ModeHTML})
	})

	h.bot.Send(sender, "Publish?", &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{publishButton, cancelButton}}})
}

func (h *Handler) preview(sender *tb.User, image []byte, post string) {
	photo := &tb.Photo{
		File:    tb.FromReader(bytes.NewReader(image)),
		Caption: post,
	}

	h.bot.Send(sender, photo, &tb.SendOptions{ParseMode: tb.ModeHTML})
}

func (h *Handler) publish(sender *tb.User, image []byte, finalPost string) {
	channelID := os.Getenv("CHANNEL_ID")

	h.bot.Send(sender, "Successfully published.")
	chatID, _ := strconv.ParseInt(channelID, 10, 64)
	photo := &tb.Photo{
		File:    tb.FromReader(bytes.NewReader(image)),
		Caption: finalPost,
	}
	h.bot.Send(tb.ChatID(chatID), photo, &tb.SendOptions{ParseMode: tb.ModeHTML})
}
