package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/teamteamdev/ugractf-2024-x/tasks/moderator/internal/store"
	"github.com/teamteamdev/ugractf-2024-x/tasks/moderator/internal/voting"
)

var (
	Debug          = flag.Bool("debug", false, "use it to enable debug logging")
	VotingDuration = flag.Duration("voting-duration", 7*time.Second, "time before voting end")
)

func banUser(db *store.DB, c tele.Context) error {
	if *Debug {
		log.Printf("banning user: %v", c.Sender().ID)
	}

	if err := store.Ban(db, c.Sender().ID); err != nil {
		return errors.Join(
			err,
			c.RespondText("Очень жаль..."),
			c.Send("Я хотел отправить тебя в бан. Но не смог. Так что попробуй еще раз!"),
		)
	}
	c.Bot().Ban(
		&tele.Chat{ID: cfg.VoteChatId},
		&tele.ChatMember{User: c.Sender()},
	)

	err1 := c.RespondAlert("Забанен.")
	err2 := c.Send("Я с тобой больше не разговариваю.")
	time.Sleep(1500 * time.Millisecond)
	err3 := c.Send("Совсем.")
	return errors.Join(err1, err2, err3)
}

func updateActivity(db *store.DB, c tele.Context) error {
	amount := 0
	msg := c.Message()
	if msg == nil {
		return nil
	}

	if msg.Text != "" || msg.Dice != nil {
		amount++
	}
	if msg.Voice != nil || msg.Poll != nil {
		amount += 25
	}
	if msg.VideoNote != nil {
		amount += 50
	}
	if msg.Audio != nil || msg.Photo != nil || msg.Video != nil {
		amount += 10
	}
	if msg.Animation != nil || msg.Document != nil || msg.Contact != nil || msg.Location != nil || msg.Venue != nil {
		amount += 5
	}
	if msg.Sticker != nil {
		amount--
	}

	if *Debug {
		log.Printf("update activity: id=%v, amount = %v", c.Sender().ID, amount)
	}

	_, err := db.Exec(`
		UPDATE chat_members
		SET activity = activity + ?
		WHERE id = ?
	`, amount, c.Sender().ID)
	return err
}

func resetActicity(db *store.DB) error {
	_, err := db.Exec(`UPDATE chat_members SET activity = 0`)
	return err
}

func becameModerator(bot *tele.Bot, db *store.DB, id int64) {
	link, err := bot.CreateInviteLink(
		tele.ChatID(cfg.ModeratorChatId),
		&tele.ChatInviteLink{MemberLimit: 1, JoinRequest: true},
	)
	if err != nil {
		log.Printf("failed create invite link to moderator chat: %v", err)
		bot.Send(
			tele.ChatID(cfg.VoteChatId),
			"На меня напали пчелы, поэтому я не смог поздравить победителя. Пусть победит еще раз.",
		)
		return
	}

	_, err = db.Exec(`UPDATE users SET invited = 2 WHERE id = ?`, id)
	if err != nil {
		log.Printf("failed update moderator status: %v", err)
		bot.Send(
			tele.ChatID(cfg.VoteChatId),
			"На меня напали пчелы, поэтому я не смог поздравить победителя. Пусть победит еще раз.",
		)
	}

	_, err = bot.Send(
		tele.ChatID(id),
		"Добро пожаловать в команду!\n"+link.InviteLink,
	)
	if err != nil {
		log.Printf("failed send moderator invite: %v", err)
	}
}

func votingEpochEnd(bot *tele.Bot, db *store.DB) error {
	candidates, err := voting.EndVoting(db)
	if err != nil || len(candidates) == 0 {
		return err
	}

	sb := strings.Builder{}
	sb.WriteString("Голосование окончено! Результаты:\n")
	for i, c := range candidates {
		fmt.Fprintf(&sb, "%d. %s --- %d\n", i+1, c.Header(), c.Votes)
	}
	sb.WriteByte('\n')

	validateResult := voting.Validate(candidates)
	if validateResult == voting.ErrAlmostEqual {
		sb.WriteString("Вы оказались слишком хороши и среди вас не удалось выбрать лучшего. ")
		sb.WriteString("Поэтому в этом голосовании никто не победил)")
	} else if validateResult == voting.ErrInvalidVotes {
		sb.WriteString("Результаты кажутся нам мошенническими, поэтому мы аннулируем голосование. ")
		sb.WriteString("Ждем всех в следующий раз!")
	} else {
		sb.WriteString("У нас есть победитель! Это " + candidates[0].Header())
		go becameModerator(bot, db, candidates[0].ID)
	}

	_, err = bot.Send(
		tele.ChatID(cfg.VoteChatId),
		sb.String(),
	)

	return err
}

