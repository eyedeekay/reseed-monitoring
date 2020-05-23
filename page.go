package monitor

import (
	"fmt"
	"github.com/yosssi/gohtml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var headline = `<!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="utf-8">
      <title> I2P Reseed Monitoring </title>
      <link rel="stylesheet" href="/styles.css">
	  <script src="/script.js"></script>
    </head>
    <body>
`

var footline = `
    </body>
</html>
`

func GeneratePageData() []error {
	config, err := SortedMap("config.json")
	if err != nil {
		log.Fatal(err)
	}
	return SortedMonitor(config)
}

func GeneratePage() (string, error) {
	var ret string
	ret += headline
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(info.Name(), ".json") {
				if info.Name() != "config.json" {
					ret += "\n" + `  <div class="` + TrimDir(path) + `">` + "\n"
					f, e := ioutil.ReadFile(path)
					if e == nil {
						pre := string(f)
						pre = strings.TrimPrefix(pre, "{")
						pre = strings.TrimSuffix(pre, "}")
						pre = strings.Replace(pre, `"`, "", -1)
						pre2 := strings.Split(pre, ",")
						for _, v := range pre2 {
							ky, vy := Split2(v)
							if ky != "Content" {
								ret += `    <div class="` + TrimDir(path) + ` keyvalue">`
								ret += "\n" + `      <span class="` + TrimDir(path) + " " + ky + ` key">` + ky + "\n"
								ret += "\n     </span>\n"
								ret += "\n" + `     <span class="` + TrimDir(path) + " " + ky + ` value">` + vy + "\n"
								ret += "\n    </span>\n"
							}
							ret += `    </div>`
						}

					} else {
						ret = e.Error()
					}
					ret += "\n  </div>\n"
					fmt.Println(path, info.Size())
				}
			}
			return nil
		})
	if err != nil {
		return "", err
	}
	ret += footline
	return gohtml.Format(ret), nil
}

func TrimDir(path string) string {
	return strings.Replace(strings.Replace(strings.Split(filepath.Dir(path), ":")[0], ".", "", -1), "/", "", -1)
}
