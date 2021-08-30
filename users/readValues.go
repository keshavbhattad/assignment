package users

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
	"strings"
	"time"
)

func readValues() []int{
	topic := "leaderboard"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10*time.Second))
	batch := conn.ReadBatch(50, 1e6) //  min, max
	var values string
	for {
		b := make([]byte, 50) // max per message
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		values = string(b)
	}

	err = batch.Close()
	if err != nil{
		log.Fatal(err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
	var vals = []int{}
	var value = strings.Split(values,"\n")
	value = strings.Fields(value[0])

	for _,val := range value{
		j,_ := strconv.Atoi(val)
		if err != nil{
			panic("Shares should be integers")
		}
		vals = append(vals,j)
	}
	return vals
}