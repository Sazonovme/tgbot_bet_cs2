package handler

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func validateMatches(input string) string {
	//parts := strings.Split(input, "#")

	// Регулярка: Team1vsTeam2_2025-08-16 15:00
	re := regexp.MustCompile(`^([A-Za-zА-Яа-я0-9]+)vs([A-Za-zА-Яа-я0-9]+)_(\d{4}-\d{2}-\d{2} \d{2}:\d{2})$`)

	// Убираем завершающий "#", если он есть
	input = strings.TrimSuffix(input, "#")

	// Разбиваем строку на блоки
	blocks := strings.Split(input, "#")
	if len(blocks) == 0 {
		return "строка пуста"
	}

	for i, part := range blocks {
		matches := re.FindStringSubmatch(part)
		if matches == nil {
			return fmt.Sprintf("ошибка в элементе %d: '%s' не соответствует формату", i+1, part)
		}

		// matches[3] = dateStr
		_, err := time.Parse("2006-01-02 15:04", matches[3])
		if err != nil {
			return fmt.Sprintf("ошибка в дате элемента %d: '%s'", i+1, matches[3])
		}
	}

	return ""
}

func validateMatchesResults(input string) string {
	// Регулярка для проверки одного блока
	blockPattern := regexp.MustCompile(`^\[(\d+)\]_(1|2|0-2|2-0|1-2|2-1)$`)

	// Убираем завершающий "#", если он есть
	input = strings.TrimSuffix(input, "#")

	// Разбиваем строку на блоки
	blocks := strings.Split(input, "#")
	if len(blocks) == 0 {
		return "строка пуста"
	}

	for i, block := range blocks {
		if !blockPattern.MatchString(block) {
			return fmt.Sprintf("ошибка в блоке %d: '%s' не соответствует формату", i+1, block)
		}
	}

	return ""
}
