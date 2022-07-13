package tgbot

import (
	"fmt"
	"strings"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/scrumtable/ds/mysql"
)

func sprintIssueSettingsEditState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var sprintIssueID int64

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in sprintIssueSettingsEdit state handler")
	}

	e, err := sess.SlotGet("sprintIssueID", &sprintIssueID)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}
	if e == false {
		return tg.StateHandlerRes{
			NextState:    tg.SessState("sprint"),
			StickMessage: true,
		}, nil
	}

	sprintIssue, err := bCtx.m.SprintIssueGetByID(sprintIssueID, sess.UserIDGet())
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message:      fmt.Sprintf("Enter new sprint issue text \n\nCurrent text:\n`%s`", sprintIssue.Text),
		StickMessage: true,
	}, nil
}

func sprintIssueSettingsEditMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

	var sprintIssueID int64

	e, err := sess.SlotGet("sprintIssueID", &sprintIssueID)
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

	text := strings.Join(sess.UpdateChain().MessageTextGet(), "; ")

	if _, err := bCtx.m.SprintIssueUpdate(mysql.SprintIssueUpdateData{
		ID:          sprintIssueID,
		TlgrmChatID: sess.UserIDGet(),
		Text:        &text,
	}); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	if err := sess.SlotDel("sprintIssueID"); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	return tg.MessageHandlerRes{
		NextState: tg.SessState("sprint"),
	}, nil
}
