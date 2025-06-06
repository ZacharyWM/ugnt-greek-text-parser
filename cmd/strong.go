package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sashabaranov/go-openai"
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

	workEntries := make([]WordEntry, 0, len(fileContents))

	for _, content := range fileContents {
		// response := queryLLM(instructions + " " + content)
		// if response == "" {
		// 	fmt.Println("No response from Ollama for content:", content)
		// 	continue
		// }

		response, err := ParseMarkdownToWordEntry(content)
		if err != nil {
			fmt.Printf("Error parsing markdown content: %v\n", err)
			continue
		}

		// var wordEntry WordEntry
		// err = json.Unmarshal([]byte(response), &wordEntry)
		// if err != nil {
		// 	fmt.Printf("Error unmarshalling JSON response: %v\n", err)
		// 	continue
		// }

		workEntries = append(workEntries, response)

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
	err = encoder.Encode(workEntries)
	if err != nil {
		fmt.Printf("Error encoding JSON to file: %v\n", err)
		return
	}

	fmt.Printf("Successfully written %d entries to %s\n", len(workEntries), outputFile)

}

func queryLLM(message string) string {
	retryClient := retryablehttp.NewClient()
	httpClient := retryClient.StandardClient()

	client := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:    "http://localhost:11434/v1",
		HTTPClient: httpClient,
	})

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "llama3.2",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
	})
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return resp.Choices[0].Message.Content
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

	defIdxs := FindAllOccurrences(md, "Definition:")
	glossIdxs := FindAllOccurrences(md, "Glosses:")

	for _, defIdx := range defIdxs {
		str := md[defIdx+len("Definition:"):]
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

	for i, glossIdx := range glossIdxs {
		if definitions[i] != "" {
			continue
		}
		str := md[glossIdx+len("Glosses:"):]
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

	for _, d := range definitions {
		sense := Sense{
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
	Senses []Sense `json:"senses"`
}

type Sense struct {
	Number     string `json:"number"`
	Definition string `json:"definition"`
}
