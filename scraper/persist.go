package scraper

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func initJSONAndStartVisiting(startingLink string, visitFn func(string) error) {
	file, err := os.Create(PSGTECH_JSON_FILE_PATH)
	if err != nil {
		fmt.Println("JSON file couldn't be created: ", err)
		return
	}
	file.WriteString("{")
	file.Close()

	visitFn(startingLink)

	file, err = os.OpenFile(PSGTECH_JSON_FILE_PATH, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("JSON file couldn't be opened: ", err)
		return
	}
	file.WriteString("}")
	file.Close()
}

func appendToJSON(pageDocument PageDocument) {
	file, err := os.OpenFile(PSGTECH_JSON_FILE_PATH, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("JSON file couldn't be opened:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("JSON file info was not gettable:", err)
	}
	if fileInfo.Size() > 1 {
		// Remove the trailing '}' to add the new entry
		file.Seek(fileInfo.Size()-1, 0)
		file.WriteString(",")
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(map[string]PageDocument{pageDocument.Url: pageDocument})
	if err != nil {
		fmt.Println("Couldn't encode data to JSON:", err)
	}

	// Add back the closing '}'
	file.WriteString("}")
}
