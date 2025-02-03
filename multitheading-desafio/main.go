package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const apiViaCep = "viacep"
const apiBrasilCep = "brasilcep"

type CepInterface interface {
}

type CEPBrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type CEPViaCEPAPI struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {

	var cep string
	var cepVia CEPViaCEPAPI
	var cepBrasil CEPBrasilAPI

	cepChannel1 := make(chan CEPViaCEPAPI)
	cepChannel2 := make(chan CEPBrasilAPI)

	fmt.Println("----- Busca CEP ----- ")
	fmt.Println("Digite um CEP válido - Sem pontuações")
	fmt.Scan(&cep)

	go func() {
		url := buildUrl(cep, "viacep")
		GetCep(url, &cepVia)
		cepChannel1 <- cepVia

	}()

	go func() {
		url := buildUrl(cep, "brasilcep")
		GetCep(url, &cepBrasil)
		cepChannel2 <- cepBrasil

	}()

	select {

	case cepVia := <-cepChannel1:
		fmt.Println("Resultado proveniente da API Via CEP: ")
		fmt.Println(cepVia)

	case cepBrasil := <-cepChannel2:
		fmt.Println("Resultado proveniente da Brasil API: ")
		fmt.Println(cepBrasil)

	case <-time.After(time.Second * 1):
		fmt.Println("TimeOut")

	}
}

func GetCep(url string, cep CepInterface) error {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&cep); err != nil {
		log.Fatalf("Erro no parse da resposta")
		return err
	}

	return nil
}

func buildUrl(cep, nome string) string {
	if nome == apiViaCep {
		return "http://viacep.com.br/ws/" + cep + "/json"

	} else if nome == apiBrasilCep {
		return "https://brasilapi.com.br/api/cep/v1/" + cep
	}

	return ""

}
