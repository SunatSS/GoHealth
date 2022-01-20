package medicines

import (
	"context"
	"errors"
	"fmt"
	"log"

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
func (s *Service) GetSomeMedicines(ctx context.Context, column string, value string, limit int) ([]*types.Medicine, error) {
	items := make([]*types.Medicine, 0)
	var sql string
	if column == "" || value == "" {
		sql = fmt.Sprintf("SELECT id, name, manafacturer, description, components, recipe_needed, price, qty, pharmacy_name, active, created, image, file FROM medicines WHERE active = true ORDER BY id LIMIT %v",limit)
	} else {
		sql = fmt.Sprintf("SELECT id, name, manafacturer, description, components, recipe_needed, price, qty, pharmacy_name, active, created, image, file FROM medicines WHERE %v = '%v' AND active = true ORDER BY id LIMIT %v", column, value, limit)
	}
	rows, err := s.pool.Query(ctx, sql)
	if err == pgx.ErrNoRows {
		log.Println("Medicines s.pool.Query No rows:", err)
		return nil, ErrNotFound
	}
	if err != nil {
		log.Println("Medicines s.pool.Query ERROR:", err)
		return nil, ErrInternal
	}
	defer rows.Close()
	for rows.Next() {
		item := &types.Medicine{}
		err = rows.Scan(&item.ID, &item.Name, &item.Manafacturer, &item.Description, &item.Components, &item.Recipe_needed, &item.Price, &item.Qty, &item.PharmacyName, &item.Active, &item.Created, &item.Image, &item.File)
		if err != nil {
			log.Println("Medicines rows.Scan ERROR:", err)
			return nil, ErrInternal
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Medicines rows.Err ERROR:", err)
		return nil, ErrInternal
	}

	return items, nil
}

// Save saves/updates medicine in database
func (s *Service) Save(ctx context.Context, item *types.Medicine) (*types.Medicine, error) {
	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `
			INSERT INTO medicines (name, manafacturer, description, components, recipe_needed, price, qty, pharmacy_name, active, image, file)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id
		`, &item.Name, &item.Manafacturer, &item.Description, &item.Components, &item.Recipe_needed, &item.Price, &item.Qty, &item.PharmacyName, &item.Active, &item.Image, &item.File).Scan(&item.ID)
		if err != nil {
			log.Println("Medicines s.pool.QueryRow ERROR:", err)
			return nil, ErrInternal
		}
		return item, nil
	} else {
		_, err := s.pool.Exec(ctx, `
			UPDATE medicines SET name = $1, manafacturer = $2, description = $3, components = $4, recipe_needed = $5, price = $6, qty = $7, pharmacy_name = $8, active = $9, image = $11, file = $12 WHERE id = $10
		`, &item.Name, &item.Manafacturer, &item.Description, &item.Components, &item.Recipe_needed, &item.Price, &item.Qty, &item.PharmacyName, &item.Active, &item.ID, &item.Image, &item.File)
		if err != nil {
			log.Println("Medicines s.pool.Exec ERROR:", err)
			return nil, ErrInternal
		}
		return item, nil
	}
}