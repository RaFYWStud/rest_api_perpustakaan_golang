package migrations

import "database/sql"

type createPerpusTable struct{}

func (m *createPerpusTable) SkipProd() bool {
    return false
}

func getCreatePerpusTable() migration {
    return &createPerpusTable{}
}

func (m *createPerpusTable) Name() string {
    return "create-perpus"
}

func (m *createPerpusTable) Up(conn *sql.Tx) error {
    // Create ENUM type for book availability status
    _, err := conn.Exec(`
        CREATE TYPE ketersediaan AS ENUM ('available', 'not_available');
    `)
    if err != nil {
        return err
    }

    // Create perpus table
    _, err = conn.Exec(`
        CREATE TABLE perpus (
            id SERIAL PRIMARY KEY,
            judul VARCHAR(255) NOT NULL,
            penulis VARCHAR(100) NOT NULL,
            status_ketersediaan ketersediaan NOT NULL DEFAULT 'available',
            link_buku VARCHAR(500),
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP NOT NULL DEFAULT NOW()
        );
    `)
    if err != nil {
        return err
    }

    // Insert sample book data
    _, err = conn.Exec(`
        INSERT INTO perpus (judul, penulis, status_ketersediaan, link_buku) VALUES
        ('Animal Farm', 'George Orwell', 'available', '/books/book_1.pdf'),
        ('Metamorphosis', 'Franz Kafka', 'available', '/books/book_2.pdf'),
        ('The Alchemist', 'Paulo Coelho', 'available', '/books/book_3.pdf'),
        ('To Kill A Mocking Bird', 'Harper Lee', 'available', '/books/book_4.pdf'),
        ('White Nights', 'Fyodor Dostoyevsky', 'available', '/books/book_5.pdf')
    `)
    if err != nil {
        return err
    }

    return nil
}

func (m *createPerpusTable) Down(conn *sql.Tx) error {
    _, err := conn.Exec("DROP TABLE IF EXISTS perpus;")
    if err != nil {
        return err
    }

    _, err = conn.Exec("DROP TYPE IF EXISTS ketersediaan;")
    if err != nil {
        return err
    }

    return nil
}