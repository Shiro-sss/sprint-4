package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов
const (
 mInKm                      = 1000   // количество метров в километре
 minInH                     = 60     // количество минут в часе
 stepLengthCoefficient      = 0.45   // коэффициент для расчета длины шага на основе роста
 walkingCaloriesCoefficient = 0.5    // коэффициент для расчета калорий при ходьбе
)

// parseTraining разбирает строку формата "шаги,активность,длительность"
func parseTraining(data string) (int, string, time.Duration, error) {
 parts := strings.Split(data, ",")
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

// distance рассчитывает дистанцию в километрах
func distance(steps int, height float64) float64 {
 stepLen := height * stepLengthCoefficient
 distanceMeters := float64(steps) * stepLen
 return distanceMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
 if duration <= 0 {
  return 0
 }
 dist := distance(steps, height)
 hours := duration.Hours()
 return dist / hours
}

// TrainingInfo формирует строку с информацией о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
 steps, activity, duration, err := parseTraining(data)
 if err != nil {
  log.Println(err)
  return "", err
 }

 distKm := distance(steps, height)
 speed := meanSpeed(steps, height, duration)
 var calories float64

 switch activity {
 case "Бег":
  calories, err = RunningSpentCalories(steps, weight, height, duration)
  if err != nil {
   log.Println(err)
   return "", err
  }
 case "Ходьба":
  calories, err = WalkingSpentCalories(steps, weight, height, duration)
  if err != nil {
   log.Println(err)
   return "", err
  }
 default:
  return "", errors.New("неизвестный тип тренировки")
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

// RunningSpentCalories считает калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
 if steps <= 0  weight <= 0  height <= 0 || duration <= 0 {
  return 0, errors.New("invalid input parameters")
 }

 speed := meanSpeed(steps, height, duration)
 durationInMinutes := duration.Minutes()
 calories := (weight * speed * durationInMinutes) / minInH
 return calories, nil
}

// WalkingSpentCalories считает калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
 if steps <= 0  weight <= 0  height <= 0 || duration <= 0 {
  return 0, errors.New("invalid input parameters")
 }

 speed := meanSpeed(steps, height, duration)
 durationInMinutes := duration.Minutes()
 calories := (weight * speed * durationInMinutes) / minInH
 calories *= walkingCaloriesCoefficient
 return calories, nil
}