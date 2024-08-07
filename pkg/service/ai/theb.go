package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

type Theb struct {
}

type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

const (
	completionUrl = "https://api.theb.ai/v1/chat/completions"
)

func (ai *Theb) Completion(content string) (string, error) {
	thebAiAPIKey := os.Getenv("THEB_API_KEY")

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": viper.GetString("ai.theb.model"),
		"messages": []map[string]string{
			{"role": "user", "content": content},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", completionUrl, bytes.NewBuffer(requestBody))
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

	var gptResponse GPTResponse
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&gptResponse); err != nil {
		return "", err
	}

	if len(gptResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in GPT response")
	}

	return gptResponse.Choices[0].Message.Content, nil
}
