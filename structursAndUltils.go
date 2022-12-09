package main

const SiteAdress = "http://localhost:8080/"

type User struct {
	Id       uint32 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	IsAdmin  bool   `json:"isAdmin"`
}

type OrderType struct {
	Id          uint32
	TypeName    string
	Price       int
	Discribtion string
}

type Order struct {
	Id            uint32
	TypeName      string
	Address       string
	UsefulTime    string
	Discribtion   string
	RequestStatus string
	StartTime     string
	CloseTime     string
	UserPhone     string
}
type ViewAdminData struct {
	Types []OrderType
	Data  []Order
}
