package utils

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/lib/pq"
)

func ConvertDbArrayToUnnestArrayString(dbArray pq.StringArray) string {
	byteArray, err := json.Marshal(dbArray)

	if err != nil {
		log.Fatalln(err)
	}

	stringifiedArray := string(byteArray)
	removedBrackets := strings.Trim(stringifiedArray, "[]")
	unnestArrayString := strings.Replace(removedBrackets, "\"", "'", 4)

	return unnestArrayString
}
