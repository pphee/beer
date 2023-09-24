package beersleo

import (
	"mime/multipart"
	"time"
)

type Beersleo struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Category  string     `db:"category" json:"category"`
	Detail    string     `db:"detail" json:"detail"`
	Image     string     `db:"image" json:"image"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type BeersleoFilter struct {
	Name string `query:"name"`
}

type BeerDTO struct {
	ID       int
	Name     string `form:"name" binding:"required"`
	Category string `form:"category" binding:"required"`
	Detail   string `form:"detail" binding:"required"`
	Image    string `form:"image"`
}

type BeerleoPagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

type BeerCreationRequest struct {
	Name     string                `form:"name" binding:"required"`
	Category string                `form:"category" binding:"required"`
	Detail   string                `form:"detail" binding:"required"`
	Image    *multipart.FileHeader `form:"image"`
}

type UpdateBeer struct {
	Name     string                `form:"name"`
	Category string                `form:"category"`
	Detail   string                `form:"detail"`
	Image    *multipart.FileHeader `form:"image"`
}
