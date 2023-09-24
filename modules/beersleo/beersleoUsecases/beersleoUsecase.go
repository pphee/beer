package beersleoUsecases

import (
	"errors"
	"github.com/peedans/beerleo/modules/beersleo"
	"github.com/peedans/beerleo/modules/beersleo/beersleoRepositories"
)

type IBeersleoUsecase interface {
	GetBeerByID(id int) (*beersleo.Beersleo, error)
	DeleteBeer(id int) error
	FilterBeersByName(req *beersleo.BeersleoFilter) ([]*beersleo.BeerDTO, error)
	GetAllBeersPagination(page, limit int) ([]*beersleo.Beersleo, int, error)
	CreateBeer(beer *beersleo.BeerDTO) (int, error)
	UpdateBeer(beer *beersleo.Beersleo) error
}

type beersleoUsecase struct {
	beersleoRepository beersleoRepositories.IBeersleoRepository
}

func BeersleoUsecase(beersleoRepository beersleoRepositories.IBeersleoRepository) IBeersleoUsecase {
	return &beersleoUsecase{
		beersleoRepository: beersleoRepository,
	}
}

var ErrInvalidBeerID = errors.New("invalid beer ID provided")

func (bu *beersleoUsecase) GetBeerByID(id int) (*beersleo.Beersleo, error) {
	return bu.beersleoRepository.GetByID(id)
}

func (bu *beersleoUsecase) DeleteBeer(id int) error {
	return bu.beersleoRepository.Delete(id)
}

func (bu *beersleoUsecase) GetAllBeersPagination(page, limit int) ([]*beersleo.Beersleo, int, error) {

	beerResponses, total, err := bu.beersleoRepository.GetAllBeersWithPagination(page, limit)

	if err != nil {
		// ถ้ามีข้อผิดพลาด ส่งคืนค่า nil, 0, และ err
		return nil, 0, err
	}

	return beerResponses, total, nil
}

func (bu *beersleoUsecase) CreateBeer(beer *beersleo.BeerDTO) (int, error) {
	return bu.beersleoRepository.Create(beer)
}

func (bu *beersleoUsecase) UpdateBeer(beer *beersleo.Beersleo) error {
	if beer.ID == 0 {
		return ErrInvalidBeerID
	}
	return bu.beersleoRepository.Update(beer)
}

func (bu *beersleoUsecase) FilterBeersByName(req *beersleo.BeersleoFilter) ([]*beersleo.BeerDTO, error) {

	beersByName, err := bu.beersleoRepository.FilterBeersByName(req)

	if err != nil {
		return nil, err
	}

	var beerResponses []*beersleo.BeerDTO

	for _, b := range beersByName {
		beerResponse := &beersleo.BeerDTO{
			ID:       b.ID,
			Name:     b.Name,
			Category: b.Category,
			Detail:   b.Detail,
			Image:    b.Image,
		}
		beerResponses = append(beerResponses, beerResponse)
	}

	return beerResponses, nil
}
