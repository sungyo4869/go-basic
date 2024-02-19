package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LogFormat struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency"`
	Path      string    `json:"Path"`
	OS        string    `json:"os"`
}

func Log(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessTime := time.Now()
		log.Print(accessTime)

		h.ServeHTTP(w, r)

		endTime := time.Now()

		latency := endTime.Sub(accessTime).Microseconds()
		path := r.URL.String()
		os := r.Context().Value(ctxKeyOS{}).(string)

		logData := LogFormat{
			Timestamp: accessTime,
			Latency:   int64(latency),
			Path:      path,
			OS:        os,
		}
		
		jsonData, err := json.Marshal(logData)
		if err != nil {
			fmt.Println("log: unable to convert struct to JSON, err=", err)
			return
		}

		fmt.Println(string(jsonData))
	}
	return http.HandlerFunc(fn)
}
