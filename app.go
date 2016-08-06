package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zemirco/papertrail"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	collection = os.Getenv("COLLECTION")
	mongoDB    = os.Getenv("MONGO_DB")
	mongoURI   = os.Getenv("MONGO_URI")
	port, _    = strconv.Atoi(os.Getenv("PT_PORT"))
	remote     = papertrail.Writer{
		Port:    port,
		Network: papertrail.UDP,
		Server:  os.Getenv("PT_SERVER"),
	}
	logger = log.New(&remote, "", log.LstdFlags)
)

type record struct {
	id   bson.ObjectId `bson:"_id"`
	hits int
	url  string
	uuid string
}

func redirect(ctx *gin.Context) {
	collection := ctx.MustGet("database").(*mgo.Collection)
	var result record
	if err := collection.Find(bson.M{"uuid": ctx.Param("uuid")}).One(&result); err != nil {
		defer logger.Printf("jre[dot]villas/%s - 404 not found", ctx.Param("uuid"))
		ctx.String(http.StatusNotFound, "404 not found")
	} else {
		defer logger.Printf("jre[dot]villas/%s - %s", ctx.Param("uuid"), result.url)
		ctx.Redirect(http.StatusFound, result.url)
	}
}
