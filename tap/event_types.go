package tap

import (
	"github.com/kiso/math"
)

type PointerType int

const (
	PointerMouseLeft PointerType = iota
	PointerMouseMiddle
	PointerMouseRight
	PointerFinger
	PointerStylusTip
	PointerStylusEraser
)

type InputModifier int
const (
	InputModifierNone InputModifier = 0
	InputModifierCommand = 1
	InputModifierOption = 2
	InputModifierControl = 4
	InputModifierShift = 8
)

type InputEventType int
const (
	InputTypeDown InputEventType = iota
	InputTypeUp 
	InputTypeMove
	InputTypeCancel
)

type InputEvent struct {
	Target interface{}
	Type InputEventType
	modifier InputModifier
}

func (e *InputEvent) Command() bool {
	return e.modifier & InputModifierCommand != 0
}

func (e *InputEvent) Option() bool {
	return e.modifier & InputModifierOption != 0
}

func (e *InputEvent) Control() bool {
	return e.modifier & InputModifierControl != 0
}

func (e *InputEvent) Shift() bool {
	return e.modifier & InputModifierShift != 0
}

type PointerEvent struct {
	InputEvent
	
	Position math.Vec2
	Index int
	PointerType PointerType
}

type KeyEvent struct {
	InputEvent
	Key rune
}

type TransformEvent struct {
	InputEvent
	PositionDelta math.Vec2
	RotationDelta float32
	ScaleDelta float32
}