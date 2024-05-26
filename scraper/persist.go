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
		fmt.Println("JSON file couldn't be created:", err)
		return
	}
	_, err = file.WriteString("{")
	if err != nil {
		fmt.Println("JSON file couldn't be initialized:", err)
	}
	file.Close()

	visitFn(startingLink)

	// Close the JSON object
	file, err = os.OpenFile(PSGTECH_JSON_FILE_PATH, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("JSON file couldn't be opened:", err)
		return
	}
	_, err = file.WriteString("}")
	if err != nil {
		fmt.Println("JSON file couldn't be closed properly:", err)
	}
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
		return
	}

	if fileInfo.Size() > 1 { // Check if the file already contains entries
		// Move back two bytes to overwrite the closing }
		file.Seek(fileInfo.Size()-1, 0)
		file.WriteString(",")
	}

	entry := fmt.Sprintf(`"%s":`, pageDocument.Url)
	_, err = file.WriteString(entry)
	if err != nil {
		fmt.Println("Couldn't write URL to JSON:", err)
		return
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(pageDocument)
	if err != nil {
		fmt.Println("Couldn't encode data to JSON:", err)
		return
	}
}
