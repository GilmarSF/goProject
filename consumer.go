package main

import (
	//"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"bytes"
)

type Usuario struct {
    //"omitempty", quando usado omite este atributo caso seja 'null'
    ID      string  `json:"id,omitempty"`
    Nome    string  `json:"nome,omitempty"`
    Email   string  `json:"email,omitempty"`
    Senha   string  `json:"senha,omitempty"`
}

func main() {

	// Grava na API
	payload := Usuario{Nome:"testPost",Email:"foi@certo.com", Senha:"uruu"}
	// convert struct para json
	jsonValue, _ := json.Marshal(payload)

	_, err := http.Post("http://localhost:8080/cadastros/", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}

	/*// coleta da API 
	respostaAPI, err := http.Get("http://localhost:8080/cadastros/2")
	if err != nil {
		log.Fatal(err)
	}
	
	dadosResposta, err := ioutil.ReadAll(respostaAPI.Body)
	if err != nil {
		log.Fatal(err)
	}

	var usuario Usuario
	// coverte de json para struct
	json.Unmarshal(dadosResposta, &usuario)

	fmt.Println(usuario.Nome)*/
}