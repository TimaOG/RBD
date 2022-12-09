package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	Login    = "root"
	Password = "root"
	DBip     = "127.0.0.1:3306"
	DBName   = "test_bd"
)

func addUser(userEmail string, userPassword string, phoneNumber string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	sqlRequest := `INSERT INTO users (Password, Email, PhoneNumber, IsAdmin) VALUES (?, ?, ?, ?)`
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	res, err := db.Query(sqlRequest, hashedPassword, userEmail, phoneNumber, false)
	if err != nil {
		panic(err.Error())
	}
	defer res.Close()
	defer db.Close()
}

func checkForUserInSystem(email string) int {
	sqlRequest := "SELECT COUNT(*) FROM users WHERE email = ?;"
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	var answ int
	res := db.QueryRow(sqlRequest, email)
	err := res.Scan(&answ)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return answ
}

func getUserByEmail(email string) User {
	sqlRequest := `SELECT pkId, password, email, isAdmin FROM users WHERE email = ?;`
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	var user User
	res := db.QueryRow(sqlRequest, email)
	err := res.Scan(&user.Id, &user.Password, &user.Email, &user.IsAdmin)
	if err == sql.ErrNoRows {
		return User{}
	}
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return user
}
func getUserById(id uint32) User {
	sqlRequest := `SELECT pkId, password, email FROM users WHERE pkId = ?;`
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	var user User
	res := db.QueryRow(sqlRequest, id)
	err := res.Scan(&user.Id, &user.Password, &user.Email)
	if err == sql.ErrNoRows {
		return User{}
	}
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return user
}

func addReqToUser(id uint32, taskType uint8, discribtion string, adress string, usefultime string) {
	sqlRequest := `INSERT INTO orders (fkCodeClient, fkOrderType, discribtion, address, usefulTime, requestStatus) VALUES (?,?,?,?,?,?)`
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	res, err := db.Query(sqlRequest, id, taskType, discribtion, adress, usefultime, "На рассмотрении")
	if err != nil {
		panic(err.Error())
	}
	defer res.Close()
	defer db.Close()
}

func getOrdersByUserId(orderID uint32) []Order {
	sqlRequest := "SELECT t1.pkId, t3.typeName, t1.address, t1.usefulTime, t1.discribtion, t1.requestStatus, t1.startTime, t1.closeTime, t2.PhoneNumber FROM orders AS t1 LEFT JOIN users AS t2 ON t1.fkCodeClient = t2.pkId LEFT JOIN ordersType AS t3 ON t1.fkOrderType = t3.pkId WHERE t1.fkCodeClient = ?;"
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	res, err := db.Query(sqlRequest, orderID)
	if err != nil {
		panic(err.Error())
	}
	orders := []Order{}
	var tmpStartTime sql.NullString
	var tmpCloseTime sql.NullString
	for res.Next() {
		order := Order{}
		err := res.Scan(&order.Id, &order.TypeName, &order.Address, &order.UsefulTime, &order.Discribtion, &order.RequestStatus, &tmpStartTime, &tmpCloseTime, &order.UserPhone)
		if tmpStartTime.Valid {
			order.StartTime = tmpStartTime.String
		}
		if tmpCloseTime.Valid {
			order.CloseTime = tmpCloseTime.String
		}
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
	}
	defer db.Close()
	return orders
}

func getOrdersByUserIdAndType(orderID uint32, oType uint8) []Order {
	sqlRequest := "SELECT t1.pkId, t3.typeName, t1.address, t1.usefulTime, t1.discribtion, t1.requestStatus, t1.startTime, t1.closeTime, t2.PhoneNumber FROM orders AS t1 LEFT JOIN users AS t2 ON t1.fkCodeClient = t2.pkId LEFT JOIN ordersType AS t3 ON t1.fkOrderType = t3.pkId WHERE t1.fkCodeClient = ? AND t1.fkOrderType = ?;"
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	res, err := db.Query(sqlRequest, orderID, oType)
	if err != nil {
		panic(err.Error())
	}
	orders := []Order{}
	var tmpStartTime sql.NullString
	var tmpCloseTime sql.NullString
	for res.Next() {
		order := Order{}
		err := res.Scan(&order.Id, &order.TypeName, &order.Address, &order.UsefulTime, &order.Discribtion, &order.RequestStatus, &tmpStartTime, &tmpCloseTime, &order.UserPhone)
		if tmpStartTime.Valid {
			order.StartTime = tmpStartTime.String
		}
		if tmpCloseTime.Valid {
			order.CloseTime = tmpCloseTime.String
		}
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
	}
	defer db.Close()
	return orders
}
func getOrdersAll() []Order {
	sqlRequest := "call getAllOrders()"
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	res, err := db.Query(sqlRequest)
	if err != nil {
		panic(err.Error())
	}
	var tmpStartTime sql.NullString
	var tmpCloseTime sql.NullString
	orders := []Order{}
	for res.Next() {
		order := Order{}
		err := res.Scan(&order.Id, &order.TypeName, &order.Address, &order.UsefulTime, &order.Discribtion, &order.RequestStatus, &tmpStartTime, &tmpCloseTime, &order.UserPhone)
		if tmpStartTime.Valid {
			order.StartTime = tmpStartTime.String
		}
		if tmpCloseTime.Valid {
			order.CloseTime = tmpCloseTime.String
		}
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
	}
	defer db.Close()
	return orders
}

func changeOrderStatus(status string, orderID uint32) {
	sqlRequest := "UPDATE orders SET requestStatus = ? WHERE pkId = ?;"
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	res, err := db.Query(sqlRequest, status, orderID)
	if err != nil {
		panic(err.Error())
	}
	defer res.Close()
	defer db.Close()
	return
}

func getOrdersType() []OrderType {
	sqlRequest := `SELECT * FROM ordersType`
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", Login, Password, DBip, DBName))
	res, err := db.Query(sqlRequest)
	if err != nil {
		panic(err.Error())
	}
	orders := []OrderType{}
	for res.Next() {
		order := OrderType{}
		err := res.Scan(&order.Id, &order.TypeName, &order.Price, &order.Discribtion)

		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
	}
	defer db.Close()
	return orders
}
