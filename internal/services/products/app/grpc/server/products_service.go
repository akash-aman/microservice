package server

import (
	"context"
	"pkg/logger"
	products_service "products/app/grpc/server/proto"
	"products/conf"
)

type ProductGrpcServerService struct {
	products_service.UnimplementedProductServiceServer
	log logger.Zapper
	cfg conf.Config
}

func NewProductGrpcServerService(log logger.Zapper, cfg conf.Config) *ProductGrpcServerService {
	return &ProductGrpcServerService{
		log: log,
		cfg: cfg,
	}
}

func (s *ProductGrpcServerService) GetProduct(ctx context.Context, req *products_service.GetProductRequest) (*products_service.GetProductResponse, error) {
	s.log.Infof(ctx, "Getting product with ID: %s", req.Id)

	// TODO: Implement database query
	product := &products_service.Product{
		Id:          req.Id,
		Name:        "Sample Product",
		Description: "This is a sample product",
		Price:       99.99,
	}

	return &products_service.GetProductResponse{
		Product: product,
	}, nil
}

func (s *ProductGrpcServerService) ListProducts(ctx context.Context, req *products_service.ListProductsRequest) (*products_service.ListProductsResponse, error) {
	s.log.Info(ctx, "Listing all products")

	// TODO: Implement database query with pagination
	products := []*products_service.Product{
		{
			Id:          "1",
			Name:        "Product 1",
			Description: "First product",
			Price:       99.99,
		},
		{
			Id:          "2",
			Name:        "Product 2",
			Description: "Second product",
			Price:       149.99,
		},
	}

	return &products_service.ListProductsResponse{
		Products: products,
		Total:    int32(len(products)),
	}, nil
}

func (s *ProductGrpcServerService) CreateProduct(ctx context.Context, req *products_service.CreateProductRequest) (*products_service.CreateProductResponse, error) {
	s.log.Infof(ctx, "Creating new product: %s", req.Name)

	// TODO: Implement database insert
	product := &products_service.Product{
		Id:          "new-id", // Should be generated
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	return &products_service.CreateProductResponse{
		Product: product,
	}, nil
}

func (s *ProductGrpcServerService) UpdateProduct(ctx context.Context, req *products_service.UpdateProductRequest) (*products_service.UpdateProductResponse, error) {
	s.log.Infof(ctx, "Updating product with ID: %s", req.Id)

	// TODO: Implement database update
	product := &products_service.Product{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	return &products_service.UpdateProductResponse{
		Product: product,
	}, nil
}

func (s *ProductGrpcServerService) DeleteProduct(ctx context.Context, req *products_service.DeleteProductRequest) (*products_service.DeleteProductResponse, error) {
	s.log.Infof(ctx, "Deleting product with ID: %s", req.Id)

	// TODO: Implement database delete
	return &products_service.DeleteProductResponse{
		Id: req.Id,
	}, nil
}
