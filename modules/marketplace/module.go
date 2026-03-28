package marketplace

import (
	"strings"

	"backend/internal/pkg/security"
	"backend/modules/eskiz"
	"backend/modules/marketplace/domain"
	"backend/modules/marketplace/handler"
	"backend/modules/marketplace/repository"
	"backend/modules/marketplace/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string, jwtExpireHours int) error {
	if err := db.AutoMigrate(&domain.User{}, &domain.VerificationCode{}, &domain.DeliveryArea{}); err != nil {
		return err
	}
	// Legacy schema compatibility: old versions had password NOT NULL.
	// Marketplace auth is SMS-only now, so this column must not block inserts.
	if db.Migrator().HasColumn(&domain.User{}, "password") {
		if err := db.Exec(`ALTER TABLE marketplace_users ALTER COLUMN password DROP NOT NULL`).Error; err != nil {
			return err
		}
	}

	repo := repository.NewAuthRepository(db)
	eskizSvc, _ := eskiz.NewFromEnv()
	authSvc := service.NewAuthService(repo, eskizSvc, jwtSecret, jwtExpireHours)
	authHandler := handler.NewAuthHandler(authSvc)
	regionRepo := repository.NewRegionRepository(db)
	regionSvc := service.NewRegionService(regionRepo)
	regionHandler := handler.NewRegionHandler(regionSvc)
	deliveryRepo := repository.NewDeliveryAreaRepository(db)
	deliverySvc := service.NewDeliveryAreaService(deliveryRepo)
	deliveryHandler := handler.NewDeliveryAreaHandler(deliverySvc)

	productRepo := repository.NewMarketplaceProductRepository(db)
	productSvc := service.NewMarketplaceProductService(productRepo)
	productHandler := handler.NewMarketplaceProductHandler(productSvc)
	categoryRepo := repository.NewMarketplaceCategoryRepository(db)
	categorySvc := service.NewMarketplaceCategoryService(categoryRepo)
	categoryHandler := handler.NewMarketplaceCategoryHandler(categorySvc)
	searchRepo := repository.NewMarketplaceSearchRepository(db)
	searchSvc := service.NewMarketplaceSearchService(searchRepo, productRepo)
	contragentBrowseRepo := repository.NewMarketplaceContragentBrowseRepository(db)
	contragentBrowseSvc := service.NewMarketplaceContragentBrowseService(contragentBrowseRepo, productRepo)
	searchHandler := handler.NewMarketplaceSearchHandler(searchSvc, contragentBrowseSvc)

	api := router.Group("/api/v1")
	authMiddleware := marketplaceAuthMiddleware(jwtSecret)
	authHandler.RegisterRoutes(api, authMiddleware)
	regionHandler.RegisterRoutes(api, authMiddleware)
	deliveryHandler.RegisterRoutes(api, authMiddleware)
	productHandler.RegisterRoutes(api)
	categoryHandler.RegisterRoutes(api)
	searchHandler.RegisterRoutes(api)
	return nil
}

func marketplaceAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"message": "Authorization header topilmadi"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(401, gin.H{"message": "Authorization formati noto'g'ri"})
			c.Abort()
			return
		}
		claims, err := security.ParseMarketplaceToken(secret, parts[1])
		if err != nil {
			c.JSON(401, gin.H{"message": "token yaroqsiz"})
			c.Abort()
			return
		}
		c.Set("marketplace_user_id", claims.UserID)
		c.Next()
	}
}
