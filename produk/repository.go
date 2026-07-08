package produk

import "database/sql"

type ProdukRepository struct {
	DB *sql.DB
}

func NewProdukRepository(db *sql.DB) *ProdukRepository {
	return &ProdukRepository{DB: db}
}

/*
Inisialisasi mysql:
// Driver: github.com/go-sql-driver/mysql
// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=true"
// db, err := sql.Open("mysql", dsn)
*/
