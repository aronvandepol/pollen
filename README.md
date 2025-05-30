# AccuWeather Pollen Checker

A simple Go command-line tool to get current pollen levels for Leiden, Netherlands, by scraping the AccuWeather website.

**Disclaimer:** This script scrapes a website and is not an official API client. Website structure changes may break it.

## Why this script?

Get quick pollen info in your terminal without needing an API key or paying for a service.

## How it Works

Fetches the AccuWeather page for Leiden and parses the HTML to extract pollen types and their status (Low, Moderate, High). Uses Lipgloss for a clean, color-coded output.

**Important:** Currently hardcoded for Leiden, Netherlands.

## Prerequisites

- Go (1.16+ for building)
- Internet connection

## Installation & Usage

You have a couple of options:

**1. Build from source:**

- Clone this repository:
  ```bash
  git clone https://github.com/aronvandepol/pollen
  cd pollen
  ```
- Go modules will handle dependencies. Build the executable:
  ```bash
  go build -o pollen
  ```
- Run the executable:
  ```bash
  ./pollen
  ```

**2. Use a pre-built binary (if available):**

- Download the appropriate binary for your operating system from the [releases page](link-to-your-releases).
- Make the binary executable (if needed on Linux/macOS):
  ```bash
  chmod +x ./pollen
  ```
- Run the binary:
  ```bash
  ./pollen
  ```

Example output:

The output will display the current pollen levels:
```bash
🌿 Pollen Levels - Leiden
─────────────────────────────────
📍 Leiden, Netherlands
🕐 2025-05-30 19:47:26

🌳 Tree Pollen      LOW
🌾 Ragweed Pollen   LOW
🍄 Mold             HIGH
🌱 Grass Pollen     LOW
💨 Dust & Dander    HIGH
```

## Future Improvements

- **Command-line arguments for location:** Allow users to specify a city or location ID instead of being hardcoded to Leiden. This would likely involve finding a way to look up location IDs on AccuWeather or another source.
- **More robust HTML parsing:** Using a dedicated HTML parsing library to make the script more resilient to website changes compared to regular expressions.
- **Error handling for network issues:** Add more specific error handling for different types of network errors (e.g., timeout, connection refused).
- **Caching:** Implement a simple caching mechanism to avoid hitting the website on every request, especially if checking frequently.
- **Different data sources:** Investigate other weather websites or public APIs that might provide pollen data.
- **Configuration file:** Allow setting the location and other options via a configuration file.

## Contributing

Contributions are welcome! If you have suggestions for improvements or want to contribute code, please feel free to open an issue or pull request on the GitHub repository.

## License

This project is open source. Feel free to copy, download or otherwise adjust.
