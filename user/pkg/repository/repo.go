package repository

import "github.com/social-media-timeline/user/pkg/mongoClient"


func get() {

	mongoClient.InitMongoClient("mongodb://localhost:27017")
}
