package main

import (
	"encoding/json"
	"fmt"
	"Booking/booking"
	"Booking/middleware"
	"Booking/models"
	"Booking/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"

	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"strconv"
)

type application struct {
	servicePort   string
	faqRepository booking.Booking
	s             *mux.Router
	pr            *prometheus.Prometheus
}

var conf models.Config

func init() {
	models.LoadConfig(&conf)
}

func main() {
	models.LoadConfig(&conf)
	app := NewApplication(conf)
	app.initServer()
	log.Fatal(http.ListenAndServe(app.servicePort, app.s))

}

func NewApplication(conf models.Config) *application {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable",
		conf.SQLDataBase.Server, conf.SQLDataBase.UserID, conf.SQLDataBase.Password, conf.SQLDataBase.Port, conf.SQLDataBase.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	log.Println("Connect with database")
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil
	}
	ords := booking.NewBookingRepository(db)

	return &application{
		servicePort:   ":9000",
		faqRepository: ords,
		pr:            prometheus.New("booking-api"),
	}
}

func (app *application) initServer() {
	app.s = mux.NewRouter().StrictSlash(true)

	app.s.Use(middleware.Metrics(app.pr))
	app.s.HandleFunc("/health", app.HealthHandler)
	app.s.Handle("/metrics", promhttp.Handler())
	app.s.HandleFunc("/rooms/list", app.GetRoomHandler).Name("get-Room").
		Methods("GET")
	app.s.HandleFunc("/bookings/list", app.GetBookingHandler).Name("get-Booking").
		Methods("GET")
	app.s.HandleFunc("/rooms", app.AddRoomHandler).Name("add-Room").
		Methods("POST")
	app.s.HandleFunc("/bookings/create", app.AddBookingHandler).Name("add-Booking").
		Methods("POST")
	app.s.HandleFunc("/rooms/{id:[0-9]+}", app.UpdateRoomHandler).Name("update-Room").
		Methods("PUT")
	app.s.HandleFunc("/bookings/{id:[0-9]+}", app.UpdateBookingHandler).Name("update-Booking").
		Methods("PUT")
	app.s.HandleFunc("/rooms/{id:[0-9]+}", app.DeleteRoomHandler).Name("delete-Room").
		Methods("DELETE")
	app.s.HandleFunc("/bookings/{id:[0-9]+}", app.DeleteBookingHandler).Name("delete-Booking").
		Methods("DELETE")
}

func (app *application) HealthHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}


func (app *application) GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("room_id")
	if id == "" {
		byteValue, err := app.faqRepository.ListRooms(0)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(byteValue)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		p, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		byteValue1, err := app.faqRepository.ListRooms(p)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(byteValue1)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (app *application) GetBookingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("room_id")
	if id == "" {
		byteValue, err := app.faqRepository.ListBooking(0)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(byteValue)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		p, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		byteValue1, err := app.faqRepository.ListBooking(p)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(byteValue1)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (app *application) AddBookingHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ques models.Booking
	err = json.Unmarshal(body, &ques)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = app.faqRepository.AddBooking(ques)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) AddRoomHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ques models.Room
	err = json.Unmarshal(body, &ques)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = app.faqRepository.AddRoom(ques)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) UpdateRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var rout models.Room
	err = json.Unmarshal(body, &rout)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.faqRepository.UpdateRoom(rout, id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) UpdateBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var rout models.Booking
	err = json.Unmarshal(body, &rout)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.faqRepository.UpdateBooking(rout, id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) DeleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.faqRepository.DeleteRoom(id)
	{
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) DeleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = app.faqRepository.DeleteBooking(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
