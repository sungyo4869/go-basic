package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"os/signal"
	"fmt"
	"context"

	"github.com/sungyo4869/go-basic/db"
	"github.com/sungyo4869/go-basic/handler/router"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

	err := realMain(ctx)
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain(ctx context.Context) error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

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
		Addr: port,
		Handler: mux,
	}

	err = srv.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
    if err := srv.Shutdown(ctx); err != nil {
        fmt.Println("main: Failed to shutdown server, err=", err)
	}
    
	return nil
}
