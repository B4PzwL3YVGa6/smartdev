package api

import (
	"net/http"

	"github.com/anyandrea/smartdev/lib/web"
)

func Error(rw http.ResponseWriter, err error) {
	web.Render().JSON(rw, http.StatusInternalServerError, err.Error())
}
