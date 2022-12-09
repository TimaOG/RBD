package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func startServerAndAddHandles() {
	r := mux.NewRouter()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/", startPage)      //
	r.HandleFunc("/login", loginPage) //
	r.HandleFunc("/registration", registrationPage)

	r.HandleFunc("/buy/{type}", buyPage)

	r.HandleFunc("/account", accountPage)
	r.HandleFunc("/account/exit", exitPage)

	r.HandleFunc("/admin", adminPage)
	r.HandleFunc("/change/{order}", adminChangeOrderStatus)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("http://localhost:8080")
	startServerAndAddHandles()
}
