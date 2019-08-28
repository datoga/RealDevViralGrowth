package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type AccountController struct {
	Port    int
	service *AccountService
}

func NewAccountController(port int, service *AccountService) *AccountController {
	return &AccountController{
		Port:    port,
		service: service,
	}
}

func (controller AccountController) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/register", controller.RegisterHandler).Methods("POST")
	router.HandleFunc("/userProfile", controller.UserProfileHandler).Queries("id", "{id}")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	http.ListenAndServe(":"+strconv.Itoa(controller.Port), loggedRouter)
}

type RegisterRequest struct {
	UserName   string `json:"userName"`
	InviteCode string `json:"inviteCode",omitempty`
}

func (controller AccountController) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	rDump, err := httputil.DumpRequest(r, false)

	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(rDump))
	}

	decoder := json.NewDecoder(r.Body)

	var registerRequest RegisterRequest

	err = decoder.Decode(&registerRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var referer *Profile

	if registerRequest.InviteCode != "" {
		var found bool

		referer, found = controller.service.GetReferer(registerRequest.InviteCode)

		if !found {
			http.Error(w, "The referer code has not been found", http.StatusBadRequest)
			return
		}
	}

	profile := controller.service.AddProfile(registerRequest.UserName, referer)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strconv.Itoa(profile.Id)))
	w.Write([]byte("\n"))
}

func (controller AccountController) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	sID := r.FormValue("id")

	if sID == "" {
		http.Error(w, "The ID should not be void", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(sID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile, found := controller.service.GetProfile(id)

	if !found {
		http.Error(w, "Profile not found", http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(profile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	w.Write([]byte("\n"))
}
