package util

import "log"

func PrintError(function, llmName string, chatID int64, message string, err error) {
	log.Printf("%s: llmName: %s, chatID: %d, %s: %s, ", function, llmName, chatID, message, err.Error())
}

func PrintInfo(function, llmName string, chatID int64, message string) {
	log.Printf("%s: llmName: %s, chatID: %d, %s, ", function, llmName, chatID, message)
}
