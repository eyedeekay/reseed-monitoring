package monitor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eyedeekay/i2p-tools-1/su3"
)

//func MakeReseedData(url, cert string) (map[string]string, error) {
func MakeReseedData(url, cert string) (*su3.File, error) {
	//	var m map[string]string
	url = PrepURL(url)
	su3bytes, err := FetchReseed(url, cert)
	if err != nil {
		return nil, err
	}
	if su3bytes == nil {
		return nil, fmt.Errorf("Reseed file fetched was null, you may have been rate-limited")
	}
	dir := TrimURL(url)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "i2pseeds.su3")
	err = ioutil.WriteFile(path, su3bytes, 0644)
	if err != nil {
		return nil, err
	}
	su3file := su3.New()
	err = su3file.UnmarshalBinary(su3bytes)
	if err != nil {
		return nil, err
	}
	return su3file, nil
}

func MakeReseedDataMap(url, cert string) (map[string]string, error) {
	su3data, err := MakeReseedData(url, cert)
	if err != nil {
		return nil, err
	}
	data := ToMap(su3data)
	bytemap := StatToBytes(data)

	dir := TrimURL(url)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "reseed.json")
	err = ioutil.WriteFile(path, bytemap, 0644)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ToMap(s *su3.File) map[string]string {
	return map[string]string{
		"Format":        fmt.Sprintf("%d", s.Format),
		"SignatureType": fmt.Sprintf("%d", s.SignatureType),
		"FileType":      fmt.Sprintf("%d", s.FileType),
		"ContentType":   fmt.Sprintf("%d", s.ContentType),
		"Version":       fmt.Sprintf("%x", s.Version),
		"SignerID":      fmt.Sprintf("%x", s.SignerID),
		"Content":       fmt.Sprintf("%x", s.Content),
		"Signature":     fmt.Sprintf("%x", s.Signature),
		"SignedBytes":   fmt.Sprintf("%x", s.SignedBytes),
	}
}

func StatToBytes(m map[string]string) []byte {
	var s = "{\n"
	for key, value := range m {
		s += `  "` + key + `":` + `"` + value + `",`
	}
	s += `}`
	return []byte(strings.Replace(strings.Replace(s, ",}", "\n}", -1), ",", ",\n", -1))
}

func TrimURL(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	return url
}

func PrepURL(url string) string {
	if strings.HasPrefix(url, "https://") {
		return url
	}
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix("http://", url)
		return "https://" + url
	}
	return "https://" + url
}
