package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LogFormat struct {
	Timestamp 	time.Time	`json:"timestamp"`
	Latency		int64		`json:"latency"`
	Path		string		`json:"Path"`
	OS			string		`json:"os"`
}

func Log(h http.Handler) http.Handler {
fn := func(w http.ResponseWriter, r *http.Request) {
	accessTime := time.Now()
	
	h.ServeHTTP(w, r)

	endTime := time.Now()

	latency := endTime.Sub(accessTime)
	path := r.URL.String()
	os := r.Context().Value(keyOS).(string)

	logData := LogFormat{
		Timestamp: accessTime,
		Latency: int64(latency),
		Path: path,
		OS: os,
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		fmt.Println("構造体をJSONに変換できませんでした:", err)
		return
	}
	// fmt.Println(logData)
	fmt.Println(string(jsonData))

}
return http.HandlerFunc(fn)
}