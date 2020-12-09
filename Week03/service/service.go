package service

import (
	"net/http"
)

func Service(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Service1"))
}

