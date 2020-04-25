package apihandlers

import (
	"net/http"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"time"
	"github.com/go-errors/errors"
	"encoding/json"
)

type infunc func(http.ResponseWriter, *http.Request)
func Recover( fn infunc ) infunc{
	return func(w http.ResponseWriter, r *http.Request){

		w.Header().Set("Access-Control-Allow-Origin", "*")

		defer func(r *http.Request) {
			ctx := appengine.NewContext(r)
			if err := recover(); err != nil {
				err := errors.Wrap(err, 1)
				fmt.Fprintf(w,"ERROR: %v", err)
				log.Criticalf(ctx,"ERROR: %v", err.ErrorStack())
			}
		}(r)

		fn(w,r)
	}
}

func RecoverApi( fn infunc ) infunc {
	return func(w http.ResponseWriter, r *http.Request){
		defer func(r *http.Request) {
			ctx := appengine.NewContext(r)
			if errr := recover(); errr != nil {

				err := errors.New(errr)
				mp := ErrorType{
					Error: err.Error(),
				}
				jsonstr , _ := json.Marshal(mp)
				fmt.Fprintf(w,"%s", string(jsonstr))
				ms := time.Now().UnixNano()
				log.Infof(ctx,"ERROR(%d): ------------ ", ms)
				log.Infof(ctx,"ERROR(%d): %s", ms, err.Error())
				log.Criticalf(ctx,"ERROR(%d): %v", ms, err.ErrorStack())
			}
		}(r)

		fn(w,r)
	}
}

func PanicIfNotNil(err error){
	if err != nil {
		panic(errors.Wrap(err, 1))
	}
}

func PanicWithMsg( str string){
	panic(fmt.Errorf("%s",str))
}

type ErrorType struct{
	Error string
}

func WriteAsJson(w http.ResponseWriter, obj interface{}){
	PanicIfNotNil(json.NewEncoder(w).Encode(obj))
}
func WriteAsJsonList(w http.ResponseWriter, obj interface{}){

	jsonstr, err := json.Marshal(obj)
	PanicIfNotNil(err)

	w.Write(jsonstr)
}