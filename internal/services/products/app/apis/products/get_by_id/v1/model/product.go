package model

type GetProductById struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func NewGetProductById(id, name, description string, price float64, stock int) *GetProductById {
	return &GetProductById{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}
}