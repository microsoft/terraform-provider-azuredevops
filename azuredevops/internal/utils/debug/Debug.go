package debug

import (
	"os"
	"time"
)

var debugWaitPassed = false

// Wait wait to attach a debugger
func Wait(force ...bool) {
	bForce := force != nil && force[0]
	if (!debugWaitPassed || bForce) && "1" == os.Getenv("AZDO_PROVIDER_DEBUG") {
		time.Sleep(20 * time.Second)
		debugWaitPassed = true
	}
}
