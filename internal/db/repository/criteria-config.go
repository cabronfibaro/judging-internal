package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type CriteriaConfig struct {
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	Categories CriteriaCategories `json:"categories"`
}

type CriteriaCategories struct {
	Common CriteriaCategory `json:"common"`
	Lead   CriteriaCategory `json:"lead"`
	Chase  CriteriaCategory `json:"chase"`
}

type CriteriaCategory struct {
	Label string        `json:"label"`
	Tags  []CriteriaTag `json:"tags"`
}

type CriteriaTag struct {
	ID           string   `json:"id"`
	Label        string   `json:"label"`
	Type         string   `json:"type"`
	ValueType    string   `json:"valueType"`
	Max          *float64 `json:"max,omitempty"`
	Multiplier   *float64 `json:"multiplier,omitempty"`
	DefaultValue *float64 `json:"defaultValue,omitempty"`
}

type CriteriaConfigRepository interface {
	Get() ([]CriteriaConfig, error)
	GetByID(int) (CriteriaConfig, error)
	Create(CriteriaConfig) error
	CopyByID(int) (CriteriaConfig, error)
	Delete(int) error
	Update(int, CriteriaConfig) error
}

type SQLiteCriteriaConfigRepository struct {
	DB *sql.DB
}

func (r *SQLiteCriteriaConfigRepository) Get() ([]CriteriaConfig, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, config
		FROM criteria_configs
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("get criteria configs: %w", err)
	}
	defer rows.Close()

	configs := make([]CriteriaConfig, 0)

	for rows.Next() {
		var config CriteriaConfig
		var raw string

		if err := rows.Scan(&config.ID, &config.Name, &raw); err != nil {
			return nil, fmt.Errorf("scan criteria config: %w", err)
		}

		if err := json.Unmarshal([]byte(raw), &config.Categories); err != nil {
			return nil, fmt.Errorf("decode criteria config %d: %w", config.ID, err)
		}

		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate criteria configs: %w", err)
	}

	return configs, nil
}

func (r *SQLiteCriteriaConfigRepository) GetByID(id int) (CriteriaConfig, error) {
	var config CriteriaConfig
	var rawConfig []byte

	err := r.DB.QueryRow(`
        SELECT id, name, config
        FROM criteria_configs
        WHERE id = ?
    `, id).Scan(&config.ID, &config.Name, &rawConfig)
	if err != nil {
		return CriteriaConfig{}, err
	}

	if err := json.Unmarshal(rawConfig, &config.Categories); err != nil {
		return CriteriaConfig{}, fmt.Errorf("unmarshal criteria config: %w", err)
	}

	return config, nil
}

func (r *SQLiteCriteriaConfigRepository) Create(config CriteriaConfig) error {
	raw, err := json.Marshal(config.Categories)
	if err != nil {
		return fmt.Errorf("encode criteria config: %w", err)
	}

	now := time.Now().Unix()

	_, err = r.DB.Exec(`
		INSERT INTO criteria_configs (name, config, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, config.Name, string(raw), now, now)
	if err != nil {
		return fmt.Errorf("create criteria config: %w", err)
	}

	return nil
}

func (r *SQLiteCriteriaConfigRepository) CopyByID(id int) (CriteriaConfig, error) {
	config, err := r.GetByID(id)
	if err != nil {
		return CriteriaConfig{}, err
	}

	config.ID = 0
	config.Name += " (Copy)"

	if err := r.Create(config); err != nil {
		return CriteriaConfig{}, err
	}

	return config, nil
}

func (r *SQLiteCriteriaConfigRepository) Delete(id int) error {
	result, err := r.DB.Exec(`
		DELETE FROM criteria_configs
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("delete criteria config: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted criteria config: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("criteria config %d not found", id)
	}

	return nil
}

func (r *SQLiteCriteriaConfigRepository) Update(id int, config CriteriaConfig) error {
	raw, err := json.Marshal(config.Categories)
	if err != nil {
		return fmt.Errorf("encode criteria config: %w", err)
	}

	result, err := r.DB.Exec(`
		UPDATE criteria_configs
		SET name = ?, config = ?, updated_at = ?
		WHERE id = ?
	`, config.Name, string(raw), time.Now().Unix(), id)
	if err != nil {
		return fmt.Errorf("update criteria config: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated criteria config: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("criteria config %d not found", id)
	}

	return nil
}
