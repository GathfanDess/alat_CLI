package payment

import "database/sql"

type PaymentRepository struct {
	DB *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

/*
Konfigurasi Driver postgres:
// Driver: github.com/lib/pq
// dsn := "host=localhost user=postgres password=secret dbname=mydb sslmode=disable"
// db, err := sql.Open("postgres", dsn)
*/
