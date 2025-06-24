package contract

import "golang-tutorial/entity"

type Repository struct {
	Perpus PerpusRepository
}

type PerpusRepository interface {
    GetAllBooks() ([]entity.Perpus, error)
    GetBookByID(id int) (*entity.Perpus, error)
    SearchBooksByTitle(title string) ([]entity.Perpus, error)
    UpdateBookStatus(id int, status entity.StatusKetersediaan) error
    SaveDownloadToken(bookID int, token string, expiresAt int64) error
    GetDownloadToken(bookID int) (string, int64, error)
    DeleteDownloadToken(bookID int) error
}
