// Package abb provides functionality for working with ABB robots and RAPID programming
package abb

// QuickReference contains common ABB robot programming concepts and snippets
// organized by topic for easy reference
var QuickReference = map[string]string{
	"coordinate_system": `RAPID Coordinate Systems Guide:
- World: Global reference system (default)
- Base: Robot base coordinate system
- Tool: Defined at tool center point (TCP)
- WorkObject: Local coordinate system for workpiece

Examples:
1. Define tool frame:
   PERS tooldata tool1:=[TRUE,[[100,0,100],[1,0,0,0]],[0.5,[0,0,0.1],[1,0,0,0],0,0,0]];

2. Define work object:
   PERS wobjdata wobj1:=[FALSE,TRUE,"",[[0,0,0],[1,0,0,0]],[[0,0,0],[1,0,0,0]]];`,

	"data_types": `RAPID Data Types Reference:
Basic Types:
- num: Numeric value (float) | Example: VAR num distance := 50.5;
- bool: TRUE/FALSE         | Example: VAR bool isReady := TRUE;
- string: Text string      | Example: VAR string message := "Ready";

Position Types:
- pos: Position [x,y,z]   | Example: VAR pos p1 := [100,200,300];
- orient: Quaternion      | Example: VAR orient rot1 := [1,0,0,0];
- pose: Position+orient   | Example: VAR pose target := [[x,y,z],[q1,q2,q3,q4]];
- robtarget: Full target  | Example: See 'robtarget' command reference`,

	"motion_types": `Robot Motion Types Guide:
1. MoveJ (Joint motion)
   - Fastest point-to-point movement
   - Non-linear TCP path
   - Best for large movements

2. MoveL (Linear motion)
   - Straight line TCP path
   - Constant velocity
   - Good for precise paths

3. MoveC (Circular motion)
   - Circular TCP path
   - Requires circle point
   - Perfect for curved paths

4. SearchL/SearchC
   - Motion with sensor input
   - Stops on sensor trigger
   - Used for part detection`,

	"speed_settings": `Speed Settings Reference:
Standard Speeds:
v5    - 5mm/s    | Very slow, precise movements
v50   - 50mm/s   | Careful movements
v100  - 100mm/s  | Normal operation speed
v500  - 500mm/s  | Fast movements
v1000 - 1000mm/s | Very fast movements
v2000 - 2000mm/s | Maximum speed for light tools
vmax  - Maximum possible speed

Custom Speed:
[Speeddata]
v100 := [100, 500, 5000, 1000];
  - TCP linear speed (mm/s)
  - TCP reorientation speed (deg/s)
  - External axis speed
  - Tool reorientation speed`,

	"zone_data": `Zone Data (Path Accuracy) Guide:
fine - Exact positioning (0mm)
z0   - 0.3mm path radius
z1   - 1mm path radius
z5   - 5mm path radius
z10  - 10mm path radius
z20  - 20mm path radius
z50  - 50mm path radius
z100 - 100mm path radius

Usage Tips:
- Use 'fine' for precise operations (picking, placing)
- Use z1-z5 for normal operations
- Use z10-z50 for fast movements
- Larger zones = smoother motion but less accuracy`,

	"io_handling": `I/O Handling Reference:
Digital I/O:
1. Outputs (DO):
   SetDO do_name, value;      | Example: SetDO do_Gripper, 1;
   PulseDO do_name;           | Example: PulseDO do_Reset;
   
2. Inputs (DI):
   value := GetDI di_name;    | Example: IF GetDI di_PartPresent = 1 THEN
   WaitDI di_name, value;     | Example: WaitDI di_Ready, 1;
   
Analog I/O:
1. Outputs (AO):
   SetAO ao_name, value;      | Example: SetAO ao_Speed, 50.5;
   
2. Inputs (AI):
   value := GetAI ai_name;    | Example: pressure := GetAI ai_Pressure;

Group I/O:
   SetGO go_name, value;      | Example: SetGO go_Status, 2;
   value := GetGI gi_name;    | Example: mode := GetGI gi_OpMode;`,

	"error_handling": `Error Handling Guide:
1. Basic Error Handler:
   ERROR
       IF ERRNO = ERR_PATH_STOP THEN
           StopMove;
           ClearPath;
           StartMove;
           RETRY;
       ELSE
           Stop;
       ENDIF

2. Common Error Types:
   ERR_PATH_STOP    - Motion path interrupted
   ERR_COLL_STOP   - Collision detected
   ERR_OUTOFBND    - Position out of range
   ERR_REFUNKDAT   - Undefined data used
   
3. Recovery Actions:
   RETRY           - Retry from error point
   TRYNEXT         - Skip to next instruction
   RETURN          - Exit routine
   EXIT            - Exit program
   
4. Error Logging:
   ErrLog error_msg;        | Logs error to system
   TPWrite "Error: " + error_msg;  | Display on FlexPendant`,

	"basic_program": `Basic Program Structure:
MODULE MainModule
    ! Variable declarations
    PERS tooldata currentTool := [...];
    PERS wobjdata currentWobj := [...];
    VAR robtarget homePos;
    
    ! Main procedure
    PROC main()
        ! Initialize
        TPWrite "Program Starting...";
        MoveJ homePos, v1000, z50, currentTool;
        
        ! Main loop
        WHILE running DO
            ! Check conditions
            IF GetDI di_StartCycle = 1 THEN
                Cycle;
            ENDIF
            
            ! Error checking
            ERROR
                IF ERRNO = ERR_PATH_STOP THEN
                    TPWrite "Path was stopped";
                    RETRY;
                ENDIF
        ENDWHILE
    ENDPROC
    
    ! Subroutines
    PROC Cycle()
        MoveJ p10, v1000, z50, currentTool;
        SetDO do_Gripper, 1;
        WaitTime 0.5;
        MoveL p20, v500, fine, currentTool;
    ENDPROC
ENDMODULE`,

	"common_patterns": `Common Programming Patterns:
1. Pick and Place:
   PROC PickAndPlace()
       MoveJ approach, v1000, z10, tool1;
       MoveL pick, v100, fine, tool1;
       SetDO do_Gripper, 1;
       WaitTime 0.2;
       MoveL approach, v100, z10, tool1;
       MoveJ place_approach, v1000, z10, tool1;
       MoveL place, v100, fine, tool1;
       SetDO do_Gripper, 0;
   ENDPROC

2. Palletizing:
   PROC Palletize()
       FOR layer FROM 1 TO 3 DO
           FOR row FROM 1 TO 2 DO
               FOR col FROM 1 TO 3 DO
                   current_pos := Offs(base_pos, 
                       col*100, row*100, layer*50);
                   MoveL current_pos, v500, z10, tool1;
               ENDFOR
           ENDFOR
       ENDFOR
   ENDPROC

3. Search Pattern:
   PROC SearchObject()
       SearchL \Stop \Tool:=tool1 
           \MaxTime:=5 
           \PoseOffs:=offs 
           start_pos, 
           search_pos, 
           v100, 
           tool1;
       IF FOUND THEN
           TPWrite "Object found!";
       ENDIF
   ENDPROC`,
}
