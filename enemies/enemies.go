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

type EnemyData struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CurseAmount int    `json:"curse_amount"`
	StartLevel  int    `json:"start_level"`
}

type EnemyComparison struct {
	Id          bool       `json:"id"`
	Name        bool       `json:"name"`
	CurseAmount Comparison `json:"curse_amount"`
	StartLevel  Comparison `json:"start_level"`
}

var LoadedEnemies map[int]*EnemyData
var EnemyList map[int]string

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

	LoadedEnemies = make(map[int]*EnemyData)
	EnemyList = make(map[int]string)

	for _, enemy := range enemies {
		LoadedEnemies[enemy.Id] = &enemy
		EnemyList[enemy.Id] = enemy.Name
	}

	logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "InitEnemies"}).Info("Succesfully loaded all enemies!")
}

func GetEnemyOfTheDay() *EnemyData {
	if len(EnemyList) == 0 {
		logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "GetEnemyOfTheDay"}).Warn("enemy list is empty, returning fake enemy!")
		return nil
	}
	dayHash := hashString(time.Now().Format("2006-1-2"))

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

func CompareEnemies(selectedEnemy EnemyData, baseEnemy EnemyData) EnemyComparison {
	return EnemyComparison{
		Id:          selectedEnemy.Id == baseEnemy.Id,
		Name:        selectedEnemy.Name == baseEnemy.Name,
		CurseAmount: compareInts(selectedEnemy.CurseAmount, baseEnemy.CurseAmount),
		StartLevel:  compareInts(selectedEnemy.StartLevel, baseEnemy.StartLevel),
	}
}

func GetEnemyList() map[int]string {
	if EnemyList == nil {
		logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "GetEnemyList"}).Warn("enemy list is empty, returning fake enemy!")
		return make(map[int]string)
	}
	return EnemyList
}
