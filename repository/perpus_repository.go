package repository

import (
	"golang-tutorial/contract"
	"golang-tutorial/entity"
	"time"

	"gorm.io/gorm"
)

type perpusRepo struct {
    db *gorm.DB
}

func implPerpusRepository(db *gorm.DB) contract.PerpusRepository {
    return &perpusRepo{
        db: db,
    }
}

func (r *perpusRepo) GetAllBooks() ([]entity.Perpus, error) {
    var books []entity.Perpus
    err := r.db.Table("perpus").Find(&books).Error
    return books, err
}

func (r *perpusRepo) GetBookByID(id int) (*entity.Perpus, error) {
    var book entity.Perpus
    err := r.db.Table("perpus").Where("id = ?", id).First(&book).Error
    if err != nil {
        return nil, err
    }
    return &book, nil
}

func (r *perpusRepo) SearchBooksByTitle(title string) ([]entity.Perpus, error) {
    var books []entity.Perpus
    err := r.db.Table("perpus").Where("judul ILIKE ?", "%"+title+"%").Find(&books).Error
    return books, err
}

func (r *perpusRepo) UpdateBookStatus(id int, status entity.StatusKetersediaan) error {
    return r.db.Table("perpus").Where("id = ?", id).Update("status_ketersediaan", status).Error
}

func (r *perpusRepo) SaveDownloadToken(bookID int, token string, expiresAt int64) error {
    r.db.Exec("DELETE FROM download_tokens WHERE book_id = ?", bookID)
    
    return r.db.Exec(`
        INSERT INTO download_tokens (book_id, token, expires_at, created_at) 
        VALUES (?, ?, ?, ?)
    `, bookID, token, time.Unix(expiresAt, 0), time.Now()).Error
}

func (r *perpusRepo) GetDownloadToken(bookID int) (string, int64, error) {
    var token string
    var expiresAt time.Time
    
    err := r.db.Raw(`
        SELECT token, expires_at FROM download_tokens 
        WHERE book_id = ? AND expires_at > NOW()
    `, bookID).Row().Scan(&token, &expiresAt)
    
    if err != nil {
        return "", 0, err
    }
    
    return token, expiresAt.Unix(), nil
}

func (r *perpusRepo) DeleteDownloadToken(bookID int) error {
    return r.db.Exec("DELETE FROM download_tokens WHERE book_id = ?", bookID).Error
}