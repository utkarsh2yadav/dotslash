package service

import (
	. "dotslash/model"
	"fmt"
	"regexp"
)

func GetLanguage(name string) *Language {
	l := &Language{
		Name:     name,
		Filename: "main",
	}
	if l.getExtension() == "" {
		return nil
	} else {
		return l
	}
}

func (l *Language) getCommand(body *WsBody) []string {
	if l.Name == "java" {
		pattern := regexp.MustCompile(`(?s)class\s+(\w+).*?public\s+static\s+void\s+main\s*\(String(?:\s*\[\]\s+\w+|\s+\w+\s*\[\])\)`)
		matches := pattern.FindStringSubmatch(body.Code)
		if len(matches) >= 2 {
			l.Filename = matches[1]
		}
	}

	file := fmt.Sprintf("%v.%v", l.Filename, l.getExtension())
	command := []string{}
	switch l.Name {
	case "c":
		command = append(command, "gcc", file, "-o", "main", "&&", "./main")
	case "cpp":
		command = append(command, "g++", file, "-o", "main", "&&", "./main")
	case "golang":
		command = append(command, "go", "run", file)
	case "java":
		command = append(command, "javac", file, "&&", "java", l.Filename)
	case "javascript":
		command = append(command, "node", file)
	case "python2":
		command = append(command, "python2", file)
	case "python3":
		command = append(command, "python3", file)
	case "typescript":
		command = append(command, "tsc", file, "&&", "node", l.Filename+".js")
	}
	return command
}

func (l *Language) getExtension() string {
	switch l.Name {
	case "c":
		return "c"
	case "cpp":
		return "cpp"
	case "golang":
		return "go"
	case "java":
		return "java"
	case "javascript":
		return "js"
	case "python2", "python3":
		return "py"
	case "typescript":
		return "ts"
	default:
		return ""
	}
}
