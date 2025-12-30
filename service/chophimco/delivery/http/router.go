package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
	"github.com/leehai1107/chophimco-server/pkg/middleware/auth"
)

type Router interface {
	Register(routerGroup gin.IRouter)
}

type routerImpl struct {
	handler    IHandler
	jwtService auth.IJWTService
}

func NewRouter(handler IHandler, jwtService auth.IJWTService) Router {
	return &routerImpl{
		handler:    handler,
		jwtService: jwtService,
	}
}

func (p *routerImpl) Register(r gin.IRouter) {
	// routes for chophimco service
	api := r.Group("api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			apiwrapper.SendSuccess(c, gin.H{
				"message":   "Chophimco service is running",
				"timestamp": time.Now(),
			})
		})
	}

	// Create auth middleware instance
	authMiddleware := auth.AuthMiddleware(p.jwtService)
	adminMiddleware := auth.RoleMiddleware("admin")
	sellerMiddleware := auth.RoleMiddleware("seller", "admin")

	// User routes
	userApi := api.Group("user")
	{
		userApi.POST("/login", p.handler.Login)
		userApi.POST("/register", p.handler.Register)
		userApi.GET("/profile", authMiddleware, p.handler.GetProfile) // Protected
	}

	// Product routes
	productApi := api.Group("product")
	{
		// Public routes
		productApi.GET("/all", p.handler.GetAllProducts)
		productApi.GET("/:id", p.handler.GetProductByID)
		productApi.GET("/category", p.handler.GetProductsByCategory)
		productApi.GET("/brand", p.handler.GetProductsByBrand)

		// Admin routes (protected)
		productApi.POST("/create", authMiddleware, adminMiddleware, p.handler.CreateProduct)
		productApi.PUT("/update", authMiddleware, adminMiddleware, p.handler.UpdateProduct)
		productApi.DELETE("/:id", authMiddleware, adminMiddleware, p.handler.DeleteProduct)

		// Variant routes (admin only)
		productApi.POST("/variant/create", authMiddleware, adminMiddleware, p.handler.CreateProductVariant)
		productApi.PUT("/variant/update", authMiddleware, adminMiddleware, p.handler.UpdateProductVariant)
	}

	// Cart routes (all require authentication)
	cartApi := api.Group("cart", authMiddleware)
	{
		cartApi.GET("", p.handler.GetCart)
		cartApi.POST("/add", p.handler.AddToCart)
		cartApi.PUT("/update", p.handler.UpdateCartItem)
		cartApi.DELETE("/:id", p.handler.RemoveFromCart)
		cartApi.DELETE("/clear", p.handler.ClearCart)
	}

	// Order routes (all require authentication)
	orderApi := api.Group("order", authMiddleware)
	{
		orderApi.POST("/create", p.handler.CreateOrder)
		orderApi.GET("/:id", p.handler.GetOrderByID)
		orderApi.GET("/my-orders", p.handler.GetMyOrders)
		orderApi.PUT("/status", adminMiddleware, p.handler.UpdateOrderStatus) // Admin only
	}

	// Voucher routes
	voucherApi := api.Group("voucher")
	{
		// Public routes
		voucherApi.GET("/active", p.handler.GetActiveVouchers)
		voucherApi.GET("", p.handler.GetVoucherByCode)
		voucherApi.GET("/validate", p.handler.ValidateVoucher)

		// Admin routes (protected)
		voucherApi.GET("/all", authMiddleware, adminMiddleware, p.handler.GetAllVouchers)
		voucherApi.POST("/create", authMiddleware, adminMiddleware, p.handler.CreateVoucher)
		voucherApi.PUT("/update", authMiddleware, adminMiddleware, p.handler.UpdateVoucher)
		voucherApi.DELETE("/:id", authMiddleware, adminMiddleware, p.handler.DeleteVoucher)
	}

	// Review routes
	reviewApi := api.Group("review")
	{
		reviewApi.GET("", p.handler.GetProductReviews)
		reviewApi.POST("/create", authMiddleware, p.handler.CreateReview) // Protected
	}

	// Seller routes
	sellerApi := api.Group("seller")
	{
		// Public seller info
		sellerApi.GET("/:id", p.handler.GetSellerByID)
		sellerApi.GET("/reviews", p.handler.GetSellerReviews)

		// Seller profile management (requires seller or admin role)
		sellerApi.POST("/profile", authMiddleware, sellerMiddleware, p.handler.CreateSellerProfile)
		sellerApi.GET("/profile", authMiddleware, sellerMiddleware, p.handler.GetSellerProfile)
		sellerApi.PUT("/profile", authMiddleware, sellerMiddleware, p.handler.UpdateSellerProfile)

		// Seller product management (requires seller or admin role)
		sellerApi.POST("/product", authMiddleware, sellerMiddleware, p.handler.CreateSellerProduct)
		sellerApi.GET("/product", authMiddleware, sellerMiddleware, p.handler.GetSellerProducts)
		sellerApi.PUT("/product", authMiddleware, sellerMiddleware, p.handler.UpdateSellerProduct)
		sellerApi.DELETE("/product/:id", authMiddleware, sellerMiddleware, p.handler.DeleteSellerProduct)

		// Product image management (requires seller or admin role)
		sellerApi.POST("/product/image", authMiddleware, sellerMiddleware, p.handler.UploadProductImage)
		sellerApi.DELETE("/product/image/:id", authMiddleware, sellerMiddleware, p.handler.DeleteProductImage)
		sellerApi.PUT("/product/image/primary", authMiddleware, sellerMiddleware, p.handler.SetPrimaryImage)

		// Seller reviews (requires authentication)
		sellerApi.POST("/reviews", authMiddleware, p.handler.CreateSellerReview)
	}

	// Admin routes (all require admin role)
	adminApi := api.Group("admin", authMiddleware, adminMiddleware)
	{
		// Seller verification
		adminApi.GET("/seller/pending", p.handler.GetPendingSellers)
		adminApi.POST("/seller/verify", p.handler.VerifySeller)
		adminApi.POST("/seller/reject", p.handler.RejectSeller)

		// Product approval
		adminApi.GET("/product/pending", p.handler.GetPendingProducts)
		adminApi.POST("/product/approve", p.handler.ApproveProduct)
		adminApi.POST("/product/reject", p.handler.RejectProduct)
	}
}
