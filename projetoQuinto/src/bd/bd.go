package bd

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type Denuncias_Struct struct {
	ID    string `json:"id,omitempty"`
	Nome  string `json:"nome,omitempty"`
	Total string `json:"total,omitempty"`
}

type DenunciasPorCategoria_Struct struct {
	ID     string `json:"id,omitempty"`
	Nome   string `json:"nome,omitempty"`
	Regiao string `json:"regiao,omitempty"`
	Total  string `json:"total,omitempty"`
}

type NovaDenuncia_Struct struct {
	Categoria  string `json:"categoria,omitempty"`
	Localidade string `json:"localidade,omitempty"`
}

/* VARIAVEIS GLOBAIS PODEM SER USADAS EM QUALQUER PARTE DO CÓDiGO*/

// array usado para enviar o total de cada categoria
var Denuncias []Denuncias_Struct

// array usado para enviar o total de denuncias por regiao
var DenunciasPorCategoria []DenunciasPorCategoria_Struct

// Usado para armazenar o ultimo 'id' do banco de dados
var proximoIdParaGravarNoBanco int

// Atualiza a struct Categorias com os dados do Banco
func AtualizarDenuncias() {

	// Inicia conexão com Banco de dados no Azure
	bancoDeDados, erro := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
	// Caso ocorrer algum erro na conexão, o erro será salvo em 'erro'
	// Caso não houver erro o retorno será vazio 'nil'
	if erro != nil {
		log.Println("erro em atualizarDenuncias(): ", erro.Error())
	}

	// Zera o Array antes de buscar novos dados no banco.
	// Dessa forma evita dados repetidos
	Denuncias = Denuncias[:0]
	DenunciasPorCategoria = DenunciasPorCategoria[:0]

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
		addCategoria := Denuncias_Struct{}

		// retornoSelectBanco.Scan varre o objeto retornoSelectBanco e salva os valores na variaveis citadas abaixo.
		// As variaveis são salvas na mesma ordem que são coletados do Banco de dados
		if erro := retornoSelectBanco.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Total); erro != nil {
			log.Println("erro ao salvar as categoriasFull retornados do Banco:", erro.Error())
		}

		// adiciona na struct principal os dados do banco
		// Sera está struct 'categorias' que será convertida para o formato JSON
		Denuncias = append(Denuncias, addCategoria)
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
		addCategoria := DenunciasPorCategoria_Struct{}
		if erro := retornoSelectBanco.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Regiao, &addCategoria.Total); erro != nil {
			log.Println("erro ao salvar as categoriasEach retornados do Banco:", erro.Error())
		}

		DenunciasPorCategoria = append(DenunciasPorCategoria, addCategoria)
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
