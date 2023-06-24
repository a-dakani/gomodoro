package model

func GetAllGomodoros() ([]Gomodoro, error) {
	var gomodoros []Gomodoro

	tx := db.Find(&gomodoros)
	if tx.Error != nil {
		return []Gomodoro{}, tx.Error
	}

	return gomodoros, nil
}

func GetGomodoroByName(name string) (Gomodoro, error) {
	var gomodoro Gomodoro

	tx := db.Where("name = ?", name).First(&gomodoro)
	if tx.Error != nil {
		return Gomodoro{}, tx.Error
	}

	return gomodoro, nil
}

func GetGomodoroByID(id uint) (Gomodoro, error) {
	var gomodoro Gomodoro

	tx := db.First(&gomodoro, id)
	if tx.Error != nil {
		return Gomodoro{}, tx.Error
	}

	return gomodoro, nil
}

func CreateGomodoro(name string) (Gomodoro, error) {
	gomodoro := &Gomodoro{
		Name: name,
	}

	tx := db.Create(gomodoro)
	if tx.Error != nil {
		return Gomodoro{}, tx.Error
	}

	return *gomodoro, nil
}

func DeleteGomodoroByID(id uint) error {
	tx := db.Delete(&Gomodoro{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func DeleteGomodoroByName(name string) error {
	tx := db.Where("name = ?", name).Delete(&Gomodoro{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func UpdateGomodoro(id uint, gomodoro *Gomodoro) error {
	tx := db.Model(&Gomodoro{}).Where("id = ?", id).Updates(gomodoro)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
