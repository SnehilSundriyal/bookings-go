package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/SnehilSundriyal/bookings-go/internal/config"
	"github.com/SnehilSundriyal/bookings-go/internal/driver"
	"github.com/SnehilSundriyal/bookings-go/internal/forms"
	"github.com/SnehilSundriyal/bookings-go/internal/helpers"
	"github.com/SnehilSundriyal/bookings-go/internal/models"
	"github.com/SnehilSundriyal/bookings-go/internal/render"
	"github.com/SnehilSundriyal/bookings-go/internal/repository"
	"github.com/SnehilSundriyal/bookings-go/internal/repository/dbrepo"
	"log"
	"net/http"
)

var minLength int

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo is the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(a, db.SQL),
	}
}

// New sets the repository for the handlers
func New(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// Reservation renders the make a reservation page and displays the form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	err := render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
	if err != nil {
		return
	}
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	minLength = 3
	form := forms.New(r.PostForm)
	form.Required("first-name", "last-name", "email", "phone")
	form.MinLength("first-name", minLength)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		err := render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		if err != nil {
			return
		}

		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	_, err := w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
	if err != nil {
		return
	}
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		return
	}
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("can't get error from session")
		log.Println("cannot get item from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	err := render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		return
	}
}
