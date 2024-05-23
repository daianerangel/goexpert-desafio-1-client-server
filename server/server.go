package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ExchangeRate struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func fetchExchangeRate(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var rate ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&rate); err != nil {
		return "", err
	}
	return rate.USDBRL.Bid, nil
}

func saveQuotationToDB(ctx context.Context, db *sql.DB, quotation string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO quotations (quotation) VALUES (?)", quotation)
	return err
}

func quotationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()

		quotation, err := fetchExchangeRate(ctx)
		if err != nil {
			log.Println("Error fetching exchange rate:", err)
			http.Error(w, "Failed to fetch exchange rate", http.StatusInternalServerError)
			return
		}

		dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer dbCancel()

		if err := saveQuotationToDB(dbCtx, db, quotation); err != nil {
			log.Println("Error saving quotation to DB:", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"bid": quotation})
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/currency_exchange")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", quotationHandler(db))
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
