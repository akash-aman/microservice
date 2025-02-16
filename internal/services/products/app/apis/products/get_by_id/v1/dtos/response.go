package dtos

type GetProductByIdResponseDto struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func NewGetProductByIdResponseDto(id, name, description string, price float64, stock int) *GetProductByIdResponseDto {
	return &GetProductByIdResponseDto{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}
}
