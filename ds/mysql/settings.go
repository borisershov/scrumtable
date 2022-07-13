package mysql

import "github.com/nixys/scrumtable/misc"

type Settings struct {
	TlgrmChatID int64  `json:"tlgrm_chat_id" gorm:"primaryKey"`
	CurrentDate string `json:"current_date"`
}

type SettingsSetData struct {
	TlgrmChatID int64  `json:"tlgrm_chat_id" gorm:"primaryKey"`
	CurrentDate string `json:"current_date"`
}

func (Settings) TableName() string {
	return "settings"
}

func (SettingsSetData) TableName() string {
	return "settings"
}

// SettingsSet adds new settings record into DB or update if exist
func (m *MySQL) SettingsSet(s SettingsSetData) (Settings, error) {

	r := m.client.
		Assign(SettingsSetData{CurrentDate: s.CurrentDate}).
		FirstOrCreate(&s)
	if r.Error != nil {
		return Settings{}, r.Error
	}

	return Settings{
		TlgrmChatID: s.TlgrmChatID,
		CurrentDate: s.CurrentDate,
	}, nil
}

// SettingsGet gets settings for specified user
func (m *MySQL) SettingsGet(tID int64) (Settings, error) {

	s := Settings{
		TlgrmChatID: tID,
	}

	r := m.client.
		Find(&s)
	if r.Error != nil {
		return Settings{}, r.Error
	}

	if r.RowsAffected == 0 {
		return Settings{}, misc.ErrNotFound
	}

	return s, nil
}
