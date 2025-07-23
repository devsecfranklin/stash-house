package minecraft

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func Check_Minecraft_Install() {
	var minecraftPath string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		minecraftPath = filepath.Join(appData, ".minecraft")
	case "darwin": // macOS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		minecraftPath = filepath.Join(homeDir, "Library", "Application Support", "minecraft")
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		minecraftPath = filepath.Join(homeDir, ".minecraft")
	default:
		fmt.Println("Unsupported operating system:", runtime.GOOS)
		return
	}

	info, err := os.Stat(minecraftPath)
	if os.IsNotExist(err) {
		fmt.Println("Minecraft is not installed ('.minecraft' folder not found).")
	} else if err != nil {
		fmt.Println("Error checking Minecraft installation:", err)
	} else if info.IsDir() {
		fmt.Println("Minecraft is installed at:", minecraftPath)
	} else {
		fmt.Println("Found a file named '.minecraft', but it's not a directory.")
	}
}

