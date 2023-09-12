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
	".cr3":  true,
	".cr2":  true,
}

var videoExtensions = map[string]bool{
	".mp4": true,
	".avi": true,
	".mov": true,
	".mkv": true,
}

var wg sync.WaitGroup
var mu sync.Mutex
var counter int

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
		dt, err := time.Parse("2006:01:02 15:04:05", strings.Split(dateStr, ": ")[1])
		if err != nil {
			return 0, 0, err
		}
		return dt.Year(), int(dt.Month()), nil
	}

	return 0, 0, fmt.Errorf("Unsupported file extension: %s", ext)
}

func processFile(path string, info os.FileInfo) {
	defer wg.Done()

	if isPhotoOrVideo(path) {
		year, month, err := getDateTaken(path)
		if err != nil {
			fmt.Printf("Failed to get date taken for %s: %s\n", path, err)
			return
		}

		dstDir := filepath.Join("./media", fmt.Sprintf("%d", year), fmt.Sprintf("%02d", month))
		if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory %s: %s\n", dstDir, err)
			return
		}

		dstPath := filepath.Join(dstDir, info.Name())
		if err := os.Rename(path, dstPath); err != nil {
			fmt.Printf("Failed to move %s to %s: %s\n", path, dstPath, err)
			return
		} else {
			mu.Lock()
			counter++
			mu.Unlock()
			fmt.Printf("Moved %s to %s\n", path, dstPath)
		}
	}
}

func main() {
	srcDir := "./media"

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			wg.Add(1)
			go processFile(path, info)
		}

		return nil
	})

	wg.Wait()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("\nMoved %d files.\n", counter)
	fmt.Println("Press any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
