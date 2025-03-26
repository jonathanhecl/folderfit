package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var version = "0.0.1"

func printUsage() {
	fmt.Println("\nUsage: folderfit <sources> -size=<totalsize> [-verbose]")
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
			sources = append(sources, arg)
		}
	}

	if totalSize == 0 {
		fmt.Println("Invalid size argument")
		printUsage()
		return
	}

	if verbose {
		fmt.Println()
		fmt.Println("Analyzing...")
	}
	folderSizes := make(map[string]int)
	for _, source := range sources {
		folderSizes[source] = calculateSize(source)
	}

	if verbose {
		for name, size := range folderSizes {
			fmt.Printf("%s - %s\n", name, formatSize(size))
		}
	}

	if verbose {
		fmt.Printf("\nTotal target size: %s\n\n", formatSize(totalSize))
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
	} else if size >= 1024 && size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	}
}
