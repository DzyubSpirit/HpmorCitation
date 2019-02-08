package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	maxID        = 4e18
	citationKind = "Citation"
)

// Citation is a struct for storing in datastore.
type Citation struct {
	ID   *datastore.Key
	Text string
}

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	botToken := os.Getenv("TELEGRAM_APITOKEN")

	http.HandleFunc("/"+botToken, func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		log.Debugf(ctx, "request")

		var update tgbotapi.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Errorf(ctx, "reading update from the request body: %v", err)
			return
		}

		if update.Message == nil || !update.Message.IsCommand() {
			return
		}
		log.Debugf(ctx, "command: %s", update.Message.Command())

		dsClient, err := datastore.NewClient(ctx, projectID)
		if err != nil {
			log.Errorf(ctx, "making datastore client: %v", err)
			return
		}

		bot, err := tgbotapi.NewBotAPI(botToken)
		if err != nil {
			log.Errorf(ctx, "making telegram bot api: %v", err)
			return
		}

		switch update.Message.Command() {
		case "usage":
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "/add <citation> - adds citation\n/cite - sends random citation")
			if _, err := bot.Send(msg); err != nil {
				log.Errorf(ctx, "sending message about saved citation: %v", err)
				return
			}
		case "add":
			v := update.Message.CommandArguments()
			k := datastore.IDKey(citationKind, rand.Int63n(maxID), nil)
			cit := &Citation{k, v}
			if _, err := dsClient.Put(ctx, k, cit); err != nil {
				log.Errorf(ctx, "puting new citation into the datastore: %v", err)
				return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Citation's saved")
			if _, err := bot.Send(msg); err != nil {
				log.Errorf(ctx, "sending message about saved citation: %v", err)
				return
			}
		case "cite":
			q := datastore.NewQuery(citationKind)
			n, err := dsClient.Count(ctx, q)
			if err != nil {
				log.Errorf(ctx, "getting count of citations from the datastore: %v", err)
				return
			}
			q = q.Offset(rand.Intn(n)).Limit(1)
			it := dsClient.Run(ctx, q)
			if err != nil {
				log.Errorf(ctx, "running the query for one citation from the datastore: %v", err)
				return
			}

			cit := new(Citation)
			_, err = it.Next(cit)
			if err != nil {
				log.Errorf(ctx, "getting citation from the iterator: %v", err)
				return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, cit.Text)
			if _, err := bot.Send(msg); err != nil {
				log.Errorf(ctx, "sending message about saved citation: %v", err)
				return
			}
		}

		w.Write([]byte("ok"))
	})
	appengine.Main()
}
