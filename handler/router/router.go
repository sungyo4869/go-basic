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
	mux.HandleFunc("/healthz", handler.NewHealthzHandler().ServeHTTP)
	mux.HandleFunc("/todos", handler.NewTODOHandler(service.NewTODOService(todoDB)).ServeHTTP)
	mux.HandleFunc("/do-panic", middleware.Recovery(handler.NewPanicHandler()).ServeHTTP)
	mux.HandleFunc("/test", middleware.StoreOSName(middleware.Log(handler.NewHealthzHandler())).ServeHTTP)

	return mux
}
