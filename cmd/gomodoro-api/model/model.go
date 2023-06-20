package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

type TimerType string
type TimerStatus string

const (
	WorkDuration       = 25 * time.Minute
	ShortBreakDuration = 5 * time.Minute
	LongBreakDuration  = 15 * time.Minute
)

const (
	WorkTimer       TimerType = "work"
	ShortBreakTimer TimerType = "shortBreak"
	LongBreakTimer  TimerType = "longBreak"
)

const (
	Idle    TimerStatus = "idle"
	Running TimerStatus = "running"
	Paused  TimerStatus = "paused"
)

// TODO change this to get timer with gomodoro instead of gomodoro with timer in the response
// TODO Change the relation

type Gomodoro struct {
	ID          uint          `gorm:"primarykey"`
	Name        string        `json:"name" gorm:"not null;unique"`
	Work        time.Duration `json:"work" gorm:"not null; default:1500000000000"`
	ShortBreak  time.Duration `json:"shortBreak" gorm:"not null; default:300000000000"`
	LongBreak   time.Duration `json:"longBreak" gorm:"not null; default:900000000000"`
	Repetitions int8          `json:"repetitions" gorm:"not null; default:4"`
	AutoStart   bool          `json:"autoStart" gorm:"not null; default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Timer struct {
	ID         uint          `gorm:"primarykey"`
	GomodoroID uint          `json:"gomodoroID" gorm:"not null; unique"`
	Gomodoro   Gomodoro      `gorm:"foreignKey:GomodoroID;constraint:OnDelete:CASCADE"`
	Type       TimerType     `json:"type" gorm:"not null"`
	Status     TimerStatus   `json:"status" gorm:"not null; default:'idle'"`
	Duration   time.Duration `json:"duration" gorm:"not null"`
	Remaining  time.Duration `json:"remaining" gorm:"not null"`
	Repetition int8          `json:"repetition" gorm:"not null; default:1"`
	StartedAt  time.Time     `json:"startedAt" gorm:"default:null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func ConnectDB(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		return err
	}

	log.Println("Connected to database")

	db.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})

	log.Println("Running migrations")

	err = db.AutoMigrate(&Gomodoro{}, &Timer{})
	if err != nil {
		return err
	}

	DB = db

	return nil
}
