package models

import "reflect"

type ModelStruct struct {
	TableName   string
	StructField []reflect.StructField
}
