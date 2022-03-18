package tgbot

import (
	"fmt"
	"time"

	tg "github.com/nixys/nxs-go-telegram"
)

func issueSettingsState(t *tg.Telegram) (tg.StateHandlerRes, error) {

	var r tg.StateHandlerRes

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in issueSettings state handler")
	}

	issueID, e, err := t.SlotGet("issueID")
	if err != nil {
		return r, err
	}
	if e == false {
		return tg.StateHandlerRes{
			NextState:    tg.SessState("schedule"),
			StickMessage: true,
		}, nil
	}

	issue, err := bCtx.m.IssueGetByID(int64(issueID.(float64)), t.UserIDGet())
	if err != nil {
		return r, err
	}

	createdAtDate, err := time.Parse("2006-01-02", issue.CreatedAt)
	if err != nil {
		return r, err
	}

	dueDate, err := time.Parse("2006-01-02", issue.Date)
	if err != nil {
		return r, err
	}

	msg := ""
	if createdAtDate.Truncate(24*time.Hour).Equal(dueDate.Truncate(24*time.Hour)) == true {
		msg = fmt.Sprintf("Issue: `%s`\nDue date: %s", issue.Text, dueDate.Format("Monday, 02-Jan-06"))
	} else {
		msg = fmt.Sprintf("Issue: `%s`\nStart date: %s\nDue date: %s", issue.Text, createdAtDate.Format("Monday, 02-Jan-06"), dueDate.Format("Monday, 02-Jan-06"))
	}

	return tg.StateHandlerRes{
		Message: msg,
		Buttons: [][]tg.Button{
			{
				{
					Text:       "✅",
					Identifier: "done",
				},
			},
			{
				{
					Text:       "📅",
					Identifier: "calendar",
				},
			},
			{
				{
					Text:       "🖋",
					Identifier: "edit",
				},
			},
			{
				{
					Text:       "🗑",
					Identifier: "del",
				},
			},
			{
				{
					Text:       "🔙",
					Identifier: "back",
				},
			},
		},
		StickMessage: true,
	}, nil
}

func issueSettingsCallback(t *tg.Telegram, uc tg.UpdateChain, identifier string) (tg.CallbackHandlerRes, error) {

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.CallbackHandlerRes{}, fmt.Errorf("can not extract user context in issueSettings callback handler")
	}

	issueID, e, err := t.SlotGet("issueID")
	if err != nil {
		return tg.CallbackHandlerRes{}, err
	}
	if e == false {
		return tg.CallbackHandlerRes{
			NextState: tg.SessState("schedule"),
		}, nil
	}

	id := int64(issueID.(float64))

	switch identifier {
	case "done":

		issue, err := bCtx.m.IssueGetByID(id, t.UserIDGet())
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		done := true
		if issue.Done == true {
			done = false
		}

		if err := bCtx.m.IssueSetDone(id, t.UserIDGet(), done); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

	case "calendar":
		return tg.CallbackHandlerRes{
			NextState: tg.SessState("issueSettingsCal"),
		}, nil
	case "edit":
		return tg.CallbackHandlerRes{
			NextState: tg.SessState("issueSettingsEdit"),
		}, nil
	case "del":
		if err := bCtx.m.IssueDel(id, t.UserIDGet()); err != nil {
			return tg.CallbackHandlerRes{}, err
		}
	}

	if err := t.SlotDel("issueID"); err != nil {
		return tg.CallbackHandlerRes{}, err
	}

	return tg.CallbackHandlerRes{
		NextState: tg.SessState("schedule"),
	}, nil
}
