package mzerolog

import (
	"net/http"
	"strings"
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

	var sb strings.Builder
	sb.WriteString(c.Request.Method)
	sb.WriteString(" ")
	sb.WriteString(c.Request.RequestURI)
	if err != nil {
		sb.WriteString(" finished (or aborted) with error")
	} else {
		sb.WriteString(" finished with no error")
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
		Int("http_status_code", httph.ErrorGetStatusCode(c.Request)).
		Msg(sb.String())
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
