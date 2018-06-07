package postgres

type Rate struct {
	ID       int
	Rate     int
	UserId   int
	RecipiId int
}

func GetRatingByRecipeID(exec Executer, recipeId string) (*Rate, error) {
	var rate *Rate
	if err := exec.Get(&rate, "SELECT * FROM ratings WHERE recipe_id = $1", recipeId); err != nil {
		return nil, err
	}
	return rate, nil
}

func UpdateRate(exec Executer, recipeId string) error {
	rate, err := GetRatingByRecipeID(exec, recipeId)
	if err != nil {
		return err
	}
	rateNewNum := rate.addRate()
	if _, err := exec.Exec("UPDATE ratings SET rate = $1 WHERE recipe_id = $2", rateNewNum, recipeId); err != nil {
		return err
	}
	return nil
}

func (rate *Rate) addRate() int {
	return rate.Rate + 1
}
