package beersleoRepositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/peedans/beerleo/modules/beersleo"
)

type IBeersleoRepository interface {
	FilterBeersByName(req *beersleo.BeersleoFilter) ([]*beersleo.Beersleo, error)
	GetByID(id int) (*beersleo.Beersleo, error)
	Delete(id int) error
	GetAllBeersWithPagination(page, limit int) ([]*beersleo.Beersleo, int, error)
	Create(beer *beersleo.BeerDTO) (int, error)
	Update(beer *beersleo.Beersleo) error
}

type beersleoRepository struct {
	db *sqlx.DB
}

func BeersleoRepository(db *sqlx.DB) IBeersleoRepository {
	return &beersleoRepository{
		db: db,
	}
}

func (r *beersleoRepository) GetByID(id int) (*beersleo.Beersleo, error) {
	var beer beersleo.Beersleo
	err := r.db.Get(&beer, "SELECT id, name, category, detail, image FROM beers WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("Beer with ID %d not found: %v", id, err)
	}
	return &beer, nil
}

func (r *beersleoRepository) Create(beer *beersleo.BeerDTO) (int, error) {
	result, err := r.db.NamedExec("INSERT INTO beers(name, category, detail, image) VALUES (:name, :category, :detail, :image)", beer)

	if err != nil {
		// ถ้ามีข้อผิดพลาด ส่งคืนค่า 0 และ err
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		// ถ้ามีข้อผิดพลาด ส่งคืนค่า 0 และ err
		return 0, err
	}

	return int(lastInsertID), nil
}

func (r *beersleoRepository) Update(beer *beersleo.Beersleo) error {
	_, err := r.db.NamedExec("UPDATE beers SET name=:name, category=:category, detail=:detail, image=:image WHERE id=:id",
		&beer)
	return err
}

func (r *beersleoRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM beers WHERE id=?", id)
	return err
}

func (r *beersleoRepository) GetAllBeersWithPagination(page, limit int) ([]*beersleo.Beersleo, int, error) {

	var beers []*beersleo.Beersleo

	offset := (page - 1) * limit

	var total int

	err := r.db.Get(&total, "SELECT COUNT(*) FROM beers")
	if err != nil {
		// ถ้าเกิดข้อผิดพลาด ส่งคืนค่า nil, 0, และ err
		return nil, 0, err
	}

	err = r.db.Select(&beers, "SELECT * FROM beers LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		// ถ้าเกิดข้อผิดพลาด ส่งคืนค่า nil, 0, และข้อผิดพลาดที่มีการระบุเพิ่มเติม
		return nil, 0, fmt.Errorf("error fetching beers with pagination: %w", err)
	}

	fmt.Printf("Limit: %d, Offset: %d\n", limit, offset)

	return beers, total, nil
}
func (r *beersleoRepository) FilterBeersByName(req *beersleo.BeersleoFilter) ([]*beersleo.Beersleo, error) {
	var beerList []*beersleo.Beersleo
	query := `SELECT id, name, category, detail, image, created_at, updated_at, deleted_at FROM beers WHERE name LIKE ?`
	rows, err := r.db.Query(query, "%"+req.Name+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var beer beersleo.Beersleo
		if err := rows.Scan(&beer.ID, &beer.Name, &beer.Category, &beer.Detail, &beer.Image, &beer.CreatedAt, &beer.UpdatedAt, &beer.DeletedAt); err != nil {
			return nil, err
		}

		beerList = append(beerList, &beer)
	}

	return beerList, nil
}
