package user

import (
	"database/sql"
	"fmt"

	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// scanRowIntoUser maps a single database row into a User struct
// This is similar to Django's model hydration from database rows
func scanRowIntoUser(row interface {
	Scan(dest ...interface{}) error
}) (*types.User, error) {
	user := new(types.User)

	// Scan reads the row's columns into our struct fields
	// Order must match the SELECT statement columns
	err := row.Scan(
		&user.ID,
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

// GetUserByEmail retrieves a user by their email address
// Similar to Django's User.objects.get(email=email)
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	// Use explicit column names instead of SELECT *
	query := `
		SELECT id, first_name, last_name, email, password, created_at 
		FROM users 
		WHERE email = ?
	`

	// QueryRow is better for single-row queries (returns only one row)
	row := s.db.QueryRow(query, email)

	user, err := scanRowIntoUser(row)
	if err == sql.ErrNoRows {
		// User not found - this is like Django's DoesNotExist
		utils.S.Debugf("User not found for email: %s", email)
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		utils.S.Errorf("Database error getting user by email %s: %v", email, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	utils.S.Debugf("Found user %d for email: %s", user.ID, email)
	return user, nil
}

// GetUserByID retrieves a user by their ID
// Similar to Django's User.objects.get(id=id)
func (s *Store) GetUserByID(id int) (*types.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at 
		FROM users 
		WHERE id = ?
	`

	row := s.db.QueryRow(query, id)

	user, err := scanRowIntoUser(row)
	if err == sql.ErrNoRows {
		utils.S.Debugf("User not found for ID: %d", id)
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		utils.S.Errorf("Database error getting user by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	utils.S.Debugf("Found user %d (%s)", user.ID, user.Email)
	return user, nil
}

// CreateUser creates a new user record in the database
// Similar to Django's User.objects.create()
func (s *Store) CreateUser(user types.User) error {
	// Use snake_case column names to match database schema
	query := `
		INSERT INTO users (first_name, last_name, email, password, created_at) 
		VALUES (?, ?, ?, ?, NOW())
	`

	res, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		utils.S.Errorf("Failed to create user with email %s: %v", user.Email, err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Get the auto-generated ID (optional but useful for logging)
	id, err := res.LastInsertId()
	if err == nil {
		utils.S.Successf("Created user %d with email: %s", id, user.Email)
	}

	return nil
}

// GetUsers retrieves all users from the database
// Similar to Django's User.objects.all()
func (s *Store) GetUsers() ([]*types.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at 
		FROM users 
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		utils.S.Errorf("Failed to query all users: %v", err)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*types.User
	for rows.Next() {
		user := new(types.User)
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		)
		if err != nil {
			utils.S.Warnf("Failed to scan user row: %v", err)
			continue // Skip this row but continue with others
		}
		users = append(users, user)
	}

	// Check for errors that occurred during iteration
	if err = rows.Err(); err != nil {
		utils.S.Errorf("Error iterating user rows: %v", err)
		return nil, fmt.Errorf("error reading users: %w", err)
	}

	utils.S.Debugf("Retrieved %d users from database", len(users))
	return users, nil
}

// UpdateUser updates an existing user's information
// Similar to Django's user.save() or User.objects.filter(id=id).update()
func (s *Store) UpdateUser(user types.User) error {
	query := `
		UPDATE users 
		SET first_name = ?, last_name = ?, email = ?, updated_at = NOW() 
		WHERE id = ?
	`

	res, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		utils.S.Errorf("Failed to update user %d: %v", user.ID, err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found", user.ID)
	}

	utils.S.Successf("Updated user %d", user.ID)
	return nil
}

// DeleteUser deletes a user by their ID
// Similar to Django's User.objects.filter(id=id).delete()
func (s *Store) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`

	res, err := s.db.Exec(query, id)
	if err != nil {
		utils.S.Errorf("Failed to delete user %d: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found", id)
	}

	utils.S.Successf("Deleted user %d", id)
	return nil
}
