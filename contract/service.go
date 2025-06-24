package contract

import (
	"golang-tutorial/dto"
	"mime/multipart"
)

type Service struct {
    Perpus PerpusService
}

type PerpusService interface {
    GetAllBooks() (*dto.PerpusListResponse, error)
    GetBookByID(id int) (*dto.PerpusData, error) // Fixed: changed return type
    SearchBooksByTitle(title string) (*dto.PerpusListResponse, error)
    BorrowBook(id int) (*dto.BorrowBookResponse, error)
    ReturnBook(payload *dto.ReturnBookRequest) (*dto.ReturnBookResponse, error) // Fixed: changed return type
    ReturnBookWithFile(payload *dto.ReturnBookRequest, bookFile *multipart.FileHeader) (*dto.ReturnBookResponse, error)
    GetBookFileWithToken(bookID int, token string) (string, error) // Added: missing method
    ValidateDownloadToken(bookID int, token string) bool // Added: missing method
}