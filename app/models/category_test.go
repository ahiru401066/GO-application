package models_test

import (
	"album/app/models"
	"album/pkg/tester"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CategoryTestSuite struct {
	tester.DBSQLiteSuite
}

func TestCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryTestSuite))
}

func (suite *CategoryTestSuite) TestCategory() {
	category, err := models.GetOrCreateCategory("test")
	suite.Assert().Nil(err)
	suite.Assert().NotNil(category.ID)
	suite.Assert().Equal("test", category.Name)

	category2, err := models.GetOrCreateCategory("test")
	suite.Assert().Nil(err)
	suite.Assert().Equal("test", category2.Name)
	suite.Assert().Equal(category.ID, category2.ID)
}
