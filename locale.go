package ksana

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var locales = make(map[string]map[string]string, 0)

func LoadLocales(path string) error {
	logger.Info("Loading i18n from " + path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		fn := f.Name()
		logger.Info("Find locale file " + fn)
		lang := fn[0:(len(fn) - 5)]

		ss := make(map[string]string, 0)
		fd, err := os.Open(path + "/" + fn)
		if err != nil {
			return err
		}
		defer fd.Close()

		err = json.NewDecoder(fd).Decode(&ss)
		if err != nil {
			return err
		}
		locales[lang] = ss

	}
	return nil
}

func T(locale, key string, args ...interface{}) string {
	if l, ok := locales[locale]; ok {
		if v, ok := l[key]; ok {
			return fmt.Sprintf(v, args...)
		}
	}
	return fmt.Sprintf("Translation [%s] not found", key)
}
