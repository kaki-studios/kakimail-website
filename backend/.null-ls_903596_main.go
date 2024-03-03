package main

import (
	"html/template"
	"log"
	"net/http"
)

func main()  {
  http.HandleFunc("/", indexHandler)
  log.Fatal(http.ListenAndServe(":8000", nil));

}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
  tmpl := template.Must(template.ParseFiles("../frontend/index.html"));
  tmpl.Execute(writer, nil)


}
