package models

type Category struct {
	ID   int
	Name string
}

func GetOrCreateCategory(name string) (*Category, error) {
	var category Category
	tx := DB.FirstOrCreate(&category, Category{Name: name})
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &category, nil
}
