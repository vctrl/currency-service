package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vctrl/currency-service/currency/internal/dto"
	"github.com/vctrl/currency-service/currency/internal/pkg/config"

	_ "github.com/lib/pq"
)

type ExchangeRateRepository struct {
	DB *sql.DB
}

type ExchangeRate struct {
	ID            int
	Date          time.Time
	BaseCurrency  string
	CurrencyRates map[string]float32
	CreatedAt     time.Time
}

type CurrencyRate struct {
	Date time.Time
	Rate float32
}

// TODO не должно быть частью репозитория!
func NewDatabaseConnection(cfg config.DatabaseConfig) (*sql.DB, string, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, "", fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, dsn, nil
}

func NewExchangeRate(db *sql.DB) (ExchangeRateRepository, error) {
	return ExchangeRateRepository{
		DB: db,
	}, nil
}

func (repo *ExchangeRateRepository) Save(ctx context.Context, date time.Time, baseCurrency string, rates map[string]float64) error {
	ratesJSON, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("failed to marshal currency rates: %w", err)
	}

	_, err = repo.DB.ExecContext(
		ctx,
		`INSERT INTO exchange_rates (date, base_currency, currency_rates) VALUES ($1, $2, $3)`,
		date, baseCurrency, ratesJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to save exchange rates: %w", err)
	}
	return nil
}

func (repo *ExchangeRateRepository) FindInInterval(ctx context.Context, dto *dto.CurrencyRequestDTO) ([]CurrencyRate, error) {
	query := `
		SELECT date, (currency_rates ->> $1)::float 
		FROM exchange_rates
		WHERE date::date BETWEEN $2 AND $3 AND base_currency = $4
	`

	rows, err := repo.DB.QueryContext(
		ctx,
		query,
		dto.TargetCurrency,
		dto.DateFrom.Format("2006-01-02"),
		dto.DateTo.Format("2006-01-02"),
		dto.BaseCurrency,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query exchange rates: %w", err)
	}
	defer rows.Close()

	var rates []CurrencyRate
	for rows.Next() {
		var rate CurrencyRate
		if err := rows.Scan(&rate.Date, &rate.Rate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rates = append(rates, rate)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return rates, nil
}
