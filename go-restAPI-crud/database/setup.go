package databases

import (
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	RedisDB *redis.Client
)

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/valorant_leaderboard"))
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
