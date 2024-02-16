package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"errors"

	"github.com/sungyo4869/go-basic/model"
	"github.com/sungyo4869/go-basic/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (t *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method{
	case http.MethodPost:
		var req model.CreateTODORequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if req.Subject == "" || err != nil{
			// Subjectが空またはデコードに失敗した場合
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		
		// DBに格納
		res, err := t.Create(r.Context(), &req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

	case http.MethodPut:
		var req model.UpdateTODORequest

		err := json.NewDecoder(r.Body).Decode(&req);
		if req.ID == 0 || req.Subject == "" || err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		res, err := t.Update(r.Context(), &req)
		if err != nil{
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

	case http.MethodGet:
		var req model.ReadTODORequest
		var err error
		max_row := int64(5)

		params := r.URL.Query()
		// クエリパラメータの値を取得
		pramStr := params.Get("prev_id")
		if pramStr != "" {
			// idが指定されていた場合、intにパース
			req.PrevID, err = strconv.ParseInt(pramStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		// クエリパラメータの値を取得
		pramStr = params.Get("size")
		if pramStr != "" {
			// サイズが指定されていた場合
			req.Size, err = strconv.ParseInt(pramStr, 10, 64)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
		}else {
			// 指定されていない場合はデフォルト値を設定
			req.Size = max_row
		}

		res, err := t.Read(r.Context(), &req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

	case http.MethodDelete:
		var req model.DeleteTODORequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if len(req.IDs) == 0 || err != nil{
			// idが指定されていないまたはデコードに失敗した場合
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		res, err := t.Delete(r.Context(), &req)
		if err != nil{
			http.Error(w, "NotFound", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	var res model.CreateTODOResponse

	result , err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	res.TODO = *result

	return &res, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {

	var res model.ReadTODOResponse
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}

	res.TODOs = []model.TODO{}
	for _, todo := range todos {
		res.TODOs = append(res.TODOs, *todo)
	}

	return &res, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	var res model.UpdateTODOResponse

	result, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	res.TODO = *result

	return &res, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	var res model.DeleteTODOResponse

	err := h.svc.DeleteTODO(ctx, req.IDs)
	if errors.Is(err, &model.ErrNotFound{}){
		return nil, model.ErrNotFound{}
	}

	return &res, nil
}