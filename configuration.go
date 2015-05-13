package ksana

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
)

type property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type bean struct {
	Name       string     `xml:"name,attr"`
	Class      string     `xml:"class,attr"`
	Properties []property `xml:"property"`
}

func (b *bean) getString(name string, def string) string {
	for _, p := range b.Properties {
		if p.Name == name {
			return p.Value
		}
	}
	return def
}

func (b *bean) getInt(name string, def int) int {
	for _, p := range b.Properties {
		if p.Name == name {
			i, err := strconv.Atoi(p.Value)
			if err == nil {
				return i
			}
			break
		}
	}
	return def
}

type configuration struct {
	XMLName xml.Name `xml:"ksana"`

	Name string `xml:"name,attr"`
	Mode string `xml:"mode,attr"`

	Port int `xml:"port,attr"`

	Beans []bean `xml:"bean"`
}

func loadConfiguration(file string, cfg *configuration) error {
	xf, err := os.Open(file)
	if err != nil {
		return err
	}
	defer xf.Close()

	var data []byte
	data, err = ioutil.ReadAll(xf)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}
