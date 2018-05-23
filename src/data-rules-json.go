package beholder

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

const srdURL = "https://raw.githubusercontent.com/BTMorton/dnd-5e-srd/3d3afe305e0178d08baf674f57411a588c8f4ac5/5esrd.json"

var ignoreSections = map[string]bool{
	"Legal Information":                      true,
	"Feats":                                  true,
	"Magic Items":                            true,
	"Spell Descriptions":                     true,
	"Appendix PH-A: Conditions":              true,
	"Appendix MM-A: Miscellaneous Creatures": true,
}

var headerReplacements = map[string]string{
	// (sic): the source json has a typo
	"Proficieny Bonus":  "Prof.",
	"Proficiency Bonus": "Prof.", // in case it gets fixed
	"Spells Known":      "Spells",
	"Cantrips Known":    "Cantrips",
}

const tableColSpace = "  "
const maxColWidth = 30

// NewJSONRulesSource constructs a DataSource that generates
// rules from JSON downloaded externally
func NewJSONRulesSource() (DataSource, error) {
	localPath, err := homedir.Expand("~/.config/beholder/srd.json")
	if err != nil {
		return nil, err
	}

	return &networkDataSource{
		URL:       srdURL,
		localPath: localPath,
		delegate: &jsonRulesSource{
			localPath: localPath,
		},
	}, nil
}

type jsonRulesSource struct {
	localPath string
}

func (d *jsonRulesSource) GetEntities() ([]Entity, error) {
	f, err := os.Open(d.localPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseRulesJSON(bufio.NewReader(f))
}

func parseRulesJSON(reader *bufio.Reader) ([]Entity, error) {
	// NOTE we parse everything by hand because the format is wacky
	// and assumes json maps are ordered (they are not)

	decoder := json.NewDecoder(reader)

	rootRule, err := readRuleSection(decoder, "")
	if err != nil {
		return nil, err
	}

	_, generated := generateEntities(RuleEntity, rootRule, true, ignoreSections)
	return generated, nil
}

func readRuleSection(decoder *json.Decoder, header string) (*ruleParts, error) {
	part := &ruleParts{
		name: header,
	}

	tok, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	switch v := tok.(type) {
	case string:
		// single string content
		readString(v, part)
		return part, nil

	case json.Delim:
		switch v {
		case '{':
			// normal case
			break

		case '[':
			err := readContent(decoder, part, 1)
			return part, err

		default:
			return nil, fmt.Errorf("(in %s) Invalid json: %s", header, tok)
		}

	default:
		return nil, fmt.Errorf("(in %s) Invalid json: %s", header, tok)
	}

	// this should perhaps be refactored into its own function:
	for {
		tok, err := decoder.Token()
		if tok == nil {
			return part, nil
		} else if err != nil {
			return nil, err
		}

		switch v := tok.(type) {
		case string:
			switch v {
			case "content":
				if err := readContent(decoder, part, 0); err != nil {
					return nil, err
				}

			case "table":
				if err := readTable(decoder, part); err != nil {
					return nil, err
				}

			default:
				section, err := readRuleSection(decoder, v)
				if err != nil {
					return nil, err
				} else if section == nil {
					return nil, fmt.Errorf("Read nil section %s", v)
				}

				part.parts = append(part.parts, section)
			}

		case json.Delim:
			if v == '}' {
				// done!
				return part, nil
			}
		}
	}
}

func readContent(decoder *json.Decoder, part *ruleParts, initialNest int) error {

	nest := initialNest

	for {
		tok, err := decoder.Token()
		if tok == nil {
			break
		} else if err != nil {
			return err
		}

		switch v := tok.(type) {
		case json.Delim:
			switch v {
			case '[':
				nest++
			case ']':
				nest--
				if nest == 0 {
					return nil
				}

			case '{':
				// table, probably
				if tableTok, err := decoder.Token(); err != nil {
					return err
				} else if tableTok.(string) != "table" {
					return fmt.Errorf("Expected a table, but: %s", tableTok)
				}
				readTable(decoder, part)

				if endTok, err := decoder.Token(); err != nil {
					return err
				} else if endTok.(json.Delim) != '}' {
					return fmt.Errorf("Expected table to end, but: %s", endTok)
				}

			case '}':
				return fmt.Errorf("Unexpected } in %s", part.name)
			}

		case string:
			readString(v, part)

			if nest == 0 {
				// if nest == 0, content was a single string value
				return nil
			}
		}
	}

	return nil
}

func readTable(decoder *json.Decoder, part *ruleParts) error {
	var builder bytes.Buffer
	builder.WriteString("[::b]")

	tableObj := make(map[string][]string)
	headers := make([]string, 0, 3)
	headerFormats := make([]string, 0, 3)
	values := 0

	if tok, err := decoder.Token(); err != nil {
		return err
	} else if tok.(json.Delim) != '{' {
		return fmt.Errorf("Expected table to start, but: %s", tok)
	}
	for {
		tok, err := decoder.Token()
		if err != nil {
			return err
		} else if d, ok := tok.(json.Delim); ok && d == '}' {
			break
		}
		header := tok.(string)
		if replacement, ok := headerReplacements[header]; ok {
			header = replacement
		}

		var rows []string
		if err := decoder.Decode(&rows); err != nil {
			return err
		}

		// calculate the widest row
		widestRow := len(header)
		for _, row := range rows {
			l := len(row)
			if l > widestRow {
				widestRow = l
			}
		}
		if widestRow > maxColWidth {
			widestRow = maxColWidth
		}

		format := fmt.Sprintf("%%%ds", widestRow)

		tableObj[header] = rows
		values = len(rows)
		headers = append(headers, header)
		headerFormats = append(headerFormats, format)

		builder.WriteString(fmt.Sprintf(format, header))
		builder.WriteString(tableColSpace)
	}

	builder.WriteString("[::-]\n")
	for i := 0; i < values; i++ {
		for j, header := range headers {
			builder.WriteString(fmt.Sprintf(
				headerFormats[j], trimmed(tableObj[header][i], maxColWidth)),
			)
			builder.WriteString(tableColSpace)
		}
		builder.WriteString("\n")
	}

	part.parts = append(part.parts, builder.String())

	return nil
}

func readString(s string, part *ruleParts) {
	if strings.HasPrefix(s, "***") {
		part.parts = append(part.parts, "")
	}

	part.parts = append(part.parts, formatText(s))
}

type formatter struct {
	regex       *regexp.Regexp
	replacement string
}

var replacements = []formatter{
	formatter{
		regexp.MustCompile(`\*\*\*([^*]+)\*\*\*`),
		"<b>$1</b>",
	},

	formatter{
		regexp.MustCompile(`\*\*([^*]+)\*\*`),
		"<i>$1</i>",
	},

	formatter{
		regexp.MustCompile(`\*([^*]+)\*`),
		"<u>$1</u>",
	},
}

func formatText(s string) string {
	for _, r := range replacements {
		s = r.regex.ReplaceAllString(s, r.replacement)
	}
	return s
}

func trimmed(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return fmt.Sprintf("%s…", s[:maxLen-1])
}
