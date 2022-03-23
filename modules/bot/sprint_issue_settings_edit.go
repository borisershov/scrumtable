package tgbot

import (
	"fmt"
	"strings"

	tg "github.com/nixys/nxs-go-telegram"
)

func sprintIssueSettingsEditState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in sprintIssueSettingsEdit state handler")
	}

	sprintIssueID, e, err := sess.SlotGet("sprintIssueID")
	if err != nil {
		return tg.StateHandlerRes{}, err
	}
	if e == false {
		return tg.StateHandlerRes{
			NextState:    tg.SessState("sprint"),
			StickMessage: true,
		}, nil
	}

	sprintIssue, err := bCtx.m.SprintIssueGetByID(int64(sprintIssueID.(float64)), sess.UserIDGet())
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message:      fmt.Sprintf("Enter new sprint issue text \n\nCurrent text:\n`%s`", sprintIssue.Text),
		StickMessage: true,
	}, nil
}

func sprintIssueSettingsEditMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

	sprintIssueID, e, err := sess.SlotGet("sprintIssueID")
	if err != nil {
		return tg.MessageHandlerRes{}, err
	}
	if e == false {
		return tg.MessageHandlerRes{
			NextState: tg.SessState("sprint"),
		}, nil
	}

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.MessageHandlerRes{}, fmt.Errorf("can not extract user context in sprintIssueSettingsEdit message handler")
	}

	if err := bCtx.m.SprintIssueUpdateText(int64(sprintIssueID.(float64)), sess.UserIDGet(), strings.Join(sess.UpdateChain().MessageTextGet(), "; ")); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	if err := sess.SlotDel("sprintIssueID"); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	return tg.MessageHandlerRes{
		NextState: tg.SessState("sprint"),
	}, nil
}
