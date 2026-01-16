package infrastructure

import "github.com/jackc/pgx/v5/pgxpool"

type MacroDataProvider struct {
	pool *pgxpool.Pool
}

func NewMacroDataProvider(pool *pgxpool.Pool) *MacroDataProvider {
	return &MacroDataProvider{pool: pool}
}

func GetRubleDollarForecast() {}

func GetCurrentCentralBankRate() {}
