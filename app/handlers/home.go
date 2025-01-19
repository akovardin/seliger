package handlers

import "github.com/pocketbase/pocketbase/core"

type Home struct {
}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) Home(e *core.RequestEvent) error {
	return nil
}
