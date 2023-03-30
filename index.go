package default_session

import (
	"github.com/infrago/infra"
	"github.com/infrago/session"
)

func Driver() session.Driver {
	return &defaultDriver{}
}

func init() {
	infra.Register("default", Driver())
}
