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
			FilePath: "../strong/G00680/01.md",
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
			FilePath: "../strong/G00740/01.md",
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
		{
			FilePath: "../strong/G01220/01.md",
			Expected: WordEntry{
				Strong: "G01220",
				Senses: []Sense{
					{
						Number:     "1.0",
						Definition: "relating to a goat",
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
						Number:     "1.0",
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
						Number:     "1.0",
						Definition: "A section of land used for growing different types of plants and trees and often to grow plants and trees that produce fruits and vegetables",
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
						Number:     "1.0",
						Definition: "to break, crush",
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
						Number:     "1.0",
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
						Number:     "1.0",
						Definition: "to agree on a course of action",
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
						Number:     "1.0",
						Definition: "Abijah, a son of Rehoboam who is in the ancestral line of Jesus",
					},
					{
						Number:     "2.0",
						Definition: "Abijah, the founded of a division of priests of which Zacharaias was a part.  This division of priestly service of which Abijah was a part is described in 1 Chronicles 24 [1Chr 24:3, 10](1ch 24:3, 10)",
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
