package postgres

import (
	"time"
)

type Recipe struct {
	ID          int
	Topic       string
	Description string
	UserId      int
	PrepTime    int
	Difficulty  int
	Vegetarian  bool

	CreatedAt time.Time
}

const timeFormat = "2006-01-02 15:04:05"

func (recipe *Recipe) CreatedAtDate() string {
	return recipe.CreatedAt.Format(timeFormat)
}

func (recipe *Recipe) CreateRecipe(exec Executer) error {
	statement := "INSERT INTO recipes (topic, user_id, description, prep_time, difficulty, vegetarian) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := exec.Exec(statement, recipe.Topic, recipe.UserId, recipe.Description, recipe.PrepTime, recipe.Difficulty, recipe.Vegetarian)
	if err != nil {
		return err
	}
	return nil
}

func GetRecipeByID(exec Executer, recipeId string) (*Recipe, error) {
	var recipe *Recipe
	if err := exec.Get(&recipe, "SELECT * FROM recipes WHERE id = $1 AND is_deleted = false", recipeId); err != nil {
		return nil, err
	}
	return recipe, nil
}

func GetAllRecipes(exec Executer) ([]*Recipe, error) {
	var recipes []*Recipe
	if err := exec.Select(&recipes, "SELECT * FROM recipes WHERE is_deleted = false ORDER BY created_at DESC"); err != nil {
		return nil, err
	}
	return recipes, nil
}

//GetAuthor returns the author of the recipe
func (recipe *Recipe) GetAuthor(exec Executer) (*User, error) {
	var user *User
	if err := exec.Get(user, "SELECT * FROM users WHERE user_id = $1", recipe.UserId); err != nil {
		return nil, err
	}
	return user, nil
}

func (recipe *Recipe) UpdateRecipe(exec Executer) error {
	if _, err := exec.Exec("UPDATE recipes SET topic = $1, description = $2, prep_time = $3, difficulty = $4, vegetarian = $5 WHERE id = $6 ", recipe.Topic, recipe.Description, recipe.PrepTime, recipe.Difficulty, recipe.Vegetarian, recipe.ID); err != nil {
		return err
	}
	return nil
}

func (recipe *Recipe) DeleteRecipe(exec Executer) error {
	if _, err := exec.Exec("UPDATE recipes SET is_deleted = true WHERE id = $1", recipe.ID); err != nil {
		return err
	}
	return nil
}
