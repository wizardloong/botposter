package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/wizardloong/botposter/pkg/handler"
	"github.com/wizardloong/botposter/pkg/service"
	tb "gopkg.in/tucnak/telebot.v2"
)

type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	telegramToken := os.Getenv("BOT_TOKEN")
	channelID := os.Getenv("CHANNEL_ID") // Убедитесь, что это значение является числовым ID канала

	log.Println(telegramToken)
	log.Println(channelID)

	b, err := tb.NewBot(tb.Settings{
		Token:  telegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	services := service.NewServices()
	handler.NewHandler(services, b)

	b.Handle(tb.OnText, func(m *tb.Message) {
		url := m.Text
		if strings.Contains(url, "shazoo.ru") {
			content, title, imageURL, err := fetchArticleData(url)
			if err != nil {
				b.Send(m.Sender, "Не удалось получить контент статьи.")
				return
			}

			// Отправка текста статьи в нейросеть
			rewrittenText, err := sendToAI(content)
			if err != nil {
				b.Send(m.Sender, "Не удалось переписать текст статьи.")
				return
			}

			// Скачивание изображения из статьи
			image, err := downloadImage(imageURL)
			if err != nil {
				b.Send(m.Sender, "Не удалось скачать изображение.")
				return
			}

			// Формирование поста для Телеграма
			post := fmt.Sprintf("👾 <b>%s</b>\n<i>%s</i>\n", title, rewrittenText)
			b.Send(m.Sender, post, &tb.SendOptions{ParseMode: tb.ModeHTML})
			b.Send(m.Sender, "Добавьте несколько слов от себя.", &tb.SendOptions{ParseMode: tb.ModeHTML})

			b.Handle(tb.OnText, func(m *tb.Message) {
				var finalPost string
				if len(m.Text) != 0 {
					finalPost = post + "\n--------\n" + m.Text
				} else {
					finalPost = post
				}

				b.Send(m.Sender, "Вот итоговый пост:\n"+finalPost, &tb.SendOptions{ParseMode: tb.ModeHTML})

				publishButton := tb.InlineButton{
					Unique: "publish",
					Text:   "Опубликовать",
				}
				editButton := tb.InlineButton{
					Unique: "edit",
					Text:   "Исправить",
				}
				b.Handle(&publishButton, func(c *tb.Callback) {
					b.Send(m.Sender, "Пост опубликован.")
					chatID, _ := strconv.ParseInt(channelID, 10, 64)
					photo := &tb.Photo{
						File:    tb.FromReader(bytes.NewReader(image)),
						Caption: finalPost,
					}
					b.Send(tb.ChatID(chatID), photo, &tb.SendOptions{ParseMode: tb.ModeHTML})
				})
				b.Handle(&editButton, func(c *tb.Callback) {
					b.Send(m.Sender, "Напишите, что нужно исправить.")
					b.Handle(tb.OnText, func(m *tb.Message) {
						finalPost = post + "\n" + m.Text
						b.Send(m.Sender, "Вот итоговый пост:\n"+finalPost, &tb.SendOptions{ParseMode: tb.ModeHTML})
					})
				})

				b.Send(m.Sender, "Опубликовать или Исправить?", &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{publishButton, editButton}}})
			})
		} else {
			b.Send(m.Sender, "Пожалуйста, отправьте ссылку на статью с сайта Shazoo.ru.")
		}
	})

	b.Start()
}

func fetchArticleData(url string) (string, string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	content := doc.Find("section.Entry__content").Text()
	title := doc.Find("article section h1").Text()
	imageURL, exists := doc.Find("figure img.w-full").Attr("src")
	if !exists {
		return "", "", "", fmt.Errorf("image not found")
	}

	return content, title, imageURL, nil
}

func sendToAI(content string) (string, error) {
	thebAiAPIKey := os.Getenv("THEB_API_KEY")

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": fmt.Sprintf("Перепиши этот текст, сократив его до %s символов или меньше: %s", os.Getenv("POST_LEN"), content)},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.theb.ai/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", thebAiAPIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Отладочная информация
	fmt.Println("Response Body:", string(body))

	var gptResponse GPTResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&gptResponse); err != nil {
		return "", err
	}

	if len(gptResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in GPT response")
	}

	return gptResponse.Choices[0].Message.Content, nil
}

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
