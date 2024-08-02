package handler

import tb "gopkg.in/tucnak/telebot.v2"

func (h *Handler) hello(m *tb.Message) {
	h.bot.Send(m.Sender, "Привет! Отправь мне ссылку на статью с сайта Shazoo.ru.")
}
