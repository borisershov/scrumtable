package mysql

import "github.com/nixys/scrumtable/misc"

// Issue contains issue data
type Issue struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	TlgrmChatID int64  `json:"tlgrm_chat_id"`
	CreatedAt   string `json:"created_at"`
	Date        string `json:"date"`
	Done        bool   `json:"done"`
	Text        string `json:"text"`
}

// IssueCreateData contains data to create new issue
type IssueCreateData struct {
	TlgrmChatID int64  `json:"tlgrm_chat_id"`
	CreatedAt   string `json:"created_at"`
	Date        string `json:"date"`
	Done        bool   `json:"done"`
	Text        string `json:"text"`
}

// IssueUpdateData contains data to update issue
type IssueUpdateData struct {
	ID          int64   `json:"id" gorm:"->"`
	TlgrmChatID int64   `json:"tlgrm_chat_id" gorm:"->"`
	CreatedAt   *string `json:"created_at"`
	Date        *string `json:"date"`
	Done        *bool   `json:"done"`
	Text        *string `json:"text"`
}

func (Issue) TableName() string {
	return "issues"
}

func (IssueCreateData) TableName() string {
	return "issues"
}

func (IssueUpdateData) TableName() string {
	return "issues"
}

// IssueCreate creates new issue
func (m *MySQL) IssueCreate(issue IssueCreateData) (Issue, error) {

	i := Issue{
		TlgrmChatID: issue.TlgrmChatID,
		CreatedAt:   issue.CreatedAt,
		Date:        issue.Date,
		Done:        issue.Done,
		Text:        issue.Text,
	}

	r := m.client.
		Create(&i)
	if r.Error != nil {
		return Issue{}, r.Error
	}

	return i, nil
}

// IssuesGetByDate gets all issues for specified user and date
func (m *MySQL) IssuesGetByDate(tID int64, date string) ([]Issue, error) {

	var i []Issue

	r := m.client.
		Where(
			m.client.
				Where(Issue{Date: date}).
				Or(Issue{CreatedAt: date}),
		).
		Where(Issue{TlgrmChatID: tID}).
		Find(&i)
	if r.Error != nil {
		return []Issue{}, r.Error
	}

	return i, nil
}

// IssueGetByID gets issue by specified ID and user
func (m *MySQL) IssueGetByID(id, tID int64) (Issue, error) {

	var i Issue

	r := m.client.
		Where(Issue{ID: id}).
		Where(Issue{TlgrmChatID: tID}).
		Find(&i)
	if r.Error != nil {
		return Issue{}, r.Error
	}

	return i, nil
}

// IssueUpdate updates issue by specified ID and user
func (m *MySQL) IssueUpdate(issue IssueUpdateData) (Issue, error) {

	var i Issue

	r := m.client.
		Where(Issue{TlgrmChatID: issue.TlgrmChatID}).
		Updates(issue).Find(&i)
	if r.Error != nil {
		return Issue{}, r.Error
	}

	if r.RowsAffected == 0 {
		return Issue{}, misc.ErrNotFound
	}

	return i, nil
}

// IssueDel deletes specified issue
func (m *MySQL) IssueDel(id, tID int64) error {

	r := m.client.
		Where(Issue{ID: id}).
		Where(Issue{TlgrmChatID: tID}).
		Delete(Issue{})
	if r.Error != nil {
		return r.Error
	}

	return nil
}
