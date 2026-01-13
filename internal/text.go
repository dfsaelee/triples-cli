package internal

import (
    "regexp" 
    "strings"
)

// function to Clean the Title
func cleanTitle(title string) string {
	// Remove hashtags and the text immediately following them
	re := regexp.MustCompile(`#\S+`)
	cleaned := re.ReplaceAllString(title, "")
	// Remove extra spaces
	cleaned = strings.TrimSpace(cleaned)
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	return cleaned
}
