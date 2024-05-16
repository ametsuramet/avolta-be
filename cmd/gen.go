package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/gommon/log"

	"github.com/iancoleman/strcase"
)

func GenerateFeature(args []string) {

	featureName := strcase.ToCamel(args[1])
	featureNameLower := strcase.ToLowerCamel(args[1])

	// HANDLER
	b, err := os.ReadFile("./handler/sample.go") // just pass the file name
	if err != nil {
		log.Error(err)
	}

	err = os.WriteFile(fmt.Sprintf("./handler/%s.go", args[1]), []byte(strings.ReplaceAll(string(b), "Sample", featureName)), 0644)
	if err != nil {
		log.Error(err)
	}
	fmt.Println("GENERATE FEATURE =>", featureName, "SUCCEED")

	// MODEL
	b, err = os.ReadFile("./model/sample.go") // just pass the file name
	if err != nil {
		log.Error(err)
	}

	err = os.WriteFile(fmt.Sprintf("./model/%s.go", args[1]), []byte(strings.ReplaceAll(string(b), "Sample", featureName)), 0644)
	if err != nil {
		log.Error(err)
	}
	fmt.Println("GENERATE MODEL =>", featureName, "SUCCEED")

	// RESP
	b, err = os.ReadFile("./object/resp/sample_response.go") // just pass the file name
	if err != nil {
		log.Error(err)
	}

	err = os.WriteFile(fmt.Sprintf("./object/resp/%s_response.go", args[1]), []byte(strings.ReplaceAll(string(b), "Sample", featureName)), 0644)
	if err != nil {
		log.Error(err)
	}
	fmt.Println("GENERATE RESP =>", featureName, "SUCCEED")

	// ROUTE
	b, err = os.ReadFile("./router/sample.txt") // just pass the file name
	if err != nil {
		log.Error(err)
	}

	router := strings.ReplaceAll(string(b), "Sample", featureName)
	router = strings.ReplaceAll(router, "sample", featureNameLower)

	b, err = os.ReadFile("./router/apiv1.go") // just pass the file name
	if err != nil {
		log.Error(err)
	}

	newRouter := strings.ReplaceAll(string(b), "// DONT REMOVE THIS LINE", router+"\n\n\t\t// DONT REMOVE THIS LINE")

	err = os.WriteFile("./router/apiv1.go", []byte(newRouter), 0644)
	if err != nil {
		log.Error(err)
	}
	fmt.Println("GENERATE ROUTE =>", featureName, "SUCCEED")

}
