package tap

import (
	"reflect"
	)

var recognizers []GestureRecognizer
var possibleRecognizers []GestureRecognizer
var activeRecognizers []GestureRecognizer

var PointerDown *EventManager
var PointerUp *EventManager
var PointerMove *EventManager
var KeyDown *EventManager
var KeyUp *EventManager

func AddRecognizer(recognizer GestureRecognizer) {
	recognizers = append(recognizers, recognizer)
}

func applyEvent(e InputEvent) {

	// Some recognizers are detecting gestures and should be 
	// fed events
	if len(activeRecognizers) > 0 {
		newActives := make([]GestureRecognizer, 0)
		for _, r := range activeRecognizers {
			switch r.OnChange(e) {
			case RecognizerStatePossible:
				// Error case, but we'll just warn and fix the issue
				print("WARN: Active recognizer returned possible. Resetting...")
				possibleRecognizers = append(possibleRecognizers, r)
			case RecognizerStateActive:
				newActives = append(newActives, r)
			}
		}
		activeRecognizers = newActives

	// Otherwise feed all events to the recognizers and update
	// the active list accordingly
	} else {
		for _, r := range recognizers {
			switch r.OnPossible(e) {
				case RecognizerStateActive:
					activeRecognizers = append(activeRecognizers, r)
			}
		}
	}
}

func init() {
	PointerDown = NewEventManager(reflect.TypeOf(PointerEvent{}))
	PointerUp = NewEventManager(reflect.TypeOf(PointerEvent{}))
	PointerMove = NewEventManager(reflect.TypeOf(PointerEvent{}))
	KeyDown = NewEventManager(reflect.TypeOf(KeyEvent{}))
	KeyUp = NewEventManager(reflect.TypeOf(KeyEvent{}))

	PointerDown.On(func(e PointerEvent) {
		applyEvent((interface{})(e).(InputEvent))
	})	

	PointerUp.On(func(e PointerEvent) {
		applyEvent((interface{})(e).(InputEvent))
	})

	PointerMove.On(func(e PointerEvent) {
		applyEvent((interface{})(e).(InputEvent))
	})

	KeyDown.On(func(e KeyEvent) {
		applyEvent((interface{})(e).(InputEvent))
	})

	KeyUp.On(func(e KeyEvent) {
		applyEvent((interface{})(e).(InputEvent))
	})
}