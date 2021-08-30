package main

import (
	"Leaderboard/users"
	"fmt"
	"log"
	"net/http"
)

var names = make(map[int]users.Info)
var idx = 0

func fun(res http.ResponseWriter,req *http.Request){
	if req.Method == "GET"{
		users.GetLeaderboard(names,idx,res,req)
	} else if req.Method == "POST"{
		users.Register(names,&idx,res,req)
	}
}

func updateShares(res http.ResponseWriter,req *http.Request){
	users.UpdateShares(names,res,req)
}

func main() {

	http.HandleFunc("/",func(res http.ResponseWriter,req *http.Request){
		fmt.Fprint(res,"Go to \nlocalhost:9000/users to Register and access the leaderboard")
	})

	http.HandleFunc("/leaderboard/v1/users", fun)

	http.HandleFunc("/leaderboard/v1/users/",updateShares)

	log.Fatal(http.ListenAndServe(":9000", nil))
}
