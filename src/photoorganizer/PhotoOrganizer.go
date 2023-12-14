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
	//"github.com/rwcarlsen/goexif/exif"
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
	".srw":  true,
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
    cmd := exec.Command("exiftool", "-DateTimeOriginal", path)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return 0, 0, err
    }

    dateStr := strings.TrimSpace(out.String())
    if strings.Contains(dateStr, ":") {
        parts := strings.Split(dateStr, ": ")
        if len(parts) < 2 {
            return 0, 0, fmt.Errorf("DateTimeOriginal not found in %s", path)
        }
        datePart := parts[1]
        dt, err := time.Parse("2006:01:02 15:04:05", datePart)
        if err != nil {
            return 0, 0, err
        }
        return dt.Year(), int(dt.Month()), nil
    }
    return 0, 0, fmt.Errorf("Failed to extract DateTimeOriginal for %s", path)
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
			// Move file to "undated" directory if there's an error
			dstDir := filepath.Join("./media", "undated")
			if filepath.Dir(path) == dstDir {
				fmt.Printf("%s is already in the undated location\n", path)
				return
			}

			if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
				fmt.Printf("Failed to create directory %s: %s\n", dstDir, err)
				return
			}

			dstPath := filepath.Join(dstDir, filepath.Base(path))
			if err := tryMoveFile(path, dstPath); err != nil {
				fmt.Printf("Failed to move %s: %s\n", path, err)
				return
			} else {
				mu.Lock()
				counter++
				printProgress(counter, totalFiles)
				mu.Unlock()
			}
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
		if err := tryMoveFile(path, dstPath); err != nil {
			fmt.Printf("Failed to move %s: %s\n", path, err)
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

func tryMoveFile(src, dst string) error {
    originalDst := dst
    counter := 1
    for {
        err := os.Rename(src, dst)
        if err != nil {
            if os.IsExist(err) {
                // File already exists, try a new name
                ext := filepath.Ext(originalDst)
                name := strings.TrimSuffix(filepath.Base(originalDst), ext)
                dst = filepath.Join(filepath.Dir(originalDst), fmt.Sprintf("%s_%d%s", name, counter, ext))
                counter++
            } else {
                // Some other error
                return err
            }
        } else {
            // File moved successfully
            return nil
        }
    }
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
