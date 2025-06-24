## Base URL

```
http://localhost:8080
```

### Available Operations

-   **Browse Books**: View and search available books
-   **Borrow Books**: Get secure download tokens
-   **Download Books**: Access files with valid tokens
-   **Return Books**: Simple or with file upload logging

---

## 1. 📖 Browse Books

### Get All Books

**Endpoint:** `GET /perpus/books`  
**Description:** Retrieve all books in the library database

**Postman Setup:**

-   Method: `GET`
-   URL: `http://localhost:8080/perpus/books`
-   Headers: None required

**Response Example:**

```json
{
    "status_code": 200,
    "message": "Berhasil mendapatkan daftar buku",
    "data": [
        {
            "id": 1,
            "judul": "Animal Farm",
            "penulis": "George Orwell",
            "status_ketersediaan": "available",
            "created_at": "2024-12-01 10:30:25",
            "updated_at": "2024-12-01 10:30:25"
        },
        {
            "id": 2,
            "judul": "Metamorphosis",
            "penulis": "Franz Kafka",
            "status_ketersediaan": "not_available",
            "created_at": "2024-12-01 10:30:25",
            "updated_at": "2024-12-01 14:20:10"
        }
    ]
}
```

### Get Book by ID

**Endpoint:** `GET /perpus/books/{id}`  
**Description:** Get detailed information about a specific book

**Postman Setup:**

-   Method: `GET`
-   URL: `http://localhost:8080/perpus/books/1`
-   Path Variable: `id = 1`

**Response Example:**

```json
{
    "status_code": 200,
    "message": "Berhasil mendapatkan data buku",
    "data": {
        "id": 1,
        "judul": "Animal Farm",
        "penulis": "George Orwell",
        "status_ketersediaan": "available",
        "created_at": "2024-12-01 10:30:25",
        "updated_at": "2024-12-01 10:30:25"
    }
}
```

### Search Books by Title

**Endpoint:** `GET /perpus/search?title={keyword}`  
**Description:** Search books by title using partial matching

**Postman Setup:**

-   Method: `GET`
-   URL: `http://localhost:8080/perpus/search`
-   Query Parameters: `title = Animal`

**Response Example:**

```json
{
    "status_code": 200,
    "message": "Berhasil mencari buku",
    "data": [
        {
            "id": 1,
            "judul": "Animal Farm",
            "penulis": "George Orwell",
            "status_ketersediaan": "available",
            "created_at": "2024-12-01 10:30:25",
            "updated_at": "2024-12-01 10:30:25"
        }
    ]
}
```

---

## 2. 📋 Borrow Book (Required for Download)

### Borrow Book

**Endpoint:** `POST /perpus/borrow/{id}`  
**Description:** Borrow a book and receive a secure download token (expires in 1 hour)

**Postman Setup:**

-   Method: `POST`
-   URL: `http://localhost:8080/perpus/borrow/1`
-   Path Variable: `id = 1`
-   Headers: None required
-   Body: None required

**Success Response:**

```json
{
    "status_code": 200,
    "message": "Berhasil meminjam buku. Gunakan link download yang disediakan (valid 1 jam).",
    "data": {
        "id": 1,
        "judul": "Animal Farm",
        "penulis": "George Orwell",
        "link_buku": "/perpus/download/1?token=a1b2c3d4e5f6789abcdef1234567890",
        "download_token": "a1b2c3d4e5f6789abcdef1234567890",
        "status_ketersediaan": "not_available",
        "created_at": "2024-12-01 10:30:25",
        "updated_at": "2024-12-01 14:30:25"
    }
}
```

**Error Responses:**

```json
// Book not found
{
    "status_code": 404,
    "message": "Buku tidak ditemukan"
}

// Book already borrowed
{
    "status_code": 400,
    "message": "Buku tidak tersedia untuk dipinjam"
}

// Book file missing
{
    "status_code": 400,
    "message": "File buku tidak tersedia"
}
```

---

## 3. 📥 Download Book (Secure)

### Download Book File

**Endpoint:** `GET /perpus/download/{id}?token={token}`  
**Description:** Download book file using valid token from borrow response

**Postman Setup:**

-   Method: `GET`
-   URL: `http://localhost:8080/perpus/download/1?token=a1b2c3d4e5f6789abcdef1234567890`
-   Path Variable: `id = 1`
-   Query Parameter: `token = a1b2c3d4e5f6789abcdef1234567890`

**Success Response:** File download starts automatically

**Error Responses:**

```json
// Missing token
{
    "error": "Download token diperlukan"
}

// Invalid/expired token
{
    "error": "token tidak valid atau sudah expired"
}

// Book not borrowed
{
    "error": "buku harus dipinjam terlebih dahulu"
}
```

**⚠️ Important Notes:**

-   Token expires in 1 hour after borrowing
-   Must borrow book first before downloading
-   Each borrow generates a new unique token

---

## 4. 📤 Return Book

### Simple Return (JSON)

**Endpoint:** `POST /perpus/return`  
**Description:** Return borrowed book without file upload

**Postman Setup:**

-   Method: `POST`
-   URL: `http://localhost:8080/perpus/return`
-   Headers: `Content-Type: application/json`
-   Body (raw JSON):

```json
{
    "book_id": 1
}
```

**Response Example:**

```json
{
    "status_code": 200,
    "message": "Berhasil mengembalikan buku",
    "data": {
        "id": 1,
        "judul": "Animal Farm",
        "penulis": "George Orwell",
        "status_ketersediaan": "available",
        "created_at": "2024-12-01 10:30:25",
        "updated_at": "2024-12-01 14:30:25"
    }
}
```

