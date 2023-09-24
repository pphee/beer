package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/peedans/beerleo/modules/beersleo/beersleoHandlers"
	"github.com/peedans/beerleo/modules/beersleo/beersleoRepositories"
	"github.com/peedans/beerleo/modules/beersleo/beersleoUsecases"
	monitorHandlers "github.com/peedans/beerleo/modules/monitorHandlers/handlers"
)

type IModuleFactory interface {
	monitorModule()
	beersleoModule()
}

type moduleFactory struct {
	r *gin.RouterGroup
	s *server
}

func InitModule(r *gin.RouterGroup, s *server) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
	}
}

func (mf *moduleFactory) monitorModule() {
	handler := monitorHandlers.MonitorHandler(mf.s.cfg)
	mf.r.GET("/", handler.HealthCheck)
}

func (mf *moduleFactory) beersleoModule() {
	repo := beersleoRepositories.BeersleoRepository(mf.s.db)
	usecases := beersleoUsecases.BeersleoUsecase(repo)
	handler := beersleoHandlers.BeersleoHandler(usecases)
	beerRouter := mf.r.Group("/beers")
	beerRouter.GET("/filter", handler.FilterBeersByName)
	beerRouter.GET("/", handler.GetAllBeersPagination)
	beerRouter.GET("/:id", handler.GetBeerByID)
	beerRouter.DELETE("/:id", handler.DeleteBeer)
	beerRouter.POST("/", handler.CreateBeer)
	beerRouter.PUT("/:id", handler.UpdateBeer)
}
