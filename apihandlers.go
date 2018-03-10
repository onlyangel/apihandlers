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

		defer func(r *http.Request) {
			ctx := appengine.NewContext(r)
			if err := recover(); err != nil {
				fmt.Fprintf(w,"ERROR: %v", err)
				ctx.Infof("ERROR: %v", err.(*errors.Error).ErrorStack())
			}
		}(r)

		fn(w,r)
	}
}

