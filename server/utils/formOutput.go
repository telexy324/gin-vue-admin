package utils

import "strings"

func FormOutput(msg string) (formedMsg string) {
	if len(msg) <= 0 {
		return
	}
	sliceString := strings.Split(msg, "\n")
	for i, s := range sliceString {
		if i == len(sliceString)-1 {
			formedMsg = formedMsg + s
		} else {
			formedMsg = formedMsg + s + "<br>"
		}
	}
	return
}
