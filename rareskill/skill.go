package rareskill

import (
	"log"
	"time"
)

type Skills struct{}

func (e Skills) Run() {
	Takanome()
	log.Println(time.Now(), "鷹の目")
	Register()
	log.Println(time.Now(), "レジスタ")
}
