package rareskill

import (
	"log"
	"time"
)

type Jobs struct{}

func (e Jobs) Run() {
	Takanome()
	log.Println(time.Now(), "鷹の目")
	Register()
	log.Println(time.Now(), "レジスタ")
}
