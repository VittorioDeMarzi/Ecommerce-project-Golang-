package products

import (
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 2. Return JSON in a HTTP response
	json.Write(w, http.StatusOK, products)
}

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request, id int64) {
	// 1. call the service -> GetProduct
	// 2. Return JSON in a HTTP response
	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		json.Write(w, http.StatusNotFound, map[string]string{"error": "Product not found"})
		return
	}
	json.Write(w, http.StatusOK, product)
}
