package apihandlers

import (
	"net/http"
	"fmt"
	"appengine"
	"github.com/go-errors/errors"
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
				ctx.Infof("ERROR: %v", err.ErrorStack())
			}
		}(r)

		fn(w,r)
	}
}

func PanicIfNil(err error){
	if err != nil {
		panic(err)
	}
}

func PanicWithMsg( str string){
	panic(fmt.Errorf("%s",str))
}