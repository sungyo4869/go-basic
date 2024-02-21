package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sungyo4869/go-basic/db"
	"github.com/sungyo4869/go-basic/handler/router"
)

func main() {
	wg := &sync.WaitGroup{}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	wg.Add(1)
	err := realMain(ctx, wg)
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}

}

func realMain(ctx context.Context, wg *sync.WaitGroup) error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	errChan := make(chan error, 1)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	// TODO: サーバーをlistenする
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		errChan <- srv.ListenAndServe()
	}()

	err = <-errChan
	if err != nil {
		return err
	}

	<-ctx.Done()

	err = srv.Shutdown(ctx)
	if err != nil {
		fmt.Println("main: Failed to shutdown server, err=", err)
		return err
	}
	
	wg.Done()

	return nil
}
