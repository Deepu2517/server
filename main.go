package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Deepu2517/server/lib"
	"github.com/gorilla/mux"
)

func main() {
	log.SetFlags(log.Lshortfile)
	prodClientID := flag.String("production-client-id", "v7i8VAzo9hUK6paPBbBQelzdfGdGoSju3NBC6wwz", "Production Client ID")
	prodClientSecret := flag.String("production-client-secret", "cYDWhlgGPL2fq4q7m2IElpA58vDvCDBx8pZdRuYSYL4z1mhsV7tvtdhgfvXA5KZ4udMK4C8HBFaLUnozDm2CqOhrmoQ8d7o4fiuTYeTndPXGxRlBAhIwV6R9fDnESdiz", "Production Client Secret")
	testClientID := flag.String("test-client-id", "", "Test Client ID")
	testClientSecret := flag.String("test-client-secret", "", "Test Client Secret")
	flag.Parse()

	// if *prodClientID == "" {
	// 	log.Fatal("Production Client ID is missing")
	// }
	//
	// if *prodClientSecret == "" {
	// 	log.Fatal("Production Client secret is missing")
	// }
	//
	// if *testClientID == "" {
	// 	log.Fatal("Test Client ID is missing")
	// }
	//
	// if *testClientSecret == "" {
	// 	log.Fatal("Test Client Secret is missing")
	// }

	lib.SetCredentials(*prodClientID, *prodClientSecret, *testClientID, *testClientSecret)

	router := mux.NewRouter()
	router.HandleFunc("/create/", createOrderTokens).Methods("POST")
	router.HandleFunc("/create", createOrderTokens).Methods("POST")
	router.HandleFunc("/status", statusHandler).Methods("GET")
	router.HandleFunc("/refund/", refundHandler).Methods("POST")
	router.HandleFunc("/ping", pingHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(serverAddr, LoggingHandler(router)))
}

func createOrderTokens(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := lib.CreateOrderTokens(r.FormValue("env"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	env := r.FormValue("env")
	orderID := r.FormValue("order_id")
	transactionID := r.FormValue("transaction_id")

	data, err := lib.GetOrderStatus(env, r.Header.Get("Authorization"), orderID, transactionID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func refundHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	env := r.FormValue("env")
	transactionID := r.FormValue("transaction_id")
	amount := r.FormValue("amount")
	refundType := r.FormValue("type")
	body := r.FormValue("body")

	statusCode, err := lib.InitiateRefund(env, r.Header.Get("Authorization"), transactionID, amount, refundType, body)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(statusCode)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
