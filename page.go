package monitor

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eyedeekay/i2p-tools-1/su3"
	"github.com/yosssi/gohtml"
)

var headline = `<!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="utf-8">
      <title> I2P Reseed Monitoring </title>
      <link rel="stylesheet" href="style.css">
	  <script src="script.js"></script>
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
	var menu string
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(info.Name(), ".json") {
				if info.Name() != "config.json" {
					ret += "\n" + `  <div class="` + TrimDir(path) + ` Reseed" id="` + TrimDir(path) + `">` + "\n"
					ret += "\n" + `    <h4><a href="#` + TrimDir(path) + `">` + filepath.Dir(path) + "</a></h4>\n"
					menu += "\n" + `  <h3> + <a class="` + TrimDir(path) + `" href="#` + TrimDir(path) + `">` + filepath.Dir(path) + "</a></h3>\n"
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
								ret += "\n      </span>\n"
								ret += "\n" + `      <span class="` + TrimDir(path) + " " + ky + ` value">` + vy + "\n"
								ret += "\n      </span>\n"
								ret += `    </div>`
							}
						}
						su3bytes, e := ioutil.ReadFile(filepath.Join(filepath.Dir(path), "i2pseeds.su3"))
						if e == nil {
							su3file := su3.New()
							err = su3file.UnmarshalBinary(su3bytes)
							if err != nil {
								return err
							}
							Valid, Errored := CheckKeys(su3file, "./reseed")
							if Errored != nil {
								ret += `<div class="` + TrimDir(path) + ` Invalid">` + Errored.Error() + `</div>`
							} else {
								ret += `<div class="` + TrimDir(path) + ` Valid">` + Valid + `</div>`
							}
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
	return gohtml.Format(headline + menu + ret + footline), nil
}

func CheckKeys(su3file *su3.File, dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, v := range files {
		file, err := ioutil.ReadFile(filepath.Join(dir, v.Name()))
		if err != nil {
			return "", err
		}
		pemfile, _ := pem.Decode(file)
		crt, err := x509.ParseCertificate(pemfile.Bytes)
		if err == nil {
			err = su3file.VerifySignature(crt)
			if err == nil {
				return "Reseed verified by certificate" + v.Name(), nil
			}
		} else {
			return "", err
		}
	}
	return "", fmt.Errorf("No reseed certs were found")
}

func TrimDir(path string) string {
	return strings.Replace(strings.Replace(strings.Split(filepath.Dir(path), ":")[0], ".", "", -1), "/", "", -1)
}
