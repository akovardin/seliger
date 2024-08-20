package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

type Ads struct {
	app *pocketbase.PocketBase
}

func NewAds(app *pocketbase.PocketBase) *Ads {
	return &Ads{
		app: app,
	}
}

type AdsResponse struct {
	Ads []Ad `json:"ads"`
}

type Ad struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Img         string `json:"img"`
	Url         string `json:"url"`
	Info        string `json:"info"`
}

func (a *Ads) List(c echo.Context) error {
	banners, err := a.app.Dao().FindRecordsByFilter(
		"banners",
		"enabled = true",
		"-created",
		100,
		0,
	)

	if err != nil {
		return err
	}

	response := AdsResponse{
		Ads: []Ad{},
	}
	for _, b := range banners {

		img := b.GetString("image")
		base := a.app.App.Settings().Meta.AppUrl

		response.Ads = append(response.Ads, Ad{
			Name:        b.GetString("name"),
			Description: b.GetString("description"),
			Url:         b.GetString("url"),
			Info:        b.GetString("info"),
			Img:         base + "/api/files/banners/" + b.Id + "/" + img,
		})
	}

	c.JSON(200, response)

	return nil
}
