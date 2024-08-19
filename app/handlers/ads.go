package handlers

import "github.com/labstack/echo/v5"

type Ads struct {
}

func NewAds() *Ads {
	return &Ads{}
}

func (a *Ads) List(c echo.Context) error {
	return nil
}
