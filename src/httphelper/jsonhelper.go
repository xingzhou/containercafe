
package httphelper

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrintJson(b []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	if err != nil {
		fmt.Printf("\nError parsing JSON\n%s\n", b)
	}else {
		//out.WriteTo(os.Stdout)
		fmt.Printf("\n%s\n", out.Bytes())
	}
}
