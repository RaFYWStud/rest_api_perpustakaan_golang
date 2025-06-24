package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang-tutorial/contract"
	"golang-tutorial/dto"
	"golang-tutorial/entity"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PerpusService struct {
    perpusRepository contract.PerpusRepository
}

func implPerpusService(repo *contract.Repository) contract.PerpusService {
    return &PerpusService{
        perpusRepository: repo.Perpus,
    }
}

func (s *PerpusService) GetAllBooks() (*dto.PerpusListResponse, error) {
    books, err := s.perpusRepository.GetAllBooks()
    if err != nil {
        return nil, err
    }

    var booksData []dto.PerpusData
    for _, book := range books {
        booksData = append(booksData, dto.PerpusData{
            ID:                 book.ID,
            Judul:              book.Judul,
            Penulis:            book.Penulis,
            StatusKetersediaan: book.StatusKetersediaan,
            CreatedAt:          book.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:          book.UpdatedAt.Format("2006-01-02 15:04:05"),
        })
    }

    response := &dto.PerpusListResponse{
        StatusCode: http.StatusOK,
        Message:    "Berhasil mendapatkan daftar buku",
        Data:       booksData,
    }

    return response, nil
}

func (s *PerpusService) GetBookByID(id int) (*dto.PerpusData, error) {
    book, err := s.perpusRepository.GetBookByID(id)
    if err != nil {
        return nil, err
    }

    return &dto.PerpusData{
        ID:                 book.ID,
        Judul:              book.Judul,
        Penulis:            book.Penulis,
        StatusKetersediaan: book.StatusKetersediaan,
        CreatedAt:          book.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:          book.UpdatedAt.Format("2006-01-02 15:04:05"),
    }, nil
}

func (s *PerpusService) SearchBooksByTitle(title string) (*dto.PerpusListResponse, error) {
    books, err := s.perpusRepository.SearchBooksByTitle(title)
    if err != nil {
        return nil, err
    }

    var booksData []dto.PerpusData
    for _, book := range books {
        booksData = append(booksData, dto.PerpusData{
            ID:                 book.ID,
            Judul:              book.Judul,
            Penulis:            book.Penulis,
            StatusKetersediaan: book.StatusKetersediaan,
            CreatedAt:          book.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:          book.UpdatedAt.Format("2006-01-02 15:04:05"),
        })
    }

    response := &dto.PerpusListResponse{
        StatusCode: http.StatusOK,
        Message:    "Berhasil mencari buku",
        Data:       booksData,
    }

    return response, nil
}

func (s *PerpusService) BorrowBook(id int) (*dto.BorrowBookResponse, error) {
    book, err := s.perpusRepository.GetBookByID(id)
    if err != nil {
        return &dto.BorrowBookResponse{
            StatusCode: http.StatusNotFound,
            Message:    "Buku tidak ditemukan",
        }, nil
    }

    if book.StatusKetersediaan != entity.Available {
        return &dto.BorrowBookResponse{
            StatusCode: http.StatusBadRequest,
            Message:    "Buku tidak tersedia untuk dipinjam",
        }, nil
    }

    bookPath := fmt.Sprintf("list_buku/book_%d.pdf", book.ID)
    if _, err := os.Stat(bookPath); os.IsNotExist(err) {
        return &dto.BorrowBookResponse{
            StatusCode: http.StatusBadRequest,
            Message:    "File buku tidak tersedia",
        }, nil
    }

    downloadToken, err := s.generateDownloadToken(book.ID)
    if err != nil {
        return &dto.BorrowBookResponse{
            StatusCode: http.StatusInternalServerError,
            Message:    "Gagal membuat token download",
        }, nil
    }

    err = s.perpusRepository.UpdateBookStatus(id, entity.NotAvailable)
    if err != nil {
        return nil, err
    }

    downloadLink := fmt.Sprintf("/perpus/download/%d?token=%s", book.ID, downloadToken)
    
    response := &dto.BorrowBookResponse{
        StatusCode: http.StatusOK,
        Message:    "Berhasil meminjam buku. Gunakan link download yang disediakan (valid 1 jam).",
        Data: dto.PerpusDataWithLink{
            ID:                 book.ID,
            Judul:              book.Judul,
            Penulis:            book.Penulis,
            StatusKetersediaan: entity.NotAvailable,
            LinkBuku:           downloadLink,
            DownloadToken:      downloadToken,
            CreatedAt:          book.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:          book.UpdatedAt.Format("2006-01-02 15:04:05"),
        },
    }

    return response, nil
}

func (s *PerpusService) generateDownloadToken(bookID int) (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    
    token := hex.EncodeToString(bytes)
    
    expiresAt := time.Now().Add(1 * time.Hour).Unix()
    
    err := s.perpusRepository.SaveDownloadToken(bookID, token, expiresAt)
    if err != nil {
        return "", err
    }
    
    return token, nil
}

func (s *PerpusService) ValidateDownloadToken(bookID int, token string) bool {
    storedToken, expiresAt, err := s.perpusRepository.GetDownloadToken(bookID)
    if err != nil {
        return false
    }
    
    if storedToken != token || time.Now().Unix() > expiresAt {
        return false
    }
    
    return true
}

func (s *PerpusService) GetBookFileWithToken(bookID int, token string) (string, error) {
    if !s.ValidateDownloadToken(bookID, token) {
        return "", fmt.Errorf("token tidak valid atau sudah expired")
    }
    
    book, err := s.perpusRepository.GetBookByID(bookID)
    if err != nil {
        return "", err
    }

    if book.StatusKetersediaan == entity.Available {
        return "", fmt.Errorf("buku harus dipinjam terlebih dahulu")
    }

    return fmt.Sprintf("list_buku/book_%d.pdf", book.ID), nil
}

