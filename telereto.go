package main

import (
	"fmt"
	"golang.org/x/exp/slices"
	tele "gopkg.in/telebot.v3"
	"log"
	"net/http"
	"net/url"
)

var ACCEPTABLE_TYPES = []string{"image/png", "image/gif", "image/jpeg", "image/bmp", "image/webp"}[:]
var CHEVERETO_HOST_API_URL, CHEVERETO_API_KEY string

func main() {
	config_file, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	config, err := NewConfig(config_file)

	// Setup tg webhook
	end_point := url.URL{
		Scheme: config.BotServer.Schema,
		Host:   config.BotServer.Host,
		Path:   config.BotServer.ApiKey,
	}

	che_url := url.URL{
		Scheme: config.CheveretoServer.Schema,
		Host:   config.CheveretoServer.Host,
		Path:   "/api/1/upload",
	}

	CHEVERETO_HOST_API_URL = che_url.String()
	CHEVERETO_API_KEY = config.CheveretoServer.ApiKey

	webhook := &tele.Webhook{
		Listen:   config.BotServer.Listen,
		Endpoint: &tele.WebhookEndpoint{PublicURL: end_point.String()},
	}
	pref := tele.Settings{
		Token:  config.BotServer.ApiKey,
		Poller: webhook,
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle(tele.OnText, func(context tele.Context) error {
		// If a text is an url, try to fetch and upload that image.
		text := context.Text()
		parsed_url, err := url.ParseRequestURI(text)
		if err != nil {
			return context.Send("请上传一张图片")
		}

		url_string := parsed_url.String()
		resp, err := http.Get(url_string)
		if err != nil {
			return context.Send("请上传一张图片")
		}
		if slices.Contains(ACCEPTABLE_TYPES, resp.Header.Get("Content-Type")) {
			return context.Send(upload_photo(url_string))
		}
		return context.Send("请上传一张图片")
	})

	bot.Handle(tele.OnPhoto, func(context tele.Context) error {
		photo := context.Message().Photo
		f, err := bot.FileByID(photo.FileID)
		if err != nil {
			return context.Send("Failed to fetch photo url from Telegram, please try again later!")
		}

		tg_url := fmt.Sprint("https://api.telegram.org/file/bot", config.BotServer.ApiKey, "/", f.FilePath)
		return context.Send(upload_photo(tg_url))
	})

	bot.Handle(tele.OnDocument, func(context tele.Context) error {
		doc := context.Message().Document
		fmt.Println(doc.MIME)
		if !slices.Contains(ACCEPTABLE_TYPES, doc.MIME) {
			return context.Send(fmt.Sprint("我们不接受 ", doc.MIME, " 类型的文件"))
		}

		f, err := bot.FileByID(doc.FileID)
		if err != nil {
			return context.Send("Failed to fetch the file url from Telegram, please try again later!")
		}

		tg_url := fmt.Sprint("https://api.telegram.org/file/bot", config.BotServer.ApiKey, "/", f.FilePath)
		return context.Send(upload_photo(tg_url))
	})

	bot.Handle("/start", func(context tele.Context) error {
		return context.Send("您好，欢迎使用 Repost.ink 图床。\n给我发送的照片会自动上传到图床。")
	})

	bot.Start()
}

func upload_photo(tg_url string) string {
	chevereto_resp, err := call_chevereto_api(tg_url)
	if err != nil {
		return fmt.Sprint("上传失败了，以下信息对我们很有帮助\n", err)
	}
	reply := fmt.Sprint(
		"成功上传图片，您可以通过以下地址访问原图:\n",
		chevereto_resp.OriginUrl,
		"\n也可以点击 ", chevereto_resp.SiteUrl, " 在我们的网站查看。\n",
		"或者，您也可以通过以下链接删除图片:\n",
		chevereto_resp.DeleteUrl,
	)
	return reply
}
