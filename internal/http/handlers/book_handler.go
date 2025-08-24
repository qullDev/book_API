package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qullDev/book_API/internal/domain/book"
	"gorm.io/gorm"
)

const (
	minReleaseYear = 1980
	maxReleaseYear = 2024
)

type BookHandler struct {
	db *gorm.DB
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	return &BookHandler{db: db}
}

func (h *BookHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.POST("", h.Create)
	rg.GET("/:id", h.Detail)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}

type createBookReq struct {
	Title       string    `json:"title" binding:"required,max=200"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ReleaseYear int       `json:"release_year" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	TotalPage   int       `json:"total_page" binding:"required"`
}

type updateBookReq struct {
	Title       *string    `json:"title" binding:"omitempty,max=200"`
	CategoryID  *uuid.UUID `json:"category_id"`
	Description *string    `json:"description"`
	ImageURL    *string    `json:"image_url"`
	ReleaseYear *int       `json:"release_year"`
	Price       *float64   `json:"price"`
	TotalPage   *int       `json:"total_page"`
}

// @Summary List all books
// @Description Get list of all books
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]book.Book
// @Router /api/books [get]
func (h *BookHandler) List(c *gin.Context) {
	var items []book.Book
	if err := h.db.Order("title asc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data buku"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// @Summary Create new book
// @Description Create a new book
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param book body createBookReq true "Book data" example({"title":"The Go Programming Language","category_id":"550e8400-e29b-41d4-a716-446655440000","description":"Comprehensive guide to Go","image_url":"https://example.com/book.jpg","release_year":2020,"price":59.99,"total_page":150})
// @Success 201 {object} map[string]book.Book "example={'data':{'id':'550e8400-e29b-41d4-a716-446655440000','title':'The Go Programming Language','category_id':'550e8400-e29b-41d4-a716-446655440000','description':'Comprehensive guide to Go','release_year':2020,'price':59.99,'total_page':150,'thickness':'tebal'}}"
// @Failure 400 {object} gin.H "example={'message':'payload tidak valid'}"
// @Router /api/books [post]
func (h *BookHandler) Create(c *gin.Context) {
	var req createBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload tidak valid", "error": err.Error()})
		return
	}
	if req.ReleaseYear < minReleaseYear || req.ReleaseYear > maxReleaseYear {
		c.JSON(http.StatusBadRequest, gin.H{"message": "release_year harus antara 1980 sampai 2024"})
		return
	}

	item := book.Book{
		Title:       req.Title,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		ReleaseYear: req.ReleaseYear,
		Price:       req.Price,
		TotalPage:   req.TotalPage,
		// Thickness akan diisi otomatis oleh hook BeforeCreate
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menambahkan buku", "error": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

// @Summary Get book detail
// @Description Get detail of a book
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]book.Book
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/books/{id} [get]
func (h *BookHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	var item book.Book
	if err := h.db.First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "buku tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil detail buku"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// @Summary Update book
// @Description Update an existing book
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Book ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Param book body updateBookReq true "Updated book data" example({"title":"Updated Title","price":49.99})
// @Success 200 {object} map[string]book.Book
// @Failure 400,404 {object} gin.H "example={'message':'buku tidak ditemukan'}"
// @Router /api/books/{id} [put]
func (h *BookHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	var existing book.Book
	if err := h.db.First(&existing, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "buku tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data buku"})
		return
	}

	var req updateBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload tidak valid", "error": err.Error()})
		return
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.CategoryID != nil {
		existing.CategoryID = *req.CategoryID
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.ImageURL != nil {
		existing.ImageURL = *req.ImageURL
	}
	if req.ReleaseYear != nil {
		if *req.ReleaseYear < minReleaseYear || *req.ReleaseYear > maxReleaseYear {
			c.JSON(http.StatusBadRequest, gin.H{"message": "release_year harus antara 1980 sampai 2024"})
			return
		}
		existing.ReleaseYear = *req.ReleaseYear
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.TotalPage != nil {
		existing.TotalPage = *req.TotalPage
		// update thickness saat total_page berubah
		if existing.TotalPage > 100 {
			existing.Thickness = "tebal"
		} else {
			existing.Thickness = "tipis"
		}
	}

	if err := h.db.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengupdate buku"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": existing})
}

// @Summary Delete book
// @Description Delete a book
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/books/{id} [delete]
func (h *BookHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	res := h.db.Delete(&book.Book{}, "id = ?", id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menghapus buku"})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "buku tidak tersedia"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "buku berhasil dihapus"})
}

