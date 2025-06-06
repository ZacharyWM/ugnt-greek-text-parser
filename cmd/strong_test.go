package main

import (
	"os"
	"testing"
)

func TestThing(t *testing.T) {
	type StrongParseTest struct {
		FilePath string
		Expected WordEntry
	}

	tests := []StrongParseTest{
		{
			FilePath: "/Users/zachm/Documents/github/ugnt-greek-text-parser/strong/G01590/01.md",
			Expected: WordEntry{
				Strong: "G01590",
				Senses: []Sense{
					{
						Number:     "1.0",
						Definition: "pertaining to the cause of something",
					},
					{
						Number:     "2.0",
						Definition: "the source",
					},
					{
						Number:     "3.0",
						Definition: "guilty, having a valid basis for a charge",
					},
				},
			},
		},
		{
			FilePath: "/Users/zachm/Documents/github/ugnt-greek-text-parser/strong/G00680/01.md",
			Expected: WordEntry{
				Strong: "G00680",
				Senses: []Sense{
					{
						Number:     "1.0",
						Definition: "the countryside as distinct from settled towns and villages",
					},
					{
						Number:     "2.0",
						Definition: "fields in which plants grow",
					},
					{
						Number:     "3.0",
						Definition: "a plot of ground",
					},
				},
			},
		},
		{
			FilePath: "/Users/zachm/Documents/github/ugnt-greek-text-parser/strong/G00740/01.md",
			Expected: WordEntry{
				Strong: "G00740",
				Senses: []Sense{
					{
						Number:     "1.0",
						Definition: "a contest, wrestling;",
					},
					{
						Number:     "2.0",
						Definition: "great fear, agony, anguish;",
					},
				},
			},
		},
	}

	for _, test := range tests {
		// if test.Expected.Strong != "G00740" {
		// 	continue
		// }
		t.Run(test.FilePath, func(t *testing.T) {
			fileContent, err := os.ReadFile(test.FilePath)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", test.FilePath, err)
			}

			entry, err := ParseMarkdownToWordEntry(string(fileContent))
			if err != nil {
				t.Fatalf("Failed to parse file %s: %v", test.FilePath, err)
			}

			if entry.Strong != test.Expected.Strong {
				t.Errorf("Expected Strong number %s; got %s", test.Expected.Strong, entry.Strong)
			}

			if len(entry.Senses) != len(test.Expected.Senses) {
				t.Errorf("Expected %d senses; got %d", len(test.Expected.Senses), len(entry.Senses))
			}

			for i, sense := range entry.Senses {
				if sense.Definition != test.Expected.Senses[i].Definition {
					t.Errorf("Expected sense %d definition '%s'; got '%s'", i+1, test.Expected.Senses[i].Definition, sense.Definition)
				}
			}
		})
	}
}
