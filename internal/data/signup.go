package data

import (
	"context"
	"database/sql"
	"errors" // Import fmt for printing
	"strings"
	"time"

	"github.com/mtechguy/quiz3/internal/validator"
)

type Signup struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FName    string `json:"fname"`
	MName    string `json:"mname"`
	LName    string `json:"lname"`
	FullName string `json:"fullname"`
	Version  int32  `json:"version"`
}

type SignupModel struct {
	DB *sql.DB
}

func (u *Signup) GetFullName() string {
	return u.FName + " " + u.MName + " " + u.LName
}

func ValidateSignup(v *validator.Validator, signup *Signup) {
	v.Check(strings.TrimSpace(signup.Email) != "", "email", "must be provided")
	v.Check(strings.TrimSpace(signup.FName) != "", "fname", "must be provided")
	v.Check(strings.TrimSpace(signup.MName) != "", "mname", "must be provided")
	v.Check(strings.TrimSpace(signup.LName) != "", "lname", "must be provided")

	v.Check(len(signup.Email) <= 100, "email", "must not be more than 100 bytes long")
	v.Check(len(signup.FName) <= 25, "fname", "must not be more than 25 bytes long")
	v.Check(len(signup.MName) <= 25, "mname", "must not be more than 25 bytes long")
	v.Check(len(signup.LName) <= 25, "lname", "must not be more than 25 bytes long")
}

func (c SignupModel) Insert(signup *Signup) error {
	fullName := signup.GetFullName() // Use GetFullName to generate the full name

	query := `
		INSERT INTO signup (email, full_name)
		VALUES ($1, $2)
		RETURNING id, version
	`
	args := []any{signup.Email, fullName}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(
		&signup.ID,
		&signup.Version)
}

func (c SignupModel) Get(id int64) (*Signup, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, email, full_name, version
		FROM signup
		WHERE id = $1
	`
	var signup Signup
	var fullName string

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&signup.ID,
		&signup.Email,
		&fullName,
		&signup.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Split full_name into FName and LName
	names := strings.SplitN(fullName, " ", 3)
	signup.FName = names[0]
	if len(names) > 1 {
		signup.MName = names[1]
	}
	if len(names) > 2 {
		signup.LName = names[2]
	}

	// Assign FullName
	signup.FullName = signup.GetFullName()

	return &signup, nil
}

func (c SignupModel) Update(signup *Signup) error {
	fullName := signup.GetFullName() // Use GetFullName to generate the full name

	query := `
		UPDATE signup
		SET email = $1, full_name = $2, version = version + 1
		WHERE id = $3
		RETURNING version
	`

	args := []any{signup.Email, fullName, signup.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&signup.Version)
	if err != nil {
		return err
	}

	return nil
}

func (c SignupModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM signup
        WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (c SignupModel) GetAll() ([]*Signup, error) {
	query := `
		SELECT id, email, full_name, version
		FROM signup
		ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	signups := []*Signup{}

	for rows.Next() {
		var signup Signup
		var fullName string

		err := rows.Scan(
			&signup.ID,
			&signup.Email,
			&fullName,
			&signup.Version,
		)
		if err != nil {
			return nil, err
		}

		// Assign the full name directly to FullName
		signup.FullName = fullName

		names := strings.SplitN(fullName, " ", 3)
		signup.FName = names[0]
		if len(names) > 1 {
			signup.MName = names[1]
		}
		if len(names) > 2 {
			signup.LName = names[2]
		}

		signups = append(signups, &signup)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return signups, nil
}
