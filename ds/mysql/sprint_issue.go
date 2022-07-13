package mysql

import "github.com/nixys/scrumtable/misc"

// Issue contains sprint issue data
type SprintIssue struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	TlgrmChatID int64  `json:"tlgrm_chat_id"`
	Date        string `json:"date"`
	Goal        bool   `json:"goal"`
	Done        bool   `json:"done"`
	Text        string `json:"text"`
}

// IssueCreateData contains data to create new sprint issue
type SprintIssueCreateData struct {
	TlgrmChatID int64  `json:"tlgrm_chat_id"`
	Date        string `json:"date"`
	Goal        bool   `json:"goal"`
	Done        bool   `json:"done"`
	Text        string `json:"text"`
}

// SprintIssueUpdateData contains data to update sprint issue
type SprintIssueUpdateData struct {
	ID          int64   `json:"id" gorm:"->"`
	TlgrmChatID int64   `json:"tlgrm_chat_id" gorm:"->"`
	Date        *string `json:"date"`
	Goal        *bool   `json:"goal"`
	Done        *bool   `json:"done"`
	Text        *string `json:"text"`
}

func (SprintIssue) TableName() string {
	return "sprint_issues"
}

func (SprintIssueCreateData) TableName() string {
	return "sprint_issues"
}

func (SprintIssueUpdateData) TableName() string {
	return "sprint_issues"
}

// SprintIssueCreate creates new sprint issue
func (m *MySQL) SprintIssueCreate(si SprintIssueCreateData) (SprintIssue, error) {

	i := SprintIssue{
		TlgrmChatID: si.TlgrmChatID,
		Date:        si.Date,
		Goal:        si.Goal,
		Done:        si.Done,
		Text:        si.Text,
	}

	r := m.client.
		Create(&i)
	if r.Error != nil {
		return SprintIssue{}, r.Error
	}

	return i, nil
}

// SprintIssuesGetByDate gets all sprint issues for specified user and date
func (m *MySQL) SprintIssuesGetByDate(tID int64, date string) ([]SprintIssue, error) {

	var i []SprintIssue

	r := m.client.
		Where(SprintIssue{Date: date}).
		Where(SprintIssue{TlgrmChatID: tID}).
		Find(&i)
	if r.Error != nil {
		return []SprintIssue{}, r.Error
	}

	return i, nil
}

// SprintIssueGetByID gets sprint issue by specified ID and user
func (m *MySQL) SprintIssueGetByID(id, tID int64) (SprintIssue, error) {

	var i SprintIssue

	r := m.client.
		Where(SprintIssue{TlgrmChatID: tID}).
		Where(SprintIssue{ID: id}).
		Find(&i)
	if r.Error != nil {
		return SprintIssue{}, r.Error
	}

	return i, nil
}

// IssueUpdate updates issue by specified ID and user
func (m *MySQL) SprintIssueUpdate(si SprintIssueUpdateData) (SprintIssue, error) {

	var i SprintIssue

	r := m.client.
		Where(SprintIssue{TlgrmChatID: si.TlgrmChatID}).
		Updates(si).Find(&i)
	if r.Error != nil {
		return SprintIssue{}, r.Error
	}

	if r.RowsAffected == 0 {
		return SprintIssue{}, misc.ErrNotFound
	}

	return i, nil
}

// SprintIssueDel deletes specified sprint issue
func (m *MySQL) SprintIssueDel(id, tID int64) error {

	r := m.client.
		Where(SprintIssue{ID: id}).
		Where(SprintIssue{TlgrmChatID: tID}).
		Delete(SprintIssue{})
	if r.Error != nil {
		return r.Error
	}

	return nil
}
