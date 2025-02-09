package dtos

type ProductResponseDto struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func NewProductResponseDto(id, name, description string, price float64, stock int) *ProductResponseDto {
	return &ProductResponseDto{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}
}
