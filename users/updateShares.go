package users

import (
	"encoding/json"
	"github.com/golang/gddo/httputil/header"
	"net/http"
	"strconv"
	"strings"
)

func UpdateShares(names map[int]Info, res http.ResponseWriter, req *http.Request){

	if req.Method == "POST"{
		//Checks the content type of request
		if req.Header.Get("Content-Type") != "" {
			value, _ := header.ParseValueAndParams(req.Header, "Content-Type")
			if value != "application/json" {
				msg := "Content-Type header is not application/json"
				http.Error(res, msg, http.StatusUnsupportedMediaType)
				return
			}
		}

		header:=res.Header()
		header.Set("Content-Type","application/json")

		id,err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/leaderboard/v1/users/"));if err != nil{
			msg:= "ID can only be integer"
			http.Error(res,msg,http.StatusUnprocessableEntity)
			return
		}
		if names[id].Name == ""{
			msg := "User not registered"
			http.Error(res,msg,http.StatusBadRequest)
		}else {

			var name Info
			err = json.NewDecoder(req.Body).Decode(&name)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}

			values := readValues()
			names[id] = Info{names[id].Name,name.Shares,compute(names,id,values)}
			res.WriteHeader(http.StatusOK)
		}

	} else if req.Method == "GET"{

		header := res.Header()
		header.Set("Content-Type","application/json")

		id,err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/leaderboard/v1/users/"));if err != nil{
			msg := "ID can only be integer"
			http.Error(res,msg,http.StatusUnprocessableEntity)
			return
		}
		if names[id].Name == "" {
			msg := "User not registered"
			http.Error(res,msg,http.StatusBadRequest)
		}else {
			err = json.NewEncoder(res).Encode(names[id])
		}
	}
}