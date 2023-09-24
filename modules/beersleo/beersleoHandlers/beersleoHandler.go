package beersleoHandlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/peedans/beerleo/modules/beersleo"
	"github.com/peedans/beerleo/modules/beersleo/beersleoUsecases"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type IBeersleoHandler interface {
	GetBeerByID(c *gin.Context)
	DeleteBeer(c *gin.Context)
	FilterBeersByName(c *gin.Context)
	GetAllBeersPagination(c *gin.Context)
	UpdateBeer(c *gin.Context)
	CreateBeer(c *gin.Context)
}

type beersleoHandler struct {
	beersleoUsecase beersleoUsecases.IBeersleoUsecase
}

func BeersleoHandler(beersleoUsecase beersleoUsecases.IBeersleoUsecase) IBeersleoHandler {
	return &beersleoHandler{
		beersleoUsecase: beersleoUsecase,
	}
}

func (h *beersleoHandler) GetBeerByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid beer ID"})
		return
	}

	beer, err := h.beersleoUsecase.GetBeerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve beer"})
		return
	}
	c.JSON(http.StatusOK, beer)
}

func (h *beersleoHandler) GetAllBeersPagination(c *gin.Context) {
	page, limit, err := getPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	beersData, total, err := h.beersleoUsecase.GetAllBeersPagination(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve beers"})
		return
	}

	pagination, err := getPagination(c, total)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := struct {
		Data   []*beersleo.Beersleo          `json:"data"`
		Paging *beersleo.BeerleoPagingResult `json:"paging"`
	}{
		Data:   beersData,
		Paging: pagination,
	}

	c.JSON(http.StatusOK, response)
}

func getPaginationParams(c *gin.Context) (page, limit int, err error) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err = strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid page number")
	}
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid limit number")
	}

	return page, limit, nil
}

func getPagination(c *gin.Context, total int) (*beersleo.BeerleoPagingResult, error) {
	page, limit, err := getPaginationParams(c)
	if err != nil {
		return nil, err
	}
	totalPages := (total + limit - 1) / limit

	return &beersleo.BeerleoPagingResult{
		Page:      page,
		Limit:     limit,
		PrevPage:  max(1, page-1),
		NextPage:  min(totalPages, page+1),
		Count:     total,
		TotalPage: totalPages,
	}, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h *beersleoHandler) CreateBeer(c *gin.Context) {
	var beerCreate beersleo.BeerCreationRequest

	if err := c.ShouldBind(&beerCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind request"})
		return
	}

	beerData := beersleo.BeerDTO{
		Name:     beerCreate.Name,
		Category: beerCreate.Category,
		Detail:   beerCreate.Detail,
	}

	id, err := h.beersleoUsecase.CreateBeer(&beerData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create beer"})
		return
	}

	beerResponse := beersleo.Beersleo{
		ID:       id,
		Name:     beerCreate.Name,
		Category: beerCreate.Category,
		Detail:   beerCreate.Detail,
	}

	imagePath, err := setBeerImage(c, &beerResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set beer image"})
		return
	}

	beerResponse.Image = imagePath

	err = h.beersleoUsecase.UpdateBeer(&beerResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beer created successfully"})
}

func (h *beersleoHandler) UpdateBeer(c *gin.Context) {
	var beerUpdate beersleo.UpdateBeer

	if err := c.ShouldBind(&beerUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind request"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid beer ID"})
		return
	}

	beerResponse, err := h.beersleoUsecase.GetBeerByID(id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch beer details"})
		return
	}

	beerResponse.Name = beerUpdate.Name
	beerResponse.Category = beerUpdate.Category
	beerResponse.Detail = beerUpdate.Detail

	imagePath, err := setBeerImage(c, beerResponse)
	fmt.Println(imagePath)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set beer image"})
		return
	}
	beerResponse.Image = imagePath

	err = h.beersleoUsecase.UpdateBeer(beerResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update beer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beer updated successfully"})
}

func setBeerImage(c *gin.Context, beer *beersleo.Beersleo) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		// ถ้ามีข้อผิดพลาด ส่งคืนค่าข้อผิดพลาด
		return "", err
	}

	if beer.Image != "" {
		beer.Image = strings.Replace(beer.Image, getHost(c), "", 1)
		pwd, err := os.Getwd()
		if err != nil {
			// ถ้ามีข้อผิดพลาด ส่งคืนค่าข้อผิดพลาด
			return "", err
		}

		if err := os.Remove(pwd + beer.Image); err != nil {
			return "", err
		}
	}
	path := "uploads/beers/" + strconv.Itoa(int(beer.ID))

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	fileName := path + "/" + file.Filename

	if err := c.SaveUploadedFile(file, fileName); err != nil {
		return "", err
	}

	beer.Image = getHost(c) + "/" + fileName

	return beer.Image, nil
}

func getHost(c *gin.Context) string {
	if c.Request.URL.Scheme == "" {
		// ถ้าว่าง จะถือว่าเป็น http และเติม "http://" หน้า host
		return "http://" + c.Request.Host
	}
	return "https://" + c.Request.Host
}

func (h *beersleoHandler) DeleteBeer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid beer ID"})
		return
	}

	err = h.beersleoUsecase.DeleteBeer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete beer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beer deleted successfully"})
}

func (h *beersleoHandler) FilterBeersByName(c *gin.Context) {

	nameQuery := c.DefaultQuery("name", "")

	if nameQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name query parameter is required"})
		return
	}

	filter := &beersleo.BeersleoFilter{
		Name: nameQuery,
	}

	beersData, err := h.beersleoUsecase.FilterBeersByName(filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter beers by name"})
		return
	}
	c.JSON(http.StatusOK, beersData)
}
