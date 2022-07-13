package tgbot

import (
	"fmt"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/scrumtable/ds/mysql"
)

func sprintIssueSettingsState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var (
		sprintDate    string
		sprintIssueID int64
	)

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in sprintIssueSettings state handler")
	}

	e, err := sess.SlotGet("sprint", &sprintDate)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}
	if e == false {
		return tg.StateHandlerRes{
			NextState:    tg.SessState("sprint"),
			StickMessage: true,
		}, nil
	}

	e, err = sess.SlotGet("sprintIssueID", &sprintIssueID)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}
	if e == false {
		return tg.StateHandlerRes{
			NextState:    tg.SessState("sprint"),
			StickMessage: true,
		}, nil
	}

	issue, err := bCtx.m.SprintIssueGetByID(sprintIssueID, sess.UserIDGet())
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message: fmt.Sprintf("Issue: `%s`\nSprint: %s", issue.Text, sprintDate),
		Buttons: [][]tg.Button{
			{
				{
					Text:       "✅",
					Identifier: "done",
				},
			},
			{
				{
					Text:       "🎯",
					Identifier: "goal",
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

func sprintIssueSettingsCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	var (
		sprintIssueID int64
		sprintDate    string
	)

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.CallbackHandlerRes{}, fmt.Errorf("can not extract user context in sprintIssueSettings callback handler")
	}

	e, err := sess.SlotGet("sprintIssueID", &sprintIssueID)
	if err != nil {
		return tg.CallbackHandlerRes{}, err
	}
	if e == false {
		return tg.CallbackHandlerRes{
			NextState: tg.SessState("sprint"),
		}, nil
	}

	switch identifier {
	case "done":

		sprintIssue, err := bCtx.m.SprintIssueGetByID(sprintIssueID, sess.UserIDGet())
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		done := true
		if sprintIssue.Done == true {
			done = false
		}

		if _, err := bCtx.m.SprintIssueUpdate(mysql.SprintIssueUpdateData{
			ID:          sprintIssueID,
			TlgrmChatID: sess.UserIDGet(),
			Done:        &done,
		}); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

	case "goal":

		e, err := sess.SlotGet("sprint", &sprintDate)
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}
		if e == false {
			return tg.CallbackHandlerRes{
				NextState: tg.SessState("sprint"),
			}, nil
		}

		sprintIssues, err := bCtx.m.SprintIssuesGetByDate(sess.UserIDGet(), sprintDate)
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		// Find previous sprint issue set as goal
		var idToUnGoal int64
		for _, si := range sprintIssues {
			if si.Goal == true && si.ID != sprintIssueID {
				idToUnGoal = si.ID
				break
			}
		}

		g := true
		// Set current issue as sprint goal
		if _, err := bCtx.m.SprintIssueUpdate(mysql.SprintIssueUpdateData{
			ID:          sprintIssueID,
			TlgrmChatID: sess.UserIDGet(),
			Goal:        &g,
		}); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		// Remove `goal` mark from previous sprint issue
		if idToUnGoal != 0 {
			g = false
			if _, err := bCtx.m.SprintIssueUpdate(mysql.SprintIssueUpdateData{
				ID:          idToUnGoal,
				TlgrmChatID: sess.UserIDGet(),
				Goal:        &g,
			}); err != nil {
				return tg.CallbackHandlerRes{}, err
			}
		}

	case "edit":
		return tg.CallbackHandlerRes{
			NextState: tg.SessState("sprintIssueSettingsEdit"),
		}, nil
	case "del":
		if err := bCtx.m.SprintIssueDel(sprintIssueID, sess.UserIDGet()); err != nil {
			return tg.CallbackHandlerRes{}, err
		}
	}

	if err := sess.SlotDel("sprintIssueID"); err != nil {
		return tg.CallbackHandlerRes{}, err
	}

	return tg.CallbackHandlerRes{
		NextState: tg.SessState("sprint"),
	}, nil
}
