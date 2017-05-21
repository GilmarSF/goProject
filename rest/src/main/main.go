package main
 
import (
    "encoding/json"
    "log"
    "net/http"
	"fmt"
    "github.com/gorilla/mux"
)
 
type Pessoa struct {
    ID        string   `json:"id"`
    Firstname string   `json:"firstname"`
    Lastname  string   `json:"lastname"`
	City  string `json:"city"`
    State string `json:"state"`
    
	/* Address usado como ponteiro de outra struct
	Address   *Address `json:"address"`

	"omitempty", quando usado omite este atributo caso seja 'null'
	Address   *Address `json:"address,omitempty"` */
}

/* O endereço poderia ser usado como ponteiro para "Pessoa"
type Address struct {
    City  string `json:"city"`
    State string `json:"state"`
}*/
 
var turma []Pessoa
 
// Traz apenas uma pessoa
func GetPessoaEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req) // salva o 'id' passado na requisição
    for _, item := range turma { // For sera executado de acordo com o tamanho do array 'turma'
        if item.ID == params["id"] { // Se encontrar uma pessoa com o ID enviado
			fmt.Fprintln(w, "Dados de: "+item.Firstname+"\n")
            json.NewEncoder(w).Encode(item) // imprime os dados da pessoa
            return // sai da func
        }
    }
	fmt.Fprintln(w, "Pessoa não encontrada!")
	/* Para exibir a struct em branco
    json.NewEncoder(w).Encode(&Pessoa{})*/
}
 
// traz a turma toda. Todas todas pessoas do array 'turma'
func GetturmaEndpoint(w http.ResponseWriter, req *http.Request) {
    json.NewEncoder(w).Encode(turma)
}
 
func CreatePessoaEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var pessoa Pessoa
    _ = json.NewDecoder(req.Body).Decode(&pessoa)
    pessoa.ID = params["id"]
    turma = append(turma, pessoa) // adiciona ao final de 'turma'
    json.NewEncoder(w).Encode(turma)
}
 
/* Quando o 'id' for deletado, sera preciso recriar o slice com o dados restantes, por isso o uso do 'indiceArray'
o for a cada rodada grava em 'indiceArray' o indice do array que ele esta.*/
func DeletePessoaEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    for indiceArray, item := range turma {
        if item.ID == params["id"] {
            turma = append(turma[:indiceArray], turma[indiceArray+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(turma)
}
 
func main() {
    router := mux.NewRouter()
	// caso tenho ponteiro para o Address usar linha abaixo
    //turma = append(turma, Pessoa{ID: "1", Firstname: "Nic", Lastname: "Raboy", Address: &Address{City: "Dublin", State: "CA"}})
    turma = append(turma, Pessoa{ID: "1", Firstname: "Nic", Lastname: "Raboy", City: "Dublin", State: "CA"})
    turma = append(turma, Pessoa{ID: "2", Firstname: "Maria", Lastname: "Raboy"})
    router.HandleFunc("/turma", GetturmaEndpoint).Methods("GET")
    router.HandleFunc("/turma/{id}", GetPessoaEndpoint).Methods("GET")
    router.HandleFunc("/turma/{id}", CreatePessoaEndpoint).Methods("POST")
    router.HandleFunc("/turma/{id}", DeletePessoaEndpoint).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8080", router))
}