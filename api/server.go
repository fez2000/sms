package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexgear/sms/common"
	"github.com/alexgear/sms/modem"
	//"github.com/gorilla/mux"
	db "github.com/alexgear/sms/database"
	"github.com/satori/go.uuid"
)

var err error

//response structure to /sms
type SMSResponse struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
	UUID   string `json:"uuid"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

func sendSMSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-type", "application/json")
	r.ParseForm()
	log.Printf("sendSMSHandler: %#v", r.Form)
	uuid := uuid.NewV1()
	sms := &common.SMS{
		UUID:   uuid.String(),
		Mobile: r.FormValue("to"),
		Body:   r.FormValue("text"),
		Status: "pending"}
	err = db.InsertMessage(sms)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := SMSResponse{Status: 200, Text: sms.Body, UUID: sms.UUID}
	toWrite, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(toWrite)
	return
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-type", "application/json")
	balance, err := modem.GetBalance(`*111#`)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := BalanceResponse{Balance: balance}
	toWrite, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(toWrite)
	return
}

func InitServer(host string, port int) error {
	//r := mux.NewRouter()
	http.HandleFunc("/api/sms", sendSMSHandler)
	http.HandleFunc("/api/balance", getBalanceHandler)

	//http.Handle("/", r)

	bind := fmt.Sprintf("%s:%d", host, port)
	log.Println("listening on: ", bind)
	return http.ListenAndServe(bind, nil)

}
