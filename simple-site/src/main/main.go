package main
 
import (
  "log"
  "net/http"
  "html/template"
)

type content struct{
	Nome string
	Idade string
}
 
func main() {
  http.HandleFunc("/", index)
  http.HandleFunc("/greet", greeter)
  log.Println("Listening...")
  err := http.ListenAndServe(":80", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
   }
}
 
func index(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("public/index.html")
  t.Execute(w, nil)
}
 
func greeter(w http.ResponseWriter, r *http.Request) {
  dados := content{}
  
  dados.Nome = r.FormValue("nome")
  dados.Idade = r.FormValue("idade")
  t, _ := template.ParseFiles("public/greeter.html")
  err := t.Execute(w, dados)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
