package tgbot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/scrumtable/ds/mysql"
)

func scheduleState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	buttons := [][]tg.Button{}

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in schedule state handler")
	}

	date, err := userCurDateGet(sess.UserIDGet(), bCtx.m)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	issues, err := bCtx.m.IssuesGetByDate(sess.UserIDGet(), date)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	dayPrev := d.Add(-time.Hour * 24)
	dayNext := d.Add(time.Hour * 24)

	y, w := d.ISOWeek()

	buttons = append(buttons, []tg.Button{
		{
			Text:       fmt.Sprintf("⬅️ %s", dayPrev.Format("02.01.2006")),
			Identifier: "date:" + dayPrev.Format("2006-01-02"),
		},
		{
			Text:       "Sprint issues",
			Identifier: "sprint:" + fmt.Sprintf("%d-%d", w, y),
		},
		{
			Text:       fmt.Sprintf("%s ➡️", dayNext.Format("02.01.2006")),
			Identifier: "date:" + dayNext.Format("2006-01-02"),
		},
	})

	for _, i := range issues {

		text := i.Text

		if i.CreatedAt != i.Date {
			if i.CreatedAt == date {
				text = "⏩ " + text
			} else {
				text = "⏪ " + text
			}
		}

		if i.Done == true {
			text = "✅ " + text
		}

		buttons = append(buttons, []tg.Button{
			{
				Text:       text,
				Identifier: "issue:" + strconv.Itoa(int(i.ID)),
			},
		})
	}

	isToday := ""
	if d.Truncate(24*time.Hour).Equal(time.Now().Truncate(24*time.Hour)) == true {
		isToday = " 🟢"
	}

	return tg.StateHandlerRes{
		Message:      fmt.Sprintf("Issues on %s%s\n\nEnter text to create new issue", d.Format("Monday, 02-Jan-06"), isToday),
		Buttons:      buttons,
		StickMessage: true,
	}, nil
}

func scheduleMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.MessageHandlerRes{}, fmt.Errorf("can not extract user context in schedule message handler")
	}

	date, err := userCurDateGet(sess.UserIDGet(), bCtx.m)
	if err != nil {
		return tg.MessageHandlerRes{}, err
	}

	// Create new issue for every message line
	for _, m := range strings.Split(strings.Join(sess.UpdateChain().MessageTextGet(), "\n"), "\n") {
		if _, err := bCtx.m.IssueCreate(mysql.IssueCreateData{
			TlgrmChatID: sess.UserIDGet(),
			Date:        date,
			CreatedAt:   date,
			Done:        false,
			Text:        m,
		}); err != nil {
			return tg.MessageHandlerRes{}, err
		}
	}

	return tg.MessageHandlerRes{
		NextState: tg.SessState("schedule"),
	}, nil
}

func scheduleCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	var r tg.CallbackHandlerRes

	action, value, err := buttonIdentifierParse(identifier)
	if err != nil {
		return r, err
	}

	switch action {
	case "date":

		bCtx, b := t.UsrCtxGet().(botCtx)
		if b == false {
			return r, fmt.Errorf("can not extract user context in schedule callback handler")
		}

		if _, err := bCtx.m.SettingsSet(mysql.SettingsSetData{
			TlgrmChatID: sess.UserIDGet(),
			CurrentDate: value,
		}); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		r.NextState = tg.SessState("schedule")

	case "issue":

		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return r, err
		}

		if err := sess.SlotSave("issueID", id); err != nil {
			return r, err
		}

		r.NextState = tg.SessState("issueSettings")

	case "sprint":

		if err := sess.SlotSave("sprint", value); err != nil {
			return r, err
		}

		r.NextState = tg.SessState("sprint")
	}

	return r, nil
}
