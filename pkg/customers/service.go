package customers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/SYSTEMTerror/GoHealth/pkg/types"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	//ErrNotFound is returned when a customer is not found
	ErrNotFound = errors.New("customer not found")
	//ErrInvalidPassword is returned when password is incorrect
	ErrInvalidPassword = errors.New("invalid password")
	//ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal error")
	//ErrExpired is returned when token is expired
	ErrExpired = errors.New("expired")
)

//Service is structure of customers service
type Service struct {
	pool *pgxpool.Pool
}

//NewService returns a new instance of Service
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// RegisterCustomer method registers customer
func (s *Service) RegisterCustomer(ctx context.Context, item *types.RegInfo) (*types.Customer, error) {
	customer := &types.Customer{}

	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternal
	}

	item.Password = string(hash)
	err = s.pool.QueryRow(ctx, `
			INSERT INTO customers (name, phone, password, address) VALUES ($1, $2, $3, $4)
			ON CONFLICT (phone) DO NOTHING
			RETURNING id, name, phone, password, address, active, created
		`, item.Name, item.Phone, item.Password, item.Address).Scan(
		&customer.ID, &customer.Name, &customer.Phone, &customer.Password,
		&customer.Address, &customer.Active, &customer.Created)
	if err != nil {
		log.Println("Save with id == 0 s.pool.QueryRow error:", err)
		return nil, ErrInternal
	}

	return customer, nil
}

// Token method generates token for customer
func (s *Service) Token(ctx context.Context, item *types.TokenInfo) (*types.Token, error) {
	var hash string
	token := &types.Token{}
	err := s.pool.QueryRow(ctx, `SELECT id, password FROM customers WHERE phone = $1`, item.Login).Scan(&token.CustomerID, &hash)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(item.Password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return nil, ErrInternal
	}

	token.Token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO customers_tokens (customer_id, token) VALUES ($1, $2)`, token.CustomerID, token.Token)
	if err != nil {
		return nil, ErrInternal
	}

	return token, nil
}

func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	var expires time.Time

	err := s.pool.QueryRow(ctx, `SELECT customer_id, expires FROM customers_tokens WHERE token = $1`, token).Scan(&id, &expires)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil || expires.Before(time.Now()) {
		return 0, ErrInternal
	}

	return id, nil
}

//EditCustomer method edits customer
func (s *Service) EditCustomer(ctx context.Context, item *types.Customer) (error) {
	sqlBase := "UPDATE customers SET {col} = $1 WHERE id = $2 RETURNING id"
	if item.Name != "" {
		sql := strings.ReplaceAll(sqlBase, "{col}", "name")
		err := s.pool.QueryRow(ctx, sql, item.Name, item.ID).Scan(&item.ID)
		if err == pgx.ErrNoRows {
			return ErrNotFound
		} else if err != nil {
			return ErrInternal
		}
	}
	if item.Password != "" {
		sql := strings.ReplaceAll(sqlBase, "{col}", "password")
		err := s.pool.QueryRow(ctx, sql, item.Password, item.ID).Scan(&item.ID)
		if err == pgx.ErrNoRows {
			return ErrNotFound
		} else if err != nil {
			return ErrInternal
		}
	}
	if item.Address != "" {
		sql := strings.ReplaceAll(sqlBase, "{col}", "address")
		err := s.pool.QueryRow(ctx, sql, item.Address, item.ID).Scan(&item.ID)
		if err == pgx.ErrNoRows {
			return ErrNotFound
		} else if err != nil {
			return ErrInternal
		}
	}

	return nil
}

func (s *Service) IsAdmin(ctx context.Context, id int64) (bool, error) {
	var isAdmin bool
	err := s.pool.QueryRow(ctx, `SELECT is_admin FROM customers WHERE id = $1`, id).Scan(&isAdmin)
	if err == pgx.ErrNoRows {
		return false, ErrNotFound
	} else if err != nil {
		return false, ErrInternal
	}

	return isAdmin, nil
}

func (s *Service) MakeAdmin(ctx context.Context, id int64) (error) {
	_, err := s.pool.Exec(ctx, `UPDATE customers SET is_admin = true WHERE id = $1`, id)
	if err == pgx.ErrNoRows {
		return ErrNotFound
	} else if err != nil {
		return ErrInternal
	}

	return nil
}

func (s *Service) GetCustomerByID(ctx context.Context, id int64) (*types.Customer, error) {
	customer := &types.Customer{}
	err := s.pool.QueryRow(ctx, `SELECT id, name, phone, address, active, created FROM customers WHERE id = $1`, id).Scan(
		&customer.ID, &customer.Name, &customer.Phone, &customer.Address, &customer.Active, &customer.Created)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, ErrInternal
	}

	return customer, nil
}