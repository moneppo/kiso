package tap

import (
	"reflect"
	"github.com/kiso/math"
)

const tapTolerance = 2

type TapEvent struct {
	InputEvent
	Position math.Vec2
}

type TapRecognizer struct {
	down bool
	e PointerEvent
}

var Tap *EventManager

func (t *TapRecognizer) OnPossible(e *InputEvent) RecognizerState {
	if p, ok := (interface{})(e).(PointerEvent); ok && p.Type == InputTypeDown {
		t.e = p
		t.down = true
		return RecognizerStateActive
	} else {
		t.down = false
		return RecognizerStateFail
	}
}

func (t *TapRecognizer) OnChange(e *InputEvent) RecognizerState {
	if p, ok := (interface{})(e).(PointerEvent); ok {
		switch p.Type {
		case InputTypeUp:
			result := TapEvent{ InputEvent{e.Target, InputTypeDown, e.modifier}, t.e.Position }
			Tap.Fire(result)
			t.down = false
			return RecognizerStateEnd
		case InputTypeMove:
			if math.Distance(p.Position, t.e.Position) > tapTolerance {
				t.down = false
				return RecognizerStateFail
			}
		}
	}
	return RecognizerStateFail
}

func init() {
	Tap = NewEventManager(reflect.TypeOf(TapEvent{}))
}