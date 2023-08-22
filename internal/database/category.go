package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func CreateCategoryTable(db *sql.DB) error {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS categories (
		id TEXT PRIMARY KEY,
		name TEXT,
		description TEXT
	);`)
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (c *Category) Create(name, description string) (*Category, error) {
	id := uuid.New().String()

	_, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)", id, name, description)
	if err != nil {
		return nil, err
	}

	return &Category{
		db:          c.db,
		ID:          id,
		Name:        name,
		Description: description,
	}, nil
}

func (c *Category) FindAll() ([]Category, error) {
	rows, err := c.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (c *Category) FindByCourseId(courseID string) (*Category, error) {
	var category Category
	err := c.db.QueryRow("SELECT categories.id, categories.name, categories.description FROM categories INNER JOIN courses ON categories.id = courses.category_id WHERE courses.id = $1", courseID).
		Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *Category) Find(id string) (*Category, error) {
	var category Category
	err := c.db.QueryRow("SELECT * FROM categories WHERE id = $1", id).
		Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
