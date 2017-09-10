// Todo codigo em Go precisa iniciar com este package
package main

// Importa todas as bibliotecas que serão utilizadas no código
import (

	// bibliotecas Nativas do Golang
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	// Pode ser usado bibliotecas não nativas
	// /gorilla/mux para o servidor HTTP
	// /denisenkom/go-mssqldb para a comunicação com Banco SQLServer
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

// Struct Principal criada para armazenar os dados retornados do Banco de dados
// Structs são mais comuns para converter para JSON em Golang

type Denuncias struct {
	ID    string `json:"id,omitempty"`
	Nome  string `json:"nome,omitempty"`
	Total string `json:"total,omitempty"`
}

type DenunciasPorCategoria struct {
	ID     string `json:"id,omitempty"`
	Nome   string `json:"nome,omitempty"`
	Regiao string `json:"regiao,omitempty"`
	Total  string `json:"total,omitempty"`
}

type NovaDenuncia struct {
	Categoria  string `json:"categoria,omitempty"`
	Localidade string `json:"localidade,omitempty"`
}

/* VARIAVEIS GLOBAIS PODEM SER USADAS EM QUALQUER PARTE DO CÓDiGO*/

// array usado para enviar o total de cada categoria
var denuncias []Denuncias

// array usado para enviar o total de denuncias por regiao
var denunciasPorCategoria []DenunciasPorCategoria

// Usado para armazenar o ultimo 'id' do banco de dados
var proximoIdParaGravarNoBanco int

// Adicona mais uma denuncia
func gravarNovaDenuncia(w http.ResponseWriter, req *http.Request) {
	// modelo que deve enviado
	// {"categoria":"4","localidade":"2"}
	log.Printf("Post Nova Denuncia")
	var novaD NovaDenuncia

	// grava em 'novaD' os dados enviados
	erro := json.NewDecoder(req.Body).Decode(&novaD)
	if erro != nil {
		log.Println("erro em ao gravar em novaD: ", erro.Error())
	}
	// imprime no terminal os valores recebidos
	fmt.Println(novaD)

	bancoDeDados, erro := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
	if erro != nil {
		log.Println("erro ao conectar com o Banco de dados:", erro.Error())
	}

	insert, erro := bancoDeDados.Query(`INSERT into tab_denuncia (id, id_categoria, id_localidade) 
									VALUES (?1, ?2, ?3)`, proximoIdParaGravarNoBanco, novaD.Categoria, novaD.Localidade)
	if erro != nil {
		log.Println("erro no INSERT:", erro.Error())
	} else {
		proximoIdParaGravarNoBanco++
	}

	defer insert.Close()       // fecha o comando Query
	defer bancoDeDados.Close() // fecha conexão com o Banco
	// atualiza struct no banco
	atualizarCategorias()

	//json.NewEncoder(w).Encode(categorias)
}

// função para enviar apenas uma categoria com o total por regiao
func pegarUmaCategoria(w http.ResponseWriter, req *http.Request) {
	// OBSERVAÇÂO: comentarios de como funciona esta na 'func GetUsuario'
	log.Printf("Get uma Categoria")
	params := mux.Vars(req)
	var categoriaEncontrada []DenunciasPorCategoria
	for _, item := range denunciasPorCategoria {
		if strings.ToLower(item.Nome) == strings.ToLower(params["uri"]) {
			categoriaEncontrada = append(categoriaEncontrada, item)
		}
	}
	//json.NewEncoder(w).Encode(denunciasPorCategoria)
	json.NewEncoder(w).Encode(categoriaEncontrada)
}

// envia os dados das categorias via GET
func pegarTodasCategorias(w http.ResponseWriter, req *http.Request) {
	log.Printf("Get categorias")
	json.NewEncoder(w).Encode(denuncias)
}

// Atualiza a struct Categorias com os dados do Banco
func atualizarCategorias() {

	// Inicia conexão com Banco de dados no Azure
	bancoDeDados, erro := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
	// Caso ocorrer algum erro na conexão, o erro será salvo em 'erro'
	// Caso não houver erro o retorno será vazio 'nil'
	if erro != nil {
		log.Println("erro em atualizarCategorias(): ", erro.Error())
	}

	// Zera o Array antes de buscar novos dados no banco.
	// Dessa forma evita dados repetidos
	denuncias = denuncias[:0]
	denunciasPorCategoria = denunciasPorCategoria[:0]

	// Select traz todos os Nomes das categorias e as quantidades de denuncias de cada um
	retornoSelectBanco, erro := bancoDeDados.Query(`SELECT d.id_categoria, c.categoria, COUNT(d.id_categoria) 
                                            FROM tab_denuncia d JOIN tab_categoria c 
                                            ON d.id_categoria = c.id 
                                            GROUP BY d.id_categoria, c.categoria`)

	if erro != nil {
		log.Println("erro no SELECT das CategoriasFull:", erro.Error())
	}

	// retornoSelectBanco.Next usado para varrer o objeto 'retornoSelectBanco' e pegar os valores retornados da Query
	for retornoSelectBanco.Next() {
		// Struct criada para receber os dados do banco de dados
		// Uma Struct permite receber dados de diferentes tipos: String, Int ...
		addCategoria := Denuncias{}

		// retornoSelectBanco.Scan varre o objeto retornoSelectBanco e salva os valores na variaveis citadas abaixo.
		// As variaveis são salvas na mesma ordem que são coletados do Banco de dados
		if erro := retornoSelectBanco.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Total); erro != nil {
			log.Println("erro ao salvar as categoriasFull retornados do Banco:", erro.Error())
		}

		// adiciona na struct principal os dados do banco
		// Sera está struct 'categorias' que será convertida para o formato JSON
		denuncias = append(denuncias, addCategoria)
	}
	log.Printf("Denuncias atualizadas")

	// select usado para trazer a Categoria e a quantidade de denuncias por regiao
	retornoSelectBanco, erro = bancoDeDados.Query(`SELECT d.id_categoria, c.categoria, l.regiao, COUNT(d.id_localidade) 
                                        FROM tab_denuncia d JOIN tab_categoria c 
                                        ON d.id_categoria = c.id 
                                        JOIN tab_localidade l 
                                        ON d.id_localidade = l.id 
                                        GROUP BY d.id_localidade, d.id_categoria, c.categoria, l.regiao`)

	if erro != nil {
		log.Println("erro no SELECT das denunciasPorRegiao:", erro.Error())
	}
	// select apenas para trazer o valor do ultimo ID
	// usado para um futuro POST
	ultimoIDBanco, erro := bancoDeDados.Query("select MAX(id) from tab_denuncia")
	if erro != nil {
		log.Println("erro no SELECT count categoria:", erro.Error())
	}

	// finaliza o comando para Query
	defer retornoSelectBanco.Close()
	defer ultimoIDBanco.Close()
	// fecha conexão com o Banco
	defer bancoDeDados.Close()

	for retornoSelectBanco.Next() {
		addCategoria := DenunciasPorCategoria{}
		if erro := retornoSelectBanco.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Regiao, &addCategoria.Total); erro != nil {
			log.Println("erro ao salvar as categoriasEach retornados do Banco:", erro.Error())
		}

		denunciasPorCategoria = append(denunciasPorCategoria, addCategoria)
	}
	log.Printf("Denuncias por categorias atualizadas")

	for ultimoIDBanco.Next() {

		if erro := ultimoIDBanco.Scan(&proximoIdParaGravarNoBanco); erro != nil {
			log.Println("erro ao salvar categoriasCount retornados do Banco:", erro.Error())
		} else {
			proximoIdParaGravarNoBanco++
		}
	}
	log.Printf("Ultimo ID atualizado: %d", proximoIdParaGravarNoBanco)
}

func main() {

	atualizarCategorias()
	router := mux.NewRouter()

	router.HandleFunc("/denuncias/", pegarTodasCategorias).Methods("GET")   // JSON com todas as categorias
	router.HandleFunc("/denuncias/{uri}", pegarUmaCategoria).Methods("GET") // devolve apenas uma categoria
	router.HandleFunc("/denuncias/", gravarNovaDenuncia).Methods("POST")    // adiciona nova denuncia

	log.Fatal(http.ListenAndServe(":8080", router)) // Server na porta 8080 [ localhost:8080 ]
}
