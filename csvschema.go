package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	Name string
	Test func(string) bool
}

var types = []fieldType{
	fieldType{
		"timestamp",
		func (value string) bool {
			_, err := time.Parse("01/02/2006", value)
			if err == nil {
				return true
			}
			
			_, err = time.Parse(time.RFC3339, value)
			return err == nil
		},
	},
	fieldType{
		"int",
		func (value string) bool {
			_, err := strconv.ParseInt(value, 10, 32)
			return err == nil
		},
	},
	fieldType{
		"bigint",
		func (value string) bool {
			_, err := strconv.ParseInt(value, 10, 64)
			return err == nil
		},
	},
	fieldType{
		"decimal",
		func (value string) bool {
			_, err := strconv.ParseFloat(value, 32)
			return err == nil
		},
	},
	fieldType{
		"boolean",
		func (value string) bool {
			_, err := strconv.ParseBool(value)
			return err == nil
		},
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

	_, filename := path.Split(os.Args[1])
	fileParts := strings.Split(filename, ".")
	tablename := strings.Join(fileParts[:len(fileParts)-1], ".")

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
				if discoveredTypeIndexes[fieldIndex] == len(types) || value == "" {
					break
				}
				
				match := types[discoveredTypeIndexes[fieldIndex]].Test(value)
				if match == true {
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

	fmt.Println("CREATE TABLE \"" + tablename + "\" (")
	for i, field := range fields {
		if field.FieldType == "character varying" {
			fmt.Printf("  \"" + field.FieldName + "\" " + field.FieldType + " (" + strconv.Itoa(field.MaxLength) + ")")
		} else {
			fmt.Printf("  \"" + field.FieldName + "\" " + field.FieldType)
		}

		if i < len(fields)-1 {
			fmt.Printf(",")
		}

		fmt.Printf("\n")
	}
	fmt.Println(");")
}
