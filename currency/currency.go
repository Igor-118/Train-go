package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"unicode"
)

type CurrencyClient struct {
	apiKey string
	client *http.Client
}

func NewCurrency(apiKey string) CurrencyClient {
	return CurrencyClient{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type CurrencyResponse struct {
	Valid bool               `json:"valid"`
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func (c CurrencyClient) GetCurrency(base string) (CurrencyResponse, error) {
	url := fmt.Sprintf("https://currencyapi.net/api/v2/rates?base=%s&output=json&key=%s", base, c.apiKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return CurrencyResponse{}, fmt.Errorf("error get запрос на мем: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return CurrencyResponse{}, fmt.Errorf("error get ответ для запроса http валюты: %w", err)
	}

	if res.StatusCode != 200 {
		return CurrencyResponse{}, fmt.Errorf("error fail get статус код для валюты: %d", res.StatusCode)
	}
	var currencyResponse CurrencyResponse
	err = json.NewDecoder(res.Body).Decode(&currencyResponse)
	if err != nil {
		return CurrencyResponse{}, fmt.Errorf("ошибка анмаршалинга валюты: %w", err)
	}
	// for key, value := range currencyResponse.Rates {
	// 	if key == "RUB" {
	// 		fmt.Println("значение валюты: ", value)
	// 	}
	// }
	return currencyResponse, nil
}

func ParseNumber(number string) (string, error) {
	if number == "" {
		return number, fmt.Errorf("ну ты сумму то напиши пёс")
	}
	// сам ты пёс, список доступных валют покажи
	if len(number) > 10 {
		return number, fmt.Errorf("ну и фигню ты написал, пошёл нахрен КАЗЬЁЛЬ")
	}
	var value string
	for _, digit := range number {
		if digit < '0' || digit > '9' {
			return number, fmt.Errorf("ТОЛЬКО ЦИФРЫ В НАЗВАНИИ ПИШИ")
		}
		value += string(digit)
	}
	return value, nil
}

func ParseCurrencyName(name string) (string, error) {
	if name == "" {
		return name, fmt.Errorf("ну ты название валюты то напиши пёс")
	}
	if len(name) != 3 {
		return name, fmt.Errorf("Название валюты должно быть длинною в 3 символа")
	}
	var currency string
	for _, letter := range name {
		Up := unicode.ToUpper(letter)
		if Up < 'A' || Up > 'Z' {
			return name, fmt.Errorf("ТОЛЬКО БУКВЫ В НАЗВАНИИ ПИШИ")
		}
		currency += string(letter)
	}
	return currency, nil
}

// func ParseCoord(s string) (Position, error) {
// 	if len(s) > 3 {
// 		return Position{}, fmt.Errorf("Шляпу мне не пиши тут, давай нормальное значение не длиннее 3 символов")
// 	}
// 	pos := Position{}
// 	var x rune
// 	var y string
// 	massiv := []rune{}
// 	for _, sign := range s {
// 		massiv = append(massiv, sign)
// 	}
// 	x = rune(massiv[0])
// 	Lowx := unicode.ToLower(x)
// 	LowX := rune(Lowx)
// 	if rune(LowX) < 'a' || rune(LowX) > 'j' {
// 		return Position{}, fmt.Errorf("Слыш чувырло, букву нормальную напиши")
// 	}
// 	X := int(LowX - 'a')
// 	y = string(massiv[1:])
// 	Y, err := strconv.Atoi(y)
// 	if err != nil {
// 		return Position{}, fmt.Errorf("ЛЭЭЭ проблема при конвертации ээ")
// 	}
// 	Y = Y - 1
// 	if Y > 9 {
// 		return Position{}, fmt.Errorf("ЛЭЭЭ больше 10 нельзя")
// 	}
// 	if Y < 0 {
// 		return Position{}, fmt.Errorf("ЛЭЭЭ меньше 0 нельзя")
// 	}
// 	pos = Position{X, Y}
// 	return pos, nil
// }
