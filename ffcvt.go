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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

const (
	CRF_HELP      = "Set the CRF value: 0-51. Higher CRF gives lower quality."
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
	crf     = flag.Int("crf", 24, CRF_HELP)
	sprintf = fmt.Sprintf
)

////////////////////////////////////////////////////////////////////////////
// Main

func main() {
	flag.Parse()
	startTime := time.Now()
	directory, _ := os.Getwd()
	transcodeEpisodes(scanEpisodes(scanDirectory(directory), directory))
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
		ep.transcode(index+1, files)
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

// Append the current episode to the episode list.
func appendEpisode(list *[]Episode, file os.FileInfo, directory *string) {
	*list = append(*list, Episode{
		name:         file.Name(),
		directory:    *directory,
		originalSize: file.Size() / 1000000,
	})
}

//==========================================================================
// Transcode handling

// Transcode the current episode
func (ep Episode) transcode(index, files int) {
	inputName, outputName := getFileNames(&ep.directory, &ep.name)
	fmt.Printf("Transcoding %d/%d: %s\n", index, files, ep.name)
	startTime := time.Now()
	ep.stat = encodeParameters(&inputName, &outputName).Run()
	ep.time = time.Since(startTime)
	ep.transcodedSize = transcodeSize(&outputName)
	ep.sizeDifference = ep.originalSize - ep.transcodedSize
}

// Returns the encode parameters for the episode
func encodeParameters(inputName, outputName *string) *exec.Cmd {
	x265Parameters := sprintf("crf=%d:%s", *crf, STATIC_PARAMS)
	return exec.Command("ffmpeg", "-i", *inputName, "-c:a", "libopus",
		"-c:v", "libx265", "-x265-params",
		x265Parameters, *outputName)
}

// Returns the size of the newly transcoded episode, if it exists.
func transcodeSize(transcodedEpisode *string) int64 {
	file, err := os.Open(*transcodedEpisode)
	if err == nil {
		stat, _ := file.Stat()
		return stat.Size() / 1000000
	} else {
		return 0
	}
}

// Returns the input name and output name for the file to be encoded
func getFileNames(directory, name *string) (string, string) {
	inputName := sprintf("%s/%s", *directory, *name)
	outputName := convertToMKV(sprintf("%s/encoded_%s", *directory, *name))
	return inputName, outputName
}

// Replaces the file extension from the input string and changes it to mkv
func convertToMKV(input string) string {
	for index := len(input) - 1; index >= 0; index-- {
		if input[index] == '.' {
			input = input[:index]
		}
	}
	return input + ".mkv"
}
