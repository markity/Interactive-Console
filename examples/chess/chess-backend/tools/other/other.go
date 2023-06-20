package other

import (
	"math/rand"
	"time"
)

func RandGetBool() bool {
	return rand.New(rand.NewSource(time.Now().Unix())).Int()%2 == 0
}

func Ignore(i interface{}) {

}
