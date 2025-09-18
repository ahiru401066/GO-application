package models_test

import (
	"album/app/models"
	"album/pkg/tester"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AlbumTestSuite struct {
	tester.DBSQLiteSuite
	originalDB *gorm.DB
}

func TestAlbumTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumTestSuite))
}

func (suite *AlbumTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.originalDB = models.DB
}

func (suite *AlbumTestSuite) AfterTest(suiteName, testName string) {
	models.DB = suite.originalDB
}

func Str2time(t string) time.Time {
	parsedTime, _ := time.Parse("2006-01-02", t)
	return parsedTime
}

func (suite *AlbumTestSuite) TestAlbum() {
	createdAlbum, err := models.CreateAlbum("Test", time.Now(), "sports")
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", createdAlbum.Title)
	suite.Assert().NotNil(createdAlbum.ReleaseDate)
	suite.Assert().NotNil(createdAlbum.Category.ID)
	suite.Assert().Equal(createdAlbum.Category.Name, "sports")

	getAlbum, err := models.GetAlbum(createdAlbum.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", getAlbum.Title)
	suite.Assert().NotNil(getAlbum.ReleaseDate)
	suite.Assert().NotNil(getAlbum.Category.ID)
	suite.Assert().Equal(getAlbum.Category.Name, "sports")

	getAlbum.Title = "updated"
	err = getAlbum.Save()
	suite.Assert().Nil(err)
	updatedAlbum, err := models.GetAlbum(createdAlbum.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updatedAlbum.Title)
	suite.Assert().NotNil(updatedAlbum.ReleaseDate)
	suite.Assert().NotNil(updatedAlbum.Category.ID)
	suite.Assert().Equal(updatedAlbum.Category.Name, "sports")

	err = updatedAlbum.Delete()
	suite.Assert().Nil(err)
	deletedAlbum, err := models.GetAlbum(updatedAlbum.ID)
	suite.Assert().Nil(deletedAlbum)
	suite.Assert().True(strings.Contains("record not found", err.Error()))
}

func (suite *AlbumTestSuite) TestAlbumMarshal() {
	album := models.Album{
		Title:       "Test",
		ReleaseDate: Str2time("2023-01-01"),
		Category:    &models.Category{Name: "sports"},
	}
	aniversary := time.Now().Year() - 2023
	albumJSON, err := album.MarshalJSON()
	suite.Assert().Nil(err)
	suite.Assert().JSONEq(fmt.Sprintf(`{
         "anniversary":%d,
         "category":{
             "id":0,"name":"sports"
         },
         "id":0,
         "releaseDate":"2023-01-01",
         "title":"Test"
    }`, aniversary), string(albumJSON))
}
