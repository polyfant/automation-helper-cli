package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/polyfant/automation-helper-cli/abb"
	"github.com/polyfant/automation-helper-cli/ai"
)

// Command represents an automation command with its description and implementation
type Command struct {
	Description string
	Execute     func(args []string) string
}

// commandRegistry stores all available commands
var commandRegistry = make(map[string]Command)

func init() {
	// Register commands
	commandRegistry["sensor"] = Command{
		Description: "Generate sensor code",
		Execute:     generateSensorCode,
	}

	commandRegistry["ai"] = Command{
		Description: "Get AI assistance with ABB RAPID code",
		Execute: func(args []string) string {
			if len(args) < 2 {
				return "Usage: ai help \"your question about the code\""
			}

			apiKey := os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				return "Error: OPENAI_API_KEY environment variable not set"
			}

			assistant := ai.NewAssistant(apiKey)
			response, err := assistant.GetHelp(strings.Join(args[1:], " "))
			if err != nil {
				return fmt.Sprintf("Error getting AI help: %v", err)
			}

			return response
		},
	}

	// Add ABB specific commands
	commandRegistry["abb"] = Command{
		Description: "Get ABB robot programming information and examples",
		Execute: func(args []string) string {
			if len(args) < 1 {
				return `Usage: abb <topic> [subtopic]
Available topics:
1. command  - Show RAPID command details
2. quickref - Show programming reference
3. help     - Show this help message

Examples:
  abb command move_j     - Show MoveJ command details
  abb quickref io        - Show I/O handling guide
  abb list              - List all available commands`
			}

			switch args[0] {
			case "command":
				if len(args) < 2 {
					// List all available commands
					var cmdList []string
					for cmd := range abb.Commands {
						cmdList = append(cmdList, cmd)
					}
					return "Available commands:\n" + strings.Join(cmdList, ", ")
				}
				if cmd, exists := abb.Commands[args[1]]; exists {
					return fmt.Sprintf("\nCommand: %s\nSyntax: %s\n\nExample:\n%s\n\nDescription:\n%s",
						cmd.Name, cmd.Syntax, cmd.Example, cmd.Description)
				}
				return "Unknown ABB command. Type 'abb command' to see available commands."

			case "quickref":
				if len(args) < 2 {
					// List all quick reference topics
					var topics []string
					for topic := range abb.QuickReference {
						topics = append(topics, topic)
					}
					return "Available quick reference topics:\n" + strings.Join(topics, ", ")
				}
				if info, exists := abb.QuickReference[args[1]]; exists {
					return info
				}
				return "Unknown topic. Type 'abb quickref' to see available topics."

			case "list":
				var result strings.Builder
				result.WriteString("\nABB RAPID Commands:\n")
				result.WriteString("================\n")
				for _, cmd := range abb.Commands {
					result.WriteString(fmt.Sprintf("%-10s - %s\n", cmd.Name, cmd.Description))
				}
				return result.String()

			default:
				return "Unknown ABB subcommand. Available: command, quickref, list"
			}
		},
	}
}

func generateSensorCode(args []string) string {
	if len(args) < 2 {
		return "Usage: sensor <type> <action>\nExample: sensor digital when_on"
	}

	sensorType := args[0]
	action := args[1]

	switch sensorType {
	case "digital":
		return generateDigitalSensorCode(action)
	case "analog":
		return generateAnalogSensorCode(action)
	default:
		return "Unknown sensor type. Available types: digital, analog"
	}
}

func generateDigitalSensorCode(action string) string {
	switch action {
	case "when_on":
		return `
PLC Ladder Logic:
|--[INPUT]--|--[OUTPUT]--|

ABB Robot:
IF DI_01 = 1 THEN
    ! Your action here
ENDIF

Siemens S7:
IF "Input_Bit" THEN
    // Your action here
END_IF`
	default:
		return "Unknown action for digital sensor"
	}
}

func generateAnalogSensorCode(action string) string {
	return `
PLC Ladder Logic:
|--[ANALOG_IN]--|--[SCALE]--|--[COMPARE]--|--[OUTPUT]--|

ABB Robot:
IF AI_01 > SET_POINT THEN
    ! Your action here
ENDIF

Siemens S7:
IF "Analog_Input" > "Set_Point" THEN
    // Your action here
END_IF`
}

func printHelp() {
	fmt.Println("\nAutomation Helper CLI")
	fmt.Println("====================")
	fmt.Println("\nAvailable commands:")
	for cmd, info := range commandRegistry {
		fmt.Printf("  %s: %s\n", cmd, info.Description)
	}
	fmt.Println("\nType 'exit' to quit")
}

func main() {
	fmt.Println("Welcome to Automation Helper CLI!")
	fmt.Println("Type 'help' for available commands or 'exit' to quit")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		command := strings.ToLower(args[0])

		switch command {
		case "exit":
			fmt.Println("Goodbye!")
			return
		case "help":
			printHelp()
		default:
			if cmd, exists := commandRegistry[command]; exists {
				result := cmd.Execute(args[1:])
				fmt.Println(result)
			} else {
				fmt.Printf("Unknown command: %s\nType 'help' for available commands\n", command)
			}
		}
	}
}