func votingEpoch(bot *tele.Bot, db *store.DB) error {
	candidates, err := voting.StartVoting(db)
	fewCandidates := false
	if err == voting.ErrTooFewCandidates {
		err = nil
		fewCandidates = true
	}
	if err != nil {
		return err
	}
	if err := resetActicity(db); err != nil {
		log.Printf("failed reset activity: %v", err)
	}

	if candidates == nil || fewCandidates {
		sb := strings.Builder{}
		sb.WriteString("К сожалению, у нас нет достаточно активных кандидатов и выборов не будет. В следующий раз будьте активнее!")
		fmt.Fprintf(&sb, "\n\nКоличество активных кандидатов с последних выборов: %d/%d", len(candidates), *voting.MinCandidates)
		_, err = bot.Send(
			tele.ChatID(cfg.VoteChatId),
			sb.String(),
		)
		return err
	}

	sb := strings.Builder{}
	rows := make([][]tele.InlineButton, len(candidates))

	sb.WriteString("Наступает голосование!\n")
	sb.WriteString("За последнее время лучше всего себя проявили следующие участники:\n")
	for i, c := range candidates {
		sb.WriteString("- ")
		sb.WriteString(c.Header())
		sb.WriteByte('\n')

		data := strconv.AppendInt([]byte("vote"), c.ID, 10)
		rows[i] = []tele.InlineButton{{Text: c.Header(), Data: string(data)}}
	}
	sb.WriteString("\nГолосуем среди них! И пусть голосование будет честным и недолгим")

	_, err = bot.Send(
		tele.ChatID(cfg.VoteChatId),
		sb.String(),
		&tele.ReplyMarkup{InlineKeyboard: rows},
	)
	if err != nil {
		return err
	}

	time.Sleep(*VotingDuration)

	return votingEpochEnd(bot, db)
}

