package main

import (
	"challenge/controller"
	"challenge/model"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {

	r := httprouter.New()                            // setting a new mux
	uc := controller.NewUserController(getsession()) // generatin new user
	r.GET("/", uc.Index)
	r.GET("/public_key", uc.GetPublic)
	r.PUT("/transaction", uc.Transaction)
	r.POST("/signature", uc.Signature)

	http.ListenAndServe(":8080", r)
}

/*Creating a session*/
func getsession() (map[string]string, map[int64]model.Trans) {
	k := make(map[string]string)
	t := make(map[int64]model.Trans)

	return k, t
}
