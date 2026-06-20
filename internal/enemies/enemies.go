package enemies

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"math"
	"os"

	"github.com/WeatherGod3218/nullscaple/internal/logging"
	t "github.com/WeatherGod3218/nullscaple/internal/nulltypes"
	"github.com/WeatherGod3218/nullscaple/internal/timeutil"
	"github.com/sirupsen/logrus"
)

var LoadedEnemies map[string]*t.EnemyData
var EnemyList []string

const CURSES_PARTIAL_DISTANCE int = 2
const SPAWN_PARTIAL_DISTANCE int = 5

func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func compareInts(selected int, base int, partial_yield int) t.ComparisonResults {
	if selected > base {
		if math.Abs(float64(selected-base)) <= float64(partial_yield) {
			return t.Comparison_High_Middle
		}
		return t.Comparison_Higher
	} else if selected < base {
		if math.Abs(float64(selected-base)) <= float64(partial_yield) {
			return t.Comparison_Low_Middle
		}
		return t.Comparison_Lower
	}

	return t.Comparison_Equal
}

// Initializes enemy data from the enemies.json file in the root directory. Establishes loadedServerEnemies and loadedClientEnemies as their respective types
func InitEnemies() {
	logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "InitEnemies"}).Info("Starting enemy loading!")

	rawData, err := os.ReadFile("enemies.json")
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "enemies", "method": "InitEnemies"}).Fatal("failed to load the enemies.json!")
	}

	var enemies []t.EnemyData

	//err = json.Unmarshal(rawData, &enemies)
	err = json.Unmarshal(rawData, &enemies)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "enemies", "method": "InitEnemies"}).Fatal("failed to unmarshal the json!")
	}

	LoadedEnemies = make(map[string]*t.EnemyData, len(enemies))
	EnemyList = make([]string, 0, len(enemies))

	for _, enemy := range enemies {
		e := enemy
		LoadedEnemies[e.Id] = &e
		EnemyList = append(EnemyList, e.Id)

	}

	logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "InitEnemies"}).Info("Succesfully loaded all enemies!")
}

func GetEnemyOfTheDay(mode string) (*t.EnemyData, error) {
	if len(LoadedEnemies) == 0 {
		logging.Logger.WithFields(logrus.Fields{"module": "enemies", "method": "GetEnemyOfTheDay"}).Warn("enemy list is empty!")
		return nil, errors.New("Empty Enemy List!")
	}

	currentDay := timeutil.GetFormattedTime()
	inputString := mode + os.Getenv("ENEMY_HASH_SALT") + currentDay
	dayHash := hashString(inputString)

	listIndex := int(dayHash) % len(LoadedEnemies)
	enemyId := EnemyList[listIndex]

	foundEnemy := LoadedEnemies[enemyId]
	return foundEnemy, nil
}

func GetEnemyFromId(id string) *t.EnemyData {
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

func CompareEnemies(selectedEnemy *t.EnemyData, baseEnemy *t.EnemyData) (bool, t.EnemyComparison) {
	comparison := t.EnemyComparison{
		Id:          selectedEnemy.Id == baseEnemy.Id,
		Name:        selectedEnemy.Name == baseEnemy.Name,
		KillMethod:  selectedEnemy.KillMethod == baseEnemy.KillMethod,
		CurseAmount: compareInts(selectedEnemy.CurseAmount, baseEnemy.CurseAmount, CURSES_PARTIAL_DISTANCE),
		StartLevel:  compareInts(selectedEnemy.StartLevel, baseEnemy.StartLevel, SPAWN_PARTIAL_DISTANCE),
	}
	isEqual := comparison.Id && comparison.Name && comparison.KillMethod && comparison.CurseAmount == t.Comparison_Equal && comparison.StartLevel == t.Comparison_Equal
	return isEqual, comparison
}

func GetEnemyList() map[string]*t.EnemyData {
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
