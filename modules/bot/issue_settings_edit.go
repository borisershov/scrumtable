package tgbot

import (
	"fmt"
	"strings"

	tg "github.com/nixys/nxs-go-telegram"
)

func issueSettingsEditState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	issueID, e, err := sess.SlotGet("issueID")
	if err != nil {
		return tg.StateHandlerRes{}, err
	}
	if e == false {
		return tg.StateHandlerRes{
			NextState:    tg.SessState("schedule"),
			StickMessage: true,
		}, nil
	}

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in issueSettingsEdit state handler")
	}

	issue, err := bCtx.m.IssueGetByID(int64(issueID.(float64)), sess.UserIDGet())
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message:      fmt.Sprintf("Enter new issue text\n\nCurrent text:\n`%s`", issue.Text),
		StickMessage: true,
	}, nil
}

func issueSettingsEditMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

	issueID, e, err := sess.SlotGet("issueID")
	if err != nil {
		return tg.MessageHandlerRes{}, err
	}
	if e == false {
		return tg.MessageHandlerRes{
			NextState: tg.SessState("schedule"),
		}, nil
	}

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.MessageHandlerRes{}, fmt.Errorf("can not extract user context in issueSettingsEdit message handler")
	}

	if err := bCtx.m.IssueUpdateText(int64(issueID.(float64)), sess.UserIDGet(), strings.Join(sess.UpdateChain().MessageTextGet(), "; ")); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	if err := sess.SlotDel("issueID"); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	return tg.MessageHandlerRes{
		NextState: tg.SessState("schedule"),
	}, nil
}
