package tap

import (
	"reflect"
	"github.com/kiso/math"
)

const startThreshold = 2

type PanEvent struct {
	InputEvent
	Position math.Vec2
	Delta math.Vec2
}

type PanRecognizer struct {
	down bool
	start PointerEvent
	last math.Vec2
}

var Pan *EventManager

func (r *PanRecognizer) OnPossible(e *InputEvent) RecognizerState {
	if p, ok := (interface{})(e).(PointerEvent); ok {
		if !r.down && p.Type == InputTypeDown {
			r.start = p
			r.last = p.Position
			r.down = true
			return RecognizerStatePossible
		} else if r.down && p.Type == InputTypeMove {
			if math.Distance(p.Position, r.start.Position) > startThreshold {
				r.last = p.Position
				result := PanEvent{ InputEvent{e.Target, InputTypeDown, e.modifier},
														p.Position, p.Position.Subtract(r.start.Position) }
				Pan.Fire(result)
				return RecognizerStateActive
			} else {
				return RecognizerStatePossible
			}
		} else {
		  r.down = false
		  return RecognizerStateFail
		}
	}
	return RecognizerStateFail
}

func (r *PanRecognizer) OnChange(e *InputEvent) RecognizerState {
	if p, ok := (interface{})(e).(PointerEvent); ok {
		switch p.Type {
		case InputTypeUp:
			result := PanEvent{ InputEvent{e.Target, InputTypeDown, e.modifier},
													p.Position, p.Position.Subtract(r.last) }
			Pan.Fire(result)
			r.down = false
			return RecognizerStateEnd
		case InputTypeMove:
			result := PanEvent{ InputEvent{e.Target, InputTypeDown, e.modifier},
													p.Position, p.Position.Subtract(r.last) }
			r.last = p.Position
			Pan.Fire(result)
			return RecognizerStateActive
		}
	}
	return RecognizerStateFail
}

func init() {
	Pan = NewEventManager(reflect.TypeOf(PanEvent{}))
}