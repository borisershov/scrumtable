package tgbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nixys/scrumtable/db/mysql"

	tg "github.com/nixys/nxs-go-telegram"
)

type Settings struct {
	MySQL     mysql.MySQL
	APIToken  string
	RedisHost string
}

type Bot struct {
	bot tg.Telegram
}

type botCtx struct {
	m mysql.MySQL
}

func Init(settings Settings) (Bot, error) {

	// Setup the bot
	bot, err := tg.Init(
		tg.Settings{
			BotSettings: tg.SettingsBot{
				BotAPI: settings.APIToken,
			},
			RedisHost: settings.RedisHost,
		},

		tg.Description{

			Commands: []tg.Command{
				{
					Command:     "calendar",
					Description: "Show calendar",
					Handler:     calendarCmd,
				},
				{
					Command:     "curdate",
					Description: "Go to current date",
					Handler:     calendarCurDateCmd,
				},
			},

			InitHandler: botInit,

			States: map[tg.SessionState]tg.State{

				tg.SessState("cal"): {
					StateHandler:    calendarState,
					CallbackHandler: calendarCallback,
				},

				tg.SessState("schedule"): {
					StateHandler:    scheduleState,
					MessageHandler:  scheduleMsg,
					CallbackHandler: scheduleCallback,
				},

				tg.SessState("sprint"): {
					StateHandler:    sprintState,
					MessageHandler:  sprintMsg,
					CallbackHandler: sprintCallback,
				},

				tg.SessState("issueSettings"): {
					StateHandler:    issueSettingsState,
					CallbackHandler: issueSettingsCallback,
				},

				tg.SessState("sprintIssueSettings"): {
					StateHandler:    sprintIssueSettingsState,
					CallbackHandler: sprintIssueSettingsCallback,
				},

				tg.SessState("issueSettingsCal"): {
					StateHandler:    issueSettingsCalState,
					CallbackHandler: issueSettingsCalCallback,
				},

				tg.SessState("issueSettingsEdit"): {
					StateHandler:   issueSettingsEditState,
					MessageHandler: issueSettingsEditMsg,
				},

				tg.SessState("sprintIssueSettingsEdit"): {
					StateHandler:   sprintIssueSettingsEditState,
					MessageHandler: sprintIssueSettingsEditMsg,
				},
			},
		},
		botCtx{
			m: settings.MySQL,
		})
	if err != nil {
		return Bot{}, fmt.Errorf("bot setup error: %v", err)
	}

	return Bot{
		bot: bot,
	}, nil
}

// runtimeBotUpdates checks updates at Telegram and put it into queue
func (b *Bot) UpdatesGet(ctx context.Context, ch chan error) {
	if err := b.bot.GetUpdates(ctx); err != nil {
		if err == tg.ErrUpdatesChanClosed {
			ch <- nil
		} else {
			ch <- err
		}
	} else {
		ch <- nil
	}
}

// runtimeBotQueue processes an updaates from queue
func (b *Bot) Queue(ctx context.Context, ch chan error) {
	timer := time.NewTimer(time.Millisecond * 200)
	for {
		select {
		case <-timer.C:
			if err := b.bot.Processing(); err != nil {
				ch <- err
			}
			timer.Reset(time.Millisecond * 200)
		case <-ctx.Done():
			return
		}
	}
}

func botInit(t *tg.Telegram, sess *tg.Session) (tg.InitHandlerRes, error) {

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.InitHandlerRes{}, fmt.Errorf("can not extract user context in botInit handler")
	}

	if _, err := userCurDateGet(sess.UserIDGet(), bCtx.m); err != nil {
		return tg.InitHandlerRes{}, err
	}

	return tg.InitHandlerRes{
		NextState: tg.SessState("schedule"),
	}, nil
}

// userCurDateGet returns current date for user from settings.
// If current date in settings for user is empty it will
// be set as `Now()`
func userCurDateGet(tID int64, m mysql.MySQL) (string, error) {

	// Get user current date
	d, err := m.SettingsGetCurDate(tID)
	if err != nil {
		return "", err
	}

	if len(d) > 0 {
		return d, nil
	}

	d = time.Now().Format("2006-01-02")

	// Set current date for user
	if err := m.SettingsSetCurDate(tID, d); err != nil {
		return "", err
	}

	return d, nil
}

func calendarRender(year, month int, useSprints bool) [][]tg.Button {

	firstDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	prevDate := time.Date(year, time.Month(month-1), 1, 0, 0, 0, 0, time.UTC)
	nextDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)

	buttons := [][]tg.Button{
		{
			{
				Text:       time.Now().Format("January, 2006"),
				Identifier: "month:" + time.Now().Format("2006-01-02"),
			},
		},
		{
			{
				Text:       "⬅️ " + prevDate.Format("January, 2006"),
				Identifier: "month:" + prevDate.Format("2006-01-02"),
			},
			{
				Text:       nextDate.Format("January, 2006") + " ➡️",
				Identifier: "month:" + nextDate.Format("2006-01-02"),
			},
		},
	}
	bb := []tg.Button{}

	// Calc date of first monday
	for ; firstDate.Weekday() != 1; firstDate = firstDate.Add(-time.Hour * 24) {
	}

	// Calc date of last monday
	lastDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	for ; lastDate.Weekday() != 1; lastDate = lastDate.Add(time.Hour * 24) {
	}

	y, m, d := time.Now().Date()
	currentDate := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	// Render calendar buttons
	for ; firstDate.Equal(lastDate) == false; firstDate = firstDate.Add(time.Hour * 24) {

		// If need to render sprints
		if useSprints == true && firstDate.Weekday() == 1 {

			y, w := firstDate.ISOWeek()

			bb = append(bb, tg.Button{
				Text:       fmt.Sprintf("#%d", w),
				Identifier: "sprint:" + fmt.Sprintf("%d-%d", w, y),
			})
		}

		text := firstDate.Format("02")
		if firstDate.Equal(currentDate) == true {
			text = "🟢"
		}

		bb = append(bb, tg.Button{
			Text:       text,
			Identifier: "date:" + firstDate.Format("2006-01-02"),
		})

		// Check if need to start next week
		if firstDate.Weekday() == 0 {
			buttons = append(buttons, bb)
			bb = []tg.Button{}
		}
	}

	if len(bb) > 0 {
		buttons = append(buttons, bb)
	}

	return buttons
}

func buttonIdentifierParse(identifier string) (string, string, error) {

	b := strings.Split(identifier, ":")

	if len(b) != 2 {
		return "", "", fmt.Errorf("incorrect button identifier")
	}

	return b[0], b[1], nil
}
