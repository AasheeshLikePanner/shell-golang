package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"github.com/peterh/liner"
)

func main() {
	commands := map[string]func(args []string){
		"exit":   exitCommand,
		"ping":   pingCommand,
		"pwd":    pwdCommand,
		"ls":     lsCommand,
		"cd":     cdCommand,
		"echo":   echoCommand,
		"mkdir":  mkdirCommand,
		"touch":  touchCommand,
		"rm":     rmCommand,
		"cat":    catCommand,
		"clear":  clearCommand,
		"date":   dateCommand,
		"curl":   curlGetCommand,
		"ip":     ipCommand,
		"vim":    vimCommand,
		"theme":  themeCommand,
	}

	line := liner.NewLiner()
	defer line.Close();
	lineAutoComplete(line)
	for {
		if input, err := line.Prompt("> "); err	== nil {
		
			input = strings.TrimSpace(input)
			parts := strings.Fields(input)
			cmd := parts[0]
			args := parts[1:]

			if cmdFunc, exists := commands[cmd]; exists {
				cmdFunc(args)
			} else {
				printError(fmt.Sprintf("Unknown command: %s", cmd))
			}
		}
	}
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

var currentTheme = "default"

var themes = map[string]map[string]string{
	"default": {
		"prompt": Cyan,
		"info":   Green,
		"error":  Red,
	},
	"dark": {
		"prompt": Purple,
		"info":   Blue,
		"error":  Yellow,
	},
	"light": {
		"prompt": White,
		"info":   Blue,
		"error":  Red,
	},
}

func printInfo(msg string) {
	fmt.Println(themes[currentTheme]["info"] + msg + Reset)
}

func printError(msg string) {
	fmt.Println(themes[currentTheme]["error"] + msg + Reset)
}

// --- Input AutoComplete ---

func lineAutoComplete(line *liner.State) {
	line.SetCompleter(func(line string) (c []string) {
		words := []string{"touch", "curl", "theme", "vim", "mkdir", "echo", "rm", "cat", "clear", "date", "ip", "ping", "pwd", "ls", "cd"}
		for _, w := range words {
			if strings.HasPrefix(w, strings.ToLower(line)) {
				c = append(c, w)
			}
		}
		return
	})
}

// --- Command Handlers ---


func vimCommand(args []string) {
	if len(args) == 0 {
		printError("vim: missing file name")
		return
	}

	cmd := exec.Command("vim", args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		printError(fmt.Sprintf("Failed to open vim: %v", err))
	}
}

func ipCommand(args []string) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		printError(fmt.Sprintf("Failed to get IP: %v", err))
		return
	}
	defer resp.Body.Close()
	ip, _ := io.ReadAll(resp.Body)
	printInfo(fmt.Sprintf("Your IP: %s", string(ip)))
}

func curlGetCommand(args []string) {
	if len(args) == 0 {
		printError("get: missing URL")
		return
	}
	resp, err := http.Get(args[0])
	if err != nil {
		printError(fmt.Sprintf("get error: %v", err))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	printInfo(string(body))
}

func clearCommand(args []string) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func catCommand(args []string) {
	if len(args) == 0 {
		printError("cat: missing file name")
		return
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		printError(fmt.Sprintf("cat error: %v", err))
		return
	}
	printInfo(string(data))
}

func rmCommand(args []string) {
	if len(args) == 0 {
		printError("rm: missing file name")
		return
	}
	err := os.Remove(args[0])
	if err != nil {
		printError(fmt.Sprintf("rm error: %v", err))
	}
}

func mkdirCommand(args []string) {
	if len(args) == 0 {
		printError("mkdir: missing folder name")
		return
	}
	err := os.Mkdir(args[0], 0755)
	if err != nil {
		printError(fmt.Sprintf("mkdir error: %v", err))
	}
}

func touchCommand(args []string) {
	if len(args) == 0 {
		printError("touch: missing file name")
		return
	}
	file, err := os.Create(args[0])
	if err != nil {
		printError(fmt.Sprintf("touch error: %v", err))
		return
	}
	file.Close()
}

func echoCommand(args []string) {
	printInfo(strings.Join(args, " "))
}

func exitCommand(args []string) {
	printInfo("Exiting...")
	os.Exit(0)
}

func pingCommand(args []string) {
	printInfo("pong")
}

func pwdCommand(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		printError(fmt.Sprintf("Error getting current working directory: %v", err))
		return
	}
	printInfo(cwd)
}

func dateCommand(args []string) {
	printInfo(time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
}

func cdCommand(args []string) {
	if len(args) == 0 {
		printError("cd: missing directory argument")
		return
	}
	err := os.Chdir(args[0])
	if err != nil {
		printError(fmt.Sprintf("Error changing directory: %v", err))
	}
}

func themeCommand(args []string) {
	if len(args) == 0 {
		printInfo("Available themes:")
		for theme := range themes {
			printInfo(" - " + theme)
		}
		return
	}

	theme := args[0]
	if _, exists := themes[theme]; !exists {
		printError(fmt.Sprintf("Theme not found: %s", theme))
		return
	}

	currentTheme = theme
	printInfo(fmt.Sprintf("Theme changed to: %s", theme))
}

func lsCommand(args []string) {
	files, err := os.ReadDir(".")
	if err != nil {
		printError(fmt.Sprintf("Error reading directory: %v", err))
		return
	}
	for _, file := range files {
		printInfo(file.Name())
	}
}