package main

import (	
	"log"
	"net/http"
	"os"
	"fmt"
	"time"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)


func main() {	
	
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintf(w, "Hello World")
		})

		http.ListenAndServe(":8080", nil)
	}()

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{		
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	updater.Idle()
}


func start(b *gotgbot.Bot, ctx *ext.Context) error {    
    message := ctx.Message
	
    _, err := message.Reply(b, "Hey! I'm PinterestBot. You can search for Pinterest videos or photos, and you can download them too. I can also provide Google images, Bing images, etc.", &gotgbot.SendMessageOpts{})
    if err != nil {
        return nil
    }

    return nil 
}

