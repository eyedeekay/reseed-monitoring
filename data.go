package monitor

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"
)

func Monitor(info map[string]string) []error {
	var errs []error
	for url, cert := range info {
		data, err := MakeReseedDataMap(url, cert, true)
		if err != nil {
			errs = append(errs, err)
		}
		if data == nil {
			errs = append(errs, err)
		}

	}
	return errs
}

func SortedMonitor(info []string) []error {
	var errs []error
	for _, kv := range info {
		url, cert := Split(kv)
		data, err := MakeReseedDataMap(url, cert, true)
		if err != nil {
			errs = append(errs, err)
		}
		if data == nil {
			errs = append(errs, err)
		}

	}
	return errs
}

func LoadMap(file string) (map[string]string, error) {
	byteValue, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}
	resultMap := make(map[string]string)
	for k, v := range result {
		resultMap[k] = v.(string)
	}
	return resultMap, nil
}

func SortedMap(file string) ([]string, error) {
	lm, err := LoadMap(file)
	if err != nil {
		return nil, err
	}
	var el []string
	for k := range lm {
		el = append(el, k)
	}
	sort.Strings(el)
	var ret []string
	for _, k := range el {
		ret = append(ret, k+"="+lm[k])
	}
	return ret, nil
}

func Split2(kv string) (string, string) {
	val := strings.SplitN(kv, ":", 2)
	if len(val) < 2 {
		return "", ""
	}
	return val[0], val[1]
}

func Split(kv string) (string, string) {
	val := strings.SplitN(kv, "=", 2)
	if len(val) < 2 {
		return "", ""
	}
	return val[0], val[1]
}
