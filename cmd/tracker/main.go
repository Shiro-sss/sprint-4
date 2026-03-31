package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	// Проверка что результат: шаги, тип, длительность
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid data format")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient  // Длина шага
	distanceMeters := float64(steps) * stepLen // Дистанция в метрах
	return distanceMeters / mInKm              // Перевожу метры в километры
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height) // Дистанция в километрах
	hours := duration.Hours()
	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activity, duration, err := parseTraining(data)
	if err != nil {

		log.Println(err)
		return "", err
	}

	distKm := distance(steps, height)           // Считаю дистанцию
	speed := meanSpeed(steps, height, duration) // Скорость

	var calories float64 // Считаю калории

	switch activity {
	case "Бег":

		calories, err = RunningSpentCalories(steps, weight, height, duration) // Калории при беге
		if err != nil {
			log.Println(err)
			return "", err
		}

	case "Ходьба":

		calories, err = WalkingSpentCalories(steps, weight, height, duration) // Калории при ходьбе
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки") // Неизвестный тип тренировки
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity,
		duration.Hours(),
		distKm,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 { // Проверка параметров
		return 0, errors.New("invalid input parameters")
	}

	speed := meanSpeed(steps, height, duration) // Средняя скорость км в час
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid input parameters")
	}

	speed := meanSpeed(steps, height, duration) // Средняя скорость
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	calories *= walkingCaloriesCoefficient
	return calories, nil
}

func main() {
	weight := 84.6
	height := 1.87

	// дневная активность
	input := []string{
		"678,0h50m",
		"792,1h14m",
		"1078,1h30m",
		"7830,2h40m",
		",3456",
		"12:40:00, 3456",
		"something is wrong",
	}

	fmt.Println("Активность в течение дня")

	var (
		dayActionsInfo string
		dayActionsLog  []string
	)

	for _, v := range input {
		dayActionsInfo = daysteps.DayActionInfo(v, weight, height)
		dayActionsLog = append(dayActionsLog, dayActionsInfo)
	}

	for _, v := range dayActionsLog {
		fmt.Println(v)
	}

	// тренировки
	trainings := []string{
		"3456,Ходьба,3h00m",
		"something is wrong",
		"678,Бег,0h5m",
		"1078,Бег,0h10m",
		",3456 Ходьба",
		"7892,Ходьба,3h10m",
		"15392,Бег,0h45m",
	}

	var trainingLog []string

	for _, v := range trainings {
		trainingInfo, err := spentcalories.TrainingInfo(v, weight, height)
		if err != nil {
			log.Printf("не получилось получить информацию о тренировке: %v", err)
			continue
		}
		trainingLog = append(trainingLog, trainingInfo)
	}

	fmt.Println("Журнал тренировок")

	for _, v := range trainingLog {
		fmt.Println(v)
	}
}
