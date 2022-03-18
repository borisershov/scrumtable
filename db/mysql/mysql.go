package mysql

import (
	"github.com/jmoiron/sqlx"
)

// MySQL it is a MySQL module context structure
type MySQL struct {
	client *sqlx.DB
}

// Settings contains settings for MySQL
type Settings struct {
	Host     string
	User     string
	Password string
	Database string
}

// Connect connects to MySQL
func Connect(s Settings) (MySQL, error) {

	var m MySQL

	client, err := sqlx.Connect("mysql", s.User+":"+s.Password+"@"+"tcp("+s.Host+")/"+s.Database+"?parseTime=true&charset=utf8mb4,utf8")
	if err != nil {
		return m, err
	}

	m.client = client

	return m, nil
}

// Close closes MySQL connection
func (m *MySQL) Close() error {
	return m.client.Close()
}
