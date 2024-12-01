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
		Name:    "MoveJ",
		Syntax:  "MoveJ Target [Speed] [Zone] [Tool]",
		Example: "MoveJ p10, v1000, z50, tool1;\nMoveJ [[100,200,300],[1,0,0,0],[0,0,0,0],[9E9,9E9,9E9,9E9,9E9,9E9]], v1000, fine, tool0;",
		Description: "Joint movement - Fastest way to move between points. Robot axes move independently to reach target.\n" +
			"- Speed (v): v1000 means 1000mm/s\n" +
			"- Zone (z): z50 means blend radius 50mm, 'fine' for exact positioning\n" +
			"- Common use: Home position movements, approach positions",
	},
	"move_l": {
		Name:    "MoveL",
		Syntax:  "MoveL Target [Speed] [Zone] [Tool]",
		Example: "MoveL p20, v100, fine, tool1;\nMoveL Offs(p20,100,0,50), v100, z10, tool1; ! Offset from p20",
		Description: "Linear movement - TCP moves in straight line. Essential for precise paths and process work.\n" +
			"- Use lower speeds (v100-v300) for process work\n" +
			"- Offs() function adds offset to target\n" +
			"- Common use: Welding, gluing, precise positioning",
	},
	"move_c": {
		Name:    "MoveC",
		Syntax:  "MoveC CirPoint ToPoint [Speed] [Zone] [Tool]",
		Example: "MoveC p30, p40, v500, z10, tool1;\n! Create circle: Same distance from start to circle point as circle point to end",
		Description: "Circular movement - Creates perfect circular arc through three points.\n" +
			"- Start point (current pos) → CirPoint → ToPoint\n" +
			"- Points should form isosceles triangle for smooth motion\n" +
			"- Common use: Arc welding, curved sealing paths",
	},
	"set_do": {
		Name:    "SetDO",
		Syntax:  "SetDO Signal Value",
		Example: "SetDO do_Gripper, 1;    ! Turn on gripper\nSetDO do_Valve, 0;      ! Turn off valve",
		Description: "Sets digital output signal. Controls binary equipment like grippers, valves.\n" +
			"- Value: 1/0 (on/off)\n" +
			"- Signals must be defined in I/O configuration\n" +
			"- Common use: Gripper control, process equipment",
	},
	"wait_di": {
		Name:    "WaitDI",
		Syntax:  "WaitDI Signal Value [\\MaxTime]",
		Example: "WaitDI di_PartPresent, 1;         ! Wait for part\nWaitDI di_Ready, 1 \\MaxTime:=5;  ! Wait max 5s",
		Description: "Waits for digital input signal. Essential for synchronizing with external events.\n" +
			"- Optional \\MaxTime prevents infinite waiting\n" +
			"- Throws error if MaxTime exceeded\n" +
			"- Common use: Part detection, process synchronization",
	},
	"if_statement": {
		Name:    "IF",
		Syntax:  "IF condition THEN ... {ELSEIF condition THEN ...} [ELSE ...] ENDIF",
		Example: "IF di_PartType = 1 THEN\n    MoveJ p10, v1000, z50, tool1;\nELSEIF di_PartType = 2 THEN\n    MoveJ p20, v1000, z50, tool1;\nENDIF",
		Description: "Conditional execution. Supports complex program logic.\n" +
			"- Conditions: =, <>, >, <, >=, <=, AND, OR, NOT\n" +
			"- Can be nested\n" +
			"- Common use: Part type selection, error handling",
	},
	"search_l": {
		Name:    "SearchL",
		Syntax:  "SearchL [\\Signal] Target [Speed] [Tool]",
		Example: "SearchL \\Stop:=di_Contact, p30, v100, tool1;\npos := CRobT();  ! Get position where contact was made",
		Description: "Linear search movement with sensor feedback.\n" +
			"- Stops immediately when signal changes\n" +
			"- Use CRobT() to get stop position\n" +
			"- Common use: Part location, calibration",
	},
	"tool_data": {
		Name:    "PERS tooldata",
		Syntax:  "PERS tooldata <name>:=[TRUE,[[x,y,z],[q1,q2,q3,q4]],[mass,[cx,cy,cz],[I1,I2,I3,I4],0,0,0]];",
		Example: "PERS tooldata myTool:=[TRUE,[[175,0,35],[1,0,0,0]],[0.5,[0,0,0.02],[1,0,0,0],0,0,0]];\n! Tool at x=175mm, z=35mm, 0.5kg",
		Description: "Define tool data - Critical for accurate positioning.\n" +
			"- Position: [x,y,z] from tool mounting point to TCP\n" +
			"- Orientation: [q1,q2,q3,q4] quaternion values\n" +
			"- Mass and center of gravity important for dynamics",
	},
	"wobj_data": {
		Name:    "PERS wobjdata",
		Syntax:  "PERS wobjdata <name>:=[FALSE,TRUE,\"\",[uframe],[oframe]];",
		Example: "PERS wobjdata myTable:=[FALSE,TRUE,\"\"[[800,0,500],[1,0,0,0]],[[0,0,0],[1,0,0,0]]];\n! Table 800mm in X, 500mm in Z",
		Description: "Define work object - Local coordinate system for parts.\n" +
			"- uframe: User frame relative to world\n" +
			"- oframe: Object frame relative to uframe\n" +
			"- Common use: Multiple identical fixtures, moving lines",
	},
	"for_loop": {
		Name:    "FOR",
		Syntax:  "FOR <var> FROM <start> TO <end> [STEP <step>] DO ... ENDFOR",
		Example: "FOR i FROM 1 TO 5 DO\n    MoveL Offs(p40,0,i*50,0), v500, z10, tool1;\nENDFOR\n! Creates 5 points 50mm apart",
		Description: "Counter-based loop. Perfect for repeated patterns.\n" +
			"- Variable automatically increments\n" +
			"- STEP defines increment value\n" +
			"- Common use: Pallet picking, pattern movements",
	},
	"while_loop": {
		Name:    "WHILE",
		Syntax:  "WHILE condition DO ... ENDWHILE",
		Example: "WHILE di_PartsPresent = 1 DO\n    MoveL pPickPos, v500, fine, tool1;\n    SetDO do_Gripper, 1;\nENDWHILE",
		Description: "Condition-based loop. Continues while condition is true.\n" +
			"- Check condition before each iteration\n" +
			"- Use caution to avoid infinite loops\n" +
			"- Common use: Continuous processes, conveyor tracking",
	},
	"set_ao": {
		Name:    "SetAO",
		Syntax:  "SetAO Signal Value",
		Example: "SetAO ao_WeldPower, 75;    ! Set welding power to 75%\nSetAO ao_Speed, v_SpeedRef;  ! Set from variable",
		Description: "Sets analog output signal. Controls variable equipment.\n" +
			"- Value range typically 0-100 or custom\n" +
			"- Can use variables as value\n" +
			"- Common use: Speed control, process parameters",
	},
	"wait_time": {
		Name:    "WaitTime",
		Syntax:  "WaitTime <seconds>",
		Example: "SetDO do_Glue, 1;\nWaitTime 0.5;    ! Wait 0.5 seconds\nSetDO do_Glue, 0;",
		Description: "Pauses program execution. Use for timing control.\n" +
			"- Specify time in seconds (can be decimal)\n" +
			"- Accurate to milliseconds\n" +
			"- Common use: Process timing, settling time",
	},
	"robtarget": {
		Name:    "robtarget",
		Syntax:  "CONST robtarget <name>:=[[x,y,z],[q1,q2,q3,q4],[cf1,cf4,cf6,cfx],[ex1,ex2,ex3,ex4,ex5,ex6]];",
		Example: "CONST robtarget p10:=[[500,0,400],[1,0,0,0],[0,0,0,0],[9E9,9E9,9E9,9E9,9E9,9E9]];\n! Point at x=500, z=400",
		Description: "Define robot target position - Complete robot configuration.\n" +
			"- [x,y,z]: Position in mm\n" +
			"- [q1-q4]: Orientation in quaternions\n" +
			"- [cf1,cf4,cf6,cfx]: Robot configuration\n" +
			"- [ex1-ex6]: External axes",
	},
	"offs": {
		Name:    "Offs",
		Syntax:  "Offs(robtarget,x,y,z)",
		Example: "MoveL Offs(p10,100,0,50), v500, z10, tool1;\n! Move to p10 offset 100mm in X, 50mm in Z",
		Description: "Create offset from robtarget. Useful for relative movements.\n" +
			"- Adds offset in x,y,z directions\n" +
			"- Maintains original orientation\n" +
			"- Common use: Pattern movements, approach positions",
	},
	"pulse_do": {
		Name:    "PulseDO",
		Syntax:  "PulseDO Signal [\\High] [\\Time:=0.1]",
		Example: "PulseDO do_Reset;         ! Quick pulse with default time\nPulseDO do_Trigger \\Time:=0.5;  ! 0.5s pulse",
		Description: "Generates a pulse on digital output signal.\n" +
			"- Default pulse time is 0.1 seconds\n" +
			"- \\High keeps signal high after pulse\n" +
			"- Common use: Reset signals, triggers",
	},
	"relative_pos": {
		Name:    "RelTool",
		Syntax:  "RelTool [\\Tool] Point [\\Dx] [\\Dy] [\\Dz] [\\Rx] [\\Ry] [\\Rz]",
		Example: "MoveL RelTool(pCurrent, 0, 0, 50), v100, fine, tool1;  ! Move up 50mm\nMoveL RelTool(pCurrent \\Dx:=100), v100, z10, tool1;  ! Move in X",
		Description: "Create position relative to tool coordinate system.\n" +
			"- Dx,Dy,Dz: Translation in mm\n" +
			"- Rx,Ry,Rz: Rotation in degrees\n" +
			"- Common use: Tool-relative movements",
	},
	"error_recovery": {
		Name:    "ERROR",
		Syntax:  "ERROR\n    [instruction]\n    ...\nENDERROR",
		Example: "ERROR\n    StopMove;         ! Stop robot\n    SetDO do_Error, 1;  ! Signal error\n    Stop;              ! Stop program\nENDERROR",
		Description: "Error recovery handler. Executes when error occurs.\n" +
			"- Place in TRAP routines\n" +
			"- Use with RAISE to trigger\n" +
			"- Common use: Error handling, safety",
	},
	"string_handling": {
		Name:    "StrMatch",
		Syntax:  "StrMatch String1 Pattern [\\MatchLength]",
		Example: "IF StrMatch(partID, \"A*\") THEN\n    ! Handle A-series parts\nENDIF",
		Description: "Pattern matching for strings.\n" +
			"- Supports wildcards (* and ?)\n" +
			"- Case sensitive\n" +
			"- Common use: Part identification",
	},
	"interrupt": {
		Name:    "CONNECT",
		Syntax:  "CONNECT Signal WITH Trap_Routine",
		Example: "CONNECT di_Emergency WITH Emergency_Stop;\n! In TRAP:\nTRAP Emergency_Stop\n    StopMove;\n    Stop;\nENDTRAP",
		Description: "Connect interrupt signal to trap routine.\n" +
			"- Executes trap immediately on signal\n" +
			"- Multiple connects possible\n" +
			"- Common use: Emergency stops, monitoring",
	},
}
