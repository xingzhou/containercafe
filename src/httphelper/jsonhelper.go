
package httphelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

func PrintJson(b []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	if err != nil {
		log.Printf("\nError parsing JSON\n%s\n", b)
	}else {
		//out.WriteTo(os.Stdout)
		log.Printf("\n%s\n", out.Bytes())
	}
}

func PrettyJson(b []byte) string{
	var out bytes.Buffer
	var s string
	err := json.Indent(&out, b, "", "\t")
	if err != nil {
		s = fmt.Sprintf("\n%s\n", b)
	}else {
		//out.WriteTo(os.Stdout)
		s = fmt.Sprintf("\n%s\n", out.Bytes())
	}
	return s
}
