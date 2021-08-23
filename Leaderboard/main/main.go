package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const MAX =5
var names=make(map[int]info)
var idx=0

type info struct {
	name string
	shares [MAX]int
	total int
}

type Pair struct {
	Key   int
	Value int
}
type PairList []Pair

func compute(i int,values []int) (total int){
	total=0
	for j:=0;j<MAX;j++{
		total+=names[i].shares[j]*values[j]
	}
	return
}

func register(name string,shares [MAX]int){
	names[idx]=info{name,shares,0}
	//total:=compute(idx,values)
	//names[idx]=info{name,shares,total}
	idx++
}

func read_values() []int{
	topic := "stocks"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10*time.Second))
	batch := conn.ReadBatch(200, 1e6) // fetch 200B min, 1MB max
	var values string
	for {
		b := make([]byte, 200) // 200B max per message
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		values=string(b)
	}
	//fmt.Println(values)

	err=batch.Close()
	if err!=nil{
		log.Fatal(err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}

	var vals = []int{}
	var value=strings.Split(values,"\n")
	value=strings.Fields(value[0])

	for _,val:= range value{
		j,_:=strconv.Atoi(val)
		if err!=nil{
			panic("Hello!")
		}
		vals = append(vals,j)
	}
	return vals
	//consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
	//	"bootstrap.servers":    "host1:9092,host2:9092",
	//	"group.id":             "foo",
	//	"auto.offset.reset":    "smallest"})
	//if err!=nil{
	//	return
	//}
	//run := true
	//for run == true {
	//	ev := consumer.Poll(0)
	//	switch e := ev.(type) {
	//	case *kafka.Message:
	//		// application-specific processing
	//		fmt.Println(ev)
	//	case kafka.Error:
	//		fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
	//		run = false
	//	default:
	//		fmt.Printf("Ignored %v\n", e)
	//	}
	//}
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func main(){
	var name string
	var values = []int{}
	for {

		fmt.Println("Enter choice\t1. Register\t2. Display people\tAny other key to Exit")
		var ch int
		fmt.Scanf("%d", &ch)
		switch ch {
		case 1:
			fmt.Println("Enter name: ")
			fmt.Scanf("%s", &name)
			var shares = [MAX]int{}
			fmt.Println("Enter share units for companies:")
			for i := 0; i < MAX; i++ {
				fmt.Printf("Company %d ", i)
				fmt.Scanf("%d", &shares[i])
			}
			register(name, shares)
			fmt.Println("Registered")
		case 2:
			values = read_values()
			for i := 0; i < idx; i++ {
				total := compute(i, values)
				names[i] = info{names[i].name, names[i].shares, total}
			}

			var mp = make(map[int]int)
			for i := 0; i < idx; i++ {
				mp[i] = names[i].total
			}
			p := make(PairList, len(mp))
			i := 0
			for k, v := range mp {
				p[i] = Pair{k, v}
				i++
			}
			sort.Sort(sort.Reverse(p))

			fmt.Printf("Name\t\t")
			for i := 0; i < MAX; i++ {
				fmt.Printf("Company %d\t", i)
			}
			fmt.Println("Total")
			for _, k := range p {
				fmt.Printf("%s\t\t", names[k.Key].name)
				for j := 0; j < MAX; j++ {
					fmt.Printf("%d\t\t", names[k.Key].shares[j])
				}
				fmt.Println(names[k.Key].total)
			}
		default:
			os.Exit(0)
		}
	}
}
