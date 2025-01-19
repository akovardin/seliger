package handlers

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type Ads struct {
	app *pocketbase.PocketBase
}

func NewAds(app *pocketbase.PocketBase) *Ads {
	return &Ads{
		app: app,
	}
}

type AdResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Img         string `json:"img"`
	Url         string `json:"url"`
	Info        string `json:"info"`
}

func (a *Ads) One(e *core.RequestEvent) error {
	// реализуем взвешанный рандом по посчитанному CPM

	response := AdResponse{}

	return e.JSON(200, response)
}

type AdsResponse struct {
	Ads []AdResponse `json:"ads"`
}

func (a *Ads) List(e *core.RequestEvent) error {
	banners, err := a.app.FindRecordsByFilter(
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
		Ads: []AdResponse{},
	}
	for _, b := range banners {

		img := b.GetString("image")
		base := a.app.App.Settings().Meta.AppURL

		response.Ads = append(response.Ads, AdResponse{
			Name:        b.GetString("name"),
			Description: b.GetString("description"),
			Url:         b.GetString("url"),
			Info:        b.GetString("info"),
			Img:         base + "/api/files/banners/" + b.Id + "/" + img,
		})
	}

	return e.JSON(200, response)
}
