package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func startPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	data := getOrdersType()
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	data := ""
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if CheckUserSingIn(w, r, email, password) {
			target := r.URL.Query().Get("target")
			http.Redirect(w, r, (SiteAdress + target), http.StatusSeeOther)
			return
		} else {
			data = "Неверный логин или пароль"
		}
	}
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}
func registrationPage(w http.ResponseWriter, r *http.Request) {
	data := ""
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")
		phoneNumber := r.FormValue("phone")

		createOrNot := true
		if password != password2 {
			data = "пароли не равны"
			createOrNot = false
		} else if checkForUserInSystem(email) >= 1 {
			data = "Пользователь уже есть в системе"
			createOrNot = false
		}
		if createOrNot {
			addUser(email, password, phoneNumber)
			http.Redirect(w, r, fmt.Sprintf(SiteAdress+"login"), http.StatusSeeOther)
			return
		}
	}
	t, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}
func buyPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buyType := vars["type"]
	if !CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(SiteAdress+"login?target=buy/"+buyType), http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		taskId, _ := strconv.Atoi(buyType)
		userIDstr, _ := r.Cookie("userID")
		userID, _ := strconv.Atoi(userIDstr.Value)
		discribtion := r.FormValue("discribtion")
		adress := r.FormValue("adress")
		utime := r.FormValue("utime")
		addReqToUser(uint32(userID), uint8(taskId), discribtion, adress, utime)
		http.Redirect(w, r, fmt.Sprintf(SiteAdress), http.StatusSeeOther)
	}
	t, err := template.ParseFiles("templates/createRequest.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func accountPage(w http.ResponseWriter, r *http.Request) {
	if !CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(SiteAdress+"login?target=account"), http.StatusSeeOther)
		return
	}
	userIDstr, _ := r.Cookie("userID")
	userID, _ := strconv.Atoi(userIDstr.Value)
	var orders []Order
	if r.Method == http.MethodPost {
		sel := r.FormValue("taskType")
		selInt, _ := strconv.Atoi(sel)
		orders = getOrdersByUserIdAndType(uint32(userID), uint8(selInt))
	} else {
		orders = getOrdersByUserId(uint32(userID))
	}
	data := getOrdersType()
	var dv = ViewAdminData{data, orders}
	t, err := template.ParseFiles("templates/account.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, dv)
}

func exitPage(w http.ResponseWriter, r *http.Request) {
	UserLogOut(w, r)
	http.Redirect(w, r, fmt.Sprintf(SiteAdress), http.StatusSeeOther)
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	if !CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(SiteAdress+"login?target=account"), http.StatusSeeOther)
		return
	}
	isAdmin, _ := r.Cookie("isAdmin")
	fmt.Println(isAdmin.Value)
	if isAdmin.Value != "true" {
		http.Redirect(w, r, fmt.Sprintf(SiteAdress), http.StatusSeeOther)
		return
	}
	orders := getOrdersAll()
	t, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, orders)
}

func adminChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	if !CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(SiteAdress+"login?target=account"), http.StatusSeeOther)
		return
	}
	isAdmin, _ := r.Cookie("isAdmin")
	if isAdmin.Value != "true" {
		http.Redirect(w, r, fmt.Sprintf(SiteAdress), http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		status := r.FormValue("status")
		vars := mux.Vars(r)
		orderIdStr := vars["order"]
		orderID, _ := strconv.Atoi(orderIdStr)
		changeOrderStatus(status, uint32(orderID)) //Сделать тригерр на обновление, если статус = принят, то поставить время начала
	}
	orders := getOrdersAll()
	t, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, orders)

}
