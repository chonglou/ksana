package ksana

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

var logger, _ = OpenLogger("ksana")

var beans = make(map[reflect.Type]interface{}, 0)

func Map(bean interface{}) {
	t := reflect.TypeOf(bean)
	logger.Debug("Register bean: " + t.String())
	beans[t] = bean
}

func Get(tp reflect.Type) (interface{}, bool) {
	val, ok := beans[tp]
	return val, ok
}

func Obj2bit(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Bit2obj(data []byte, obj interface{}) error {
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(data)
	err := dec.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}
