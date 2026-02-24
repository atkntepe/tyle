package layout

type Direction string

const (
	Right    Direction = "right"
	Left     Direction = "left"
	Down     Direction = "down"
	Up       Direction = "up"
	Previous Direction = "previous"
	Next     Direction = "next"
)

type StepAction string

const (
	ActionSplit    StepAction = "split"
	ActionFocus    StepAction = "focus"
	ActionEqualize StepAction = "equalize"
	ActionDelay    StepAction = "delay"
)

type LayoutStep struct {
	Action    StepAction
	Direction Direction
	DelayMs   int
}

type Layout struct {
	ID          string
	Name        string
	Description string
	Preview     []string
	Steps       []LayoutStep
	PaneCount   int
	FinalFocus  Direction
}
