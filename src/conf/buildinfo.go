package conf

import(
	"fmt"
	"io/ioutil"
	"encoding/json"
)

const Version = 1.10

func GetVerStr() string{
	buildId, buildDate := GetBuildInfo("build.info")
	s := fmt.Sprintf("ver=%4.2f  build_id=%s  build_date=%s", Version, buildId, buildDate)
	return s
}

func GetBuildInfo(fname string) (buildId string, buildDate string){

	b, err := ioutil.ReadFile(fname)
	if err != nil {
		Log.Printf("error reading build info file %s", fname)
		return
	}

	// fill map
	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil{
		Log.Printf("error in json unmarshalling, err=%v", err)
		return
	}

	m := f.(map[string]interface{})
	for k, v := range m {
		switch k {
		case "build-id":
			buildId = v.(string)
			break
		case"build-date":
			buildDate = v.(string)
			break
		default:
			break
		}
	}

	return
}
