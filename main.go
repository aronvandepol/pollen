package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type PollenData struct {
	Name   string
	Status string
}

var (
	// Styles using ANSI 16 colors
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")) // Bright Green

	headerStyle = lipgloss.NewStyle().
			BorderBottom(true).
			BorderForeground(lipgloss.Color("7")). // Light Gray
			PaddingBottom(1).
			MarginBottom(1)

	pollenItemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Status badge styles using ANSI 16 colors
	lowBadge = lipgloss.NewStyle().
			Background(lipgloss.Color("10")). // Bright Green
			Foreground(lipgloss.Color("0")).  // Black
			Padding(0, 1).
			Bold(true).
			Width(6). // Fixed width for alignment
			Align(lipgloss.Center)

	moderateBadge = lipgloss.NewStyle().
			Background(lipgloss.Color("11")). // Bright Yellow
			Foreground(lipgloss.Color("0")).  // Black
			Padding(0, 1).
			Bold(true).
			Width(6). // Fixed width
			Align(lipgloss.Center)

	highBadge = lipgloss.NewStyle().
			Background(lipgloss.Color("9")). // Bright Red
			Foreground(lipgloss.Color("0")).
			Padding(0, 1).
			Bold(true).
			Width(6). // Fixed width
			Align(lipgloss.Center)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")). // Light Gray
			PaddingTop(1).
			MarginTop(1)

	// Error and warning styles using explicit Lipgloss Color
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).Bold(true) // Bright Red

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).Bold(true) // Bright Yellow
)

func getBadge(status string) string {
	switch strings.ToLower(status) {
	case "low":
		return lowBadge.Render("LOW")
	case "moderate":
		return moderateBadge.Render("MOD")
	case "high", "very high":
		return highBadge.Render("HIGH")
	default:
		return status
	}
}

func getEmoji(pollenType string) string {
	lower := strings.ToLower(pollenType)
	switch {
	case strings.Contains(lower, "tree"):
		return "🌳"
	case strings.Contains(lower, "grass"):
		return "🌱"
	case strings.Contains(lower, "ragweed"):
		return "🌾"
	case strings.Contains(lower, "mold"):
		return "🍄"
	case strings.Contains(lower, "dust") || strings.Contains(lower, "dander"):
		return "💨"
	default:
		return "🌿"
	}
}

func parsePollenData(html string) []PollenData {
	var pollenData []PollenData

	cardRegex := regexp.MustCompile(`(?s)<a[^>]*class="[^"]*index-list-card[^"]*"[^>]*>(.*?)</a>`)
	cards := cardRegex.FindAllStringSubmatch(html, -1)

	nameRegex := regexp.MustCompile(`<div class="index-name"[^>]*>([^<]+)</div>`)
	statusRegex := regexp.MustCompile(`<div class="index-status-text">([^<]+)</div>`)

	for _, card := range cards {
		if len(card) < 2 {
			continue
		}

		cardContent := card[1]

		nameMatch := nameRegex.FindStringSubmatch(cardContent)
		statusMatch := statusRegex.FindStringSubmatch(cardContent)

		if len(nameMatch) > 1 && len(statusMatch) > 1 {
			name := strings.TrimSpace(nameMatch[1])
			status := strings.TrimSpace(statusMatch[1])
			name = strings.ReplaceAll(name, "&amp;", "&")

			lowerName := strings.ToLower(name)
			if strings.Contains(lowerName, "pollen") || strings.Contains(lowerName, "mold") || strings.Contains(lowerName, "dust") {
				pollenData = append(pollenData, PollenData{
					Name:   name,
					Status: status,
				})
			}
		}
	}

	return pollenData
}

func fetchPollenData() ([]PollenData, error) {
	url := "https://www.accuweather.com/en/nl/leiden/251527/health-activities/251527"

	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip decompression error: %v", err)
		}
		defer reader.(*gzip.Reader).Close()
	case "br":
		reader = resp.Body
	case "":
		reader = resp.Body
	default:
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	html := string(body)
	return parsePollenData(html), nil
}

func main() {
	// Header
	fmt.Println(titleStyle.Render(" 🌿 Pollen Levels - Leiden "))
	fmt.Println(headerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			"📍 Leiden, Netherlands",
			"🕐 "+time.Now().Format("2006-01-02 15:04"),
		),
	))

	pollenData, err := fetchPollenData()
	if err != nil {
		fmt.Println(errorStyle.Render("❌ Error: " + err.Error()))
		os.Exit(1)
	}

	if len(pollenData) == 0 {
		fmt.Println(warningStyle.Render("⚠️  No pollen data found"))
		os.Exit(1)
	}

	fmt.Println()
	for _, pollen := range pollenData {
		emoji := getEmoji(pollen.Name)
		badge := getBadge(pollen.Status)

		pollenLine := lipgloss.JoinHorizontal(lipgloss.Left,
			fmt.Sprintf("%s %-20s", emoji, pollen.Name),
			badge,
		)
		fmt.Println(pollenItemStyle.Render(pollenLine))
	}

	// // Legend
	// fmt.Println(footerStyle.Render(
	// 	lipgloss.JoinHorizontal(lipgloss.Left,
	// 		"Legend:",
	// 		lowBadge.Render("LOW"),
	// 		moderateBadge.Render("MOD"),
	// 		highBadge.Render("HIGH"),
	// 	),
	// ))
}
