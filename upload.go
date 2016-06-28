package main

import (

	//	"encoding/json"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"

	"github.com/dchest/uniuri"
	//	"mime/multipart"
	"net/http"
	"os"
)

const (
	LENGTH     = 6
	PORT       = ":8080"
	DIRECTORY  = "/tmp/"
	UPADDRESS  = "http://localhost"
	dbUSERNAME = ""
	dbNAME     = ""
	dbPASSWORD = ""
	DATABASE   = dbUSERNAME + ":" + dbPASSWORD + "@/" + dbNAME + "?charset=utf8"
)

type Result struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Hash string `json:"hash"`
	Size int64  `jason:"size"`
}

type Response struct {
	Success     bool     `json:"success"`
	ErrorCode   int      `json:"errorcode,omitempty"`
	Description string   `json:"description,omitempty"`
	Files       []Result `json:"files,omitempty"`
}

func generateName() string {
	name := uniuri.NewLen(LENGTH)
	db, err := db.Open("mysql", DATABSE)
	return name
}
func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()

	resp := Response{Files: []Result{}}
	if err != nil {
		resp.ErrorCode = http.StatusInternalServerError
		resp.Description = err.Error()
		resp.Success = false
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}
		s := generateName()
		extName := filepath.Ext(part.FileName())
		filename := s + extName
		dst, err := os.Create(DIRECTORY + filename)
		defer dst.Close()

		if err != nil {
			resp.ErrorCode = http.StatusInternalServerError
			resp.Description = err.Error()
			return
		}

		h := sha1.New()
		t := io.TeeReader(part, h)
		_, err = io.Copy(dst, t)

		if err != nil {
			resp.ErrorCode = http.StatusInternalServerError
			resp.Description = err.Error()
			return
		}
		hash := h.Sum(nil)
		sha1 := base64.URLEncoding.EncodeToString(hash)
		size, _ := dst.Stat()
		res := Result{
			URL:  UPADDRESS + "/" + s + extName,
			Name: part.FileName(),
			Hash: sha1,
			Size: size.Size(),
		}
		resp.Files = append(resp.Files, res)

	}
	fmt.Println(resp)
}

func main() {
	http.HandleFunc("/upload.php", uploadHandler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		panic(err)
	}

}
