package data

import "time"

type Recipe struct {
	Id         int
	Uuid       string
	Topic      string
	UserId     int
	PrepTime   int
	Difficulty int
	Vegetarian bool
	CreatedAt  time.Time
}

// type Comment struct {
// 	Id        int
// 	Uuid      string
// 	Body      string
// 	UserId    int
// 	ThreadId  int
// 	CreatedAt time.Time
// }

const timeFormat = "2006-01-02 15:04:05"

func (recipe *Recipe) CreatedAtDate() string {
	return recipe.CreatedAt.Format(timeFormat)
}

func (user *User) CreateRecipe(topic string, preptime int, difficulty int, vegetarian bool) (recipe Recipe, err error) {
	statement := "insert into recipes (uuid, topic, user_id, prep_time, difficulty, vegetarian, created_at) values ($1, $2, $3, $4, $5,$6,$7) returning id,uuid, topic, user_id, prep_time, difficulty, vegetarian, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), topic, user.Id, preptime, difficulty, vegetarian, time.Now()).Scan(&recipe.Id, &recipe.Uuid, &recipe.Topic, &recipe.UserId, &recipe.PrepTime, &recipe.Difficulty, &recipe.Vegetarian, &recipe.CreatedAt)
	return
}

func GetAllRecipes() (recipes []Recipe, err error) {
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, prep_time, difficulty, vegetarian, created_at FROM recipes ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		recipe := Recipe{}
		if err = rows.Scan(&recipe.Id, &recipe.Uuid, &recipe.Topic, &recipe.UserId, &recipe.PrepTime, &recipe.Difficulty, &recipe.Vegetarian, &recipe.CreatedAt); err != nil {
			return
		}
		recipes = append(recipes, recipe)
	}
	rows.Close()
	return
}
