package users

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
)

type Pair struct {
	Key   int
	Value int
}

type PairList []Pair

func compute(names map[int]Info,i int,values []int) (total int){
	total = 0
	for j := 0; j < MAX; j++{
		total += names[i].Shares[j]*values[j]
	}
	return
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func GetLeaderboard(names map[int]Info,idx int,res http.ResponseWriter,_ *http.Request){
	header := res.Header()
	header.Set("Content-Type","application/json")

	var values = []int{}
	values = readValues()
	for i := 0; i < idx; i++ {
		total := compute(names, i, values)
		names[i] = Info{names[i].Name, names[i].Shares, total}
	}

	var mp = make(map[int]int)
	for i := 0; i < idx; i++ {
		mp[i] = names[i].Total
	}
	p := make(PairList, len(mp))
	i := 0
	for k, v := range mp {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))

	for _, k := range p {
		name := names[k.Key].Name
		shares := names[k.Key].Shares
		total := names[k.Key].Total
		data := Info{name,shares,total}

		err := json.NewEncoder(res).Encode(data);if err != nil{
			log.Fatal(err)
		}
	}
}
