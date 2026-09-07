package memes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MemeResponse struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type MemClient struct {
	memKey string
}

func NewMem(memKey string) MemClient {
	return MemClient{
		memKey: memKey,
	}
}

func (m MemClient) Mem(keyWord string) (MemeResponse, error) {
	url := "https://api.apileague.com/retrieve-random-meme?api-key=%s&keywords=%s"
	resp, err := http.Get(fmt.Sprintf(url, m.memKey, keyWord))
	if err != nil {
		return MemeResponse{}, fmt.Errorf("не удалось получить мемчеГ %w", err)
	}

	if resp.StatusCode != 200 {
		return MemeResponse{}, fmt.Errorf("ошибка статус кода: %d", resp.StatusCode)
	}

	var memeResponse MemeResponse
	err = json.NewDecoder(resp.Body).Decode(&memeResponse)
	if err != nil {
		return MemeResponse{}, fmt.Errorf("ошибка анмаршалинга мема: %w", err)
	}

	return MemeResponse{
		Url:         memeResponse.Url,
		Description: memeResponse.Description,
	}, nil
}
