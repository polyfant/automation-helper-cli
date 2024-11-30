package abb

// ABBCommand represents a specific ABB robot command with its syntax and example
type ABBCommand struct {
	Name        string
	Syntax      string
	Example     string
	Description string
}

// Common ABB RAPID commands
var Commands = map[string]ABBCommand{
	"move_j": {
		Name:        "MoveJ",
		Syntax:      "MoveJ Target [Speed] [Zone] [Tool]",
		Description: "Joint movement - moves robot to position using axis movement",
	},
	"move_l": {
		Name:        "MoveL",
		Syntax:      "MoveL Target [Speed] [Zone] [Tool]",
		Description: "Linear movement - moves robot in straight line to position",
	},
	"set_do": {
		Name:        "SetDO",
		Syntax:      "SetDO Signal Value",
		Description: "Sets digital output signal",
	},
	"wait_di": {
		Name:        "WaitDI",
		Syntax:      "WaitDI Signal Value [\\MaxTime]",
		Description: "Waits for digital input signal to reach specified value",
	},
	"if_statement": {
		Name:        "IF",
		Syntax:      "IF condition THEN ... ENDIF",
		Description: "Conditional execution of instructions",
	},
}
