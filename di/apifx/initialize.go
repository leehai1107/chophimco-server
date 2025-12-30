package apifx

import (
	"github.com/leehai1107/chophimco-server/pkg/middleware/auth"
	"github.com/leehai1107/chophimco-server/service/chophimco/delivery/http"
	"github.com/leehai1107/chophimco-server/service/chophimco/repository"
	"github.com/leehai1107/chophimco-server/service/chophimco/usecase"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Provide(
	provideRouter,
	provideHandler,
	provideJWTService,

	// Repositories
	provideUserRepo,
	provideProductRepo,
	provideCartRepo,
	provideOrderRepo,
	provideVoucherRepo,
	provideReviewRepo,
	providePaymentRepo,
	provideSellerRepo,

	// Usecases
	provideUserUsecase,
	provideProductUsecase,
	provideCartUsecase,
	provideOrderUsecase,
	provideVoucherUsecase,
	provideReviewUsecase,
	provideSellerUsecase,
)

func provideRouter(handler http.IHandler, jwtService auth.IJWTService) http.Router {
	return http.NewRouter(handler, jwtService)
}

func provideJWTService() auth.IJWTService {
	return auth.NewJWTService()
}

func provideHandler(
	userUsecase usecase.IUserUsecase,
	productUsecase usecase.IProductUsecase,
	cartUsecase usecase.ICartUsecase,
	orderUsecase usecase.IOrderUsecase,
	voucherUsecase usecase.IVoucherUsecase,
	reviewUsecase usecase.IReviewUsecase,
	sellerUsecase usecase.ISellerUsecase,
) http.IHandler {
	handler := http.NewHandler(
		userUsecase,
		productUsecase,
		cartUsecase,
		orderUsecase,
		voucherUsecase,
		reviewUsecase,
		sellerUsecase,
	)
	return handler
}

// Repository providers
func provideUserRepo(db *gorm.DB) repository.IUserRepo {
	return repository.NewUserRepo(db)
}

func provideProductRepo(db *gorm.DB) repository.IProductRepo {
	return repository.NewProductRepo(db)
}

func provideCartRepo(db *gorm.DB) repository.ICartRepo {
	return repository.NewCartRepo(db)
}

func provideOrderRepo(db *gorm.DB) repository.IOrderRepo {
	return repository.NewOrderRepo(db)
}

func provideVoucherRepo(db *gorm.DB) repository.IVoucherRepo {
	return repository.NewVoucherRepo(db)
}

func provideReviewRepo(db *gorm.DB) repository.IReviewRepo {
	return repository.NewReviewRepo(db)
}

func providePaymentRepo(db *gorm.DB) repository.IPaymentRepo {
	return repository.NewPaymentRepo(db)
}

func provideSellerRepo(db *gorm.DB) repository.ISellerRepository {
	return repository.NewSellerRepo(db)
}

// Usecase providers
func provideUserUsecase(repo repository.IUserRepo, jwtService auth.IJWTService) usecase.IUserUsecase {
	return usecase.NewUserUsecase(repo, jwtService)
}

func provideProductUsecase(repo repository.IProductRepo) usecase.IProductUsecase {
	return usecase.NewProductUsecase(repo)
}

func provideCartUsecase(
	cartRepo repository.ICartRepo,
	productRepo repository.IProductRepo,
) usecase.ICartUsecase {
	return usecase.NewCartUsecase(cartRepo, productRepo)
}

func provideOrderUsecase(
	orderRepo repository.IOrderRepo,
	cartRepo repository.ICartRepo,
	voucherRepo repository.IVoucherRepo,
	productRepo repository.IProductRepo,
	paymentRepo repository.IPaymentRepo,
) usecase.IOrderUsecase {
	return usecase.NewOrderUsecase(orderRepo, cartRepo, voucherRepo, productRepo, paymentRepo)
}

func provideVoucherUsecase(repo repository.IVoucherRepo) usecase.IVoucherUsecase {
	return usecase.NewVoucherUsecase(repo)
}

func provideReviewUsecase(reviewRepo repository.IReviewRepo) usecase.IReviewUsecase {
	return usecase.NewReviewUsecase(reviewRepo)
}

func provideSellerUsecase(
	sellerRepo repository.ISellerRepository,
	userRepo repository.IUserRepo,
) usecase.ISellerUsecase {
	return usecase.NewSellerUsecase(sellerRepo, userRepo)
}
