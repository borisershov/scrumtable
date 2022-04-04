package tgbot

import (
	"fmt"
	"time"

	tg "github.com/nixys/nxs-go-telegram"
)

func calendarCmd(t *tg.Telegram, sess *tg.Session, cmd string, args string) (tg.CommandHandlerRes, error) {
	return tg.CommandHandlerRes{
		NextState: tg.SessState("cal"),
	}, nil
}

func calendarCurDateCmd(t *tg.Telegram, sess *tg.Session, cmd string, args string) (tg.CommandHandlerRes, error) {

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.CommandHandlerRes{}, fmt.Errorf("can not extract user context in calendar current date cmd handler")
	}

	d := time.Now().Format("2006-01-02")

	// Set current date for user
	if err := bCtx.m.SettingsSetCurDate(sess.UserIDGet(), d); err != nil {
		return tg.CommandHandlerRes{}, err
	}

	if err := sess.SlotDel("calDate"); err != nil {
		return tg.CommandHandlerRes{}, err
	}

	return tg.CommandHandlerRes{
		NextState: tg.SessState("schedule"),
	}, nil
}

func calendarState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var (
		date time.Time
		c    string
	)

	e, err := sess.SlotGet("calDate", &c)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	if e == true {
		date, err = time.Parse("2006-01-02", c)
		if err != nil {
			return tg.StateHandlerRes{}, err
		}
	} else {

		bCtx, b := t.UsrCtxGet().(botCtx)
		if b == false {
			return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in calendar state handler")
		}

		ud, err := userCurDateGet(sess.UserIDGet(), bCtx.m)
		if err != nil {
			return tg.StateHandlerRes{}, err
		}

		date, err = time.Parse("2006-01-02", ud)
		if err != nil {
			return tg.StateHandlerRes{}, err
		}
	}

	return tg.StateHandlerRes{
		Message:      "Select a day or week to watch already created your issues or make a new one",
		Buttons:      calendarRender(date.Year(), int(date.Month()), true),
		StickMessage: true,
	}, nil
}

func calendarCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	var r tg.CallbackHandlerRes

	action, value, err := buttonIdentifierParse(identifier)
	if err != nil {
		return r, err
	}

	switch action {
	case "date":

		bCtx, b := t.UsrCtxGet().(botCtx)
		if b == false {
			return r, fmt.Errorf("can not extract user context in calendar callback handler")
		}

		if err := bCtx.m.SettingsSetCurDate(sess.UserIDGet(), value); err != nil {
			return r, err
		}

		if err := sess.SlotDel("calDate"); err != nil {
			return r, err
		}

		r.NextState = tg.SessState("schedule")

	case "sprint":

		if err := sess.SlotSave("sprint", value); err != nil {
			return r, err
		}

		r.NextState = tg.SessState("sprint")

	case "month":

		if err := sess.SlotSave("calDate", value); err != nil {
			return r, err
		}

		r.NextState = tg.SessState("cal")
	}

	return r, nil
}