func initBot(bot *tele.Bot, db *store.DB) {
	startMenu := &tele.ReplyMarkup{}
	startQuizButton := startMenu.WebApp("Давай", &tele.WebApp{URL: cfg.RootURL + "/quiz"})
	noQuizButton := startMenu.Data("Нет", "no")
	startMenu.Inline(
		startMenu.Row(startQuizButton),
		startMenu.Row(noQuizButton),
	)

	bot.SetMyDescription("Я милый и умный HR из TMK Network. Напишите мне, поверьте, вам надо.", "")
	bot.SetMyShortDescription("Свяжитесь со мной, если хотите присоединиться к TMK Network.", "")
	bot.SetMyName("TMK Network HR", "")

	bot.Send(
		tele.ChatID(cfg.VoteChatId),
		"Доброе утро, дорогие коллеги!",
	)

	bot.Use(func(inner tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			id := c.Sender().ID
			if banned, err := store.IsBanned(db, id); err != nil {
				log.Printf("failed check banned status: %v", err)
				return inner(c)
			} else if !banned {
				return inner(c)
			} else {
				return nil
			}
		}
	})

	bot.Handle("/start", func(c tele.Context) error {
		user := c.Sender()

		if *Debug {
			log.Printf("registering: %v", user.ID)
		}

		_, err := db.Exec(
			`INSERT INTO users (id, name, username) VALUES (?, ?, ?) ON CONFLICT DO NOTHING`,
			user.ID,
			user.FirstName,
			user.Username,
		)
		if err != nil {
			return errors.Join(
				err,
				c.Send("Увы, у меня не получилось вас зарегистрировать. Наверное, так быть не должно, попробуйте в другой фазе луны"),
			)
		}
		return c.Send(
			"Привет. Рад, что ты заинтересовался)\nДавай сначала выполним небольшое тестовое задание?",
			startMenu,
		)
	})

	bot.Handle(&noQuizButton, func(c tele.Context) error {
		bot.EditReplyMarkup(c.Message(), nil)
		return banUser(db, c)
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		if *Debug {
			log.Printf("text message from: chat = %v, user = %v", c.Chat().ID, c.Sender().ID)
		}

		if c.Chat().ID == c.Sender().ID {
			return c.Send("Я даже не знаю, что с тобой делать")
		}
		if c.Chat().ID == cfg.VoteChatId {
			return updateActivity(db, c)
		}
		return nil
	})

	activities := []string{
		tele.OnSticker,
		tele.OnVideo,
		tele.OnAudio,
		tele.OnPhoto,
		tele.OnVoice,
		tele.OnVideoNote,
		tele.OnAnimation,
		tele.OnDocument,
		tele.OnPoll,
		tele.OnLocation,
		tele.OnVenue,
		tele.OnDice,
		tele.OnContact,
	}
	activityHandler := func(c tele.Context) error {
		if c.Chat().ID == cfg.VoteChatId {
			return updateActivity(db, c)
		}
		return nil
	}
	for _, activity := range activities {
		bot.Handle(activity, activityHandler)
	}

	bot.Handle(tele.OnChatMember, func(c tele.Context) error {
		if c.Chat().ID != cfg.VoteChatId {
			return nil
		}

		member := c.ChatMember().NewChatMember
		switch member.Role {
		case tele.Creator:
			fallthrough
		case tele.Administrator:
			fallthrough
		case tele.Member:
			if *Debug {
				log.Printf("user join vote chat: %v (%v)", member.User.ID, member.Role)
			}
			_, err := db.Exec(
				`INSERT INTO chat_members (id) VALUES (?) ON CONFLICT DO NOTHING`,
				member.User.ID,
			)
			return err
		case tele.Left:
			fallthrough
		case tele.Kicked:
			if *Debug {
				log.Printf("user left vote chat: %v", member.User.ID)
			}
			_, err := db.Exec(
				`DELETE FROM chat_members WHERE id = ?`,
				member.User.ID,
			)
			return err
		default:
			return nil
		}
	})

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		cb := c.Callback()
		if *Debug {
			log.Printf("callback query: %v (data = %q)", cb, cb.Data)
		}
		if idString, ok := strings.CutPrefix(cb.Data, "vote"); ok {
			id, err := strconv.ParseInt(idString, 10, 64)
			if err != nil {
				log.Printf("invalid callback data")
				return nil
			}

			_, err = db.Exec(
				`UPDATE chat_members SET vote = ? WHERE id = ?`,
				id, c.Sender().ID,
			)
			if err != nil {
				log.Printf("failed save vote: %v", err)
				return err
			}

			return c.RespondText("Г-о-л-о-о-с принят")
		} else {
			return c.Respond()
		}
	})

	bot.Handle(tele.OnChatJoinRequest, func(c tele.Context) error {
		joinRequest := c.ChatJoinRequest()
		if joinRequest.Chat.ID != cfg.ModeratorChatId {
			return c.Bot().DeclineJoinRequest(
				joinRequest.Chat,
				joinRequest.Sender,
			)
		}

		row := db.QueryRow(`SELECT invited FROM users WHERE id = ?`, joinRequest.Sender.ID)
		var invited int
		if err := row.Scan(&invited); err != nil || invited != 2 {
			return errors.Join(
				err,
				c.Bot().DeclineJoinRequest(joinRequest.Chat, joinRequest.Sender),
			)
		} else {
			return c.Bot().ApproveJoinRequest(joinRequest.Chat, joinRequest.Sender)
		}
	})
}

func quizDoneCallback(bot *tele.Bot, db *store.DB) func(id int64) {
	return func(id int64) {
		row := db.QueryRow(`SELECT invited FROM users WHERE id = ?`, id)
		var invited int
		if err := row.Scan(&invited); err != nil {
			bot.OnError(err, nil)
			return
		}
		if invited > 0 {
			return
		}

		if *Debug {
			log.Printf("quiz done: %v", id)
		}

		link, err := bot.CreateInviteLink(
			tele.ChatID(cfg.VoteChatId),
			&tele.ChatInviteLink{MemberLimit: 1},
		)
		if err != nil {
			bot.OnError(err, nil)
			return
		}

		_, err = bot.Send(
			tele.ChatID(id),
			strings.Join([]string{
				"Ты справился с первым этапом. Теперь переходим ко второму)",
				"В нем тебе предстоит доказать, что ты лучше других кандидатов.",
				"Ты попадешь в чат, в котором время от времени будут происходить выборы.",
				"Если выиграешь, то станешь новым членом команды Teemoorka Network!",
				"",
				"Ссылка на чат: " + link.InviteLink,
			}, "\n"),
		)
		if err != nil {
			bot.OnError(err, nil)
			return
		}

		_, err = db.Exec(`UPDATE users SET invited = 1 WHERE id = ?`, id)
		if err != nil {
			bot.OnError(err, nil)
		}
	}
}
