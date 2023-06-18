package model

import "time"

func CreateDefaultTimer(gomodoroID uint) error {
	return CreateTimer(gomodoroID, WorkTimer, WorkDuration, 1)
}

func CreateTimer(gomodoroID uint, timerType TimerType, duration time.Duration, repetition int8) error {
	timer := &Timer{
		GomodoroID: gomodoroID,
		Type:       timerType,
		Duration:   duration,
		Repetition: repetition,
	}

	tx := DB.Create(timer)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func GetTimer(gomodoroID uint) (Timer, error) {
	var timer Timer

	tx := DB.Where("gomodoro_id = ?", gomodoroID).First(&timer)
	if tx.Error != nil {
		return Timer{}, tx.Error
	}

	return timer, nil
}

func UpdateTimer(id uint, timer *Timer) error {
	tx := DB.Model(&Timer{}).Where("id = ?", id).Updates(timer)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
