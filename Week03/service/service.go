package service

import (
	"net/http"
)

func Service1(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Service1"))
}

func Service2(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Service2"))
}

