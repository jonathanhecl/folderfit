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
		fmt.Println("Calculating selection...")
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

func selectBestFolders(folderSizes map[string]int, totalSize int) map[string]int {
	names := make([]string, 0, len(folderSizes))
	sizes := make([]int, 0, len(folderSizes))
	for name, size := range folderSizes {
		names = append(names, name)
		sizes = append(sizes, size)
	}

	n := len(names)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, totalSize+1)
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= totalSize; j++ {
			if sizes[i-1] <= j {
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-sizes[i-1]]+sizes[i-1])
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}

	selected := make(map[string]int)
	j := totalSize
	for i := n; i > 0 && dp[i][j] != 0; i-- {
		if dp[i][j] != dp[i-1][j] {
			selected[names[i-1]] = sizes[i-1]
			j -= sizes[i-1]
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
