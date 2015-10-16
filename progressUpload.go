package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func uploadPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("<html><title>Go upload</title><body><form action='http://localhost:8080/receive' method='post' enctype='multipart/form-data'><label for='file'>Filename:</label><input type='file' name='file' id='file'><input type='submit' value='Upload' ></form></body></html>")))
}

func uploadProgress(w http.ResponseWriter, r *http.Request) {

	mr, err := r.MultipartReader()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	length := r.ContentLength

	//ticker := time.Tick(time.Millisecond) // <-- use this in production
	ticker := time.Tick(time.Second) // this is for demo purpose with longer delay

	for {

		var read int64
		var p float32
		part, err := mr.NextPart()

		if err == io.EOF {
			fmt.Printf("\nDone!")
			break
		}

		dst, err := os.OpenFile("upload.jpg", os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			return
		}

		for {

			buffer := make([]byte, 100000)
			cBytes, err := part.Read(buffer)
			if err == io.EOF {
				fmt.Printf("\nLast buffer read!")
				break
			}
			read = read + int64(cBytes)

			//fmt.Printf("\r read: %v  length : %v \n", read, length)

			if read > 0 {
				p = float32(read*100) / float32(length)
				//fmt.Printf("progress: %v \n", p)
				<-ticker
				fmt.Printf("\rUploading progress %v", p) // for console
				dst.Write(buffer[0:cBytes])
			} else {
				break
			}

		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", uploadPage)
	mux.HandleFunc("/receive", uploadProgress)

	http.ListenAndServe(":8080", mux)
}
