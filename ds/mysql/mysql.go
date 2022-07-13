package mysql

import (
	"fmt"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQL it is a MySQL module context structure
type MySQL struct {
	client *gorm.DB
}

// Settings contains settings for MySQL
type Param struct {
	Host     string
	User     string
	Password string
	Database string
}

// Connect connects to MySQL
func Connect(p Param) (MySQL, error) {

	client, err := gorm.Open(gmysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		p.User,
		p.Password,
		p.Host,
		p.Database)), &gorm.Config{})
	if err != nil {
		return MySQL{}, err
	}

	return MySQL{
		client: client,
	}, nil
}

// Close closes MySQL connection
func (m *MySQL) Close() error {
	db, _ := m.client.DB()
	return db.Close()
}
