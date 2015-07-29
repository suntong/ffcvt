////////////////////////////////////////////////////////////////////////////
// Porgram: FfCvt
// Purpose: ffmpeg convert wrapper tool
// Authors: Tong Sun (c) 2015, All rights reserved
////////////////////////////////////////////////////////////////////////////

/*

Transcodes all episodes in the given directory and all of it's subdirectories
using ffmpeg.

Initial version based (but also heavily hacked) on
https://gist.github.com/mmstick/3182c1c8596c1f830c7e
by Michael Murphy (mmstick)

*/

////////////////////////////////////////////////////////////////////////////
// Program start

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

const (
	STATIC_PARAMS = "me=star:subme=7:bframes=16:b-adapt=2:ref=16:rc-lookahead=60:max-merge=5:tu-intra-depth=4:tu-inter-depth=4"
)

// Contains information about each episode
type Episode struct {
	name           string
	directory      string
	originalSize   int64
	transcodedSize int64
	sizeDifference int64
	time           time.Duration
	stat           error
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	sprintf = fmt.Sprintf
)

////////////////////////////////////////////////////////////////////////////
// Main

func main() {
	flag.Usage = Usage
	flag.Parse()

	// One mandatory arguments, either -d or -f
	if len(Opts.Directory)+len(Opts.File) < 1 {
		Usage()
	}

	startTime := time.Now()
	if Opts.Directory != "" {
		transcodeEpisodes(scanEpisodes(scanDirectory(Opts.Directory), Opts.Directory))
	} else if Opts.File != "" {
		outputName := getOutputName(Opts.File)
		fmt.Printf("\n== Transcoding: %s\n", Opts.File)
		transcodeFile(Opts.File, outputName)
	}
	fmt.Printf("Transcoding completed in %s\n", time.Since(startTime))
}

////////////////////////////////////////////////////////////////////////////
// Function definitions

//==========================================================================
// Directory & Episodes handling

// Transcodes all episodes in the episode list
func transcodeEpisodes(episodeList *[]Episode) {
	files := len(*episodeList)
	for index, ep := range *episodeList {
		ep.transcodeEpisode(index+1, files)
		ep.status()
	}
}

// Print the status of the transcoded episode
func (episode Episode) status() {
	if episode.stat != nil {
		fmt.Println("Failed to transcode", episode.name)
	} else {
		fmt.Println("Transcoded", episode.name)
		fmt.Println("Original Size:", episode.originalSize, "MB")
		fmt.Println("New Size:", episode.transcodedSize, "MB")
		fmt.Println("Difference:", episode.sizeDifference, "MB")
		fmt.Println("Time:", episode.time)
	}
}

// Recurse through each subdirectory and adds each episode to the episode list
func scanEpisodes(directoryList []os.FileInfo, directory string) *[]Episode {
	list := []Episode{}
	for _, file := range directoryList {
		if file.IsDir() {
			recurseDirectory(&directory, file.Name())
		} else {
			appendEpisode(&list, file, &directory)
		}
	}
	return &list
}

// Returns a list of files in the current directory
func scanDirectory(path string) []os.FileInfo {
	directory, _ := ioutil.ReadDir(path)
	return directory
}

// If the file is a directory, recurse through the directory.
func recurseDirectory(directory *string, filename string) {
	subdirectory := sprintf("%s/%s", *directory, filename)
	scanEpisodes(scanDirectory(subdirectory), subdirectory)
}

// Append the current episode to the episode list, unless it's encoded already
func appendEpisode(list *[]Episode, file os.FileInfo, directory *string) {
	fname := file.Name()
	if fname[len(fname)-5:] == "_.mkv" {
		return
	}

	*list = append(*list, Episode{
		name:         fname,
		directory:    *directory,
		originalSize: file.Size() / 1000000,
	})
}

//==========================================================================
// Transcode handling

// Transcode the current episode
func (ep Episode) transcodeEpisode(index, files int) {
	inputName := sprintf("%s/%s", ep.directory, ep.name)
	outputName := getOutputName(inputName)
	fmt.Printf("\n== Transcoding [%d/%d]: %s\n", index, files, ep.name)
	ep.time, ep.stat = transcodeFile(inputName, outputName)
	ep.transcodedSize = transcodeSize(outputName)
	ep.sizeDifference = ep.originalSize - ep.transcodedSize
}

func transcodeFile(inputName, outputName string) (time.Duration, error) {
	startTime := time.Now()

	args := encodeParametersV(encodeParametersA(
		[]string{"-i", inputName}))
	args = append(args, os.Args...)
	args = append(args, outputName)
	debug(Opts.FFMpeg)
	debug(strings.Join(args, " "))

	cmd := exec.Command(Opts.FFMpeg, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("%s: Exec error - %s", progname, err.Error())
	}
	//fmt.Printf("\n== Out:\n%s\n", out.String())
	time := time.Since(startTime)
	return time, err
}

// Returns the encode parameters for Audio
func encodeParametersA(args []string) []string {
	if Opts.AC {
		args = append(args, "-c:a", "copy")
		return args
	}
	return args
}

// Returns the encode parameters for Video
func encodeParametersV(args []string) []string {
	if Opts.VC {
		args = append(args, "-c:v", "copy")
		return args
	}
	return args
}

// Returns the size of the newly transcoded episode, if it exists.
func transcodeSize(transcodedEpisode string) int64 {
	file, err := os.Open(transcodedEpisode)
	if err == nil {
		stat, _ := file.Stat()
		return stat.Size() / 1000000
	} else {
		log.Printf("%s: Open error - %s", progname, err.Error())
		return 0
	}
}

// Replaces the file extension from the input string with _.mkv
func getOutputName(input string) string {
	for index := len(input) - 1; index >= 0; index-- {
		if input[index] == '.' {
			input = input[:index]
		}
	}
	return input + "_.mkv"
}

func debug(input string) {
	if Opts.Debug == 0 {
		return
	}
	print("] ")
	print(input)
	print("\n")
}
