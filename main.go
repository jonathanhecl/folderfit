package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var version = "1.0.2"

func printUsage() {
	fmt.Println("Usage: folderfit <sources> -size=<totalsize> [-verbose]")
	fmt.Println("- Sources can be a list of files and folders or a single * to include all files and folders in the current directory")
	fmt.Println("- Size is the total size in bytes")
	fmt.Println("- Verbose is optional and will print more information")
	fmt.Println("Example: folderfit * -size=1024")
}

func main() {
	fmt.Println("FolderFit v", version)

	if len(os.Args) < 3 {
		printUsage()
		return
	}

	var totalSize int = 0
	var sources []string
	var verbose bool

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-size=") {
			sizeStr := strings.TrimPrefix(arg, "-size=")
			var err error
			totalSize, err = strconv.Atoi(sizeStr)
			if err != nil {
				fmt.Println("Invalid size argument")
				return
			}
		} else if arg == "-verbose" {
			verbose = true
		} else {
			if arg == "*" {
				sources = append(sources, getAllFilesAndFolders(".")...)
			} else {
				sources = append(sources, arg)
			}
		}
	}

	if totalSize == 0 {
		fmt.Println("Invalid size argument")
		printUsage()
		return
	}

	if verbose {
		fmt.Println()
		fmt.Println("Calculating sizes...")
	}
	folderSizes := make(map[string]int)
	for _, source := range sources {
		folderSizes[source] = calculateSize(source)
	}

	if verbose {
		totalSizeSource := 0
		for name, size := range folderSizes {
			fmt.Printf("%s - %s\n", name, formatSize(size))
			totalSizeSource += size
		}
		fmt.Printf("\nTotal source size: %s (%d files)", formatSize(totalSizeSource), len(folderSizes))
	}

	if verbose {
		fmt.Printf("\nTotal target size: %s\n", formatSize(totalSize))
		fmt.Println()
		fmt.Println("Calculating selection... (", estimateTime(len(folderSizes), totalSize), " estimated time)")
	}

	selected := selectBestFolders(folderSizes, totalSize)

	if len(selected) == 0 {
		fmt.Println("No selection possible")
		return
	}

	if verbose {
		fmt.Println("Selected:")
	}
	for name, size := range selected {
		fmt.Printf("%s - %s\n", name, formatSize(size))
	}
	fmt.Printf("\nSelection size: %s / %s\n", formatSize(calculateTotalSize(selected)), formatSize(totalSize))
	fmt.Printf("Free space: %s\n", formatSize(totalSize-calculateTotalSize(selected)))
}

func getAllFilesAndFolders(path string) []string {
	var files []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return files
	}
	for _, entry := range entries {
		files = append(files, filepath.Join(path, entry.Name()))
	}

	return files
}

func calculateSize(source string) int {
	var totalSize int

	info, err := os.Stat(source)
	if err != nil {
		return 0
	}

	if info.IsDir() {
		files, err := os.ReadDir(source)
		if err != nil {
			return 0
		}
		for _, file := range files {
			totalSize += calculateSize(filepath.Join(source, file.Name()))
		}
	} else {
		totalSize = int(info.Size())
	}

	return totalSize
}

func estimateTime(folderCount int, totalSizeBytes int) string {
	// Ajustar la constante para reflejar el nuevo algoritmo más rápido
	const k = 0.000001

	estimatedMs := k * float64(folderCount) * float64(totalSizeBytes)

	if estimatedMs < 1000 {
		return fmt.Sprintf("%.2f ms", estimatedMs)
	} else if estimatedMs < 60000 {
		return fmt.Sprintf("%.2f s", estimatedMs/1000)
	} else {
		return fmt.Sprintf("%.2f m", estimatedMs/60000)
	}
}

func selectBestFolders(folderSizes map[string]int, totalSize int) map[string]int {
	totalSizeOfAllFiles := 0
	for _, size := range folderSizes {
		totalSizeOfAllFiles += size
	}

	if totalSizeOfAllFiles <= totalSize {
		return folderSizes
	}

	scalingFactor := 1
	scaledTotalSize := totalSize

	if totalSize > 100000 {
		if totalSize > 1000000000 {
			scalingFactor = 1000000
		} else if totalSize > 1000000 {
			scalingFactor = 1000
		} else {
			scalingFactor = 100
		}
		scaledTotalSize = totalSize / scalingFactor
	}

	names := make([]string, 0, len(folderSizes))
	sizes := make([]int, 0, len(folderSizes))
	scaledSizes := make([]int, 0, len(folderSizes))

	for name, size := range folderSizes {
		names = append(names, name)
		sizes = append(sizes, size)
		scaledSize := size / scalingFactor
		if scalingFactor > 1 && size > 0 && scaledSize == 0 {
			scaledSize = 1
		}
		scaledSizes = append(scaledSizes, scaledSize)
	}

	n := len(names)
	dp := make([]int, scaledTotalSize+1)

	keep := make([][]bool, n+1)
	for i := range keep {
		keep[i] = make([]bool, scaledTotalSize+1)
	}

	for i := 1; i <= n; i++ {
		for j := scaledTotalSize; j >= 1; j-- {
			if scaledSizes[i-1] <= j {
				prev := dp[j]
				taken := dp[j-scaledSizes[i-1]] + scaledSizes[i-1]

				if taken > prev {
					dp[j] = taken
					keep[i][j] = true
				}
			}
		}
	}

	selected := make(map[string]int)
	j := scaledTotalSize

	for i := n; i > 0 && j > 0; i-- {
		if keep[i][j] {
			selected[names[i-1]] = sizes[i-1]
			j -= scaledSizes[i-1]
		}
	}

	for {
		total := calculateTotalSize(selected)
		if total <= totalSize {
			break
		}

		minValueRatio := float64(1 << 60)
		var elementToRemove string
		for name, size := range selected {
			scaledSize := size / scalingFactor
			if scaledSize == 0 {
				scaledSize = 1
			}

			ratio := float64(size) / float64(scaledSize)
			if ratio < minValueRatio {
				minValueRatio = ratio
				elementToRemove = name
			}
		}

		if elementToRemove != "" {
			delete(selected, elementToRemove)
		} else {
			selected = make(map[string]int)
			break
		}
	}

	for {
		addMore := false
		currentTotal := calculateTotalSize(selected)
		remaining := totalSize - currentTotal

		if remaining <= 0 {
			break
		}

		for name, size := range folderSizes {
			_, exists := selected[name]
			if !exists && size <= remaining {
				selected[name] = size
				addMore = true
				break
			}
		}

		if !addMore {
			break
		}
	}

	return selected
}

func calculateTotalSize(folderSizes map[string]int) int {
	total := 0
	for _, size := range folderSizes {
		total += size
	}
	return total
}

func formatSize(size int) string {
	if size < 1024 {
		return fmt.Sprintf("%d bytes", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%d KB", size/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}
