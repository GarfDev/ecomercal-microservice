package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gopkg.in/square/go-jose.v2"
)

type Session struct {
	IP   string `json:"id"`   // real_ip
	UUID string `json:"uuid"` // user_id
	CID  string `json:"cid"`  // card_id
}

var validate *validator.Validate

var rcpt = jose.Recipient{
	Algorithm:  jose.PBES2_HS256_A128KW,
	Key:        "mypassphrase",
	PBES2Count: 4096,
	PBES2Salt:  []byte{},
}

func e(c *gin.Context) {
	session_obj := Session{}

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	json.Unmarshal(jsonData, &session_obj)

	err = validate.Struct(session_obj)

	if err != nil {
		fmt.Printf("%+v\n", err)
		fmt.Printf("%+v\n", session_obj)
		c.String(400, err.Error())
		return
	}

	enc, err := jose.NewEncrypter(jose.A128CBC_HS256, rcpt, nil)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	jewPlaintextToken, err := enc.Encrypt(jsonData)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	serialized, err := jewPlaintextToken.CompactSerialize()
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Header("sid", serialized)
	c.String(200, string(jsonData))
}

func v(c *gin.Context) {
	session_obj := Session{}
	token := c.GetHeader("sid")

	if token == "" {
		return
	}

	jwe, err := jose.ParseEncrypted(token)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	decrypted, err := jwe.Decrypt("mypassphrase")
	if err != nil {
		c.String(400, err.Error())
		return
	}

	json.Unmarshal(decrypted, &session_obj)
	err = validate.Struct(session_obj)
	v := reflect.ValueOf(session_obj)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		c.Header(typeOfS.Field(i).Name, v.Field(i).String())
	}

}

func main() {

	r := gin.Default()
	validate = validator.New()

	r.POST("/encrypt", e)
	r.GET("/validate", v)

	r.Run("0.0.0.0:4000")
}
