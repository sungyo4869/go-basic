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
	if r.Method == http.MethodPost {
		var req model.CreateTODORequest
		var res model.CreateTODOResponse

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		if req.Subject == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		result , err := t.svc.CreateTODO(r.Context(), req.Subject, req.Description)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		res.TODO = *result

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
    }else if r.Method == http.MethodPut {
		var req model.UpdateTODORequest
		var res model.UpdateTODOResponse
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		result, err := t.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
		if err != nil{
			http.Error(w, "Bad Request", 400)
			return
		}

		res.TODO = *result
		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
	}else if r.Method == http.MethodGet {
		var req model.ReadTODORequest
		var res model.ReadTODOResponse
		var err error
		max_row := int64(5)

		params := r.URL.Query()

		// 特定のクエリパラメータの値を取得
		pramStr := params.Get("prev_id")

		if pramStr != "" {
			req.PrevID, err = strconv.ParseInt(pramStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}

		pramStr = params.Get("size")

		if pramStr != "" {
			req.Size, err = strconv.ParseInt(pramStr, 10, 64)
			if err != nil {
				http.Error(w, "Bad Request", 400)
				return
			}
		}else {
			req.Size = max_row
		}

		todos, err := t.svc.ReadTODO(r.Context(), req.PrevID, req.Size)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		res.TODOs = []model.TODO{}

		for _, todo := range todos {
			res.TODOs = append(res.TODOs, *todo)
		}

		err = json.NewEncoder(w).Encode(&res)

		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
	}else if r.Method == http.MethodDelete{
		var req model.DeleteTODORequest
		var res model.DeleteTODOResponse

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		if len(req.IDs) == 0{
			http.Error(w, "Bad Request", 400)
			return
		}

		result := t.svc.DeleteTODO(r.Context(), req.IDs)

		if errors.Is(result, &model.ErrNotFound{}){
			http.Error(w, "NotFound", 404)
			return
		}

		err := json.NewEncoder(w).Encode(&res)

		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
	}
}


// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}