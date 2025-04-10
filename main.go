package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"runtime"
	"os/exec"
	"time"
	"net/http"
	"io"
)

func main() {
	commands := map[string]func(args []string, curDir string){
		"exit": exitCommand,
		"ping": pingCommand,
		"pwd":  pwdCommand,
		"ls":   lsCommand,
		"cd": cdCommand,
		"echo": echoCommand,
		"mkdir": mkdirCommand,
		"touch": touchCommand,
		"rm": rmCommand,
		"cat": catCommand,
		"clear": clearCommand,
		"date": dateCommand,
		"curl": curlGetCommand,
		"ip": ipCommand,
		"vim": vimCommand,
	}

	curDir, err := os.Getwd()
	if err != nil {	
		fmt.Println("Error getting current working directory:", err)
		return;
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Split input into command and args
		parts := strings.Fields(line)
		cmd := parts[0]
		args := parts[1:]

		if cmdFunc, exists := commands[cmd]; exists {
			cmdFunc(args, curDir)
		} else {
			fmt.Println("Unknown command:", cmd);
		}
	}
}

// --- Command Handlers ---

func vimCommand(args []string, curDir string) {
	if len(args) == 0 {
		fmt.Println("vim: missing file name")
		return
	}

	cmd := exec.Command("vim", args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to open vim:", err)
	}
}

func ipCommand(args []string, curDir string) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Println("Failed to get IP:", err)
		return
	}
	defer resp.Body.Close()
	ip, _ := io.ReadAll(resp.Body)
	fmt.Println("Your IP:", string(ip))
}


func curlGetCommand(args []string, curDir string) {
	if len(args) == 0 {
		fmt.Println("get: missing URL")
		return
	}
	resp, err := http.Get(args[0])
	if err != nil {
		fmt.Println("get error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func clearCommand(args []string, curDir string) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func catCommand(args []string, curDir string) {
	if len(args) == 0 {
		fmt.Println("cat: missing file name")
		return
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("cat error:", err)
		return
	}
	fmt.Println(string(data))
}


func rmCommand(args []string, curDir string) {
	if len(args) == 0 {
		fmt.Println("rm: missing file name")
		return
	}
	err := os.Remove(args[0])
	if err != nil {
		fmt.Println("rm error:", err)
	}
}


func mkdirCommand(args []string, curDir string) {
	if len(args) == 0 {
		fmt.Println("mkdir: missing folder name")
		return
	}
	err := os.Mkdir(args[0], 0755)
	if err != nil {
		fmt.Println("mkdir error:", err)
	}
}

func touchCommand(args []string, curDir string) {
	if len(args) == 0 {
		fmt.Println("touch: missing file name")
		return
	}
	file, err := os.Create(args[0])
	if err != nil {
		fmt.Println("touch error:", err)
		return
	}
	file.Close()
}


func echoCommand(args []string, curDir string) {
	fmt.Println(strings.Join(args, " "))
}

func exitCommand(args [] string, curDir string) {
	fmt.Println("Exiting...")
	os.Exit(0)
}

func pingCommand(args [] string, curDir string) {
	fmt.Println("pong")
}

func pwdCommand(args [] string, curDir string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}
	fmt.Println(cwd)
}

func dateCommand(args []string, curDir string) {
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
}

func cdCommand(args []string, curDir string){
	err := os.Chdir(args[0]);
	if err != nil {
		fmt.Println("Error changing directory:", err);
		return
	}

}

func lsCommand(args []string, curDir string) {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
