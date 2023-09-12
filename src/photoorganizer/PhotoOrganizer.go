package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var photoExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".cr2":  true,
	".cr3":  true,
	".tif":  true,
	".tiff": true,
	".nef":  true,
	".orf":  true,
	".dng":  true,
	".arw":  true,
}

var videoExtensions = map[string]bool{
	".mp4": true,
	".avi": true,
	".mov": true,
	".mkv": true,
	".flv": true,
	".wmv": true,
	".m4v": true,
}

var wg sync.WaitGroup
var mu sync.Mutex
var counter, totalFiles int
var semaphore = make(chan struct{}, 2)

func isPhotoOrVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return photoExtensions[ext] || videoExtensions[ext]
}

func getDateTaken(path string) (int, int, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if photoExtensions[ext] {
		file, err := os.Open(path)
		if err != nil {
			return 0, 0, err
		}
		defer file.Close()

		x, err := exif.Decode(file)
		if err != nil {
			return 0, 0, err
		}

		dt, err := x.DateTime()
		if err != nil {
			return 0, 0, err
		}
		return dt.Year(), int(dt.Month()), nil
	} else if videoExtensions[ext] {
		cmd := exec.Command("exiftool", "-DateTimeOriginal", path)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return 0, 0, err
		}

		dateStr := strings.TrimSpace(out.String())
		parts := strings.Split(dateStr, ": ")
		if len(parts) < 2 {
			return 0, 0, fmt.Errorf("Unexpected exiftool output format: %s", dateStr)
		}

		dt, err := time.Parse("2006:01:02 15:04:05", parts[1])
		if err != nil {
			return 0, 0, err
		}
		return dt.Year(), int(dt.Month()), nil
	}

	return 0, 0, fmt.Errorf("Unsupported file extension: %s", ext)
}
func processFile(path string) {
	defer wg.Done()
	defer func() {
		<-semaphore // release a slot
	}()

	if isPhotoOrVideo(path) {
		year, month, err := getDateTaken(path)
		if err != nil {
			fmt.Printf("Failed to get date taken for %s: %s\n", path, err)
			return
		}

		dstDir := filepath.Join("./media", fmt.Sprintf("%d", year), fmt.Sprintf("%02d", month))
		if filepath.Dir(path) == dstDir {
			fmt.Printf("%s is already in the correct location\n", path)
			return
		}

		if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory %s: %s\n", dstDir, err)
			return
		}

		dstPath := filepath.Join(dstDir, filepath.Base(path))
		if err := os.Rename(path, dstPath); err != nil {
			fmt.Printf("Failed to move %s to %s: %s\n", path, dstPath, err)
			return
		} else {
			mu.Lock()
			counter++
			printProgress(counter, totalFiles)
			mu.Unlock()
		}
	}
}

func printProgress(count, total int) {
	const barLength = 40
	progress := float64(count) / float64(total)
	bars := int(progress * barLength)
	fmt.Printf("\r[%s%s] %d/%d", strings.Repeat("=", bars), strings.Repeat(" ", barLength-bars), count, total)
}

func main() {
	// Welcome message
	fmt.Println("=============================================")
	fmt.Println("   Mike's Photo & Video Sorting Tool!")
	fmt.Println("=============================================")
	fmt.Println("This tool organizes your photos and videos by extracting the 'date taken' metadata and sorting them into folders by year and month.")
	fmt.Println("For accurate results, please ensure that 'exiftool' is installed on your system.")
	fmt.Println("\nCo-coded with ChatGPT from OpenAI! 🚀")
	fmt.Println("\nLet's get started!")
	fmt.Println("---------------------------------------------\n")

	fmt.Println("Which folder would you like to sort?")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	srcDir := scanner.Text()

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isPhotoOrVideo(path) {
			totalFiles++
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			semaphore <- struct{}{} // acquire a slot
			wg.Add(1)
			go processFile(path)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("\nMoved %d files out of %d.\n", counter, totalFiles)
	fmt.Println("Press any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