### Return with File Upload

**Endpoint:** `POST /perpus/return/upload`  
**Description:** Return book with file upload and activity logging

**Postman Setup:**

-   Method: `POST`
-   URL: `http://localhost:8080/perpus/return/upload`
-   Body Type: `form-data`
-   Form Fields:
    -   `book_id`: `1` (Text)
    -   `book_file`: [Select File] (File - PDF/DOC/DOCX/TXT, max 100MB)

**Response Example:**

```json
{
    "status_code": 200,
    "message": "Berhasil mengembalikan buku dengan file upload",
    "data": {
        "id": 1,
        "judul": "Animal Farm",
        "penulis": "George Orwell",
        "status_ketersediaan": "available",
        "created_at": "2024-12-01 10:30:25",
        "updated_at": "2024-12-01 14:30:25"
    },
    "log_info": {
        "log_file_name": "log_1_20241201_143025.txt",
        "book_file_name": "book_1_20241201_143025_yourfile.pdf",
        "log_path": "log_pengembalian_buku/log_1_20241201_143025.txt",
        "book_path": "log_pengembalian_buku/book_1_20241201_143025_yourfile.pdf",
        "returned_at": "2024-12-01 14:30:25"
    }
}
```

**File Upload Restrictions:**

-   Maximum file size: 100MB
-   Allowed formats: PDF, DOC, DOCX, TXT
-   Files are stored in `log_pengembalian_buku/` directory

---

## 🔄 Complete Usage Workflow

### Standard Library Process:

1. **Browse Available Books**

    ```
    GET /perpus/books
    ```

2. **Search Specific Book** (Optional)

    ```
    GET /perpus/search?title=Animal
    ```

3. **Borrow Book** (Required for download)

    ```
    POST /perpus/borrow/1
    ```

    → Receive secure download link + token

4. **Download Book** (Using token from step 3)

    ```
    GET /perpus/download/1?token=your_token_here
    ```

5. **Return Book** (Choose one method)

    **Method A - Simple Return:**

    ```
    POST /perpus/return
    Body: {"book_id": 1}
    ```

    **Method B - Return with File Upload:**

    ```
    POST /perpus/return/upload
    Form-data: book_id=1, book_file=[your_file]
    ```

---

## 📊 Sample Books Database

| ID  | Title                  | Author             | Status    |
| --- | ---------------------- | ------------------ | --------- |
| 1   | Animal Farm            | George Orwell      | available |
| 2   | Metamorphosis          | Franz Kafka        | available |
| 3   | The Alchemist          | Paulo Coelho       | available |
| 4   | To Kill A Mocking Bird | Harper Lee         | available |
| 5   | White Nights           | Fyodor Dostoyevsky | available |

---

## 🔒 Security Features

### Token-Based Download System

-   **Secure Access**: Download only available through valid tokens
-   **Time-Limited**: Tokens expire after 1 hour
-   **Unique Tokens**: Each borrow generates cryptographically secure token
-   **Auto Cleanup**: Tokens automatically deleted on book return

### File Upload Security

-   **Size Validation**: Maximum 100MB per file
-   **Type Validation**: Only PDF, DOC, DOCX, TXT allowed
-   **Safe Storage**: Files stored with timestamp and unique naming

### Activity Logging

-   **Complete Tracking**: All return activities logged with metadata
-   **JSON Format**: Structured logging for easy parsing
-   **File Preservation**: Uploaded files preserved for audit trail

---

## ❌ Common Error Responses

```json
// Book not found
{
    "status_code": 404,
    "message": "Buku tidak ditemukan"
}

// Book already borrowed
{
    "status_code": 400,
    "message": "Buku tidak tersedia untuk dipinjam"
}

// Book already available (cannot return)
{
    "status_code": 400,
    "message": "Buku sudah tersedia, tidak perlu dikembalikan"
}

// File validation error
{
    "status_code": 400,
    "message": "File terlalu besar, maksimal 100MB"
}

// Invalid request format
{
    "error": "Key: 'ReturnBookRequest.BookID' Error:Field validation for 'BookID' failed on the 'required' tag"
}
```

---

## 🛠️ Postman Collection Setup

### Environment Variables

Create environment in Postman with these variables:

-   `base_url`: `http://localhost:8080`
-   `book_id`: `1` (for testing)
-   `download_token`: (will be populated from borrow response)

### Collection Structure

```
Library Management API/
├── 📚 Browse Books/
│   ├── Get All Books
│   ├── Get Book by ID
│   └── Search Books
├── 📋 Borrow Book/
│   └── Borrow Book by ID
├── 📥 Download Book/
│   └── Download with Token
└── 📤 Return Book/
    ├── Simple Return
    └── Return with File Upload
```

### Testing Tips

1. Always browse books first to see available options
2. Copy `download_token` from borrow response for download requests
3. Use environment variables for consistent testing
4. Test error scenarios (invalid IDs, expired tokens, etc.)
5. Verify file uploads work with different file types

---

## 📁 File System Structure

```
project_root/
├── list_buku/                  # Original book files
│   ├── book_1.pdf             # Animal Farm
│   ├── book_2.pdf             # Metamorphosis
│   ├── book_3.pdf             # The Alchemist
│   ├── book_4.pdf             # To Kill A Mocking Bird
│   └── book_5.pdf             # White Nights
└── log_pengembalian_buku/      # Return activity logs
    ├── book_1_20241201_143025_upload.pdf
    ├── log_1_20241201_143025.txt
    └── ...
```

This secure library system ensures proper book lending workflow while maintaining detailed activity logs and
