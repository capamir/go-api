package user

import (
	"database/sql"
	"fmt"

	"github.com/capamir/go-api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// maps a single database row into a User struct
// rows.Scan() is like unpacking database columns into struct fields
func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	// rows.Scan() reads the current row's columns into our struct fields
	// The order must match the SELECT statement columns
	err := rows.Scan(
		&user.ID,        // Pass memory addresses so Scan can modify the values
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row into user: %w", err)
	}

	return user, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	// The ? placeholder prevents SQL injection (like Django's parameterized queries)
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for email %s: %w", email, err)
	}
	
	// IMPORTANT: Always close rows to free database resources
	// Similar to closing a file handle or database cursor in other languages
	defer rows.Close()

	var user *types.User
	
	// Returns true if a row is available, false if no more rows
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user for email %s: %w", email, err)
		}
		// Since we expect only one user per email, we could break here
	}

	// This catches errors that happen during rows.Next()
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows for email %s: %w", email, err)
	}

	// Check if we actually found a user
	// this would be like catching User.DoesNotExist
	if user == nil || user.ID == 0 {
		return nil, fmt.Errorf("user not found for email: %s", email)
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
// Similar to Django's User.objects.get(id=id)
func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for user ID %d: %w", id, err)
	}
	defer rows.Close()

	var user *types.User
	
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user for ID %d: %w", id, err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows for user ID %d: %w", id, err)
	}

	if user == nil || user.ID == 0 {
		return nil, fmt.Errorf("user not found for ID: %d", id)
	}

	return user, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}