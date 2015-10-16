package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"log"
	"io/ioutil"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("html/upload.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		display(w, "upload", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//get the multipart reader for the request.
		reader, err := r.MultipartReader()
//		r.ParseMultipartForm(0)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			if part.FormName() == "path" {
				j, err := ioutil.ReadAll(part)
				if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)//do something
				return
				}
			log.Println(string(j))
			//log.Println(part)
			}

			//if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}
			dst, err := os.Create("uploaded/" + part.FileName())
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		//display success message.
		display(w, "upload", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	//static file handler.
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("uploaded"))))

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
