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
			FilePath: "../strong/G01590/01.md",
			Expected: WordEntry{
				Strong: "G01590",
				Senses: []Sense{
					{
						Definition: "causative of",
					},
					{
						Definition: "the source",
					},
					{
						Definition: "blameworthy",
					},
				},
			},
		},
		{
			FilePath: "../strong/G00680/01.md",
			Expected: WordEntry{
				Strong: "G00680",
				Senses: []Sense{
					{
						Definition: "the countryside",
					},
					{
						Definition: "fields",
					},
					{
						Definition: "a piece of ground",
					},
				},
			},
		},
		{
			FilePath: "../strong/G00740/01.md",
			Expected: WordEntry{
				Strong: "G00740",
				Senses: []Sense{
					{
						Definition: "a contest, wrestling;",
					},
					{
						Definition: "great fear, agony, anguish;",
					},
				},
			},
		},
		{
			FilePath: "../strong/G01220/01.md",
			Expected: WordEntry{
				Strong: "G01220",
				Senses: []Sense{
					{
						Definition: "of a goat",
					},
				},
			},
		},
		{
			FilePath: "../strong/G08495/01.md",
			Expected: WordEntry{
				Strong: "G08495",
				Senses: []Sense{
					{
						Definition: "to boast;",
					},
				},
			},
		},
		{
			FilePath: "../strong/G27790/01.md",
			Expected: WordEntry{
				Strong: "G27790",
				Senses: []Sense{
					{
						Definition: "garden",
					},
				},
			},
		},
		{
			FilePath: "../strong/G49190/01.md",
			Expected: WordEntry{
				Strong: "G49190",
				Senses: []Sense{
					{
						Definition: "to break in pieces, crush",
					},
				},
			},
		},
		{
			FilePath: "../strong/G49286/01.md",
			Expected: WordEntry{
				Strong: "G49286",
				Senses: []Sense{
					{
						Definition: "", // No definition provided in the original file
					},
				},
			},
		},
		{
			FilePath: "../strong/G49340/01.md",
			Expected: WordEntry{
				Strong: "G49340",
				Senses: []Sense{
					{
						Definition: "to determine, agree, covenant",
					},
				},
			},
		},
		{
			FilePath: "../strong/G00070/01.md",
			Expected: WordEntry{
				Strong: "G00070",
				Senses: []Sense{
					{
						Definition: "Abia, Abijah",
					},
					{
						Definition: "Abijah",
					},
				},
			},
		},
		{
			FilePath: "../strong/G00110/01.md",
			Expected: WordEntry{
				Strong: "G00110",
				Senses: []Sense{
					{
						Definition: "Abraham",
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
