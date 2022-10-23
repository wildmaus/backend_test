package utils

import (
	"log"
	"time"
)

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		log.Printf("Try connect to postgressql, tries left: %v\n", attemtps-1)
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}
