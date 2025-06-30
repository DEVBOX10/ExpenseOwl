package storage

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Storage interface for all storage types
type Storage interface {
	Close() error
	GetConfig() (*Config, error)

	GetCategories() ([]string, error)
	UpdateCategories(categories []string) error
	GetTags() ([]string, error)
	UpdateTags(tags []string) error

	GetCurrency() (string, error)
	UpdateCurrency(currency string) error
	GetStartDate() (int, error)
	UpdateStartDate(startDate int) error

	GetRecurringExpenses() ([]RecurringExpense, error)
	GetRecurringExpense(id string) (*RecurringExpense, error)
	AddRecurringExpense(recurringExpense *RecurringExpense) error
	RemoveRecurringExpense(id string, removeAll bool) error
	UpdateRecurringExpense(id string, recurringExpense *RecurringExpense, updateAll bool) error

	GetAllExpenses() ([]*Expense, error)
	GetExpense(id string) (*Expense, error)
	AddExpense(expense *Expense) error
	RemoveExpense(id string) error
	UpdateExpense(id string, expense *Expense) error
}

// config for expense data
type Config struct {
	Categories        []string           `json:"categories"`
	Currency          string             `json:"currency"`
	StartDate         int                `json:"startDate"`
	Tags              []string           `json:"tags"`
	RecurringExpenses []RecurringExpense `json:"recurringExpenses"`
}

type RecurringExpense struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Amount      float64   `json:"amount"`
	StartDate   time.Time `json:"startDate"`   // date of the first occurrence
	Interval    string    `json:"interval"`    // daily, weekly, monthly, yearly
	Occurrences int       `json:"occurrences"` // 0 for 10 years (heuristic), 10 for 10 occurrences
}

type BackendType string

const (
	BackendTypeJSON     BackendType = "json"
	BackendTypePostgres BackendType = "postgres"
)

type ExpenseType string

const (
	ExpenseTypeRecurring    ExpenseType = "recurring"
	ExpenseTypeInstantiated ExpenseType = "instantiated"
)

// config for the storage backend
type SystemConfig struct {
	StorageURL  string
	StorageType BackendType
	StorageUser string
	StoragePass string
}

// expense struct
type Expense struct {
	ID          string      `json:"id"`
	ExpenseType ExpenseType `json:"expenseType"`
	Name        string      `json:"name"`
	Tags        []string    `json:"tags"`
	Category    string      `json:"category"`
	Amount      float64     `json:"amount"`
	Date        time.Time   `json:"date"`
}

func (e *Expense) Validate() error {
	if e.Name == "" {
		return errors.New("expense name is required")
	}
	if e.Category == "" {
		return errors.New("category is required")
	}
	if e.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

func (c *Config) SetBaseConfig() {
	c.Categories = defaultCategories
	c.Currency = "usd"
	c.StartDate = 1
	c.Tags = []string{}
	c.RecurringExpenses = []RecurringExpense{}
}

func (c *SystemConfig) SetStorageConfig() {
	c.StorageType = backendTypeFromEnv(os.Getenv("STORAGE_TYPE"))
	c.StorageURL = backendURLFromEnv(os.Getenv("STORAGE_URL"))
	c.StorageUser = os.Getenv("STORAGE_USER")
	c.StoragePass = os.Getenv("STORAGE_PASS")
}

func backendTypeFromEnv(env string) BackendType {
	switch env {
	case "json":
		return BackendTypeJSON
	case "postgres":
		return BackendTypePostgres
	}
	return BackendTypeJSON
}

func backendURLFromEnv(env string) string {
	switch env {
	case "json":
		return "data"
	case "postgres":
		return "" // not sure about default url
	}
	return "data"
}

// initializes the storage backend
func InitializeStorage() (Storage, error) {
	baseConfig := SystemConfig{}
	baseConfig.SetStorageConfig()
	switch baseConfig.StorageType {
	case BackendTypeJSON:
		return InitializeJsonStore(baseConfig)
	case BackendTypePostgres:
		return InitializePostgresStore(baseConfig)
	}
	return nil, fmt.Errorf("invalid data store: %s", baseConfig.StorageType)
}

// variables
var defaultCategories = []string{
	"Food",
	"Groceries",
	"Travel",
	"Rent",
	"Utilities",
	"Entertainment",
	"Healthcare",
	"Shopping",
	"Miscellaneous",
	"Income",
}

var supportedCurrencies = []string{
	"usd", // US Dollar
	"eur", // Euro
	"gbp", // British Pound
	"jpy", // Japanese Yen
	"cny", // Chinese Yuan
	"krw", // Korean Won
	"inr", // Indian Rupee
	"rub", // Russian Ruble
	"brl", // Brazilian Real
	"zar", // South African Rand
	"aed", // UAE Dirham
	"aud", // Australian Dollar
	"cad", // Canadian Dollar
	"chf", // Swiss Franc
	"hkd", // Hong Kong Dollar
	"sgd", // Singapore Dollar
	"thb", // Thai Baht
	"try", // Turkish Lira
	"mxn", // Mexican Peso
	"php", // Philippine Peso
	"pln", // Polish Złoty
	"sek", // Swedish Krona
	"nzd", // New Zealand Dollar
	"dkk", // Danish Krone
	"idr", // Indonesian Rupiah
	"ils", // Israeli New Shekel
	"vnd", // Vietnamese Dong
	"myr", // Malaysian Ringgit
}
