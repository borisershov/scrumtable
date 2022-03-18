package tgbot

import (
	"fmt"
	"time"

	tg "github.com/nixys/nxs-go-telegram"
)

func issueSettingsCalState(t *tg.Telegram) (tg.StateHandlerRes, error) {

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
			return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in issueSettingsCal state handler")
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
		Message:      fmt.Sprintf("Select issue due date"),
		Buttons:      calendarRender(date.Year(), int(date.Month()), true),
		StickMessage: true,
	}, nil
}

func issueSettingsCalCallback(t *tg.Telegram, uc tg.UpdateChain, identifier string) (tg.CallbackHandlerRes, error) {

	var r tg.CallbackHandlerRes

	issueID, e, err := t.SlotGet("issueID")
	if err != nil {
		return tg.CallbackHandlerRes{}, err
	}
	if e == false {
		return tg.CallbackHandlerRes{
			NextState: tg.SessState("schedule"),
		}, nil
	}

	action, value, err := buttonIdentifierParse(identifier)
	if err != nil {
		return r, err
	}

	switch action {
	case "date":

		// TODO: need to forbid to set issue date to past

		bCtx, b := t.UsrCtxGet().(botCtx)
		if b == false {
			return r, fmt.Errorf("can not extract user context in issueSettingsCal callback handler")
		}

		if err := bCtx.m.IssueUpdateDate(int64(issueID.(float64)), uc.UserIDGet(), value); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		if err := t.SlotDel("issueID"); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		r.NextState = tg.SessState("schedule")

	case "month":

		if err := t.SlotSave("calDate", value); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		r.NextState = tg.SessState("issueSettingsCal")
	}

	return r, nil
}
