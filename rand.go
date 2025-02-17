package pkgep

import (
	"fmt"
	"math/rand"
	"time"
)

func RandStr() string {
	return fmt.Sprintf("%v%v", time.Now().UnixNano(), rand.Intn(90000)+10000)
}
