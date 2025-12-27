package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
)

type Router interface {
	Register(routerGroup gin.IRouter)
}

type routerImpl struct {
	handler IHandler
}

func NewRouter(handler IHandler) Router {
	return &routerImpl{
		handler: handler,
	}
}

func (p *routerImpl) Register(r gin.IRouter) {
	// routes for chophimco service
	api := r.Group("chophimco/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			apiwrapper.SendSuccess(c, gin.H{
				"message":   "Chophimco service is running",
				"timestamp": time.Now(),
			})
		})
	}

	// User routes
	userApi := api.Group("user")
	{
		userApi.POST("/login", p.handler.Login)
		userApi.POST("/register", p.handler.Register)
		userApi.GET("/profile", p.handler.GetProfile) // Requires auth
	}

	// Product routes
	productApi := api.Group("product")
	{
		// Public routes
		productApi.GET("/all", p.handler.GetAllProducts)
		productApi.GET("/:id", p.handler.GetProductByID)
		productApi.GET("/category", p.handler.GetProductsByCategory)
		productApi.GET("/brand", p.handler.GetProductsByBrand)

		// Admin routes
		productApi.POST("/create", p.handler.CreateProduct)
		productApi.PUT("/update", p.handler.UpdateProduct)
		productApi.DELETE("/:id", p.handler.DeleteProduct)

		// Variant routes
		productApi.POST("/variant/create", p.handler.CreateProductVariant)
		productApi.PUT("/variant/update", p.handler.UpdateProductVariant)
	}

	// Cart routes (all require authentication)
	cartApi := api.Group("cart")
	{
		cartApi.GET("", p.handler.GetCart)
		cartApi.POST("/add", p.handler.AddToCart)
		cartApi.PUT("/update", p.handler.UpdateCartItem)
		cartApi.DELETE("/:id", p.handler.RemoveFromCart)
		cartApi.DELETE("/clear", p.handler.ClearCart)
	}

	// Order routes
	orderApi := api.Group("order")
	{
		orderApi.POST("/create", p.handler.CreateOrder)
		orderApi.GET("/:id", p.handler.GetOrderByID)
		orderApi.GET("/my-orders", p.handler.GetMyOrders)
		orderApi.PUT("/status", p.handler.UpdateOrderStatus) // Admin only
	}

	// Voucher routes
	voucherApi := api.Group("voucher")
	{
		// Public routes
		voucherApi.GET("/active", p.handler.GetActiveVouchers)
		voucherApi.GET("", p.handler.GetVoucherByCode)
		voucherApi.GET("/validate", p.handler.ValidateVoucher)

		// Admin routes
		voucherApi.GET("/all", p.handler.GetAllVouchers)
		voucherApi.POST("/create", p.handler.CreateVoucher)
		voucherApi.PUT("/update", p.handler.UpdateVoucher)
		voucherApi.DELETE("/:id", p.handler.DeleteVoucher)
	}

	// Review routes
	reviewApi := api.Group("review")
	{
		reviewApi.GET("", p.handler.GetProductReviews)
		reviewApi.POST("/create", p.handler.CreateReview) // Requires auth
	}
}
