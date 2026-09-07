package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tgtest/models"
	//"tgtest/currency"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

// func (r *Repo) CurrencyExchange(ctx context.Context, userID int64, currancyName string, amount string) error {
// 	valyuta, err := currency.ParseCurrencyName(currancyName)
// 	if err != nil {
// 		return fmt.Errorf("ошибка в парсинге имени валюты: %w", err)
// 	}
// 	summa, err := currency.ParseNumber(amount)
// 	if err != nil {
// 		return fmt.Errorf("ошибка в парсинге суммы для перевода: %w", err)
// 	}
// 	_, err := r.db.Exec(ctx, "update users set city = $1 where id = $2", city, userID)
// 	if err != nil {
// 		return fmt.Errorf("error db.Exex: %w", err)
// 	}

// 	return err
// }

func (r *Repo) IsUserExistsByUsername(ctx context.Context, userName string) (bool, error) {
	var existOrNot bool
	err := r.db.QueryRow(ctx, "select exists(select 1 from users where username = $1)",
		userName).Scan(&existOrNot)
	if err != nil {
		return false, fmt.Errorf("error checking user: %w", err)
	}
	return existOrNot, nil
}

func (r *Repo) IsUserExistsByID(ctx context.Context, userID int64) (bool, error) {
	var existOrNot bool
	err := r.db.QueryRow(ctx, "select exists(select 1 from users where id = $1)",
		userID).Scan(&existOrNot)
	if err != nil {
		return false, fmt.Errorf("error checking user: %w", err)
	}
	return existOrNot, nil
}

func (r *Repo) GetUserCity(ctx context.Context, userID int64) (string, error) {
	var city string

	row := r.db.QueryRow(ctx, "select coalesce(city, '') from users where id=$1", userID)
	err := row.Scan(&city)
	if err != nil {
		return "", fmt.Errorf("error row.Scan %w", err)
	}

	return city, nil
}

func (r *Repo) CreateUser(ctx context.Context, userID int64, username string) error {
	_, err := r.db.Exec(ctx, "insert into users(id, username) values ($1, $2)", userID, username)
	if err != nil {
		return fmt.Errorf("error db.Exex: %w", err)
	}

	return nil
}

func (r *Repo) UpdateCity(ctx context.Context, userID int64, city string) error {
	_, err := r.db.Exec(ctx, "update users set city = $1 where id = $2", city, userID)
	if err != nil {
		return fmt.Errorf("error db.Exex: %w", err)
	}

	return err
}

func (r *Repo) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	user := models.User{}

	row := r.db.QueryRow(ctx, "select id, coalesce(city, ''), created_at from users where id = $1", userID)
	err := row.Scan(&user.ID, &user.City, &user.CreatedAt)
	if err != nil {
		fmt.Printf("SCAN ERROR: %#v\n", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("error db.Scan: %w", err)
	}

	return &user, nil
}
