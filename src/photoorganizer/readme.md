# Mike's Photo & Video Sorting Tool

This is a utility tool designed to help sort photos and videos into directories based on their "date taken" metadata. It's a result of collaborative coding between Mike and ChatGPT from OpenAI.

## Code Breakdown

- **Imports**: Various packages are imported for the functionality ranging from handling file paths, running external commands, synchronization primitives, to extracting EXIF data.

- **Global Variables**: 
  - `photoExtensions` and `videoExtensions`: Maps that contain commonly used photo and video file extensions.
  - `wg`: WaitGroup to wait for all goroutines to finish.
  - `mu`: Mutex for thread-safe increments of the counter.
  - `counter` and `totalFiles`: Keep track of progress.
  - `semaphore`: A buffered channel used to limit the number of concurrent goroutines.

- **Functions**:
  - `isPhotoOrVideo`: Checks if a given filename has a photo or video extension.
  - `getDateTaken`: Extracts the date taken from a photo or video using EXIF data or the `exiftool` for videos.
  - `processFile`: Handles the logic for checking if a file is already sorted, and if not, moves it to the correct directory.
  - `printProgress`: Displays a simple progress bar in the console.
  - `main`: Main entry point which prints out a welcome message, takes user input for the source directory, and initiates the file processing.

## Usage

1. Ensure `exiftool` is installed on your system.
2. Compile and run the Go script. You'll be prompted to provide the directory containing your photos and videos.
3. The script will organize your files into folders based on the "date taken" metadata.

## Open Source License

MIT License

Copyright (c) [2023] [Michael Romero]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
