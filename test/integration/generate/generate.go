package main

import (
	_ "embed"
	"event-facematch-backend/test/integration/testlog"
	"io/fs"
	"os"
	"strings"
)

// Generate the swagger client
//go:generate docker compose -f docker-compose-swagger-gen.yaml run -T --rm swagger-gen

// Fix the generated client
//go:generate go run generate.go

//go:embed substitutions/client_call_api_before.txt
var clientCallAPIBefore []byte

//go:embed substitutions/client_call_api_after.txt
var clientCallAPIAfter []byte

//go:embed substitutions/client_import_before.txt
var clientImportsBefore []byte

//go:embed substitutions/client_import_after.txt
var clientImportsAfter []byte

// Post-generate script to fix the generated code - added logging to the client
func main() {
	fileBytes, err := os.ReadFile("swagger/client.go")
	if err != nil {
		panic(err)
	}

	fileString := string(fileBytes)
	fileString = strings.ReplaceAll(
		fileString, string(clientCallAPIBefore), string(clientCallAPIAfter))
	fileString = strings.ReplaceAll(
		fileString, string(clientImportsBefore), string(clientImportsAfter))

	writePermissions := 0600
	err = os.WriteFile("swagger/client.go", []byte(fileString), fs.FileMode(writePermissions))
	if err != nil {
		testlog.Logln("Code generation error")
	}
}
