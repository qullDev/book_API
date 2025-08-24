package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qullDev/book_API/internal/domain/book"
	"github.com/qullDev/book_API/internal/domain/category"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.POST("", h.Create)
	rg.GET("/:id", h.Detail)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
	rg.GET("/:id/books", h.ListBooks)
}

type createCategoryReq struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type updateCategoryReq struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

func (h *CategoryHandler) List(c *gin.Context) {
	var items []category.Category
	if err := h.db.Order("name asc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload tidak valid", "error": err.Error()})
		return
	}
	item := category.Category{Name: req.Name}
	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menambahkan kategori"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *CategoryHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	var item category.Category
	if err := h.db.First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kategori tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil detail kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	var req updateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload tidak valid", "error": err.Error()})
		return
	}
	var item category.Category
	if err := h.db.First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kategori tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data kategori"})
		return
	}
	item.Name = req.Name
	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengupdate kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	res := h.db.Delete(&category.Category{}, "id = ?", id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menghapus kategori"})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "kategori tidak tersedia"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil dihapus"})
}

func (h *CategoryHandler) ListBooks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	var books []book.Book
	if err := h.db.Where("category_id = ?", id).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil buku pada kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}
