package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/oxks/mbfl/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
