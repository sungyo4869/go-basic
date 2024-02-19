package router

import (
	"database/sql"
	"net/http"

	"github.com/sungyo4869/go-basic/handler"
	"github.com/sungyo4869/go-basic/handler/middleware"
	"github.com/sungyo4869/go-basic/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	// エンドポイントの定義
	mux.HandleFunc("/healthz", middleware.Recovery(middleware.StoreOSName(middleware.Log(handler.NewHealthzHandler()))).ServeHTTP)
	mux.HandleFunc("/todos", middleware.Recovery(middleware.StoreOSName(middleware.Log(handler.NewTODOHandler(service.NewTODOService(todoDB))))).ServeHTTP)
	mux.HandleFunc("/do-panic", middleware.Recovery(middleware.StoreOSName(middleware.Log(handler.NewPanicHandler()))).ServeHTTP)

	return mux
}
