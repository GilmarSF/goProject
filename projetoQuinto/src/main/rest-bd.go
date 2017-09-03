// Todo codigo em Go precisa iniciar com este package
package main

// Importa todas as bibliotecas que serão utilizadas no código
import (

	// bibliotecas Nativas do Golang
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// Pode ser usado bibliotecas não nativas
	// /gorilla/mux para o servidor HTTP
	// /denisenkom/go-mssqldb para a comunicação com Banco SQLServer
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

// Struct Principal criada para armazenar os dados retornados do Banco de dados
// Structs são mais comuns para converter para JSON em Golang
type Usuario struct {

	// "omitempty", quando usado, omite este atributo caso seja 'null'
	ID    string `json:"id,omitempty"`
	Nome  string `json:"nome,omitempty"`
	Email string `json:"email,omitempty"`
	Senha string `json:"senha,omitempty"`
}

type CategoriaFull struct {
	ID    string `json:"id,omitempty"`
	Nome  string `json:"nome,omitempty"`
	Total string `json:"total,omitempty"`
}

type CategoriaEach struct {
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
// array para salvar os Usuarios
var cadastros []Usuario

// array usado para enviar o total de cada categoria
var categorias []CategoriaFull

// array usado para enviar o total de denuncias por regiao
var categoriasPorRegiao []CategoriaEach

// Usado para armazenar o ultimo 'id' do banco de dados
var countBanco int
var proximoIdParaGravarNoBanco int

// Função utilizada enviar apenas um Usuario via Rest
func GetUsuario(w http.ResponseWriter, req *http.Request) {

	// Um aviso qualquer que será impresso no "prompt"
	log.Printf("Get Usuario")

	// Salva o 'id' (valor) passado na requisição
	params := mux.Vars(req)

	// For usado para varrer todos os Usuarios que estão salvos no Array 'cadastros'
	// 'range' significo por todo o array
	// Neste for esta sendo criado uma Struct 'item' com os dados do Usuario que possuir o mesmo ID da requisição
	// o sinal [ := ] é usado para criar uma nova variavel do tipo do valor que será adicionado a ela
	for _, item := range cadastros {

		// Verifica em cada Usuario do array 'cadastros' se o ID é igual ao id passado na requisição
		if item.ID == params["id"] {

			// Quando encontrado, converte a Struct 'item' para o
			// formato 'JSON' e envia via REST
			json.NewEncoder(w).Encode(item)
			// sai da Func
			return
		}
	}

	// Caso não encontrar, envia um JSON vazio
	json.NewEncoder(w).Encode(&Usuario{})
}

// traz a cadastros toda. Todas todas Usuarios do array 'cadastros'
func GetCadastros(w http.ResponseWriter, req *http.Request) {
	log.Printf("Get cadastros")
	json.NewEncoder(w).Encode(cadastros) //correto cadastros
}

// Adicona mais uma denuncia
func gravarNovaDenuncia(w http.ResponseWriter, req *http.Request) {
	// modelo que deve enviado
	// {"categoria":"4","localidade":"2"}
	log.Printf("Post Nova Denuncia")
	var novaD NovaDenuncia

	// grava em 'novaD' os dados enviados
	_ = json.NewDecoder(req.Body).Decode(&novaD)
	// imprime no terminal os valores recebidos
	fmt.Println(novaD)

	bancoDeDados, erro := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
	if erro != nil {
		log.Println("erro ao conectar com o Banco de dados:", erro.Error())
	}

	rows, erro := bancoDeDados.Query("INSERT into tab_denuncia (id, id_categoria, id_localidade) VALUES (?1, ?2, ?3)", proximoIdParaGravarNoBanco, novaD.Categoria, novaD.Localidade)
	if erro != nil {
		log.Println("erro no INSERT:", erro.Error())
	} else {
		proximoIdParaGravarNoBanco++
	}

	defer rows.Close()         // fecha o comando Query
	defer bancoDeDados.Close() // fecha conexão com o Banco
	// atualiza struct no banco
	atualizarCategorias()

	//json.NewEncoder(w).Encode(categorias)
}

// cria o arquivo tabela.json no server
func CriaArquivoJSON(w http.ResponseWriter, req *http.Request) {

	// Busca no banco de dados todas as categorias e os totais de reclamações de novo
	// Assim toda vez que chamar essa func os dados estarão sempre att
	atualizarCategorias()
	log.Printf("Tabela.json atualizada")

	// Lê todo conteudo do arquivo default.json e salva em 'jsonOut'
	jsonOut, erro := ioutil.ReadFile("default.json")
	if erro != nil {
		fmt.Println(erro)
	}

	// for usado para varrer todos os dados Salvos no array d structs 'categorias'
	// a cada dado encontrado, é sobrescrito no arquivo default.json
	for _, item := range categorias {

		// Salva em 'attCategoria' o conteudo copiado de 'default.json'
		// já com a palavra "Categoria" substituida por 'item.Nome'
		// 'item' recebe do 'for' citado acima os dados da posição atual do array de struct
		// Obs.: O numero 1 no final significa que é para alterar apenas a primeira
		// palavra "Categoria encontrada", se estivesse 2 substituiria as 2 primeiras
		// palavras "Categoria" encontrada, se estivesse -1 substituiria todas as palavras
		// "Categoria" para a palavra salva no atual item.Nome
		attCategoria := bytes.Replace(jsonOut, []byte("Categoria"), []byte(item.Nome), 1)

		// Nesta parte sera gravado todo o conteudo de attCategoria com a alteração acima
		// no arquivo tabela.json
		if erro = ioutil.WriteFile("../../bin/pages/tabela.json", attCategoria, 0666); erro != nil {
			fmt.Println(erro)
		}

		// Lê todo conteudo do arquivo tabela.json e salva em 'jsonOut'
		// a partir daqui cada alteração no arquivo 'tabela.json' o mesmo deve
		// ser relido e salvo em 'jsonOut' para as proximas alterações não serem
		// baseadas no arquivo 'default.json'
		jsonOut, erro = ioutil.ReadFile("../../bin/pages/tabela.json")
		if erro != nil {
			fmt.Println(erro)
		}

		// Faz o mesmo descrito acima para o numero total de incidentes
		// Obs.: apesar de ser escrito no formato string
		// o Javascript na hora de ler entendera que é um numero INT pois não tem " "
		attTotal := bytes.Replace(jsonOut, []byte("0"), []byte(item.Total), 1)
		if erro = ioutil.WriteFile("../../bin/pages/tabela.json", attTotal, 0666); erro != nil {
			fmt.Println(erro)
		}

		jsonOut, erro = ioutil.ReadFile("../../bin/pages/tabela.json")
		if erro != nil {
			fmt.Println(erro)
		}
	}

	out, _ := ioutil.ReadFile("tabela.json")
	//if erro != nil {
	//    fmt.Println(erro)
	//}

	//texto := []byte("conteudo")
	//erro = ioutil.WriteFile("../../bin/pages/teste.json.html", out, 0644)
	//if erro != nil {
	//    fmt.Println(erro)
	//}

	//fmt.Println(out)
	json.NewEncoder(w).Encode(out)
}

// função para enviar apenas uma categoria com o total por regiao
func pegarUmaCategoria(w http.ResponseWriter, req *http.Request) {
	// OBSERVAÇÂO: comentarios de como funciona esta na 'func GetUsuario'
	log.Printf("Get uma Categoria")
	params := mux.Vars(req)
	var categoriaFound []CategoriaEach
	for _, item := range categoriasPorRegiao {
		if item.ID == params["id"] {
			categoriaFound = append(categoriaFound, item)
		}
	}
	//json.NewEncoder(w).Encode(categoriasPorRegiao)
	json.NewEncoder(w).Encode(categoriaFound)
}

// envia os dados das categorias via GET
func pegarTodasCategorias(w http.ResponseWriter, req *http.Request) {
	log.Printf("Get categorias")
	json.NewEncoder(w).Encode(categorias)
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

	// Select traz todos os Nomes das categorias e as quantidades de denuncias de cada um
	denuciasPorCategoria, erro := bancoDeDados.Query(`SELECT d.id_categoria, c.categoria, COUNT(d.id_categoria) 
                                            FROM tab_denuncia d JOIN tab_categoria c 
                                            ON d.id_categoria = c.id 
                                            GROUP BY d.id_categoria, c.categoria`)

	if erro != nil {
		log.Println("erro no SELECT das CategoriasFull:", erro.Error())
	}

	// select usado para trazer a Categoria e a quantidade de denuncias por regiao
	denuciasPorRegiao, erro := bancoDeDados.Query(`SELECT d.id_categoria, c.categoria, l.regiao, COUNT(d.id_localidade) 
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
	defer denuciasPorCategoria.Close()
	defer denuciasPorRegiao.Close()
	defer ultimoIDBanco.Close()
	// fecha conexão com o Banco
	defer bancoDeDados.Close()

	// Zera o Array antes de buscar novos dados no banco.
	// Dessa forma evita dados repetidos
	categorias = categorias[:0]
	categoriasPorRegiao = categoriasPorRegiao[:0]

	// denuciasPorCategoria.Next usado para varrer o objeto 'denuciasPorCategoria' e pegar os valores retornados da Query
	for denuciasPorCategoria.Next() {
		// Struct criada para receber os dados do banco de dados
		// Uma Struct permite receber dados de diferentes tipos: String, Int ...
		addCategoria := CategoriaFull{}

		// denuciasPorCategoria.Scan varre o objeto denuciasPorCategoria e salva os valores na variaveis citadas abaixo.
		// As variaveis são salvas na mesma ordem que são coletados do Banco de dados
		if erro := denuciasPorCategoria.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Total); erro != nil {
			log.Println("erro ao salvar as categoriasFull retornados do Banco:", erro.Error())
		}

		// adiciona na struct principal os dados do banco
		// Sera está struct 'cadastros' que será convertida para o formato JSON
		categorias = append(categorias, addCategoria)

		// Pega o valor do ID do ultimo dado buscado no banco e converte de String para Int
	}

	for denuciasPorRegiao.Next() {
		addCategoria := CategoriaEach{}
		if erro := denuciasPorRegiao.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Regiao, &addCategoria.Total); erro != nil {
			log.Println("erro ao salvar as categoriasEach retornados do Banco:", erro.Error())
		}

		categoriasPorRegiao = append(categoriasPorRegiao, addCategoria)
	}

	for ultimoIDBanco.Next() {

		if erro := ultimoIDBanco.Scan(&proximoIdParaGravarNoBanco); erro != nil {
			log.Println("erro ao salvar categoriasCount retornados do Banco:", erro.Error())
		} else {
			proximoIdParaGravarNoBanco++
		}
	}

	log.Printf("Categorias atualizadas!")
}

func main() {

	atualizarCategorias()
	router := mux.NewRouter()

	router.HandleFunc("/categorias/", pegarTodasCategorias).Methods("GET")  // JSON com todas as categorias
	router.HandleFunc("/categorias/{id}", pegarUmaCategoria).Methods("GET") // devolve apenas uma categoria
	router.HandleFunc("/categorias/", gravarNovaDenuncia).Methods("POST")   // adiciona nova categoria

	log.Fatal(http.ListenAndServe(":8080", router)) // Server na porta 8080 [ localhost:8080 ]
}
