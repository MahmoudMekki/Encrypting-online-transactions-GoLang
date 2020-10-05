package controller

import (
	"challenge/model"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

/*Creating the database for the public and private keys */
/*Creating the database for every transaction */

type UserController struct {
	keysDB  map[string]string
	transDB map[int64]model.Trans
}

/*making public key global sothat we can retrieve the private key later*/
var pub string

/*creating new user session with both databases*/
func NewUserController(k map[string]string, t map[int64]model.Trans) *UserController {
	return &UserController{k, t}
}

/*starting page with handler interface */
func (uc UserController) Index(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	s := `
	<!DOCTYPE html>
	<html>
	<head>
	<meta charset = "UTF-8">
	<title>MEkki</title>
	</head>
	<body>
	<h1><a href = "/public_key">GOTO locahost:8080/publickey</a></h1>
	</body>
	</html>
`
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(res, s)
}

/*getting the keys */
func (uc UserController) GetPublic(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	pub, priv := generate() // local made function that returns both private and public keys
	if _, ok := uc.keysDB[pub]; !ok {
		uc.keysDB[pub] = priv
	}
	var key = model.PublicDaemon{
		pub,
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(key) // generating our json and send to the server

}

/*Getting the private and public keys*/
func generate() (public, private string) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	public = string(pub)
	private = string(priv)

	return public, private
}

/*handling the transactions */
func (uc UserController) Transaction(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	t := model.Trans{}     // setting structure of the transaction
	tid := model.TransID{} // Id structure
	json.NewDecoder(req.Body).Decode(&t)

	id := random() // generating random id

	tid = model.TransID{id}

	uc.transDB[id] = t // stroing at our transaction database
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(tid)

}

func random() int64 {
	x, _ := rand.Int(rand.Reader, big.NewInt(100))
	return x.Int64()
}

/*Handling the signature method post*/
func (uc UserController) Signature(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ids := model.TransIDs{}
	trans := []model.Trans{}
	signature := model.Signature{}
	json.NewDecoder(req.Body).Decode(&ids)
	for i := range ids.IDs {
		if v, ok := uc.transDB[ids.IDs[i]]; ok {
			trans[i] = v
		}
	}
	bs := make([]byte, 0, len(trans))
	for i := range trans {
		bs = append(bs, []byte(trans[i].Txn)...)
	}

	priv := ed25519.PrivateKey(uc.keysDB[pub])

	sign := ed25519.Sign(priv, bs)
	signature = model.Signature{trans, string(sign)}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(os.Stdout).Encode(signature)

	json.NewEncoder(res).Encode(signature)

}
