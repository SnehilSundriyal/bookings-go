package handlers

import (
	"encoding/json"
	"errors"
	"github.com/SnehilSundriyal/bookings-go/internal/config"
	"github.com/SnehilSundriyal/bookings-go/internal/driver"
	"github.com/SnehilSundriyal/bookings-go/internal/forms"
	"github.com/SnehilSundriyal/bookings-go/internal/helpers"
	"github.com/SnehilSundriyal/bookings-go/internal/models"
	"github.com/SnehilSundriyal/bookings-go/internal/render"
	"github.com/SnehilSundriyal/bookings-go/internal/repository"
	"github.com/SnehilSundriyal/bookings-go/internal/repository/dbrepo"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from the session"))
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	err = render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		return
	}
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from the session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		err := render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		if err != nil {
			return
		}
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// PostAvailability handles
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	layout := "2006-01-02"

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Available Rooms")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    0,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	err = render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

	if err != nil {
		return
	}
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
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
	err := render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// ReservationSummary displays the res summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	err := render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		return
	}
}

// ChooseRoom displays available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL parameters, builds a session variable and takes user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	startDate := r.URL.Query().Get("s")
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate := r.URL.Query().Get("e")
	if err != nil {
		helpers.ServerError(w, err)
	}

	var res models.Reservation
	layout := "2006-01-02"

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID
	res.StartDate, err = time.Parse(layout, startDate)
	if err != nil {
		helpers.ServerError(w, err)
	}
	res.EndDate, err = time.Parse(layout, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}
