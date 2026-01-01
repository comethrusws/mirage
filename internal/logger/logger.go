package logger

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00d4ff")).
			PaddingLeft(1).
			PaddingRight(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	methodStyleGET = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#10b981")).
			PaddingLeft(1).
			PaddingRight(1)

	methodStylePOST = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#3b82f6")).
			PaddingLeft(1).
			PaddingRight(1)

	methodStylePUT = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f59e0b")).
		PaddingLeft(1).
		PaddingRight(1)

	methodStyleDELETE = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ef4444")).
				PaddingLeft(1).
				PaddingRight(1)

	methodStyleDefault = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7c3aed")).
				PaddingLeft(1).
				PaddingRight(1)

	statusStyleOK = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#10b981"))

	statusStyleWarn = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#f59e0b"))

	statusStyleError = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ef4444"))

	mockStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00d4ff")).
			PaddingLeft(1)

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a0aec0"))

	durationStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)
)

func PrintBanner(version string) {
	banner := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00d4ff")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00d4ff")).
		Padding(0, 2).
		Align(lipgloss.Center).
		Width(50)

	title := fmt.Sprintf("MIRAGE %s", version)
	subtitle := "API Mocking Gateway & Recorder"

	fmt.Println()
	fmt.Println(banner.Render(title))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Center).Width(50).Render(subtitle))
	fmt.Println()
}

func LogRequest(method, url, body string) {
	methodStyled := getMethodStyle(method).Render(method)
	urlStyled := urlStyle.Render(url)
	
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s  %s %s\n", 
		lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(timestamp),
		methodStyled,
		urlStyled)
	
	if body != "" && len(body) < 200 {
		fmt.Printf("         %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render("→ "+body))
	}
}

func LogResponse(status int, duration time.Duration, body string) {
	statusStyled := getStatusStyle(status).Render(fmt.Sprintf("%d", status))
	durationStyled := durationStyle.Render(duration.String())
	
	fmt.Printf("         %s  %s\n", statusStyled, durationStyled)
	
	if body != "" && len(body) < 200 {
		fmt.Printf("         %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render("← "+body))
	}
}

func LogMock(scenarioName string, status int, duration time.Duration) {
	mockStyled := mockStyle.Render("MOCK")
	scenarioStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("#00d4ff")).Render(scenarioName)
	statusStyled := getStatusStyle(status).Render(fmt.Sprintf("%d", status))
	durationStyled := durationStyle.Render(duration.String())
	
	fmt.Printf("         %s %s  %s  %s\n", mockStyled, scenarioStyled, statusStyled, durationStyled)
}

func LogInfo(message string) {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#00d4ff")).Render("ℹ " + message))
}

func LogSuccess(message string) {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#10b981")).Render("✓ " + message))
}

func LogError(message string) {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render("✗ " + message))
}

func getMethodStyle(method string) lipgloss.Style {
	switch method {
	case "GET":
		return methodStyleGET
	case "POST":
		return methodStylePOST
	case "PUT", "PATCH":
		return methodStylePUT
	case "DELETE":
		return methodStyleDELETE
	default:
		return methodStyleDefault
	}
}

func getStatusStyle(status int) lipgloss.Style {
	if status >= 200 && status < 300 {
		return statusStyleOK
	} else if status >= 300 && status < 400 {
		return statusStyleWarn
	} else {
		return statusStyleError
	}
}
