package models_test

import (
	"album/app/models"
	"album/pkg/tester"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

func (suite *AlbumTestSuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	models.DB = mockGormDB
	return mock
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

func (suite *AlbumTestSuite) TestAnniversary() {
	mockedClock := tester.NewMockClock(Str2time("2022-04-01"))

	album := models.Album{ReleaseDate: Str2time("2022-04-01")}
	suite.Assert().Equal(0, album.Anniversary(mockedClock))
	album = models.Album{ReleaseDate: Str2time("2021-04-02")}
	suite.Assert().Equal(0, album.Anniversary(mockedClock))
	album = models.Album{ReleaseDate: Str2time("2021-04-01")}
	suite.Assert().Equal(1, album.Anniversary(mockedClock))
	album = models.Album{ReleaseDate: Str2time("2020-04-02")}
	suite.Assert().Equal(1, album.Anniversary(mockedClock))
	album = models.Album{ReleaseDate: Str2time("2020-04-01")}
	suite.Assert().Equal(2, album.Anniversary(mockedClock))
}

func (suite *AlbumTestSuite) TestAlbumCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`name` = ? ORDER BY `categories`.`id` LIMIT ?")).WithArgs("sports", 1).WillReturnError((errors.New("create error")))
	createdAlbum, err := models.CreateAlbum("Test", Str2time("2023-01-01"), "sports")
	suite.Assert().Nil(createdAlbum)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *AlbumTestSuite) TestAlbumGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `albums` WHERE `albums`.`id` = ? ORDER BY `albums`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("get error"))

	album, err := models.GetAlbum(1)
	suite.Assert().Nil(album)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *AlbumTestSuite) TestAlbumSaveFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`name` = ? ORDER BY `categories`.`id` LIMIT ?")).WithArgs("sports", 1).WillReturnError(errors.New("save error"))

	album := models.Album{
		Title:       "updated",
		ReleaseDate: Str2time("2023-01-01"),
		Category:    &models.Category{Name: "sports"},
	}

	err := album.Save()
	suite.Assert().NotNil(err)
	suite.Assert().Equal("save error", err.Error())
}

func (suite *AlbumTestSuite) TestAlbumDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec("DELETE FROM `albums` WHERE id = ?").WithArgs(0).WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()
	mockDB.ExpectCommit()

	album := models.Album{
		Title:       "Test",
		ReleaseDate: Str2time("2023-01-01"),
		Category:    &models.Category{Name: "sports"},
	}

	err := album.Delete()
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
