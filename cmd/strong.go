package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const instructions = `You are a helpful assistant that converts Strong's Greek definitions into JSON format.
The JSON should be structured to work with the WordEntry Go struct:

type WordEntry struct {
	Word   string 'json:"word"'
	Strong string 'json:"strong"'
	Senses []struct {
		Number     string   'json:"number"'
		Definition string   'json:"definition"'
		Citations  []string 'json:"citations"'
	} 'json:"senses"'
}

Only return valid JSON, do not include any other text or explanations.
Here is the content:
`

func StrongsToJSON() {
	outputDir := "/Users/zachm/Documents/github/ugnt-greek-text-parser/strong_output"
	sourceDir := "/Users/zachm/Documents/github/ugnt-greek-text-parser/strong"

	fileContents, err := readAllMarkdownFiles(sourceDir)
	if err != nil {
		fmt.Printf("Error reading markdown files: %v\n", err)
		return
	}

	wordEntries := make([]WordEntry, 0, len(fileContents))

	for _, content := range fileContents {
		entry, err := ParseMarkdownToWordEntry(content)
		if err != nil {
			fmt.Printf("Error parsing markdown content: %v\n", err)
			continue
		}

		wordEntries = append(wordEntries, entry)
	}

	outputFile := filepath.Join(outputDir, "strong_output.json")
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Set indentation for pretty printing
	err = encoder.Encode(wordEntries)
	if err != nil {
		fmt.Printf("Error encoding JSON to file: %v\n", err)
		return
	}

	fmt.Printf("Successfully written %d entries to %s\n", len(wordEntries), outputFile)

}

func readAllMarkdownFiles(sourceDir string) ([]string, error) {
	var contents []string
	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".md" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			contents = append(contents, string(data))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// ParseMarkdownToWordEntry parses the content of a 01.md file into a WordEntry struct.
func ParseMarkdownToWordEntry(md string) (WordEntry, error) {
	var entry WordEntry

	// Extract Strong number
	strongRe := regexp.MustCompile(`(?m)\Strongs: ([A-Z0-9]+)`)
	if match := strongRe.FindStringSubmatch(md); len(match) > 1 {
		entry.Strong = match[1]
	}

	definitions := []string{}

	glossIdxs := FindAllOccurrences(md, "Glosses:")
	expIdxs := FindAllOccurrences(md, "Explanation:")

	for _, defIdx := range glossIdxs {
		str := md[defIdx+len("Glosses:"):]
		endIdx := strings.Index(str, "#")
		if endIdx == -1 {
			endIdx = strings.Index(str, "-")
		}
		if endIdx == -1 {
			endIdx = len(str)
		}
		definition := strings.TrimSpace(str[:endIdx])

		definitions = append(definitions, definition)
	}

	for i, glossIdx := range expIdxs {
		if definitions[i] != "" {
			continue
		}
		str := md[glossIdx+len("Explanation:"):]
		endIdx := strings.Index(str, "#")
		if endIdx == -1 {
			endIdx = strings.Index(str, "-")
		}
		if endIdx == -1 {
			endIdx = len(str)
		}
		gloss := strings.TrimSpace(str[:endIdx])
		if gloss != "" {
			definitions[i] = gloss
		}
	}

	// Extract Senses
	senses := []Sense{}

	for i, d := range definitions {
		sense := Sense{
			Number:     i + 1,
			Definition: d,
		}
		senses = append(senses, sense)
	}
	entry.Senses = senses

	return entry, nil
}

func FindAllOccurrences(text, word string) []int {
	var indexes []int
	offset := 0
	for {
		i := strings.Index(text[offset:], word)
		if i == -1 {
			break
		}
		indexes = append(indexes, offset+i)
		offset += i + len(word)
	}
	return indexes
}

type WordEntry struct {
	Strong string  `json:"strong"`
	Senses []Sense `json:"definitions"`
}

type Sense struct {
	Number     int    `json:"number"`
	Definition string `json:"definition"`
}
