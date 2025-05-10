package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65     // средняя длина шага.
	mInKm                      = 1000     // количество метров в километре.
	minInH                     = 60       // количество минут в часе.
	stepLengthCoefficient      = 0.45     // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5      // коэффициент для расчета калорий при ходьбе
	activity_run               = "Бег"    // Активность типа бег
	activity_walk              = "Ходьба" // Активность типа ходьба
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",") //Разделение строки на слайс строк
	if len(parts) != 3 {              //проверка длины слайса
		return 0, "", 0, fmt.Errorf("invalid format")
	}
	firstElement := parts[0]                 //1 часть. количество шагов
	steps, err := strconv.Atoi(firstElement) // преобразование в int
	if err != nil {
		return 0, "", 0, fmt.Errorf("conversion error: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("the number of steps must be greater than 0")
	}
	activity := parts[1] //2 часть. активность
	if activity == "" {
		return 0, "", 0, fmt.Errorf("invalid activity format")
	}
	durationOfTheWalk := parts[2]                          // 3 часть строки. продолжительность прогулки
	duration, err := time.ParseDuration(durationOfTheWalk) // строка в time.Duration
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid walk duration format")
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("the duration of the walk should be more than 0")
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

	stepLength := height * float64(stepLengthCoefficient) //длина 1 шага
	distanceInM := float64(steps) * stepLength
	distanceInKm := distanceInM / mInKm
	return distanceInKm

}

// meanSpeed возвращает среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration <= 0 { //продолжительность прогулки
		return 0
	}
	distance := distance(steps, height) //вытаскиываем дистанцию
	averageSpeed := distance / duration.Hours()
	if duration == 0 {
		return 0
	}
	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	if weight <= 0 {
		return "", fmt.Errorf("weight must be greater than 0")
	}
	if height <= 0 {
		return "", fmt.Errorf("height must be greater than 0")
	}

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("data error: %w", err)
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	calories := 0.0
	switch {
	case activity == activity_walk:
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case activity == activity_run:
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("unknown training type")
	}

	if err != nil {
		return "", fmt.Errorf("data error: %w", err)
	}

	sample := "Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n"
	line := fmt.Sprintf(sample, activity, duration.Hours(), distance, speed, calories)
	return line, nil
}

// RunningSpentCalories возвращает количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 {
		return 0, fmt.Errorf("the duration of the walk should be more than 0")
	}
	if steps <= 0 {
		return 0, fmt.Errorf("the number of steps must be greater than 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be greater than 0")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationInMinutes) / minInH
	return calories, nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 {
		return 0, fmt.Errorf("the duration of the walk should be more than 0")
	}
	if steps <= 0 {
		return 0, fmt.Errorf("the number of steps must be greater than 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be greater than 0")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationInMinutes) / minInH * walkingCaloriesCoefficient

	return calories, nil
}
