package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/teamteamdev/ugractf-2024-x/tasks/moderator/internal/quiz"
	"github.com/teamteamdev/ugractf-2024-x/tasks/moderator/internal/store"
)

type configStruct struct {
	BotToken        string `json:"bot_token"`
	DBURI           string `json:"db_uri"`
	RootURL         string `json:"root_url"`
	Bind            string `json:"bind"`
	VoteChatId      int64  `json:"vote_chat_id"`
	ModeratorChatId int64  `json:"moderator_chat_id"`
}

var (
	configFile   = flag.String("config", "", "path to config json file")
	periodVoting = flag.Duration("period-voting", 30*time.Minute, "period between voting epoch")

	cfg configStruct
)

func main() {
	flag.Parse()

	if f, err := os.Open(*configFile); err != nil {
		log.Fatalf("failed open bot config file: %v", err)
	} else if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		log.Fatalf("failed read bot config: %v", err)
	} else {
		f.Close()
	}

	wh := &tele.Webhook{
		AllowedUpdates: []string{
			"message",
			"callback_query",
			"chat_member",
			"chat_join_request",
		},
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: cfg.RootURL + "/" + cfg.BotToken,
		},
	}
	http.Handle("/"+cfg.BotToken, wh)
	bot, err := tele.NewBot(tele.Settings{
		Token: cfg.BotToken,
		OnError: func(err error, c tele.Context) {
			log.Printf("failed handle update: %v", err)
		},
		Poller: wh,
	})
	if err != nil {
		log.Fatalf("failed create bot: %v", err)
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)

	log.Printf("opening database...")
	db, err := store.NewDB(cfg.DBURI)
	if err != nil {
		log.Fatalf("failed open database: %v", err)
	}
	defer db.Close()

	initBot(bot, db)
	quiz.SetupHTTPHandlers(db, quizDoneCallback(bot, db))

	server := http.Server{Addr: cfg.Bind}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Print("starting http server...")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			cancel(err)
		}
		log.Printf("stopped http server")
	}()

	go func() {
		log.Print("starting bot...")
		bot.Start()
	}()

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(*periodVoting)
		defer ticker.Stop()
		if err := votingEpoch(bot, db); err != nil {
			log.Printf("failed voting epoch: %v", err)
		}
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case <-ticker.C:
				if err := votingEpoch(bot, db); err != nil {
					log.Printf("failed voting epoch: %v", err)
				}
			}
		}
	}()

	log.Print("started")

	<-ctx.Done()
	log.Print("stopping...")
	server.Close()
	exitCode := 0
	if err := context.Cause(ctx); !errors.Is(err, context.Canceled) {
		log.Printf("cause of shutdown: %v", err)
		exitCode = 1
	}

	wg.Wait()
	log.Printf("stopped")
	db.Close()
	os.Exit(exitCode)
}
