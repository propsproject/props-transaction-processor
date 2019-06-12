package test

import (
	"github.com/propsproject/pending-props/core/state"
	"strconv"

	//"github.com/propsproject/pending-props/core/state"
	"log"
	"testing"
	"time"
)

func TestActivityTimestamp(t *testing.T) {
	blockTimestamp, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 32)
	logDate := 20190523
	log.Print(logDate)
	log.Print("test")
	log.Print(state.ValidateActivityTimestamp(logDate, blockTimestamp))
}