package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/wizardloong/botposter/pkg/service/ai"
	"github.com/wizardloong/botposter/pkg/service/parser"
)

func (s *Service) RewriteArticle(url string) (string, []byte, error) {

	proc := s.NewParser(url)
	if proc == nil {
		return "", nil, errors.New("don't understand that link")
	}

	content, title, imageURL, err := proc.FetchArticleData(url)
	if err != nil {
		return "", nil, errors.New("cannot fetch article content")
	}

	prompt := fmt.Sprintf(viper.GetString("prompt.rw")+": %s", content)

	// Отправка текста статьи в нейросеть
	rewrittenText, err := s.NewAI().Completion(prompt)
	if err != nil {
		return "", nil, errors.New("failed to rewrite article content")
	}

	// Скачивание изображения из статьи
	image, err := proc.DownloadImage(imageURL)
	if err != nil {
		return "", nil, errors.New("cannot download image")
	}
	// Формирование поста для Телеграма
	post := fmt.Sprintf("👾 <b>%s</b>\n\n<i>%s</i>\n", title, rewrittenText)

	return post, image, nil
}

func (s *Service) NewParser(url string) parser.Parser {
	if strings.Contains(url, "shazoo.ru") {
		return &parser.Shazoo{}
	}

	return nil
}

func (s *Service) NewAI() ai.AI {
	return &ai.Theb{}
}
