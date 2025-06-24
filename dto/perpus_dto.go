package dto

import "golang-tutorial/entity"

type PerpusData struct {
    ID                 int                        `json:"id"`
    Judul              string                     `json:"judul"`
    Penulis            string                     `json:"penulis"`
    StatusKetersediaan entity.StatusKetersediaan `json:"status_ketersediaan"`
    CreatedAt          string                     `json:"created_at"`
    UpdatedAt          string                     `json:"updated_at"`
}

type PerpusDataWithLink struct {
    ID                 int                        `json:"id"`
    Judul              string                     `json:"judul"`
    Penulis            string                     `json:"penulis"`
    LinkBuku           string                     `json:"link_buku"`
    DownloadToken      string                     `json:"download_token"` // Token untuk keamanan
    StatusKetersediaan entity.StatusKetersediaan `json:"status_ketersediaan"`
    CreatedAt          string                     `json:"created_at"`
    UpdatedAt          string                     `json:"updated_at"`
}

type PerpusListResponse struct {
    StatusCode int          `json:"status_code"`
    Message    string       `json:"message"`
    Data       []PerpusData `json:"data"`
}

type BorrowBookResponse struct {
    StatusCode int                `json:"status_code"`
    Message    string             `json:"message"`
    Data       PerpusDataWithLink `json:"data,omitempty"`
}

type ReturnBookRequest struct {
    BookID int `json:"book_id" form:"book_id" binding:"required"`
}

type ReturnBookResponse struct {
    StatusCode int                `json:"status_code"`
    Message    string             `json:"message"`
    Data       PerpusData         `json:"data"`
    LogInfo    *ReturnLogInfo     `json:"log_info,omitempty"`
}

type ReturnLogInfo struct {
    LogFileName  string `json:"log_file_name"`
    BookFileName string `json:"book_file_name"`
    LogPath      string `json:"log_path"`
    BookPath     string `json:"book_path"`
    ReturnedAt   string `json:"returned_at"`
}