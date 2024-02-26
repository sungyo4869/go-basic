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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
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

	mux := router.NewRouter(todoDB)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.ListenAndServe()
	}()

	select {
    case err := <-errChan:
        if err != nil && err != http.ErrServerClosed {
            return err
        }
    case <-ctx.Done():
		shutdownErrChan := make(chan error, 1)
        go func() {
        	shutdownErrChan <- srv.Shutdown(context.Background())
        }()

		err = <- shutdownErrChan;
		if err != nil {
			fmt.Println("main: Failed to shutdown server, err=", err)
		} else {
			fmt.Println("main: Server shutdown completed successfully")
		}

		wg.Done()
    }
	return nil
}
