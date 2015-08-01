////////////////////////////////////////////////////////////////////////////
// Porgram: FfCvt
// Purpose: ffmpeg convert wrapper tool
// Authors: Tong Sun (c) 2015, All rights reserved
////////////////////////////////////////////////////////////////////////////

/*

Transcodes all episodes in the given directory and all of it's subdirectories
using ffmpeg.

Forked from Michael Murphy (mmstick)'s
https://gist.github.com/mmstick/3182c1c8596c1f830c7e

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
	getDefault()

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
	if Opts.Force {
		args = append(args, "-y")
	}
	args = append(args, flag.Args()...)
	args = append(args, outputName)
	debug(Opts.FFMpeg, 2)
	debug(strings.Join(args, " "), 1)

	cmd := exec.Command(Opts.FFMpeg, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("%s: Exec error - %s", progname, err.Error())
	}
	fmt.Printf("%s\n", out.String())
	time := time.Since(startTime)
	return time, err
}

// Returns the encode parameters for Audio
func encodeParametersA(args []string) []string {
	if Opts.AC {
		args = append(args, "-c:a", "copy")
		return args
	}
	if Opts.AN {
		args = append(args, "-an")
		return args
	}
	if Opts.A2Opus {
		Opts.AES = "libopus"
	}
	if Opts.AES != "" {
		args = append(args, "-c:a", Opts.AES)
	}
	if Opts.ABR != "" {
		args = append(args, "-b:a", Opts.ABR)
	}
	if Opts.AEA != "" {
		args = append(args, strings.Fields(Opts.AEA)...)
	}
	return args
}

// Returns the encode parameters for Video
func encodeParametersV(args []string) []string {
	if Opts.VC {
		args = append(args, "-c:v", "copy")
		return args
	}
	if Opts.VN {
		args = append(args, "-vn")
		return args
	}
	if Opts.V2X265 {
		Opts.VES = "libx265"
	}
	if Opts.VES != "" {
		args = append(args, "-c:v", Opts.VES)
	}
	if Opts.CRF != "" {
		if Opts.VES[:6] == "libx26" {
			args = append(args, "-"+Opts.VES[3:]+"-params", "crf="+Opts.CRF)
		}
	}
	if Opts.VEA != "" {
		args = append(args, strings.Fields(Opts.VEA)...)
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
	return input + Opts.Suffix + "_.mkv"
}

func debug(input string, threshold int) {
	if !(Opts.Debug >= threshold) {
		return
	}
	print("] ")
	print(input)
	print("\n")
}
