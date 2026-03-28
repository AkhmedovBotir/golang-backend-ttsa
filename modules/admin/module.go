package admin

import (
	"backend/modules/admin/domain"
	"backend/modules/admin/handler"
	"backend/modules/admin/repository"
	"backend/modules/admin/service"
	contrDomain "backend/modules/contragents/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string, jwtExpireHours int) error {
	if err := db.AutoMigrate(&domain.Admin{}, &domain.Region{}, &domain.District{}, &domain.MFY{}, &domain.ContragentType{}, &domain.Contragent{}, &domain.NeighborhoodShop{}, &domain.Agent{}, &domain.Punkt{}, &domain.Manager{}, &domain.Category{}, &contrDomain.Product{}, &contrDomain.ProductImage{}); err != nil {
		return err
	}

	repo := repository.NewAdminPostgresRepository(db)
	svc := service.NewAdminService(repo)
	h := handler.NewAdminHandler(svc, jwtSecret, jwtExpireHours)
	regionRepo := repository.NewRegionRepository(db)
	regionSvc := service.NewRegionService(regionRepo)
	regionHandler := handler.NewRegionHandler(regionSvc)
	contragentTypeRepo := repository.NewContragentTypeRepository(db)
	contragentTypeSvc := service.NewContragentTypeService(contragentTypeRepo)
	contragentTypeHandler := handler.NewContragentTypeHandler(contragentTypeSvc)
	contragentRepo := repository.NewContragentRepository(db)
	contragentSvc := service.NewContragentService(contragentRepo, regionRepo, contragentTypeRepo)
	contragentHandler := handler.NewContragentHandler(contragentSvc)
	neighborhoodShopRepo := repository.NewNeighborhoodShopRepository(db)
	neighborhoodShopSvc := service.NewNeighborhoodShopService(neighborhoodShopRepo, regionRepo)
	neighborhoodShopHandler := handler.NewNeighborhoodShopHandler(neighborhoodShopSvc)
	agentRepo := repository.NewAgentRepository(db)
	agentSvc := service.NewAgentService(agentRepo, regionRepo)
	agentHandler := handler.NewAgentHandler(agentSvc)
	punktRepo := repository.NewPunktRepository(db)
	punktSvc := service.NewPunktService(punktRepo, regionRepo)
	punktHandler := handler.NewPunktHandler(punktSvc)
	managerRepo := repository.NewManagerRepository(db)
	managerSvc := service.NewManagerService(managerRepo, regionRepo)
	managerHandler := handler.NewManagerHandler(managerSvc)
	marketplaceUserRepo := repository.NewMarketplaceUserRepository(db)
	marketplaceUserSvc := service.NewMarketplaceUserService(marketplaceUserRepo)
	marketplaceUserHandler := handler.NewMarketplaceUserHandler(marketplaceUserSvc)
	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	subcategoryHandler := handler.NewSubcategoryHandler(categorySvc)
	productRepo := repository.NewAdminProductRepository(db)
	productSvc := service.NewAdminProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)
	authMiddleware := h.AuthMiddleware()
	onlyGeneral := handler.GeneralOnly()

	api := router.Group("/api/v1")
	{
		api.POST("/admins/login", h.Login)

		adminRoutes := api.Group("/admins")
		adminRoutes.Use(authMiddleware, onlyGeneral)
		{
			adminRoutes.POST("", h.Create)
			adminRoutes.GET("", h.GetAll)
			adminRoutes.GET("/:id", h.GetByID)
			adminRoutes.PUT("/:id", h.Update)
			adminRoutes.PATCH("/:id/status", h.UpdateStatus)
			adminRoutes.DELETE("/:id", h.Delete)
		}

		regionHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		contragentTypeHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		contragentHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		neighborhoodShopHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		agentHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		punktHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		managerHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		marketplaceUserHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		categoryHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		subcategoryHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
		productHandler.RegisterRoutes(api, authMiddleware, onlyGeneral)
	}
	return nil
}
