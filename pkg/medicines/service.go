package medicines

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/SYSTEMTerror/GoHealth/pkg/types"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	//ErrNotFound is returned when a customer is not found
	ErrNotFound = errors.New("medecine not found")
	//ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal error")
)

// Service is a medicines service
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a new medicines service
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// GetSomeMedicines returns needed amount of active medicines
func (s *Service) GetSomeMedicines(ctx context.Context, column string, value string, limit int) ([]*types.Medicine, int, error) {
	items := make([]*types.Medicine, 0)
	var sql string
	if column == "" || value == "" {
		sql = fmt.Sprintf("SELECT id, name, manafacturer, description, components, recipe_needed, price, qty, pharmacy_name, pharmacy_phone, pharmacy_address, active, created, image, file FROM medicines WHERE active = true ORDER BY id LIMIT %v", limit)
	} else {
		sql = fmt.Sprintf("SELECT id, name, manafacturer, description, components, recipe_needed, price, qty, pharmacy_name, pharmacy_phone, pharmacy_address, active, created, image, file FROM medicines WHERE %v = '%v' AND active = true ORDER BY id LIMIT %v", column, value, limit)
	}
	rows, err := s.pool.Query(ctx, sql)
	if err == pgx.ErrNoRows {
		log.Println("Medicines s.pool.Query No rows:", err)
		return nil, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println("Medicines s.pool.Query ERROR:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}
	defer rows.Close()
	for rows.Next() {
		item := &types.Medicine{}
		err = rows.Scan(&item.ID, &item.Name, &item.Manafacturer, &item.Description, &item.Components, &item.Recipe_needed, &item.Price, &item.Qty, &item.PharmacyName, &item.PharmacyPhone, &item.PharmacyAddress, &item.Active, &item.Created, &item.Image, &item.File)
		if err != nil {
			log.Println("Medicines rows.Scan ERROR:", err)
			return nil, http.StatusInternalServerError, ErrInternal
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Medicines rows.Err ERROR:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return items, http.StatusOK, nil
}

// Save saves/updates medicine in database
func (s *Service) Save(ctx context.Context, item *types.Medicine) (*types.Medicine, int, error) {
	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `
			INSERT INTO medicines (name, manafacturer, description, components, recipe_needed, price, qty, pharmacy_name, active, image, file)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id
		`, &item.Name, &item.Manafacturer, &item.Description, &item.Components, &item.Recipe_needed, &item.Price, &item.Qty, &item.PharmacyName, &item.Active, &item.Image, &item.File).Scan(&item.ID)
		if err != nil {
			log.Println("Medicines s.pool.QueryRow ERROR:", err)
			return nil, http.StatusInternalServerError, ErrInternal
		}
		return item, http.StatusOK, nil
	}

	_, err := s.pool.Exec(ctx, `
		UPDATE medicines SET name = $1, manafacturer = $2, description = $3, components = $4, recipe_needed = $5, price = $6, qty = $7, pharmacy_name = $8, active = $9, image = $11, file = $12 WHERE id = $10
	`, &item.Name, &item.Manafacturer, &item.Description, &item.Components, &item.Recipe_needed, &item.Price, &item.Qty, &item.PharmacyName, &item.Active, &item.ID, &item.Image, &item.File)
	if err != nil {
		log.Println("Medicines s.pool.Exec ERROR:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}
	return item, http.StatusOK, nil
}

// Order
func (s *Service) Order(ctx context.Context, item *types.Order) (*types.Order, int, error) {
	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `
			INSERT INTO orders (customer_id, medicine_id, pharmacy_name, qty, price)
			VALUES ($1, $2, $3, $4, $5) RETURNING id
		`, &item.CustomerID, &item.MedicineID, &item.PharmacyName, &item.Qty, &item.Price).Scan(&item.ID)
		if err != nil {
			log.Println("Medicines s.pool.QueryRow ERROR:", err)
			return nil, http.StatusInternalServerError, ErrInternal
		}
		return item, http.StatusOK, nil
	}

	_, err := s.pool.Exec(ctx, `
		UPDATE orders SET customer_id = $1, medicine_id = $2, pharmacy_name = $3, qty = $4, price = $5 WHERE id = $6
	`, &item.CustomerID, &item.MedicineID, &item.PharmacyName, &item.Qty, &item.Created, &item.ID)
	if err != nil {
		log.Println("Medicines s.pool.Exec ERROR:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}
	return item, http.StatusOK, nil
}

// GetOrderByID
func (s *Service) GetOrderByID(ctx context.Context, id int64) (*types.Order, int, error) {
	item := &types.Order{}
	err := s.pool.QueryRow(ctx, `
		SELECT id, customer_id, medicine_id, pharmacy_name, qty, price, status FROM orders WHERE id = $1
	`, id).Scan(&item.ID, &item.CustomerID, &item.MedicineID, &item.PharmacyName, &item.Qty, &item.Price, &item.Status)
	if err == pgx.ErrNoRows {
		log.Println("Medicines s.pool.QueryRow No rows:", err)
		return nil, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println("Medicines s.pool.QueryRow ERROR:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return item, http.StatusOK, nil
}

// GelAllOrders
func (s *Service) GetAllOrders(ctx context.Context) ([]*types.Order, int, error) {
	var items []*types.Order
	sql := `
		SELECT id, customer_id, medicine_id, pharmacy_name, qty, price, status, created FROM orders
	`
	rows, err := s.pool.Query(ctx, sql)
	if err == pgx.ErrNoRows {
		log.Println("Medicines s.pool.Query No rows:", err)
		return nil, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println("Medicines s.pool.Query ERROR:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}
	defer rows.Close()
	for rows.Next() {
		item := &types.Order{}
		err := rows.Scan(&item.ID, &item.CustomerID, &item.MedicineID, &item.PharmacyName, &item.Qty, &item.Price, &item.Status, &item.Created)
		if err != nil {
			log.Println("Medicines rows.Scan ERROR:", err)
			return nil, http.StatusInternalServerError, ErrInternal
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Medicines rows.Err ERROR:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return items, http.StatusOK, nil
}

// SetOrderStatus order status
func (s *Service) SetOrderStatus(ctx context.Context, id int64, status string) (int, error) {
	if status == "confirmed" {
		_, err := s.pool.Exec(ctx, `
			UPDATE medicines SET qty = qty - $1 WHERE id = $2
		`, &id, &id)
		if err != nil {
			log.Println("Medicines s.pool.Exec ERROR:", err)
			return http.StatusInternalServerError, ErrInternal
		}
	}

	_, err := s.pool.Exec(ctx, `
		UPDATE orders SET status = $2 WHERE id = $1
	`, &id, &status)
	if err != nil {
		log.Println("Medicines s.pool.Exec ERROR:", err)
		return http.StatusInternalServerError, ErrInternal
	}

	return http.StatusOK, nil
}
