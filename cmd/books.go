package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Book struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	ID     int     `json:"id"`
	BookID int     `json:"bookId"`
	Number int     `json:"number"`
	Verses []Verse `json:"verses"`
}

type Verse struct {
	ID        int    `json:"id"`
	ChapterID int    `json:"chapterId"`
	Number    int    `json:"number"`
	Words     []Word `json:"words"`
}

type Word struct {
	ID      int    `json:"id"`
	VerseID int    `json:"verseId"`
	Text    string `json:"text"`
	Lemma   string `json:"lemma"`
	Strong  string `json:"strong"`
	Morph   string `json:"morph"`
}

func parseBooks() {
	fmt.Printf("Parsing started\n")

	books, err := parseUGNTFiles()
	if err != nil {
		fmt.Printf("Error parsing UGNT files: %v\n", err)
		return
	}

	err = exportBooksToJSON(books)
	if err != nil {
		fmt.Printf("Error exporting books to JSON: %v\n", err)
		return
	}

	fmt.Printf("Parsing completed successfully\n")
}

func parseUGNTFiles() ([]Book, error) {
	var books []Book

	files, err := os.ReadDir("ugnt")
	if err != nil {
		return nil, fmt.Errorf("error reading ugnt directory: %v", err)
	}

	bookID := 1

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".usfm") {
			continue
		}

		filePath := filepath.Join("ugnt", file.Name())
		book, err := parseUGNTFile(filePath, bookID)
		if err != nil {
			return nil, fmt.Errorf("error parsing file %s: %v", file.Name(), err)
		}

		books = append(books, book)
		bookID++
	}

	return books, nil
}

func parseUGNTFile(filePath string, bookID int) (Book, error) {
	book := Book{
		ID:       bookID,
		Chapters: []Chapter{},
	}

	file, err := os.Open(filePath)
	if err != nil {
		return book, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var currentChapter *Chapter
	var currentVerse *Verse
	var wordID int = 1
	var verseID int = 1
	var chapterID int = 1

	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if strings.HasPrefix(line, "\\h ") {
			book.Title = strings.TrimPrefix(line, "\\h ")
			continue
		}

		if strings.HasPrefix(line, "\\c ") {
			chapterNumStr := strings.TrimPrefix(line, "\\c ")
			chapterNum, err := strconv.Atoi(chapterNumStr)
			if err != nil {
				return book, fmt.Errorf("line %d: invalid chapter number: %s", lineNum, chapterNumStr)
			}

			chapter := Chapter{
				ID:     chapterID,
				BookID: book.ID,
				Number: chapterNum,
				Verses: []Verse{},
			}
			book.Chapters = append(book.Chapters, chapter)
			currentChapter = &book.Chapters[len(book.Chapters)-1]
			chapterID++
			verseID = 1 // Reset verse ID for new chapter
			continue
		}

		// Check for new verse
		if strings.HasPrefix(line, "\\v ") {
			if currentChapter == nil {
				return book, fmt.Errorf("line %d: verse found before chapter", lineNum)
			}

			parts := strings.SplitN(line, " ", 3)
			if len(parts) < 2 {
				return book, fmt.Errorf("line %d: malformed verse line", lineNum)
			}

			verseNumStr := parts[1]
			verseNum, err := strconv.Atoi(verseNumStr)
			if err != nil {
				return book, fmt.Errorf("line %d: invalid verse number: %s", lineNum, verseNumStr)
			}

			verse := Verse{
				ID:        verseID,
				ChapterID: currentChapter.ID,
				Number:    verseNum,
				Words:     []Word{},
			}
			currentChapter.Verses = append(currentChapter.Verses, verse)
			currentVerse = &currentChapter.Verses[len(currentChapter.Verses)-1]
			verseID++

			if len(parts) == 3 && strings.Contains(parts[2], "\\w ") {
				processWords(parts[2], currentVerse, &wordID)
			}
			continue
		}

		if currentVerse != nil && strings.Contains(line, "\\w ") {
			processWords(line, currentVerse, &wordID)
		}
	}

	if err := scanner.Err(); err != nil {
		return book, fmt.Errorf("error scanning file: %v", err)
	}

	return book, nil
}

func processWords(line string, verse *Verse, wordID *int) {
	// Split on \w to find all word markers
	parts := strings.Split(line, "\\w ")

	for i, part := range parts {
		if i == 0 || part == "" {
			continue // Skip the part before the first \w or empty parts
		}

		// Find the end of the word text (up to the | character)
		pipeIndex := strings.Index(part, "|")
		if pipeIndex == -1 {
			continue // Skip if no | character
		}

		wordText := part[:pipeIndex]
		wordProps := part[pipeIndex+1:]

		// Find the end of the word entry (marked by \w*)
		wordEndIndex := strings.Index(wordProps, "\\w*")
		if wordEndIndex == -1 {
			continue // Skip if no closing \w*
		}

		wordProps = wordProps[:wordEndIndex]

		// Extract word properties
		lemma := extractProperty(wordProps, "lemma=")
		strong := extractProperty(wordProps, "strong=")
		morph := extractProperty(wordProps, "x-morph=")

		word := Word{
			ID:      *wordID,
			VerseID: verse.ID,
			Text:    wordText,
			Lemma:   lemma,
			Strong:  strong,
			Morph:   morph,
		}

		verse.Words = append(verse.Words, word)
		(*wordID)++
	}
}

func extractProperty(text, propName string) string {
	start := strings.Index(text, propName)
	if start == -1 {
		return ""
	}

	start += len(propName)
	if start >= len(text) {
		return ""
	}

	// Check if the property is quoted
	if text[start] == '"' {
		start++
		end := strings.Index(text[start:], "\"")
		if end == -1 {
			return ""
		}
		return text[start : start+end]
	}

	// If not quoted, read until space or end
	end := start
	for end < len(text) && text[end] != ' ' {
		end++
	}

	return text[start:end]
}

func exportBooksToJSON(books []Book) error {
	if err := os.MkdirAll("output", 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	for _, book := range books {
		fileName := fmt.Sprintf("%d_%s.json", book.ID, book.Title)
		filePath := filepath.Join("output", fileName)

		bookJSON, err := json.MarshalIndent(book, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal book %s to JSON: %v", book.Title, err)
		}

		if err := os.WriteFile(filePath, bookJSON, 0644); err != nil {
			return fmt.Errorf("failed to write book %s to file: %v", book.Title, err)
		}
	}

	return nil
}
