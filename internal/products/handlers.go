package products

import (
	"log"
	"net/http"

	"github.com/VittorioDeMarzi/Ecommerce-project-Golang-/internal/products/json"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// 1. call the service -> ListProducts
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println("err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 2. Return JSON in a HTTP response
	json.Write(w, http.StatusOK, products)
}
