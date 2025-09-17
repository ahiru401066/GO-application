package models

import (
	"encoding/json"
	"time"

	"album/api"
	"album/pkg"
)

type Album struct {
	ID          int
	Title       string
	ReleaseDate time.Time
	CategoryID  int
	Category    *Category
}

func (a *Album) Anniversary(clock pkg.Clock) int {
	now := clock.Now()
	years := now.Year() - a.ReleaseDate.Year()
	releaseDay := pkg.GetAdjustedReleaseDay(a.ReleaseDate, now)
	if now.YearDay() < releaseDay {
		years -= 1
	}
	return years
}

func (a *Album) MarshalJSON() ([]byte, error) {
	return json.Marshal(&api.AlbumResponse{
		Id:          a.ID,
		Title:       a.Title,
		Anniversary: a.Anniversary(pkg.RealClock{}),
		ReleaseDate: api.ReleaseDate{Time: a.ReleaseDate},
		Category: api.Category{
			Id:   &a.CategoryID,
			Name: api.CategoryName(a.Category.Name),
		},
	})
}

func CreateAlbum(title string, releaseDate time.Time, categoryName string) (*Album, error) {
	category, err := GetOrCreateCategory(categoryName)
	if err != nil {
		return nil, err
	}

	album := &Album{
		Title:       title,
		ReleaseDate: releaseDate,
		CategoryID:  category.ID,
		Category:    category,
	}

	if err := DB.Create(album).Error; err != nil {
		return nil, err
	}
	return album, nil
}

func GetAlbum(ID int) (*Album, error) {
	var album Album
	if err := DB.Preload("Category").First(&album, ID).Error; err != nil {
		return nil, err
	}
	return &album, nil
}

func (a *Album) Save() error {
	category, err := GetOrCreateCategory(a.Category.Name)
	if err != nil {
		return err
	}
	a.Category = category
	a.CategoryID = category.ID

	if err := DB.Save(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a *Album) Delete() error {
	if err := DB.Where("id = ?", &a.ID).Delete(&a).Error; err != nil {
		return err
	}
	return nil
}
