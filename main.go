package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/billitech/json-to-dart/utils"
	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var args struct {
		Src  string `arg:"positional"`
		Dist string `arg:"positional"`
	}
	arg.MustParse(&args)

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(args.Src) < 1 {
		args.Src = filepath.Join(path, "jsons")
	}

	if len(args.Dist) < 1 {
		args.Dist = filepath.Join(path, "lib", "models")
	}

	files, err := os.ReadDir(args.Src)
	if err != nil {
		panic(err)
	}

	err = utils.EnsureDir(args.Dist)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			name := strings.TrimSuffix(file.Name(), ".json")
			path := filepath.Join(args.Src, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}

			var jsonObj map[string]interface{}

			err = json.Unmarshal(data, &jsonObj)
			if err != nil {
				panic(err)
			}

			dartStr := utils.JsonToDart(jsonObj, "$", name)
			err = os.WriteFile(filepath.Join(args.Dist, fmt.Sprintf("%s.dart", strcase.ToKebab(name))), []byte(dartStr), 0644)
			if err != nil {
				panic(err)
			}

			greenInfo := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("Processed file:  %s\n", greenInfo(path))
		}
	}

	doneNotice := color.New(color.Bold, color.FgWhite, color.BgGreen).PrintlnFunc()
	doneNotice("  Done  ")
}
