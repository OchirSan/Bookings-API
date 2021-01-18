package booking

import (
	"database/sql"
	"encoding/json"
	"errors"
	"Booking/models"
)

type Booking interface {
	AddRoom(order models.Room) error
	AddBooking(order models.Booking) error
	ListRooms(id int) ([]byte, error)
	ListBooking(p int) ([]byte, error)
	DeleteRoom(id int) error
	DeleteBooking(id int) error
	UpdateRoom(s models.Room, id int) error
	UpdateBooking(q models.Booking, id int) error
}

type booking struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) Booking {
	return &booking{db: db}
}

func (r *booking) AddRoom(v models.Room) error {
	sqlStatement := `INSERT INTO rooms (description, price) VALUES ($1, $2)`
	_, err := r.db.Exec(sqlStatement, v.Description, v.Price)
	if err != nil {
		return err
	}
	return errors.New("data add in db")
}

func (r *booking) AddBooking(v models.Booking) error {
	sqlStatement := `INSERT INTO booking (date_start, date_end, room_id) VALUES ($1, $2, $3,$4)`
	_, err := r.db.Exec(sqlStatement, v.Date_start, v.Date_end, v.Room_id)
	if err != nil {
		return err
	}
	return errors.New("data add in db")
}

func (r *booking) ListBooking(p int) ([]byte, error) {
	if p != 0 {
		rows, err := r.db.Query("SELECT booking_id, date_start, date_end, room_id FROM booking where room_id = $1 ", p)
		if err != nil {
			return nil,err
		}
		defer rows.Close()
		contentfromdb := make([]models.Booking, 0)
		for rows.Next() {
			c := models.Booking{}
			err := rows.Scan(&c.Booking_id,&c.Date_start, &c.Date_end, &c.Room_id)
			if err != nil {
				return nil, err
			}
			contentfromdb = append(contentfromdb, c)
		}
		jsonContentFromDB, err := json.Marshal(contentfromdb)
		if err != nil {
			return nil, err
		}
		return jsonContentFromDB, nil
	} else {
		rows, err := r.db.Query("SELECT booking_id,date_start, date_end,room_id FROM booking order by date_start")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		contentfromdb := make([]models.Booking, 0)
		for rows.Next() {
			c := models.Booking{}
			err := rows.Scan(&c.Booking_id, &c.Date_start, &c.Date_end,&c.Room_id)
			if err != nil {
				return nil, err
			}


			contentfromdb = append(contentfromdb, c)
		}
		jsonContentFromDB, err := json.Marshal(contentfromdb)
		if err != nil {
			return nil, err
		}
		return jsonContentFromDB, nil
	}
}

func (r *booking) ListRooms(id int) ([]byte, error) {
	if id == 0 {
		rows, err := r.db.Query("SELECT room_id, description, price FROM rooms order by price")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		contentfromdb := make([]models.Room, 0)
		for rows.Next() {
			c := models.Room{}
			err := rows.Scan(&c.Room_id, &c.Description, &c.Price)
			if err != nil {
				return nil, err
			}
			contentfromdb = append(contentfromdb, c)
		}
		jsonContentFromDB, err := json.Marshal(contentfromdb)
		if err != nil {
			return nil, err
		}
		return jsonContentFromDB, nil
	} else {
		rows := r.db.QueryRow("SELECT room_id,description, price FROM rooms where room_id = $1", id)
		c := models.Room{}
		err := rows.Scan(&c.Room_id, &c.Description, &c.Price)
		if err != nil {
			return nil, err
		}
		jsonc, err2 := json.Marshal(c)
		if err2 != nil {
			return nil, err2
		}
		return jsonc, nil
	}
}

func (r *booking) DeleteRoom(id int) error {
	sqlStatement := `DELETE FROM rooms where room_id = $1`
	_, err := r.db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	sqlStatement2 := `DELETE FROM booking where room_id = $1`
	_, err= r.db.Exec(sqlStatement2, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *booking) DeleteBooking(id int) error {
	sqlStatement := `DELETE FROM booking where booking_id = $1`
	_, err := r.db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *booking) UpdateRoom(s models.Room, id int) error {
	sqlStatement := `UPDATE rooms set description = $1, price =$2  where room_id = $3`
	_, err := r.db.Exec(sqlStatement, s.Description, s.Price, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *booking) UpdateBooking(s models.Booking, id int) error {
	sqlStatement := `UPDATE booking set room_id = $1, date_start = $2, date_end = $3  where booking_id = $5`
	_, err := r.db.Exec(sqlStatement, s.Room_id, s.Date_start,s.Date_end, id)
	if err != nil {
		return err
	}
	return nil
}


