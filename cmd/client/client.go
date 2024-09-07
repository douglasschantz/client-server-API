package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Bid struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var buscaBid Bid
	err = json.Unmarshal(body, &buscaBid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer a parse da resposta:\n")
		panic(err)
	}

	err = GravarTxt(buscaBid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gravar txt: \n")
		panic(err)
	}

}

func GravarTxt(cotacao_bid Bid) (err error) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", cotacao_bid.Bid))
	if err != nil {
		return err
	}

	return nil
}
