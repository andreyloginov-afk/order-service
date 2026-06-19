package mzerolog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/andreyloginov-afk/order-service/internal/pkg/http/httph"
)

type middleware struct {
	log zerolog.Logger

	fromOptions struct {
		skipper func(r *http.Request) bool
	}
}

func (m *middleware) Callback(c *gin.Context) {
	start := time.Now()

	c.Next()

	err := httph.ErrorGet(c.Request)

	if m.fromOptions.skipper(c.Request) {
		return
	}

	var msg string
	if err != nil {
		msg = fmt.Sprintf("%s %s finished with error", c.Request.Method, c.Request.RequestURI)
	} else {
		msg = fmt.Sprintf("%s %s finished with no error", c.Request.Method, c.Request.RequestURI)
	}

	event := m.log.Debug()
	if err != nil {
		event = m.log.Error()
	}

	event.
		Err(err).
		Ctx(c.Request.Context()).
		Dur("exec_time", time.Since(start)).
		Str("client_ip", c.ClientIP()).
		Msg(msg)
}

func NewMiddleware(opts ...Option) gin.HandlerFunc {
	m := &middleware{
		log: log.Logger,
	}
	m.fromOptions.skipper = defaultSkipper

	for _, opt := range opts {
		opt(m)
	}

	return m.Callback
}

func defaultSkipper(_ *http.Request) bool {
	return false
}
