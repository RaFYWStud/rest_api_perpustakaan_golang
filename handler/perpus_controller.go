package handler

import (
	"golang-tutorial/contract"
	"golang-tutorial/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type perpusController struct {
    service contract.PerpusService
}

func (c *perpusController) getPrefix() string {
    return "/perpus"
}

func (c *perpusController) initService(service *contract.Service) {
    c.service = service.Perpus
}

func (c *perpusController) initRoute(app *gin.RouterGroup) {
    // Search Endpoint
    app.GET("/books", c.GetAllBooks)
    app.GET("/books/:id", c.GetBookByID)
    app.GET("/search", c.SearchBooksByTitle)

    // Borrowing endpoints
    app.POST("/borrow/:id", c.BorrowBook)
    app.GET("/download/:id", c.DownloadBook)

    // Return endpoints
    app.POST("/return", c.ReturnBook)
    app.POST("/return/upload", c.ReturnBookWithFile)
}

func (c *perpusController) GetAllBooks(ctx *gin.Context) {
    response, err := c.service.GetAllBooks()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(response.StatusCode, response)
}

func (c *perpusController) GetBookByID(ctx *gin.Context) {
    id := ctx.Param("id")
    bookID, err := strconv.Atoi(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
        return
    }

    book, err := c.service.GetBookByID(bookID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status_code": http.StatusOK,
        "message":     "Berhasil mendapatkan data buku",
        "data":        book,
    })
}

func (c *perpusController) SearchBooksByTitle(ctx *gin.Context) {
    title := ctx.Query("title")
    if title == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Title parameter is required"})
        return
    }

    response, err := c.service.SearchBooksByTitle(title)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(response.StatusCode, response)
}

func (c *perpusController) BorrowBook(ctx *gin.Context) {
    id := ctx.Param("id")
    bookID, err := strconv.Atoi(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
        return
    }

    response, err := c.service.BorrowBook(bookID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(response.StatusCode, response)
}

// Secure download with token validation
func (c *perpusController) DownloadBook(ctx *gin.Context) {
    id := ctx.Param("id")
    bookID, err := strconv.Atoi(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
        return
    }

    token := ctx.Query("token")
    if token == "" {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Download token diperlukan"})
        return
    }

    filePath, err := c.service.GetBookFileWithToken(bookID, token)
    if err != nil {
        ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    ctx.File(filePath)
}

func (c *perpusController) ReturnBook(ctx *gin.Context) {
    var payload dto.ReturnBookRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := c.service.ReturnBook(&payload)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(response.StatusCode, response)
}

func (c *perpusController) ReturnBookWithFile(ctx *gin.Context) {
    var payload dto.ReturnBookRequest

    if err := ctx.ShouldBind(&payload); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    file, err := ctx.FormFile("book_file")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "File buku harus diupload"})
        return
    }

    response, err := c.service.ReturnBookWithFile(&payload, file)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(response.StatusCode, response)
}