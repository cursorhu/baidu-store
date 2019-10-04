package main

import (
	"net/http"
	"server/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownLoadFile)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandle)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
