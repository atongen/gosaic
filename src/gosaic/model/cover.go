package model

import (
	"fmt"
	"time"
)

type Cover struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	AspectId int64  `db:"aspect_id"`
	Width    int    `db:"width"`
	Height   int    `db:"height"`
}

func CoverNameAspect(aspectId int64, width, height, num int) string {
	return fmt.Sprintf("type:aspect,aspectId:%d,width:%d,height:%d,num:%d",
		aspectId, width, height, num)
}

func CoverNameQuad(aspectId int64, width, height, num, maxDepth, minArea int) string {
	return fmt.Sprintf("type:quad,aspectId:%d,width:%d,height:%d,num:%d,maxDepth:%d,minArea:%d,ts:%d",
		aspectId, width, height, num, maxDepth, minArea, time.Now().Unix())
}
