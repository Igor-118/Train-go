package facts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FactCliet struct {
	apiKey string
	client *http.Client
}

func NewFact(apiKey string) FactCliet {
	return FactCliet{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type FactResponse struct {
	Fact string `json:"fact"`
}

func (f FactCliet) GetFacts() (FactResponse, error) {
	url := "https://api.api-ninjas.com/v1/facts"
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return FactResponse{}, fmt.Errorf("error get запрос на мем: %w", err)
	}

	req.Header.Set("X-Api-Key", f.apiKey) //req.Header.Set("X-Api-Key", "Твой ключ")
	res, err := f.client.Do(req)

	if err != nil {
		return FactResponse{}, fmt.Errorf("error get ответ для запроса http факта: %w", err)
	}

	if res.StatusCode != 200 {
		return FactResponse{}, fmt.Errorf("error fail get статус код для факта: %d", res.StatusCode)
	}
	fmt.Println(res.Body)
	var factResponse []FactResponse
	err = json.NewDecoder(res.Body).Decode(&factResponse)
	if err != nil {
		return FactResponse{}, fmt.Errorf("ошибка анмаршалинга факта: %w", err)
	}

	return factResponse[0], nil
}
