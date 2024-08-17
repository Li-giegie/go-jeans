package test

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type config struct {
	sliceMaxLen int
	sliceMinLen int
	strMaxLen   int
	strMinLen   int
}

var conf = config{
	sliceMaxLen: 5,
	sliceMinLen: 1,
	strMinLen:   1,
	strMaxLen:   8,
}

func genRandValue(arg interface{}, c ...config) {
	if len(c) == 0 {
		c = []config{conf}
	}
	rv, ok := arg.(reflect.Value)
	if ok {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	} else {
		rv = reflect.ValueOf(arg).Elem()
		if rv.Kind() == reflect.Invalid {
			panic("randIntValue : invalid arg")
		}
	}
	switch rv.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		rv.SetInt(rnd.Int63())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		rv.SetUint(rnd.Uint64())
	case reflect.Float32, reflect.Float64:
		rv.SetFloat(rnd.Float64())
	case reflect.String:
		rv.SetString(randStr(rnd.Intn(c[0].strMaxLen) + c[0].strMaxLen))
	case reflect.Bool:
		if rnd.Int() >= 0 {
			rv.SetBool(true)
		} else {
			rv.SetBool(false)
		}
	case reflect.Slice:
		switch rv.Interface().(type) {
		case []uint8:
			rv.SetBytes([]byte(randStr(rnd.Intn(c[0].strMaxLen) + c[0].strMinLen)))
		default:
			capLen := rnd.Intn(c[0].sliceMaxLen) + c[0].sliceMinLen
			val := reflect.MakeSlice(rv.Type(), capLen, capLen)
			for i := 0; i < val.Len(); i++ {
				genRandValue(val.Index(i), c[0])
			}
			rv.Set(val)
		}
	case reflect.Struct:
		if !rv.CanSet() {
			log.Println("private field", rv.String())
			return
		}
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
			fmt.Println("1.1", rv.Kind())
		}
		for i := 0; i < rv.NumField(); i++ {
			if rv.Field(i).Kind() == reflect.Struct {
				genRandValue(rv.Field(i))
				continue
			}
			if rv.Field(i).Kind() == reflect.Ptr {
				if rv.Field(i).IsNil() {
					rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
				}
			}
			genRandValue(rv.Field(i))
		}
	default:
		panic(fmt.Sprintf("Not supported type %s %s", rv.Type().Name(), rv.String()))
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rnd.Intn(len(letters))]
	}
	return string(b)
}
