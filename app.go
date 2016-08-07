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
	ptPort          = os.Getenv("PT_PORT")
	ptServer        = os.Getenv("PT_SERVER")
	runAddress      = os.Getenv("RUN_ADDRESS")
)

type record struct {
	ID   bson.ObjectId `bson:"_id"`
	Hits int           `json:"hits"`
	URL  string        `json:"url"`
	UUID string        `json:"uuid"`
}

func atoi(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}

func main() {
	remote := papertrail.Writer{
		Port:    atoi(ptPort),
		Network: papertrail.UDP,
		Server:  ptServer,
	}
	log.SetOutput(&remote)
	middleware := dbMiddleware()
	router := gin.Default()
	router.Use(middleware)
	router.GET("/:uuid", redirect)
	router.Run(runAddress)
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
		ctx.String(http.StatusNotFound, "404 not found")
		log.Printf("jre[dot]villas/%s - 404 not found\n", ctx.Param("uuid"))
		return
	}
	ctx.Redirect(http.StatusFound, result.URL)
	log.Printf("jre[dot]villas/%s - %s\n", ctx.Param("uuid"), result.URL)
}
