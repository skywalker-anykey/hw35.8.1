package main

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"time"
)

type Proverb struct {
	Address string
	Text    string
}

const (
	// Настройки сервера
	addr     = "localhost:12345"
	protocol = "tcp4"
	// Сайт с поговорками
	url = "https://go-proverbs.github.io/"
)

func main() {
	// Получаем массив поговорок
	proverbs := newProverbs()

	// Создаем сервер
	listener, err := net.Listen(protocol, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	// Обрабатываем подключения
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleConn(conn, &proverbs)
	}
}

func handleConn(conn net.Conn, proverb *[]Proverb) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	log.Println("New connection from:", conn.RemoteAddr())

	for {
		p := getRandomProverb(*proverb)
		_, _ = conn.Write([]byte(p.Text + "\n\r"))
		_, _ = conn.Write([]byte(p.Address + "\n\r"))
		time.Sleep(time.Second * 3)
	}
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
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Рекомендация ментора после проверки
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

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
	rand.Seed(time.Now().UnixNano()) // Рекомендация ментора после проверки, но я бы не добавлял
	return ar[rand.Intn(len(ar))]
}
