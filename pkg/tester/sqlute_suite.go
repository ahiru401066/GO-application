package tester

import (
	"album/app/models"
	"album/configs"
	"os"

	"github.com/stretchr/testify/suite"
)

type DBSQLiteSuite struct {
	suite.Suite
}

func (suite *DBSQLiteSuite) SetupSuite() {
	configs.Config.DBName = "unittest.sqlite"
	err := models.SetDatabase(models.InstanceSqlLite)
	suite.Assert().Nil(err)

	for _, model := range models.GetModels() {
		err := models.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

func (suite *DBSQLiteSuite) TearDownSuite() {
	err := os.Remove(configs.Config.DBName)
	suite.Assert().Nil(err)
}
