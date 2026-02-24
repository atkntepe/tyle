package layout

func Presets() []Layout {
	return []Layout{
		twoColumns(),
		twoRows(),
		threeColumns(),
		mainRightStack(),
		leftStackMain(),
		mainSideStack(),
		grid2x2(),
		mainTopTwoBottom(),
		twoTopOneBottom(),
		threeTopOneBottom(),
	}
}

func twoColumns() Layout {
	return Layout{
		ID:          "two-columns",
		Name:        "Two Columns",
		Description: "Two equal vertical panes side by side",
		Preview: []string{
			"┌─────┬─────┐",
			"│  A  │  B  │",
			"│     │     │",
			"└─────┴─────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Right},
			{Action: ActionFocus, Direction: Left},
			{Action: ActionEqualize},
		},
		PaneCount:  2,
		FinalFocus: Left,
	}
}

func twoRows() Layout {
	return Layout{
		ID:          "two-rows",
		Name:        "Two Rows",
		Description: "Two equal horizontal panes stacked",
		Preview: []string{
			"┌───────────┐",
			"│     A     │",
			"├───────────┤",
			"│     B     │",
			"└───────────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Down},
			{Action: ActionFocus, Direction: Up},
			{Action: ActionEqualize},
		},
		PaneCount:  2,
		FinalFocus: Up,
	}
}

func threeColumns() Layout {
	return Layout{
		ID:          "three-columns",
		Name:        "Three Columns",
		Description: "Three equal vertical panes in a row",
		Preview: []string{
			"┌───┬───┬───┐",
			"│ A │ B │ C │",
			"│   │   │   │",
			"└───┴───┴───┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Right},
			{Action: ActionSplit, Direction: Right},
			{Action: ActionFocus, Direction: Previous},
			{Action: ActionFocus, Direction: Previous},
			{Action: ActionEqualize},
		},
		PaneCount:  3,
		FinalFocus: Previous,
	}
}

func mainRightStack() Layout {
	return Layout{
		ID:          "main-right-stack",
		Name:        "Main + Right Stack",
		Description: "Large main pane with two stacked panes on the right",
		Preview: []string{
			"┌──────┬──────┐",
			"│      │  B   │",
			"│  A   ├──────┤",
			"│      │  C   │",
			"└──────┴──────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Right},
			{Action: ActionSplit, Direction: Down},
			{Action: ActionFocus, Direction: Left},
			{Action: ActionEqualize},
		},
		PaneCount:  3,
		FinalFocus: Left,
	}
}

func leftStackMain() Layout {
	return Layout{
		ID:          "left-stack-main",
		Name:        "Left Stack + Main",
		Description: "Two stacked panes on the left with a large main pane",
		Preview: []string{
			"┌──────┬──────┐",
			"│  A   │      │",
			"├──────┤  B   │",
			"│  C   │      │",
			"└──────┴──────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Right},
			{Action: ActionFocus, Direction: Left},
			{Action: ActionSplit, Direction: Down},
			{Action: ActionFocus, Direction: Right},
			{Action: ActionEqualize},
		},
		PaneCount:  3,
		FinalFocus: Right,
	}
}

func mainSideStack() Layout {
	return Layout{
		ID:          "main-side-stack",
		Name:        "Main + Side Stack",
		Description: "Wide main pane with a narrow side stack",
		Preview: []string{
			"┌──────┬─────┐",
			"│      │  B  │",
			"│  A   ├─────┤",
			"│      │  C  │",
			"└──────┴─────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Right},
			{Action: ActionSplit, Direction: Down},
			{Action: ActionFocus, Direction: Left},
			{Action: ActionEqualize},
		},
		PaneCount:  3,
		FinalFocus: Left,
	}
}

func grid2x2() Layout {
	return Layout{
		ID:          "grid-2x2",
		Name:        "Grid 2x2",
		Description: "Four equal panes in a 2x2 grid",
		Preview: []string{
			"┌─────┬─────┐",
			"│  A  │  B  │",
			"├─────┼─────┤",
			"│  C  │  D  │",
			"└─────┴─────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Right},
			{Action: ActionSplit, Direction: Down},
			{Action: ActionFocus, Direction: Left},
			{Action: ActionSplit, Direction: Down},
			{Action: ActionFocus, Direction: Up},
			{Action: ActionEqualize},
		},
		PaneCount:  4,
		FinalFocus: Up,
	}
}

func mainTopTwoBottom() Layout {
	return Layout{
		ID:          "main-top-two-bottom",
		Name:        "Main Top + Two Bottom",
		Description: "Wide main pane on top with two panes below",
		Preview: []string{
			"┌───────────┐",
			"│     A     │",
			"├─────┬─────┤",
			"│  B  │  C  │",
			"└─────┴─────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Down},
			{Action: ActionSplit, Direction: Right},
			{Action: ActionFocus, Direction: Up},
			{Action: ActionEqualize},
		},
		PaneCount:  3,
		FinalFocus: Up,
	}
}

func twoTopOneBottom() Layout {
	return Layout{
		ID:          "two-top-one-bottom",
		Name:        "Two Top + One Bottom",
		Description: "Two panes on top with a wide pane on the bottom",
		Preview: []string{
			"┌─────┬─────┐",
			"│  A  │  B  │",
			"├─────┴─────┤",
			"│     C     │",
			"└───────────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Down},
			{Action: ActionFocus, Direction: Up},
			{Action: ActionSplit, Direction: Right},
			{Action: ActionFocus, Direction: Left},
			{Action: ActionEqualize},
		},
		PaneCount:  3,
		FinalFocus: Left,
	}
}

func threeTopOneBottom() Layout {
	return Layout{
		ID:          "three-top-one-bottom",
		Name:        "Three Top + One Bottom",
		Description: "Three panes on top with a wide pane on the bottom",
		Preview: []string{
			"┌───┬───┬───┐",
			"│ A │ B │ C │",
			"├───┴───┴───┤",
			"│     D     │",
			"└───────────┘",
		},
		Steps: []LayoutStep{
			{Action: ActionSplit, Direction: Down},
			{Action: ActionFocus, Direction: Up},
			{Action: ActionSplit, Direction: Right},
			{Action: ActionSplit, Direction: Right},
			{Action: ActionFocus, Direction: Left},
			{Action: ActionFocus, Direction: Left},
			{Action: ActionEqualize},
		},
		PaneCount:  4,
		FinalFocus: Left,
	}
}
