package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
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
	// TODO: реализовать функцию
	// Разделяю строку запятой
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

	// Преобразую длительность (parts[2]) в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	// Длина шага
	stepLen := height * stepLengthCoefficient

	// Дистанция в метрах
	distanceMeters := float64(steps) * stepLen

	// Перевожу метры в километры: делю на количество метров в километре
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0
	}

	// Дистанция в километрах
	dist := distance(steps, height)

	// duration.Hours() возвращает значение в часах (float64)
	hours := duration.Hours()

	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	steps, activity, duration, err := parseTraining(data)
	if err != nil {

		log.Println(err)
		return "", err
	}

	// Считаю дистанцию и среднюю скорость
	distKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// Считаю калории
	var calories float64

	switch activity {
	case "Бег":
		// Калории при беге
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		// Калории при ходьбе
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		// Неизвестный тип тренировки
		return "", errors.New("неизвестный тип тренировки")
	}

	// Строка результата
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
	// TODO: реализовать функцию
	// Проверка входных параметров
	// Должны быть только положительные значения
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid input parameters")
	}

	// Средняя скорость (км/ч) с помощью meanSpeed()
	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid input parameters")
	}

	// Средняя скорость с помощью meanSpeed()
	speed := meanSpeed(steps, height, duration)

	// Перевод продолжительности в минуты
	durationInMinutes := duration.Minutes()

	// Колории
	calories := (weight * speed * durationInMinutes) / minInH

	calories *= walkingCaloriesCoefficient

	return calories, nil
}
