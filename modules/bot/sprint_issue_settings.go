package tgbot

import (
	"fmt"

	tg "github.com/nixys/nxs-go-telegram"
)

func sprintIssueSettingsState(t *tg.Telegram) (tg.StateHandlerRes, error) {

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.StateHandlerRes{}, fmt.Errorf("can not extract user context in sprintIssueSettings state handler")
	}

	sprintDate, e, err := t.SlotGet("sprint")
	if err != nil {
		return tg.StateHandlerRes{}, err
	}
	if e == false {
		return tg.StateHandlerRes{
			NextState:    tg.SessState("sprint"),
			StickMessage: true,
		}, nil
	}

	sprintIssueID, e, err := t.SlotGet("sprintIssueID")
	if err != nil {
		return tg.StateHandlerRes{}, err
	}
	if e == false {
		return tg.StateHandlerRes{
			NextState:    tg.SessState("sprint"),
			StickMessage: true,
		}, nil
	}

	issue, err := bCtx.m.SprintIssueGetByID(int64(sprintIssueID.(float64)), t.UserIDGet())
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message: fmt.Sprintf("Issue: `%s`\nSprint: %s", issue.Text, sprintDate.(string)),
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

func sprintIssueSettingsCallback(t *tg.Telegram, uc tg.UpdateChain, identifier string) (tg.CallbackHandlerRes, error) {

	bCtx, b := t.UsrCtxGet().(botCtx)
	if b == false {
		return tg.CallbackHandlerRes{}, fmt.Errorf("can not extract user context in sprintIssueSettings callback handler")
	}

	sprintIssueID, e, err := t.SlotGet("sprintIssueID")
	if err != nil {
		return tg.CallbackHandlerRes{}, err
	}
	if e == false {
		return tg.CallbackHandlerRes{
			NextState: tg.SessState("sprint"),
		}, nil
	}

	id := int64(sprintIssueID.(float64))

	switch identifier {
	case "done":

		sprintIssue, err := bCtx.m.SprintIssueGetByID(id, t.UserIDGet())
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		done := true
		if sprintIssue.Done == true {
			done = false
		}

		if err := bCtx.m.SprintIssueSetDone(id, t.UserIDGet(), done); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

	case "goal":

		sprintDate, e, err := t.SlotGet("sprint")
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}
		if e == false {
			return tg.CallbackHandlerRes{
				NextState: tg.SessState("sprint"),
			}, nil
		}

		sprintIssues, err := bCtx.m.SprintIssuesGetByDate(t.UserIDGet(), sprintDate.(string))
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		// Find previous sprint issue set as goal
		var idToUnGoal int64
		for _, si := range sprintIssues {
			if si.Goal == true && si.ID != id {
				idToUnGoal = si.ID
				break
			}
		}

		// Set current issue as sprint goal
		if err := bCtx.m.SprintIssueSetGoal(id, t.UserIDGet(), true); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		// Remove `goal` mark from previous sprint issue
		if idToUnGoal != 0 {
			if err := bCtx.m.SprintIssueSetGoal(idToUnGoal, t.UserIDGet(), false); err != nil {
				return tg.CallbackHandlerRes{}, err
			}
		}

	case "edit":
		return tg.CallbackHandlerRes{
			NextState: tg.SessState("sprintIssueSettingsEdit"),
		}, nil
	case "del":
		if err := bCtx.m.SprintIssueDel(id, t.UserIDGet()); err != nil {
			return tg.CallbackHandlerRes{}, err
		}
	}

	if err := t.SlotDel("sprintIssueID"); err != nil {
		return tg.CallbackHandlerRes{}, err
	}

	return tg.CallbackHandlerRes{
		NextState: tg.SessState("sprint"),
	}, nil
}
