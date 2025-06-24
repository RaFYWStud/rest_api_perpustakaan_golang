package migrations

import "database/sql"

type createDownloadTokensTable struct{}

func (m *createDownloadTokensTable) SkipProd() bool {
    return false
}

func getCreateDownloadTokensTable() migration {
    return &createDownloadTokensTable{}
}

func (m *createDownloadTokensTable) Name() string {
    return "create-download-tokens"
}

func (m *createDownloadTokensTable) Up(conn *sql.Tx) error {
    _, err := conn.Exec(`
        CREATE TABLE download_tokens (
            id SERIAL PRIMARY KEY,
            book_id INTEGER REFERENCES perpus(id) ON DELETE CASCADE,
            token VARCHAR(64) NOT NULL UNIQUE,
            expires_at TIMESTAMP NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        );
    `)
    if err != nil {
        return err
    }

    // Create indexes for better performance
    _, err = conn.Exec(`
        CREATE INDEX idx_download_tokens_book_id ON download_tokens(book_id);
        CREATE INDEX idx_download_tokens_token ON download_tokens(token);
        CREATE INDEX idx_download_tokens_expires_at ON download_tokens(expires_at);
    `)
    
    return err
}

func (m *createDownloadTokensTable) Down(conn *sql.Tx) error {
    _, err := conn.Exec("DROP TABLE IF EXISTS download_tokens;")
    return err
}