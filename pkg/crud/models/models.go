package models

import "time"

// ConsoleReponse ...
type ConsoleReponse struct {
	Response []string `json:"response"`
}

// Admin ...
type Admin struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Category of quotes
type Category struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Infrastructure ...
type Infrastructure struct {
	School   string `json:"school"`
	Hospital string `json:"hospital"`
	Church   string `json:"church"`
	Park     string `json:"park"`
}

// Houses ...
type Houses struct {
	ID             string         `json:"id"`
	CodeID         string         `json:"code_id"`
	Title          string         `json:"title"`
	Price          int            `json:"price"`
	CreatedAt      string         `json:"created_at"`
	Description    string         `json:"description"`
	Owner          Owners         `json:"owner_id"`
	Infrastructure Infrastructure `json:"infrastructure"`
	Address        Address        `json:"addr"`
	Rooms          string         `json:"rooms"`
	City           Cities         `json:"city"`
	Rayons         Rayons         `json:"rayon"`
	Category       Category       `json:"category"`
	Status         string         `json:"status"`
	Structure      Structure      `json:"structure"`
}

// Structure ...
type Structure struct {
	Square int `json:"square"`
	Floor  int `json:"floor"`
}

// Cities ...
type Cities struct {
	ID   string `json:"id"`
	City string `json:"city"`
}

// Rayons ...
type Rayons struct {
	ID     string `json:"id"`
	Rayon  string `json:"rayon"`
	CityID string `json:"city_id"`
}

// Address ...
type Address struct {
	Street    string `json:"street"`
	Number    string `json:"number"`
	Longitude string `json:"long"`
	Latittude string `json:"lati"`
}

// Agents type ...
type Agents struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Individual string `json:"individual"`
	Logo       string `json:"logo"`
	Phone      string `json:"phone"`
	Others     string `json:"others"` // jsonb
}

// Owners ...
type Owners struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Lastname string    `json:"lastname"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Image    string    `json:"image"`
	Bio      string    `json:"bio"`
	Password string    `json:"password"`
	ReqDate  time.Time `json:"reqdate"`
}

// Users ...
type Users struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	UserName string    `json:"username"`
	Image    string    `json:"image"`
	Bio      string    `json:"bio"`
	Password string    `json:"password"`
	Status   string    `json:"status"`
	ReqDate  time.Time `json:"reqdate"`
}

// SearchObject of quotes
type SearchObject struct {
	City           string `json:"city"`
	Rayon          string `json:"rayon"`
	Room           int    `json:"room"`
	Floor          int    `json:"floor"`
	Square         int    `json:"square"`
	Type           string `json:"type"`
	Category       string `json:"category"`
	Infrastructure Infrastructure
}

// Burger ...
type Burger struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	Removed bool   `json:"removed"`
}
