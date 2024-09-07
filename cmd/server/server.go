package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CotacaoDolar struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type BidOnly struct {
	Bid string `json:"bid"`
}

func main() {

	http.HandleFunc("/cotacao", BuscaDolarHandler)
	http.ListenAndServe(":8080", nil)

}

func BuscaDolarHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cot_dolar, error := BuscaCotacaodolarCtx()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	/*result, err := json.Marshal(cot_dolar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)*/

	//json.NewEncoder(w).Encode(cot_dolar)
	fmt.Println(cot_dolar)
	err := SalvarBidDB(cot_dolar)
	if err != nil {
		log.Printf("Erro ao gravar no banco de dados: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bidOnly := BidOnly{Bid: cot_dolar.Usdbrl.Bid}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bidOnly)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Request Finalizada")
}

func BuscaCotacaodolarCtx() (*CotacaoDolar, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição:\n")
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cotacao CotacaoDolar

	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil

}

func SalvarBidDB(bid *CotacaoDolar) error {
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&CotacaoDolar{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := db.WithContext(ctx).Create(&bid).Error; err != nil {
		return err
	}

	return nil
}
