package msgfmt

import (
	"strings"
)

// Usually something like
// ───────────────
// >
// ───────────────
// Used by Claude Code, Goose, and Aider.
func findGreaterThanMessageBox(lines []string) int {
	for i := len(lines) - 1; i >= max(len(lines)-6, 0); i-- {
		if strings.Contains(lines[i], ">") {
			if i > 0 && strings.Contains(lines[i-1], "───────────────") {
				return i - 1
			}
			return i
		}
	}
	return -1
}

// Usually something like
// ───────────────
// |
// ───────────────
func findGenericSlimMessageBox(lines []string) int {
	for i := len(lines) - 3; i >= max(len(lines)-9, 0); i-- {
		if strings.Contains(lines[i], "───────────────") &&
			(strings.Contains(lines[i+1], "|") || strings.Contains(lines[i+1], "│") || strings.Contains(lines[i+1], "❯")) &&
			strings.Contains(lines[i+2], "───────────────") {
			return i
		}
	}
	return -1
}

func removeMessageBox(msg string) string {
	lines := strings.Split(msg, "\n")

	messageBoxStartIdx := findGreaterThanMessageBox(lines)
	if messageBoxStartIdx == -1 {
		messageBoxStartIdx = findGenericSlimMessageBox(lines)
	}

	if messageBoxStartIdx != -1 {
		lines = lines[:messageBoxStartIdx]
	}

	return strings.Join(lines, "\n")
}

func removeCodexInputBox(msg string) string {
	lines := strings.Split(msg, "\n")
	// Remove the input box, we need to match the exact pattern, because thinking follows the same pattern of ▌ followed by text
	if len(lines) >= 2 && strings.Contains(lines[len(lines)-2], "▌ Ask Codex to do anything") {
		idx := len(lines) - 2
		lines = append(lines[:idx], lines[idx+1:]...)
	}
	return strings.Join(lines, "\n")
}

func removeOpencodeMessageBox(msg string) string {
	lines := strings.Split(msg, "\n")
	//
	//  ┃
	//  ┃
	//  ┃
	//  ┃  Build  Anthropic Claude Sonnet 4
	//  ╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
	//                                tab switch agent  ctrl+p commands
	//
	for i := len(lines) - 1; i >= 4; i-- {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "╹▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀") {
			lines = lines[:i-4]
			break
		}
	}
	return strings.Join(lines, "\n")
}

// removeOpenClawMessageBox removes the OpenClaw TUI input box.
// OpenClaw TUI has a chat input area at the bottom with borders like:
//   ┌──────────────────────────────────────────────────────────────┐
//   │ >                                                           │
//   └──────────────────────────────────────────────────────────────┘
// Or the "🦞" lobster prompt indicator.
func removeOpenClawMessageBox(msg string) string {
	lines := strings.Split(msg, "\n")
	// Search from the bottom for the input box pattern
	for i := len(lines) - 1; i >= max(len(lines)-8, 0); i-- {
		line := strings.TrimSpace(lines[i])
		// OpenClaw TUI footer with lobster emoji or border
		if strings.HasPrefix(line, "└") && strings.HasSuffix(line, "┘") {
			// Find matching top border
			for j := i - 1; j >= max(i-4, 0); j-- {
				jLine := strings.TrimSpace(lines[j])
				if strings.HasPrefix(jLine, "┌") && strings.HasSuffix(jLine, "┐") {
					return strings.Join(lines[:j], "\n")
				}
			}
			return strings.Join(lines[:i], "\n")
		}
		// Also handle the "🦞" prompt line
		if strings.Contains(line, "🦞") && (strings.Contains(line, ">") || strings.Contains(line, "│")) {
			return strings.Join(lines[:i], "\n")
		}
	}
	return msg
}

func removeAmpMessageBox(msg string) string {
	lines := strings.Split(msg, "\n")
	msgBoxEndFound := false
	msgBoxStartIdx := len(lines)
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if !msgBoxEndFound && strings.HasPrefix(line, "╰") && strings.HasSuffix(line, "╯") {
			msgBoxEndFound = true
		}
		if msgBoxEndFound && strings.HasPrefix(line, "╭") && strings.HasSuffix(line, "╮") {
			msgBoxStartIdx = i
			break
		}
	}
	formattedMsg := strings.Join(lines[:msgBoxStartIdx], "\n")
	if len(formattedMsg) == 0 {
		return "Welcome to Amp"
	}
	return formattedMsg
}
