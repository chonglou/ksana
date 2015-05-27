package ksana

import (
	"reflect"
)

var logger, _ = OpenLogger("ksana")

var beans = make(map[reflect.Type]interface{}, 0)

func Map(bean interface{}) {
  t := reflect.TypeOf(bean)
  logger.Debug("Register bean: "+t.String())
	beans[t] = bean
}

func Get(tp reflect.Type) (interface{}, bool) {
  val, ok :=beans[tp]
	return val, ok
}
