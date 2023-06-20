package model

func GetAllGomodoros() ([]Gomodoro, error) {
	var gomodoros []Gomodoro

	tx := DB.Find(&gomodoros)
	if tx.Error != nil {
		return []Gomodoro{}, tx.Error
	}

	return gomodoros, nil
}

func GetGomodoroByName(name string) (Gomodoro, error) {
	var gomodoro Gomodoro

	tx := DB.Where("name = ?", name).First(&gomodoro)
	if tx.Error != nil {
		return Gomodoro{}, tx.Error
	}

	return gomodoro, nil
}

func GetGomodoroByID(id uint) (Gomodoro, error) {
	var gomodoro Gomodoro

	tx := DB.First(&gomodoro, id)
	if tx.Error != nil {
		return Gomodoro{}, tx.Error
	}

	return gomodoro, nil
}

func CreateGomodoro(name string) (Gomodoro, error) {
	gomodoro := &Gomodoro{
		Name: name,
	}

	tx := DB.Create(gomodoro)
	if tx.Error != nil {
		return Gomodoro{}, tx.Error
	}

	return *gomodoro, nil
}

func DeleteGomodoroByID(id uint) error {
	tx := DB.Delete(&Gomodoro{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func DeleteGomodoroByName(name string) error {
	tx := DB.Where("name = ?", name).Delete(&Gomodoro{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func UpdateGomodoro(id uint, gomodoro *Gomodoro) error {
	tx := DB.Model(&Gomodoro{}).Where("id = ?", id).Updates(gomodoro)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
