package testlog

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	cases "golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var caser cases.Caser

func init() {
	caser = cases.Title(language.English)
}

func Logln(a ...any) {
	_, _ = fmt.Println(a...)
}

func Logf(format string, a ...any) {
	_, _ = fmt.Printf(format, a...)
}

func StartFlow(name string, description string) {
	Logln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	Logf("Flow: %s: start\n", caser.String(name))
	Logf("Description: %s\n\n", description)
}

func EndFlow(name string) {
	Logf("\nFlow: %s: success\n", caser.String(name))
	Logln("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}

func LogRequest(request *http.Request) {
	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = fmt.Printf("-----------------------------\nREQUEST:\n%s", string(reqDump))
}

func LogResponse(response *http.Response) {
	respDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Fatal(err)
	}

	Logf("\nRESPONSE:\n%s", string(respDump))
}
