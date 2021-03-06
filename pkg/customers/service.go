package customers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
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
	//ErrEmptyName is returned when name is empty
	ErrEmptyName = errors.New("empty name")
	//ErrEmptyPassword is returned when password is empty
	ErrEmptyPassword = errors.New("empty password")
	//ErrEmptyAddress is returned when address is empty
	ErrEmptyAddress = errors.New("empty address")
)

//Service is a customers service
type Service struct {
	pool *pgxpool.Pool
}

//NewService creates new customers service
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// RegisterCustomer registers customer
func (s *Service) RegisterCustomer(ctx context.Context, item *types.RegInfo) (*types.Customer, int, error) {
	customer := &types.Customer{}

	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Save bcrypt.GenerateFromPassword Error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
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
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return customer, http.StatusOK, nil
}

// Token generates token for customer
func (s *Service) Token(ctx context.Context, item *types.TokenInfo) (*types.Token, int, error) {
	var hash string
	token := &types.Token{}
	err := s.pool.QueryRow(ctx, `SELECT id, password FROM customers WHERE phone = $1`, item.Login).Scan(&token.CustomerID, &hash)
	if err == pgx.ErrNoRows {
		log.Println("Token s.pool.QueryRow error:", err)
		return nil, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(item.Password))
	if err != nil {
		log.Println("Token bcrypt.CompareHashAndPassword error:", err)
		return nil, http.StatusUnauthorized, ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		log.Println("Token rand.Read len : %w (must be 256), error: %w", n, err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	token.Token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO customers_tokens (customer_id, token) VALUES ($1, $2)`, token.CustomerID, token.Token)
	if err != nil {
		log.Println("Token s.pool.Exec error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return token, http.StatusOK, nil
}

// IDByToken returns customer id by token
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	var expires time.Time

	err := s.pool.QueryRow(ctx, `SELECT customer_id, expires FROM customers_tokens WHERE token = $1`, token).Scan(&id, &expires)
	if err == pgx.ErrNoRows {
		log.Println("IDByToken s.pool.QueryRow No rows:", err)
		return 0, nil
	}
	if err != nil || expires.Before(time.Now()) {
		log.Println("IDByToken s.pool.QueryRow error:", err)
		return 0, ErrInternal
	}

	return id, nil
}

// EditCustomer edits customer
func (s *Service) EditCustomer(ctx context.Context, item *types.Customer) (int, error) {
	err := validate(item)
	if err != nil {
		log.Println("EditCustomer validate error:", err)
		return http.StatusBadRequest, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("EditCustomer bcrypt.GenerateFromPassword Error:", err)
		return http.StatusInternalServerError, ErrInternal
	}
	item.Password = string(hash)
	err = s.pool.QueryRow(ctx, `
		UPDATE customers SET name = $2, password = $3, address = $4 WHERE id = $1 RETURNING id
	`,item.ID, item.Name, item.Password, item.Address).Scan(&item.ID)
	if err == pgx.ErrNoRows {
		log.Println("EditCustomer s.pool.QueryRow No rows:", err)
		return http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println("EditCustomer s.pool.QueryRow error:", err)
		return http.StatusInternalServerError, ErrInternal
	}

	return http.StatusOK, nil
}

// Validates customer data
func validate(item *types.Customer) error {
	if item.Name == "" {
		return ErrEmptyName
	}
	if item.Password == "" {
		return ErrEmptyPassword
	}
	if item.Address == "" {
		return ErrEmptyAddress
	}
	return nil
}

// IsAdmin checks if customer is admin
func (s *Service) IsAdmin(ctx context.Context, id int64) (bool, int, error) {
	var isAdmin bool
	err := s.pool.QueryRow(ctx, `SELECT is_admin FROM customers WHERE id = $1`, id).Scan(&isAdmin)
	if err == pgx.ErrNoRows {
		log.Println("IsAdmin s.pool.QueryRow No rows:", err)
		return false, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println("IsAdmin s.pool.QueryRow error:", err)
		return false, http.StatusInternalServerError, ErrInternal
	}

	return isAdmin, http.StatusOK, nil
}

// MakeAdmin makes customer admin
func (s *Service) MakeAdmin(ctx context.Context, makeAdminInfo *types.MakeAdminInfo) (int, error) {
	_, err := s.pool.Exec(ctx, `UPDATE customers SET is_admin = $2 WHERE id = $1`, makeAdminInfo.ID, makeAdminInfo.AdminStatus)
	if err != nil {
		log.Println("MakeAdmin s.pool.Exec error:", err)
		return http.StatusInternalServerError, ErrInternal
	}

	return http.StatusOK, nil
}

// GetCustomerById returns customer by id
func (s *Service) GetCustomerByID(ctx context.Context, id int64) (*types.Customer, int, error) {
	customer := &types.Customer{}
	err := s.pool.QueryRow(ctx, `SELECT id, name, phone, address, password, is_admin, active, created FROM customers WHERE id = $1`, id).Scan(
		&customer.ID, &customer.Name, &customer.Phone, &customer.Address, &customer.Password, &customer.IsAdmin, &customer.Active, &customer.Created)
	if err == pgx.ErrNoRows {
		log.Println("GetCustomerByID s.pool.QueryRow No rows:", err)
		return nil, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println("GetCustomerByID s.pool.QueryRow error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return customer, http.StatusOK, nil
}

//GetAllCustomers returns all customers
func (s *Service) GetAllCustomers(ctx context.Context) ([]*types.Customer, int, error) {
	var customers []*types.Customer
	rows, err := s.pool.Query(ctx, `SELECT id, name, phone, address, password, is_admin, active, created FROM customers`)
	if err != nil {
		log.Println("GetAllCustomers s.pool.Query error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		customer := &types.Customer{}
		err := rows.Scan(
			&customer.ID, &customer.Name, &customer.Phone, &customer.Address, &customer.Password, &customer.IsAdmin, &customer.Active, &customer.Created)
		if err != nil {
			log.Println("GetAllCustomers rows.Scan error:", err)
			return nil, http.StatusInternalServerError, ErrInternal
		}
		customers = append(customers, customer)
	}

	return customers, http.StatusOK, nil
}
