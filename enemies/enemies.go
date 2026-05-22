package enemies

import (
	"encoding/json"
	"hash/fnv"
	"os"
	"time"

	"github.com/WeatherGod3218/nullscaple/logging"
	"github.com/sirupsen/logrus"
)

type Comparison string

const (
	Higher Comparison = "higher"
	Equal  Comparison = "equal"
	Lower  Comparison = "lower"
)

// Internal Enemy Data
type EnemyData struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CurseAmount int    `json:"curse_amount"`
	StartLevel  int    `json:"start_level"`
	KillMethod  string `json:"kill_method"`
	Filename    string `json:"filename"`
}

// Internal Enemy Data
type EnemyComparison struct {
	Id          bool       `json:"id"`
	Name        bool       `json:"name"`
	KillMethod  bool       `json:"kill_method"`
	CurseAmount Comparison `json:"curse_amount"`
	StartLevel  Comparison `json:"start_level"`
}

type EnemyRequest struct {
	ID   string `json:"enemy_id"`
	Mode string `json:"gameplay_mode"`
}

var LoadedEnemies map[int]*EnemyData

func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func compareInts(selected int, base int) Comparison {
	if selected > base {
		return Higher
	} else if selected < base {
		return Lower
	}

	return Equal
}

// Initializes enemy data from the enemies.json file in the root directory. Establishes loadedServerEnemies and loadedClientEnemies as their respective types
func InitEnemies() {
	logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "InitEnemies"}).Info("Starting enemy loading!")

	rawData, err := os.ReadFile("enemies.json")
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "enemies", "method": "InitEnemies"}).Fatal("failed to load the enemies.json!")
	}

	var enemies []EnemyData

	err = json.Unmarshal(rawData, &enemies)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "enemies", "method": "InitEnemies"}).Fatal("failed to unmarshal the json!")
	}

	LoadedEnemies = make(map[int]*EnemyData, len(enemies))

	for _, enemy := range enemies {
		LoadedEnemies[enemy.Id] = &enemy
	}

	logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "InitEnemies"}).Info("Succesfully loaded all enemies!")
}

func GetEnemyOfTheDay(mode string) *EnemyData {
	if len(LoadedEnemies) == 0 {
		logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "GetEnemyOfTheDay"}).Warn("enemy list is empty, returning fake enemy!")
		return nil
	}

	currentDay := time.Now().Format("2006-1-2")
	inputString := mode + os.Getenv("ENEMY_HASH_SALT") + currentDay
	dayHash := hashString(inputString)

	index := int(dayHash) % len(LoadedEnemies)
	return LoadedEnemies[index]
}

func GetEnemyFromId(id int) *EnemyData {
	if LoadedEnemies == nil {
		logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "GetEnemyFromId"}).Warn("enemy list is empty, returning fake enemy!")
		return nil
	}

	value, ok := LoadedEnemies[id]
	if ok {
		return value
	}

	logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "GetEnemyFromId"}).Warn("attempted to fetch an empty Id")
	return nil
}

func CompareEnemies(selectedEnemy *EnemyData, baseEnemy *EnemyData) EnemyComparison {
	return EnemyComparison{
		Id:          selectedEnemy.Id == baseEnemy.Id,
		Name:        selectedEnemy.Name == baseEnemy.Name,
		KillMethod:  selectedEnemy.KillMethod == baseEnemy.KillMethod,
		CurseAmount: compareInts(selectedEnemy.CurseAmount, baseEnemy.CurseAmount),
		StartLevel:  compareInts(selectedEnemy.StartLevel, baseEnemy.StartLevel),
	}
}

func GetEnemyList() map[int]*EnemyData {
	if LoadedEnemies == nil {
		logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "GetEnemyList"}).Warn("enemy list is empty, returning fake enemy!")
		return nil
	}
	return LoadedEnemies
}

func CheckIfStringIsMode(mode string) bool {
	switch mode {
	case "casual", "standard", "extreme":
		return true
	}
	return false
}
