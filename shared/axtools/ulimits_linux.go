package axtools

import (
	"github.com/rs/zerolog/log"
	"syscall"
)

func SetupForHighLoad() {

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Warn().Err(err).Msg("Error Getting Rlimit")
	} else {
		log.Info().Interface("RLimit", rLimit).Msg("Rlimit before setup")
		rLimit.Max = 256000
		rLimit.Cur = 256000

		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			log.Warn().Err(err).Msg("Error Setting Rlimit")
		} else {
			err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
			if err != nil {
				log.Warn().Err(err).Msg("Error Setting Rlimit after change")
			} else {
				log.Info().Interface("RLimit", rLimit).Msg("Rlimit after setup")
			}
		}
	}
}
