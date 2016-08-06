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
	mongoCollection = os.Getenv("MONGO_COLLECTION")
	mongoDB         = os.Getenv("MONGO_DB")
	mongoURI        = os.Getenv("MONGO_URI")
	port, _         = strconv.Atoi(os.Getenv("PT_PORT"))
	remote          = papertrail.Writer{
		Port:    port,
		Network: papertrail.UDP,
		Server:  os.Getenv("PT_SERVER"),
	}
	logger = log.New(&remote, "", log.LstdFlags)
)

type record struct {
	ID   bson.ObjectId `bson:"_id"`
	Hits int           `json:"hits"`
	URL  string        `json:"url"`
	UUID string        `json:"uuid"`
}

func main() {
	middleware := dbMiddleware()
	router := gin.Default()
	router.Use(middleware)
	router.GET("/:uuid", redirect)
	router.Run(":8080")
}

func dbMiddleware() gin.HandlerFunc {
	if session, err := mgo.Dial(mongoURI); err != nil {
		panic(err)
	} else {
		return func(ctx *gin.Context) {
			copy := session.Copy()
			defer copy.Close()
			ctx.Set("collection", copy.DB(mongoDB).C(mongoCollection))
			ctx.Next()
		}
	}
}

func redirect(ctx *gin.Context) {
	collection := ctx.MustGet("collection").(*mgo.Collection)
	var result record
	if err := collection.Find(bson.M{"uuid": ctx.Param("uuid")}).One(&result); err != nil {
		defer logger.Printf("jre[dot]villas/%s - 404 not found\n", ctx.Param("uuid"))
		ctx.String(http.StatusNotFound, "404 not found")
	} else {
		defer logger.Printf("jre[dot]villas/%s - %s\n", ctx.Param("uuid"), result.URL)
		ctx.Redirect(http.StatusFound, result.URL)
	}
}
