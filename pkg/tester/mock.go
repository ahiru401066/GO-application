package tester

import (
	"album/pkg/logger"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MockDB() (mock sqlmock.Sqlmock, mockGormDB *gorm.DB) {
	mockDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		logger.Fatal(err.Error())
	}

	mockGormDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "mock_db",
		DriverName:                "mysql",
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	return mock, mockGormDB
}

type mockClock struct {
	t time.Time
}

func NewMockClock(t time.Time) mockClock {
	return mockClock{t}
}

func (m mockClock) Now() time.Time {
	return m.t
}
