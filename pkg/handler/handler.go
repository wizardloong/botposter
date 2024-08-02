package handler

import (
	"github.com/wizardloong/botposter/pkg/service"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	services *service.Service
	bot      *tb.Bot
}

func NewHandler(services *service.Service, bot *tb.Bot) *Handler {
	return &Handler{
		services: services,
		bot:      bot,
	}
}

func (h *Handler) InitRoutes() {
	h.bot.Handle("/start", h.hello)
}
