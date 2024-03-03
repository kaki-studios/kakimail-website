package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
  "strings"
)

func main()  {

  http.HandleFunc("/frontend/*", staticHandler);
  http.HandleFunc("/", indexHandler)
  log.Fatal(http.ListenAndServe(":8000", nil));

}

func indexHandler(writer http.ResponseWriter, request *http.Request) {

  fmt.Println("index: r is ", request.URL.Path);
  tmpl := template.Must(template.ParseFiles("../frontend/index.html"));
  tmpl.Execute(writer, nil)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
   // load the file using r.URL.Path, like /static/scripts/index.js"
  fmt.Println("static: r is ", r.URL.Path);
  data,err:=os.ReadFile(r.URL.Path)
  if err != nil {
    fmt.Println("", err)
  }
   // Figure out file type:
  if strings.HasSuffix(r.URL.Path, "js") {
    w.Header().Set("Content-Type","text/javascript")  
  } else {
    w.Header().Set("Content-Type","text/css")
  }
   w.Write(data)
}
