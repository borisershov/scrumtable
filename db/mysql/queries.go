package mysql

import "database/sql"

// Issue contains issue data
type Issue struct {
	ID        int64
	UID       int64
	CreatedAt string
	Date      string
	Done      bool
	Text      string
}

// Issue contains sprint issue data
type SprintIssue struct {
	ID   int64
	UID  int64
	Goal bool
	Done bool
	Text string
}

// SettingsGetCurDate gets current date for specified user (by Telegram ID)
func (m *MySQL) SettingsGetCurDate(tID int64) (string, error) {

	type table struct {
		TlgrmChatID sql.NullInt64  `db:"tlgrm_chat_id"`
		CurrentDate sql.NullString `db:"current_date"`
	}

	t := []table{}

	err := m.client.Select(&t, "SELECT `tlgrm_chat_id`, `current_date` FROM `settings` WHERE `tlgrm_chat_id` = ?", tID)
	if err != nil {
		return "", err
	}

	if len(t) == 0 {
		return "", nil
	}

	return t[0].CurrentDate.String, nil
}

// SettingsSetCurDate sets current date for specified user (by Telegram ID)
func (m *MySQL) SettingsSetCurDate(tID int64, curDate string) error {

	date, err := m.SettingsGetCurDate(tID)

	if len(date) == 0 {

		s, err := m.client.Prepare("INSERT INTO `settings` (`tlgrm_chat_id`, `current_date`) VALUES(?, ?)")
		if err != nil {
			return err
		}
		defer s.Close()

		_, err = s.Exec(tID, curDate)

		return err
	}

	s, err := m.client.Prepare("UPDATE `settings` SET `current_date` = ? WHERE `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(curDate, tID)

	return err
}

// IssuesGetByDate gets all issues for specified user and date
func (m *MySQL) IssuesGetByDate(tID int64, date string) ([]Issue, error) {

	var issues []Issue

	type table struct {
		ID          sql.NullInt64  `db:"id"`
		TlgrmChatID sql.NullInt64  `db:"tlgrm_chat_id"`
		CreatedAt   sql.NullString `db:"created_at"`
		Date        sql.NullString `db:"date"`
		Done        sql.NullBool   `db:"done"`
		Text        sql.NullString `db:"text"`
	}

	t := []table{}

	err := m.client.Select(&t, "SELECT `id`, `tlgrm_chat_id`, `created_at`, `date`, `done`, `text` FROM `issues` WHERE `tlgrm_chat_id` = ? AND (`created_at` = ? OR `date` = ?)", tID, date, date)
	if err != nil {
		return []Issue{}, err
	}

	if len(t) == 0 {
		return []Issue{}, nil
	}

	for _, i := range t {
		issues = append(issues, Issue{
			ID:        i.ID.Int64,
			UID:       i.TlgrmChatID.Int64,
			Date:      i.Date.String,
			CreatedAt: i.CreatedAt.String,
			Done:      i.Done.Bool,
			Text:      i.Text.String,
		})
	}

	return issues, nil
}

// IssueGetByID gets issue by specified ID and user
func (m *MySQL) IssueGetByID(id, tID int64) (Issue, error) {

	type table struct {
		ID          sql.NullInt64  `db:"id"`
		TlgrmChatID sql.NullInt64  `db:"tlgrm_chat_id"`
		CreatedAt   sql.NullString `db:"created_at"`
		Date        sql.NullString `db:"date"`
		Done        sql.NullBool   `db:"done"`
		Text        sql.NullString `db:"text"`
	}

	t := []table{}

	err := m.client.Select(&t, "SELECT `id`, `tlgrm_chat_id`, `created_at`, `date`, `done`, `text` FROM `issues` WHERE `id` = ? AND `tlgrm_chat_id` = ?", id, tID)
	if err != nil {
		return Issue{}, err
	}

	if len(t) == 0 {
		return Issue{}, nil
	}

	return Issue{
		ID:        t[0].ID.Int64,
		UID:       t[0].TlgrmChatID.Int64,
		CreatedAt: t[0].CreatedAt.String,
		Date:      t[0].Date.String,
		Done:      t[0].Done.Bool,
		Text:      t[0].Text.String,
	}, nil
}

// IssueAdd adds new issue for user
func (m *MySQL) IssueAdd(tID int64, date string, text string) (Issue, error) {

	s, err := m.client.Prepare("INSERT INTO `issues` (`tlgrm_chat_id`, `created_at`, `date`, `done`, `text`) VALUES(?, ?, ?, 0, ?)")
	if err != nil {
		return Issue{}, err
	}
	defer s.Close()

	r, err := s.Exec(tID, date, date, text)
	if err != nil {
		return Issue{}, err
	}

	l, _ := r.LastInsertId()

	return Issue{
		ID:        l,
		UID:       tID,
		CreatedAt: date,
		Date:      date,
		Done:      false,
		Text:      text,
	}, nil
}

// IssueSetDone sets specified issue as done/undone
func (m *MySQL) IssueSetDone(id, tID int64, done bool) error {

	s, err := m.client.Prepare("UPDATE `issues` SET `done` = ? WHERE `id` = ? AND `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(done, id, tID)

	return err
}

// IssueUpdateDate updates date
func (m *MySQL) IssueUpdateDate(id, tID int64, date string) error {

	s, err := m.client.Prepare("UPDATE `issues` SET `date` = ? WHERE `id` = ? AND `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(date, id, tID)

	return err
}

// IssueUpdateText updates issue text
func (m *MySQL) IssueUpdateText(id, tID int64, text string) error {

	s, err := m.client.Prepare("UPDATE `issues` SET `text` = ? WHERE `id` = ? AND `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(text, id, tID)

	return err
}

// IssueDel deletes issue by ID
func (m *MySQL) IssueDel(id, tID int64) error {

	s, err := m.client.Prepare("DELETE FROM `issues` WHERE `id` = ? AND `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(id, tID)

	return err
}

// SprintIssuesGetByDate gets sprint issues by specified dete and user
func (m *MySQL) SprintIssuesGetByDate(tID int64, date string) ([]SprintIssue, error) {

	var issues []SprintIssue

	type table struct {
		ID          sql.NullInt64  `db:"id"`
		TlgrmChatID sql.NullInt64  `db:"tlgrm_chat_id"`
		Date        sql.NullString `db:"date"`
		Goal        sql.NullBool   `db:"goal"`
		Done        sql.NullBool   `db:"done"`
		Text        sql.NullString `db:"text"`
	}

	t := []table{}

	err := m.client.Select(&t, "SELECT `id`, `tlgrm_chat_id`, `date`, `goal`, `done`, `text` FROM `sprint_issues` WHERE `tlgrm_chat_id` = ? AND `date` = ?", tID, date)
	if err != nil {
		return []SprintIssue{}, err
	}

	if len(t) == 0 {
		return []SprintIssue{}, nil
	}

	for _, i := range t {
		issues = append(issues, SprintIssue{
			ID:   i.ID.Int64,
			UID:  i.TlgrmChatID.Int64,
			Goal: i.Goal.Bool,
			Done: i.Done.Bool,
			Text: i.Text.String,
		})
	}

	return issues, nil
}

// SprintIssueGetByID gets sprint issue by specified ID and user
func (m *MySQL) SprintIssueGetByID(id, tID int64) (SprintIssue, error) {

	type table struct {
		ID          sql.NullInt64  `db:"id"`
		TlgrmChatID sql.NullInt64  `db:"tlgrm_chat_id"`
		Date        sql.NullString `db:"date"`
		Goal        sql.NullBool   `db:"goal"`
		Done        sql.NullBool   `db:"done"`
		Text        sql.NullString `db:"text"`
	}

	t := []table{}

	err := m.client.Select(&t, "SELECT `id`, `tlgrm_chat_id`, `date`, `goal`, `done`, `text` FROM `sprint_issues` WHERE `id` = ? AND `tlgrm_chat_id` = ?", id, tID)
	if err != nil {
		return SprintIssue{}, err
	}

	if len(t) == 0 {
		return SprintIssue{}, nil
	}

	return SprintIssue{
		ID:   t[0].ID.Int64,
		UID:  t[0].TlgrmChatID.Int64,
		Goal: t[0].Goal.Bool,
		Done: t[0].Done.Bool,
		Text: t[0].Text.String,
	}, nil
}

// SprintIssueAdd adds new sprint issue
func (m *MySQL) SprintIssueAdd(tID int64, date string, goal bool, text string) (SprintIssue, error) {

	s, err := m.client.Prepare("INSERT INTO `sprint_issues` (`tlgrm_chat_id`, `date`, `goal`, `done`, `text`) VALUES(?, ?, ?, 0, ?)")
	if err != nil {
		return SprintIssue{}, err
	}
	defer s.Close()

	r, err := s.Exec(tID, date, goal, text)
	if err != nil {
		return SprintIssue{}, err
	}

	l, _ := r.LastInsertId()

	return SprintIssue{
		ID:   l,
		UID:  tID,
		Goal: goal,
		Done: false,
		Text: text,
	}, nil
}

// SprintIssueSetDone sets done/undonw for specified sprint issue
func (m *MySQL) SprintIssueSetDone(id, tID int64, done bool) error {

	s, err := m.client.Prepare("UPDATE `sprint_issues` SET `done` = ? WHERE `id` = ? AND `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(done, id, tID)

	return err
}

// SprintIssueSetGoal sets specified sprint issue as goal
func (m *MySQL) SprintIssueSetGoal(id, tID int64, goal bool) error {

	s, err := m.client.Prepare("UPDATE `sprint_issues` SET `goal` = ? WHERE `id` = ? AND `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(goal, id, tID)

	return err
}

// SprintIssueUpdateText updates text for sprint issue
func (m *MySQL) SprintIssueUpdateText(id, tID int64, text string) error {

	s, err := m.client.Prepare("UPDATE `sprint_issues` SET `text` = ? WHERE `id` = ? AND `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(text, id, tID)

	return err
}

// SprintIssueDel deletes specified sprint issue
func (m *MySQL) SprintIssueDel(id, tID int64) error {

	s, err := m.client.Prepare("DELETE FROM `sprint_issues` WHERE `id` = ? AND `tlgrm_chat_id` = ?")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(id, tID)

	return err
}
