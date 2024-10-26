package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
)

type Proverb struct {
	Address string
	Text    string
}

func main() {
	// Получаем массив поговорок
	proverbs := newProverbs()
	fmt.Println(getRandomProverb(proverbs))
}

// newProverbs - конструктор
func newProverbs() []Proverb {
	html, err := getHTML()
	if err != nil {
		// Если не получили ответ, то дальше нечего показывать
		log.Fatal(err)
	}

	p := parseHTML(html)
	if len(p) < 1 {
		log.Fatal("Поговорки не получены")
	}
	return p
}

// getHTML - получает ответ web сервера с поговорками https://go-proverbs.github.io/
func getHTML() (string, error) {
	url := "https://go-proverbs.github.io/"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// parseHTML - парсит HTML и получаем объект с поговоркой и ссылкой
func parseHTML(html string) (ar []Proverb) {
	re := regexp.MustCompile(`<h3><a href="(.+)">(.+)</a></h3>`)
	find := re.FindAllStringSubmatch(html, -1)

	for _, s := range find {
		ar = append(ar, Proverb{Address: s[1], Text: s[2]})
	}
	return
}

// getRandomProverb - получить случайную поговорку
func getRandomProverb(ar []Proverb) Proverb {
	return ar[rand.Intn(len(ar))]
}
