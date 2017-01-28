package util

import (
	"fmt"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("util")

// Managed is an interface that is implemented by an object that needs to be
// managed via start/stop
type Managed interface {

	// Start starts the managed object
	Start() error

	// Stop stops the manged object
	Stop() error
}

// startManaged starts a "Managed" object
func startManaged(managed Managed) error {

	defer func() error {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			return err
		}

		return nil
	}()

	return managed.Start()
}

// stopManaged stops a "Managed" object
func stopManaged(managed Managed) error {

	defer func() error {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			return err
		}

		return nil
	}()

	return managed.Stop()
}

// StartManaged starts a Managed object, handles panics and logs details
func StartManaged(name string, managed Managed) error {

	log.Debugf("%s: Starting...", name)
	err := managed.Start()

	if err != nil {
		log.Errorf("%s: Error Starting", name)
		return err
	}

	log.Debugf("%s: Started", name)
	return nil
}

// StopManaged stops a Managed object, handles panics and logs details
func StopManaged(name string, managed Managed) error {

	log.Debugf("%s: Stopping...", name)

	err := stopManaged(managed)

	if err != nil {
		log.Errorf("Error stopping '%s': %s", name, err.Error())
		return err
	}

	log.Debugf("%s: Stopped", name)
	return nil
}
