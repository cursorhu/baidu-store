package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"server/meta"
	"server/util"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传的html的页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, fmt.Sprintf("读取文件错误==>%s", err))
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("fail to get data err :%s", err)
			return
		}
		defer file.Close()

		fileMeta := &meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("fail to create file err=>%s", err)
			return
		}
		defer newFile.Close()
		size, err := io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("fail to copy file err:%s", err)
			return
		}

		fileMeta.FileSize = size
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)
		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success")
}

//
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileHash)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		fmt.Println(err)
		return
	}
	io.WriteString(w, string(data))
}

func DownLoadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileHash)
	file, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-disposition", "attachment;filename="+fileMeta.FileName)
	_, _ = w.Write(data)
}

func FileMetaUpdateHandle(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	currentMeta := meta.GetFileMeta(fileSha1)
	currentMeta.FileName = newFileName
	meta.UpdateFileMeta(currentMeta)
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(currentMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	fileHash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileHash)
	_ = os.Remove(fileMeta.Location)
	meta.RemoveFileMeta(fileHash)
	w.WriteHeader(http.StatusOK)
}
