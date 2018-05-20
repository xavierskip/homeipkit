package main

import "net/http"
import "net/url"
import "fmt"

const APIEndpoint = "https://api.telegram.org/bot%s/%s"

type TelegramApi struct {
    Token string
    Client *http.Client
}

func NewBot(token string) *TelegramApi {
    bot := &TelegramApi{
        Token:  token,
        Client: &http.Client{},
    }
    return bot
}

func (bot TelegramApi) SendMessage(chat_id string, text string) (*http.Response, error) {
    method := "sendMessage"
    v := url.Values{}
    v.Set("chat_id", chat_id)
    v.Set("text", text)
    url := fmt.Sprintf(APIEndpoint, bot.Token, method)
    resp, err := bot.Client.PostForm(url, v)
    if err != nil {
        return nil, err
    }
    return resp, nil
}