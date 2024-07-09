package helpers

import (
	"fmt"
	"github.com/SnehilSundriyal/bookings-go/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// New sets up app config for helpers
func New(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client Error:", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}