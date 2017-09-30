package bd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type DadosDasDenuncias struct {
	ID     string `json:"id,omitempty"`
	Nome   string `json:"nome,omitempty"`
	Total  string `json:"total,omitempty"`
	Regiao string `json:"regiao,omitempty"`
}

type NovaDenuncia struct {
	Categoria  string `json:"categoria,omitempty"`
	Localidade string `json:"localidade,omitempty"`
}

// array usado para enviar o total de cada categoria
var Denuncias []DadosDasDenuncias

// array usado para enviar o total de denuncias por regiao
var DenunciasPorCategoria []DadosDasDenuncias

// Usado para armazenar o ultimo 'id' do banco de dados
var proximoIdParaGravarNoBanco int

var (
	servidor = "pwbt.database.windows.net"
	usuario  = "admin-jose"
	senha    = "123abc!@#"
	banco    = "PWBT"
	porta    = 1433
)

var stringDeConexao = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%d", servidor, usuario, senha, banco, porta)

var bancoDeDados, _ = sql.Open("tipoBanco", "stringDeConexao")
var erroBD error

func GravarNovaDenuncia(nova *NovaDenuncia) {

	AbrirConexaoBanco()
	insert, erro := bancoDeDados.Query(`INSERT into tab_denuncia (id, id_categoria, id_localidade) 
										VALUES (?1, ?2, ?3)`, proximoIdParaGravarNoBanco, &nova.Categoria, &nova.Localidade)

	defer insert.Close()       // fecha comando Query
	defer bancoDeDados.Close() // fecha conexão com o Banco
	if erro != nil {
		log.Println("erro no INSERT:", erro.Error())
	} else {
		proximoIdParaGravarNoBanco++
	}
}

func AbrirConexaoBanco() {
	bancoDeDados, erroBD = sql.Open("mssql", stringDeConexao)
	if erroBD != nil {
		log.Println("erro ao conectar ao Banco: ", erroBD.Error())
	}
}

func AtualizarTodasDenuncias() {
	AtualizarDenuncias()
	AtualizarDenunciasPorCategoria()
	AtualizarUltimoIDBanco()
}

func AtualizarDenuncias() {

	AbrirConexaoBanco()

	Denuncias = Denuncias[:0]

	query := `SELECT d.id_categoria, c.categoria, COUNT(d.id_categoria) 
			FROM tab_denuncia d JOIN tab_categoria c 
			ON d.id_categoria = c.id 
			GROUP BY d.id_categoria, c.categoria`

	consultarNoBanco(query, &Denuncias, 1)

	defer bancoDeDados.Close()

	log.Printf("Denuncias atualizadas")
}

func AtualizarDenunciasPorCategoria() {

	AbrirConexaoBanco()

	DenunciasPorCategoria = DenunciasPorCategoria[:0]

	query := `SELECT d.id_categoria, c.categoria, COUNT(d.id_localidade), l.regiao
			FROM tab_denuncia d JOIN tab_categoria c 
			ON d.id_categoria = c.id 
			JOIN tab_localidade l 
			ON d.id_localidade = l.id 
			GROUP BY d.id_localidade, d.id_categoria, c.categoria, l.regiao`

	consultarNoBanco(query, &DenunciasPorCategoria, 2)

	defer bancoDeDados.Close()

	log.Printf("Denuncias por categorias atualizadas")
}

func consultarNoBanco(query string, den *[]DadosDasDenuncias, opcao int) {

	retornoSelectBanco, erro := bancoDeDados.Query(query)
	if erro != nil {
		log.Println("erro no SELECT das Categorias:", erro.Error())
	}

	for retornoSelectBanco.Next() {
		addCategoria := DadosDasDenuncias{}
		if opcao == 1 {
			if erro := retornoSelectBanco.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Total); erro != nil {
				log.Println("erro ao salvar as categoriasFull retornados do Banco:", erro.Error())
			}
		} else {
			if erro := retornoSelectBanco.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Total, &addCategoria.Regiao); erro != nil {
				log.Println("erro ao salvar as categoriasEach retornados do Banco:", erro.Error())
			}
		}
		*den = append(*den, addCategoria)
	}
}

func AtualizarUltimoIDBanco() {

	AbrirConexaoBanco()

	ultimoIDBanco, erro := bancoDeDados.Query("select MAX(id) from tab_denuncia")
	if erro != nil {
		log.Println("erro no SELECT count categoria:", erro.Error())
	}

	defer ultimoIDBanco.Close()
	defer bancoDeDados.Close()

	for ultimoIDBanco.Next() {

		if erro := ultimoIDBanco.Scan(&proximoIdParaGravarNoBanco); erro != nil {
			log.Println("erro ao salvar categoriasCount retornados do Banco:", erro.Error())
		} else {
			proximoIdParaGravarNoBanco++
		}
	}
	log.Printf("Ultimo ID atualizado: %d", proximoIdParaGravarNoBanco)
}
