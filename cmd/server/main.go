package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	gatewayhttp "github.com/Saidharan-dev/gateway/internal/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	failFirst := int64(1)
	if value := os.Getenv("FAIL_FIRST"); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			failFirst = parsed
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("gateway ok"))
	})
	mux.Handle("/charge", gatewayhttp.NewChargeHandler(failFirst))

	log.Printf("gateway listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}