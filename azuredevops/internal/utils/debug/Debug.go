package debug

import (
	"log"
	"os"
	"strconv"
	"time"
)

var debugWaitPassed = false

const defaultWait = 20

// Wait wait to attach a debugger if the environment variable AZDO_PROVIDER_DEBUG is set to 1
// In addtion the function supports reading the wait time from the AZDO_PROVIDER_DEBUG_WAIT environment variable
// The default waiting time is 20 seconds
func Wait(force ...bool) {
	if "1" == os.Getenv("AZDO_PROVIDER_DEBUG") {
		bForce := force != nil && force[0]
		if !debugWaitPassed || bForce {
			wait := defaultWait
			if v, ok := os.LookupEnv("AZDO_PROVIDER_DEBUG_WAIT"); ok {
				if i, err := strconv.Atoi(v); err == nil {
					wait = i
				} else {
					log.Printf("[INFO] Failed to convert value %s of environment 'AZDO_PROVIDER_DEBUG_WAIT' variable to Int32", v)
				}
			}
			time.Sleep(time.Duration(wait) * time.Second)
			debugWaitPassed = true
		}
	}
}
