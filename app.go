package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zemirco/papertrail"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func redirect(ctx *gin.Context) {
	collection := ctx.MustGet("database").(*mgo.Collection)
	var result record
	if err := collection.Find(bson.M{"uuid": ctx.Param("uuid")}).One(&result); err != nil {
		defer logger.Printf("jre[dot]villas/%s - NOT FOUND", ctx.Param("uuid"))
		ctx.String(http.StatusNotFound, "404 not found")
		return
	}
	defer logger.Printf("jre[dot]villas/%s - %s", ctx.Param("uuid"), result.url)
	ctx.Redirect(http.StatusFound, result.url)
}
