package users

import (
	"encoding/json"
	"github.com/golang/gddo/httputil/header"
	"net/http"
)

const MAX =5

type Info struct {
	Name string     `json:"name"`
	Shares [MAX]int `json:"shares"`
	Total int       `json:"total"`
}

type dummy struct{
	Id int `json:"id"`
}

func Register(users map[int]Info,idx *int,res http.ResponseWriter,req *http.Request){

	//Checks the content type of request
	if req.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(req.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(res, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	header := res.Header()
	header.Set("Content-Type","application/json")

	var name Info

	err := json.NewDecoder(req.Body).Decode(&name)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	users[*idx] = name

	json.NewEncoder(res).Encode(dummy{*idx})
	*idx++
}
