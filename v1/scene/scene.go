// Tideland Go Application Support - Scene
//
// Copyright (C) 2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package scene

//--------------------
// IMPORTS
//--------------------

import (
	"time"

	"github.com/tideland/goas/v1/version"
	"github.com/tideland/goas/v2/loop"
	"github.com/tideland/goas/v3/errors"
)

//--------------------
// VERSION
//--------------------

// PackageVersion returns the version of the version package.
func PackageVersion() version.Version {
	return version.New(1, 0, 0, "alpha")
}

//--------------------
// SCENE
//--------------------

// CleanupFunc is a function for the cleanup of props after
// a scene ended.
type CleanupFunc func(key string, prop interface{}) error

// box contains a prop and a possible cleanup function.
type box struct {
	key     string
	prop    interface{}
	cleanup CleanupFunc
}

const (
	storeProp = iota
	fetchProp
	disposeProp
)

// envelope contains information transfered between client and scene.
type envelope struct {
	kind     int
	box      *box
	err      error
	respChan chan *envelope
}

// Scene is the access point to one scene. It has to be created once
// for a continuous flow of operations and then passed between all
// functions and goroutine which are actors of the scene.
type Scene interface {
	// Stop tells the scene to end and waits until it is done.
	Stop() error

	// Abort tells the scene to end due to the passed error.
	// Here only the first error will be stored for later evaluation.
	Abort(err error)

	// Wait blocks the caller until the scene ended and returns a
	// possible error or nil.
	Wait() error

	// Store stores a prop with a given key. The key must not exist.
	Store(key string, prop interface{}) error

	// StoreClean stores a prop with a given key and a cleanup
	// function called when a scene ends. The key must not exist.
	StoreClean(key string, prop interface{}, cleanup CleanupFunc) error

	// Fetch retrieves a prop.
	Fetch(key string) (interface{}, error)
}

// scene implements Scene.
type scene struct {
	props       map[string]*box
	inactivity  time.Duration
	timeout     time.Duration
	commandChan chan *envelope
	backend     loop.Loop
}

// StartScene creates and runs a new scene.
func StartScene() Scene {
	s := &scene{
		props:       make(map[string]*box),
		commandChan: make(chan *envelope, 1),
	}
	s.backend = loop.Go(s.backendLoop)
	return s
}

// Stop is specified on the Scene interface.
func (s *scene) Stop() error {
	return s.backend.Stop()
}

// Abort is specified on the Scene interface.
func (s *scene) Abort(err error) {
	s.backend.Kill(err)
}

// Wait is specified on the Scene interface.
func (s *scene) Wait() error {
	return s.backend.Wait()
}

// Store is specified on the Scene interface.
func (s *scene) Store(key string, prop interface{}) error {
	return s.StoreClean(key, prop, nil)
}

// StoreClean is specified on the Scene interface.
func (s *scene) StoreClean(key string, prop interface{}, cleanup CleanupFunc) error {
	cmd := &envelope{
		kind: storeProp,
		box: &box{
			key:     key,
			prop:    prop,
			cleanup: cleanup,
		},
		respChan: make(chan *envelope, 1),
	}
	_, err := s.command(cmd)
	return err
}

// Fetch is specified on the Scene interface.
func (s *scene) Fetch(key string) (interface{}, error) {
	cmd := &envelope{
		kind: fetchProp,
		box: &box{
			key: key,
		},
		respChan: make(chan *envelope, 1),
	}
	resp, err := s.command(cmd)
	if err != nil {
		return nil, err
	}
	return resp.box.prop, nil
}

// backendLoop runs the backend loop of the scene.
func (s *scene) backendLoop(l loop.Loop) error {
	// Init timers.
	var watchdog <-chan time.Time
	var clapperboard <-chan time.Time
	if s.timeout != 0 {
		clapperboard = time.After(s.timeout)
	}
	// Run loop.
	for {
		if s.inactivity != 0 {
			watchdog = time.After(s.inactivity)
		}
		select {
		case <-l.ShallStop():
			return nil
		case to := <-watchdog:
			return errors.New(ErrTimeout, errorMessages, "inactivity", to)
		case to := <-clapperboard:
			return errors.New(ErrTimeout, errorMessages, "absolute", to)
		case cmd := <-s.commandChan:
			s.processCommand(cmd)
		}
	}
}

// command sends a command envelope to the backend and
// waits for the response.
func (s *scene) command(cmd *envelope) (*envelope, error) {
	select {
	case s.commandChan <- cmd:
	case <-s.backend.IsStopping():
		return nil, errors.New(ErrSceneEnded, errorMessages)
	}
	select {
	case <-s.backend.IsStopping():
		return nil, s.Wait()
	case resp := <-cmd.respChan:
		if resp.err != nil {
			return nil, resp.err
		}
		return resp, nil
	}
}

// processCommand processes the sent commands.
func (s *scene) processCommand(cmd *envelope) {
	switch cmd.kind {
	case storeProp:
		// Add a new prop.
		_, ok := s.props[cmd.box.key]
		if ok {
			cmd.err = errors.New(ErrPropAlreadyExist, errorMessages, cmd.box.key)
		} else {
			s.props[cmd.box.key] = cmd.box
		}
	case fetchProp:
		// Retrieve a prop.
		box, ok := s.props[cmd.box.key]
		if !ok {
			cmd.err = errors.New(ErrPropNotFound, errorMessages, cmd.box.key)
		} else {
			cmd.box = box
		}
	case disposeProp:
		// Remove a prop.
		box, ok := s.props[cmd.box.key]
		if ok {
			delete(s.props, cmd.box.key)
			if box.cleanup != nil {
				cmd.err = box.cleanup(cmd.box.key, cmd.box.prop)
			}
		}
	default:
		panic("illegal command")
	}
	// Return the changed command as response.
	cmd.respChan <- cmd
}

// EOF