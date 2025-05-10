package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	parts := strings.Split(data, ",") //Разделение строки на слайс строк
	if len(parts) != 2 {              //проверка длины слайса
		return 0, 0, fmt.Errorf("invalid format")
	}
	firstElement := parts[0]                 //первая часть. количество шагов
	steps, err := strconv.Atoi(firstElement) // преобразование в int
	if err != nil {
		return 0, 0, fmt.Errorf("conversion error: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("the number of steps must be greater than 0")
	}
	durationOfTheWalk := parts[1]                          // вторая часть. продолжительность прогулки
	duration, err := time.ParseDuration(durationOfTheWalk) // строка в time.Duration
	if err != nil {
		return 0, 0, fmt.Errorf("conversion error: %w", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("the duration of the walk should be more than 0")
	}
	return steps, duration, nil
}

// DayActionInfo вычисляет количество шагов, дистанцию в км и потраченные калл.
func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data) //данные о шагах и продолж.прогулки
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distanceMeters := float64(steps) * stepLength //дист в метрах.шаги * длину шага
	distanceKm := distanceMeters / mInKm          //дист в км.дист в метрах/ число м в км
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	sample := "Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n"
	return fmt.Sprintf(sample, steps, distanceKm, calories)

}
