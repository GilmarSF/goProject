package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	//"bytes"
)

type Usuario struct {
    //"omitempty", quando usado omite este atributo caso seja 'null'
    ID      string  `json:"id,omitempty"`
    Nome    string  `json:"nome,omitempty"`
    Email   string  `json:"email,omitempty"`
    Senha   string  `json:"senha,omitempty"`
}

type CategoriaFull struct{
    ID      string  `json:"id,omitempty"`
    Nome    string  `json:"nome,omitempty"`
    Total   string  `json:"total,omitempty"`
}

type CategoriaEach struct{
    ID      string `json:"id,omitempty"`
    Nome    string `json:"nome,omitempty"`
    Regiao  string `json:"regiao,omitempty"`
    Total   string `json:"total,omitempty"`
}

type NovaDenuncia struct{
    Categoria string    `json:"categoria,omitempty"`
    Localidade string   `json:"localidade,omitempty"`
}

func main() {
	/*
	////////////////////////////// Grava na API ///////////////////////////////////
	
	// Informações de como deve ser enviado para salvar no Banco
	// categorias: 1 = assalto; 2 = abuso sexual; 3 = transito; 4 = violencia; 5 = assassinato
	// Localidades: 1 = sul; 2 = norte; 3 = lest; 4 = oeste; 5 = central
	
	grava := NovaDenuncia{Categoria:"3",Localidade:"5"}
	// convert struct para json
	jsonValue, _ := json.Marshal(grava)
	// Envia para a API gravar no banco
	_, err := http.Post("http://localhost:8080/categorias/", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	} */

	/////////////////////////////// Coleta da API ////////////////////////////////// 
	//respostaAPI, err := http.Get("http://localhost:8080/categorias/") //TRAZ TODAS AS CATEGORIAS, UTIL PARA O GRAFICO GERAL
	respostaAPI, err := http.Get("http://localhost:8080/categorias/5")
	if err != nil {
		log.Fatal(err)
	}
	
	dadosResposta, err := ioutil.ReadAll(respostaAPI.Body)
	if err != nil {
		log.Fatal(err)
	}

	var categEach []CategoriaEach
	// coverte de json para struct
	json.Unmarshal(dadosResposta, &categEach)

	for _, item := range categEach{
		fmt.Println(item.Nome)	
	}	
}