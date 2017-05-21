package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Usuario struct {
    //"omitempty", quando usado omite este atributo caso seja 'null'
    ID      string  `json:"id,omitempty"`
    Nome    string  `json:"nome,omitempty"`
    Email   string  `json:"email,omitempty"`
    Senha   string  `json:"senha,omitempty"`
}

func main() {
	// coleta da API 
	respostaAPI, err := http.Get("http://localhost:8080/cadastros/4")
	if err != nil {
		log.Fatal(err)
	}
	
	respostaDados, err := ioutil.ReadAll(respostaAPI.Body)
	if err != nil {
		log.Fatal(err)
	}

	var usuario Usuario
	// coverte de json para struct
	json.Unmarshal(respostaDados, &usuario)

	fmt.Println(usuario.Nome)

}