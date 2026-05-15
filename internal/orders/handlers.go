package orders

import (
	"net/http"

	"github.com/VittorioDeMarzi/Ecommerce-project-Golang-/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating an order
	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdOrder, err := h.service.CreateOrder(r.Context(), tempOrder)

	if err != nil {
		if err == ErrorProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)

}