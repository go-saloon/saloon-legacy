package grifts

import (
	"github.com/go-saloon/saloon/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
