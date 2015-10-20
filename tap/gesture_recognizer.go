package tap

type RecognizerState int
const (
	RecognizerStatePossible RecognizerState = iota
	RecognizerStateActive
	RecognizerStateEnd
	RecognizerStateFail
)

type GestureRecognizer interface {
	OnPossible(InputEvent) RecognizerState
	OnChange(InputEvent) RecognizerState
}