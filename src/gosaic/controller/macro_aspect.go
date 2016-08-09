package controller

import (
	"fmt"
	"gosaic/environment"
	"time"
)

func MacroAspect(env environment.Environment, path string, coverWidth, coverHeight, partialWidth, partialHeight, num int) {
	ts := time.Now().Format(time.RubyDate)
	name := fmt.Sprintf("%s-%s", path, ts)

	CoverAspect(env, name, coverWidth, coverHeight, partialWidth, partialHeight, num)
	Macro(env, path, name)
}
