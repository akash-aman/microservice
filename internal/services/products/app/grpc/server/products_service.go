package server

/**
 * * For Concurrency
 * - Please go through below docs to understand how concurrency plays its role.
 * - https://github.com/grpc/grpc-go/blob/master/Documentation/concurrency.md
 * 
 * * For Encoding:
 * - https://github.com/grpc/grpc-go/blob/master/Documentation/encoding.md
 */

import (
	"context"
	"io"
	"pkg/logger"
	products_service "products/app/grpc/server/proto"
	"products/conf"
	"time"
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

func (s *ProductGrpcServerService) BiDiStreamProducts(stream products_service.ProductService_BiDiStreamProductsServer) error {
	s.log.Info(stream.Context(), "Starting bidirectional stream")

	for {
		product, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// Process the product and send back modified version
		modifiedProduct := &products_service.Product{
			Id:          product.Id,
			Name:        "Modified: " + product.Name,
			Description: product.Description,
			Price:       product.Price * 1.1, // 10% price increase
		}

		time.Sleep(2 * time.Second)
		if err := stream.Send(modifiedProduct); err != nil {
			return err
		}
	}
}

func (s *ProductGrpcServerService) ClientStreamProducts(stream products_service.ProductService_ClientStreamProductsServer) error {
	s.log.Info(stream.Context(), "Starting client stream")

	var products []*products_service.Product
	for {
		product, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		products = append(products, product)
	}

	// Process all received products and send back the list
	return stream.SendAndClose(&products_service.ProductList{
		Products: products,
	})
}

func (s *ProductGrpcServerService) ServerStreamProducts(req *products_service.Product, stream products_service.ProductService_ServerStreamProductsServer) error {
	s.log.Info(stream.Context(), "Starting server stream")

	// Generate variations of the product
	variations := []*products_service.Product{
		{
			Id:          req.Id + "-1",
			Name:        req.Name + " (Small)",
			Description: req.Description,
			Price:       req.Price * 0.8,
		},
		{
			Id:          req.Id + "-2",
			Name:        req.Name + " (Medium)",
			Description: req.Description,
			Price:       req.Price,
		},
		{
			Id:          req.Id + "-3",
			Name:        req.Name + " (Large)",
			Description: req.Description,
			Price:       req.Price * 1.2,
		},
	}

	for _, product := range variations {
		time.Sleep(2 * time.Second)
		if err := stream.Send(product); err != nil {
			return err
		}
	}

	return nil
}

func (s *ProductGrpcServerService) UnaryStreamProducts(ctx context.Context, req *products_service.Product) (*products_service.Product, error) {
	s.log.Infof(ctx, "Processing unary stream product: %s", req.Name)

	// Process the product and return modified version
	return &products_service.Product{
		Id:          req.Id,
		Name:        "Processed: " + req.Name,
		Description: req.Description,
		Price:       req.Price * 1.05, // 5% price increase
	}, nil
}
