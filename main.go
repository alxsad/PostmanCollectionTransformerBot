package main

import (
	"errors"
	"flag"
	"github.com/yanzay/tbot/v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var token = flag.String("token", "", "Telegram bot token")

func main() {

	flag.Parse()
	bot := tbot.New(*token)
	c := bot.Client()

	bot.HandleMessage(".*", func(m *tbot.Message) {
		if m.Document == nil {
			return
		}
		log.Println("New message received")
		c.SendChatAction(m.Chat.ID, tbot.ActionUploadDocument)

		doc, err := c.GetFile(m.Document.FileID)
		if err != nil {
			c.SendMessage(m.Chat.ID, err.Error())
			return
		}
		url := c.FileURL(doc)
		resp, err := http.Get(url)
		if err != nil {
			c.SendMessage(m.Chat.ID, err.Error())
			return
		}
		defer resp.Body.Close()

		tmp, err := ioutil.TempFile("", m.Document.FileName)
		if err != nil {
			c.SendMessage(m.Chat.ID, err.Error())
			return
		}
		defer os.Remove(tmp.Name())
		log.Printf("Created file: %s", tmp.Name())

		io.Copy(tmp, resp.Body)
		tmp.Close()

		result, err := transform(tmp.Name())
		if err != nil {
			c.SendMessage(m.Chat.ID, err.Error())
			return
		}

		c.SendDocumentFile(m.Chat.ID, result)
	})

	log.Println("Bot started")
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func transform(filepath string) (string, error) {
	cmd := exec.Command(
		"postman-collection-transformer",
		"convert",
		"-i",
		filepath,
		"-o",
		filepath + ".v2.json",
		"-j",
		"1.0.0",
		"-p",
		"2.0.0",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	if strings.Contains(string(out), "[ERROR]") {
		return "", errors.New("incorrect file")
	}
	return filepath + ".v2.json", nil
}

