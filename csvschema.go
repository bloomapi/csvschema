package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var nonFriendlyCharacters = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
var tooManyUnderscores = regexp.MustCompile(`_+`)

func friendlyName(name string) string {
	friendly := strings.ToLower(name)
	friendly = nonFriendlyCharacters.ReplaceAllString(friendly, "_")
	friendly = tooManyUnderscores.ReplaceAllString(friendly, "_")
	friendly = strings.Trim(friendly, "_")
	return friendly
}

type fieldType struct {
	Name       string
	Expression *regexp.Regexp
}

var types = []fieldType{
	fieldType{
		"timestamp",
		regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\d)$`),
	},
	fieldType{
		"bigint",
		regexp.MustCompile(`^\-?[1-9]\d{9,17}$`),
	},
	fieldType{
		"int",
		regexp.MustCompile(`^((\-?[1-9]\d{0,8})|0)$`),
	},
	fieldType{
		"decimal",
		regexp.MustCompile(`^((\-?\d*(\.\d+)?)|0)$`),
	},
	fieldType{
		"boolean",
		regexp.MustCompile(`^(true|false|TRUE|FALSE|True|False)$`),
	},
}

type FieldInfo struct {
	FieldName string
	FieldType string
	MaxLength int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <csv>\n", os.Args[0])
		os.Exit(1)
	}

	fileReader, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	csvReader := csv.NewReader(fileReader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	columns, err := csvReader.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	discoveredTypeIndexes := make([]int, len(columns))
	discoveredMaxLengths := make([]int, len(columns))
	fields := []FieldInfo{}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for fieldIndex, _ := range columns {
			value := row[fieldIndex]
			valueLength := len(value)

			if discoveredMaxLengths[fieldIndex] < valueLength {
				discoveredMaxLengths[fieldIndex] = valueLength
			}

			if discoveredTypeIndexes[fieldIndex] == len(types) {
				continue
			}

			for {
				if discoveredTypeIndexes[fieldIndex] == len(types) {
					break
				}
				match := types[discoveredTypeIndexes[fieldIndex]].Expression.MatchString(value)
				if match == true || value == "" {
					break
				}
				discoveredTypeIndexes[fieldIndex] += 1
			}
		}
	}

	for fieldIndex, fieldName := range columns {
		var fieldType string
		if discoveredTypeIndexes[fieldIndex] == len(types) {
			fieldType = "character varying"
		} else {
			fieldType = types[discoveredTypeIndexes[fieldIndex]].Name
		}

		fields = append(fields, FieldInfo{
			friendlyName(fieldName),
			fieldType,
			discoveredMaxLengths[fieldIndex],
		})
	}

	fmt.Println("CREATE TABLE sample (")
	for i, field := range fields {
		if field.FieldType == "character varying" {
			fmt.Printf("  " + field.FieldName + " " + field.FieldType + " (" + strconv.Itoa(field.MaxLength*2) + ")")
		} else {
			fmt.Printf("  " + field.FieldName + " " + field.FieldType)
		}

		if i < len(fields)-1 {
			fmt.Printf(",")
		}

		fmt.Printf("\n")
	}
	fmt.Println(");")
}
