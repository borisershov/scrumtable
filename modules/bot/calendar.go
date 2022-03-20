package tgbot

import (
	"fmt"
	"time"

	tg "github.com/nixys/nxs-go-telegram"
)

func calendarCmd(t *tg.Telegram, uc tg.UpdateChain, cmd string, args string) (tg.CommandHandlerRes, error) {
	return tg.CommandHandlerRes{
		NextState: tg.SessState("cal"),
	}, nil
}

func calendarState(t *tg.Telegram) (tg.StateHandlerRes, error) {

	var date time.Time

	c, e, err := t.SlotGet("calDate")
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	if e == true {
		date, err = time.Parse("2006-01-02", c.(string))
		if err != nil {
			return tg.StateHandlerRes{}, err
		}
	} else {

		bCtx, b := t.UsrCtxGet().(botCtx)
		if b == false {
			return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in calendar state handler")
		}

		ud, err := userCurDateGet(t.UserIDGet(), bCtx.m)
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

func calendarCallback(t *tg.Telegram, uc tg.UpdateChain, identifier string) (tg.CallbackHandlerRes, error) {

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

		if err := bCtx.m.SettingsSetCurDate(uc.UserIDGet(), value); err != nil {
			return r, err
		}

		if err := t.SlotDel("calDate"); err != nil {
			return r, err
		}

		r.NextState = tg.SessState("schedule")

	case "sprint":

		if err := t.SlotSave("sprint", value); err != nil {
			return r, err
		}

		r.NextState = tg.SessState("sprint")

	case "month":

		if err := t.SlotSave("calDate", value); err != nil {
			return r, err
		}

		r.NextState = tg.SessState("cal")
	}

	return r, nil
}