func (s *PerpusService) ReturnBook(payload *dto.ReturnBookRequest) (*dto.ReturnBookResponse, error) {
    book, err := s.perpusRepository.GetBookByID(payload.BookID)
    if err != nil {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusNotFound,
            Message:    "Buku tidak ditemukan",
        }, nil
    }

    if book.StatusKetersediaan == entity.Available {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusBadRequest,
            Message:    "Buku sudah tersedia, tidak perlu dikembalikan",
        }, nil
    }

    s.perpusRepository.DeleteDownloadToken(payload.BookID)

    err = s.perpusRepository.UpdateBookStatus(payload.BookID, entity.Available)
    if err != nil {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusInternalServerError,
            Message:    "Gagal mengembalikan buku",
        }, nil
    }

    response := &dto.ReturnBookResponse{
        StatusCode: http.StatusOK,
        Message:    "Berhasil mengembalikan buku",
        Data: dto.PerpusData{
            ID:                 book.ID,
            Judul:              book.Judul,
            Penulis:            book.Penulis,
            StatusKetersediaan: entity.Available,
            CreatedAt:          book.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:          book.UpdatedAt.Format("2006-01-02 15:04:05"),
        },
    }

    return response, nil
}

func (s *PerpusService) ReturnBookWithFile(payload *dto.ReturnBookRequest, bookFile *multipart.FileHeader) (*dto.ReturnBookResponse, error) {
    book, err := s.perpusRepository.GetBookByID(payload.BookID)
    if err != nil {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusNotFound,
            Message:    "Buku tidak ditemukan",
        }, nil
    }

    if book.StatusKetersediaan == entity.Available {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusBadRequest,
            Message:    "Buku sudah tersedia, tidak perlu dikembalikan",
        }, nil
    }

    if bookFile == nil {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusBadRequest,
            Message:    "File buku harus diupload untuk pengembalian",
        }, nil
    }

    if err := s.validateBookFile(bookFile); err != nil {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusBadRequest,
            Message:    err.Error(),
        }, nil
    }

    logInfo, err := s.processReturnFile(book, bookFile)
    if err != nil {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusInternalServerError,
            Message:    "Gagal memproses file pengembalian: " + err.Error(),
        }, nil
    }

    s.perpusRepository.DeleteDownloadToken(payload.BookID)

    err = s.perpusRepository.UpdateBookStatus(payload.BookID, entity.Available)
    if err != nil {
        return &dto.ReturnBookResponse{
            StatusCode: http.StatusInternalServerError,
            Message:    "Gagal mengupdate status buku",
        }, nil
    }

    response := &dto.ReturnBookResponse{
        StatusCode: http.StatusOK,
        Message:    "Berhasil mengembalikan buku dengan file upload",
        Data: dto.PerpusData{
            ID:                 book.ID,
            Judul:              book.Judul,
            Penulis:            book.Penulis,
            StatusKetersediaan: entity.Available,
            CreatedAt:          book.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:          book.UpdatedAt.Format("2006-01-02 15:04:05"),
        },
        LogInfo: logInfo,
    }

    return response, nil
}

func (s *PerpusService) validateBookFile(fileHeader *multipart.FileHeader) error {
    if fileHeader.Size > 100*1024*1024 {
        return fmt.Errorf("file terlalu besar, maksimal 100MB")
    }

    allowedExts := map[string]bool{
        ".pdf": true, ".doc": true, ".docx": true, ".txt": true,
    }

    ext := filepath.Ext(fileHeader.Filename)
    if !allowedExts[ext] {
        return fmt.Errorf("tipe file tidak diizinkan: %s. Hanya diperbolehkan: pdf, doc, docx, txt", ext)
    }

    return nil
}

func (s *PerpusService) processReturnFile(book *entity.Perpus, fileHeader *multipart.FileHeader) (*dto.ReturnLogInfo, error) {
    timestamp := time.Now()
    timestampStr := timestamp.Format("20060102_150405")

    bookFileName := fmt.Sprintf("book_%d_%s_%s", book.ID, timestampStr, filepath.Base(fileHeader.Filename))
    logFileName := fmt.Sprintf("log_%d_%s.txt", book.ID, timestampStr)

    logDir := "log_pengembalian_buku"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, fmt.Errorf("gagal membuat direktori log: %v", err)
    }

    bookPath := filepath.Join(logDir, bookFileName)
    logPath := filepath.Join(logDir, logFileName)

    if err := s.saveUploadedFile(fileHeader, bookPath); err != nil {
        return nil, fmt.Errorf("gagal menyimpan file buku: %v", err)
    }

    logData := map[string]interface{}{
        "book_id":       book.ID,
        "judul":         book.Judul,
        "penulis":       book.Penulis,
        "returned_at":   timestamp.Format("2006-01-02 15:04:05"),
        "book_file":     bookFileName,
        "book_path":     bookPath,
        "file_size":     fileHeader.Size,
        "original_name": fileHeader.Filename,
    }

    if err := s.saveLogFile(logPath, logData); err != nil {
        return nil, fmt.Errorf("gagal menyimpan log: %v", err)
    }

    return &dto.ReturnLogInfo{
        LogFileName:  logFileName,
        BookFileName: bookFileName,
        LogPath:      logPath,
        BookPath:     bookPath,
        ReturnedAt:   timestamp.Format("2006-01-02 15:04:05"),
    }, nil
}

func (s *PerpusService) saveUploadedFile(fileHeader *multipart.FileHeader, dst string) error {
    src, err := fileHeader.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    return err
}

func (s *PerpusService) saveLogFile(path string, data map[string]interface{}) error {
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, jsonData, 0644)
}