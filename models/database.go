package models

// SQLDataBase struct
type SQLDataBase struct {
	Server          string   `toml:"Server"`
	Database        string   `toml:"Database"`
	Port            int      `toml:"Port"`
	ApplicationName string   `toml:"ApplicationName"`
	MaxIdleConns    int      `toml:"MaxIdleConns"`
	MaxOpenConns    int      `toml:"MaxOpenConns"`
	ConnMaxLifetime duration `toml:"ConnMaxLifetime"`
	UserID          string
	Password        string
}

type Room struct {
	Room_id     int    `json:"room_id"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type Booking struct {
	Booking_id int    `json:"booking_id"`
	Date_start string `json:"date_start"`
	Date_end   string `json:"date_end"`
	Room_id    int    `json:"room_id"`
}
