package model

type GetProductById struct {
	ID string `validate:"required,gte=2,lte=50"`
}

func NewGetProductById(id string) *GetProductById {
	return &GetProductById{
		ID: id,
	}
}
