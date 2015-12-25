////////////////////////////////////////////////////////////////////////////
// Porgram: FfCvt
// Purpose: ffmpeg convert wrapper tool
// Authors: Tong Sun (c) 2015, All rights reserved
////////////////////////////////////////////////////////////////////////////

/*

Transcodes all videos in the given directory and all of it's subdirectories
using ffmpeg.

*/

////////////////////////////////////////////////////////////////////////////
// Program start

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

const _encodedExt = "_.mkv"

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	sprintf           = fmt.Sprintf
	encodedExt string = _encodedExt
	total_org  int64  = 1
	total_new  int64  = 1
	videos     []string
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

	// Sanity check
	if Opts.WDirectory != "" {
		// To error on the safe side -- when -d is not given but -f is,
		// path.Clean(Opts.Directory) will return ".", thus forcing
		// the work directory cannot be the same as pwd
		// because the encodedExt might be conflicting with the source file
		encodedExt = encodedExt[1:] // now ".mkv"
		Opts.Directory = path.Clean(Opts.Directory)
		Opts.WDirectory = path.Clean(Opts.WDirectory)
		absd, _ := filepath.Abs(Opts.Directory)
		absw, _ := filepath.Abs(Opts.WDirectory)
		if absd == absw {
			log.Fatalf("[%s] Error: work directory (%s) cannot be the same\n\t\tas the source directory (%s).", progname, absw, absd)
		}

		// The basename of the source directory will be created under the work
		//directory, which will become the new work directory
		Opts.WDirectory += string(os.PathSeparator) + filepath.Base(absd)
		os.Mkdir(Opts.WDirectory, os.ModePerm)
		//debug(Opts.WDirectory, 2)
	} else {
		Opts.Par2C = false
	}

	startTime := time.Now()
	if Opts.Directory != "" {
		filepath.Walk(Opts.Directory, visit)
		transcodeVideos(startTime)
	} else if Opts.File != "" {
		fmt.Printf("\n== Transcoding: %s\n", Opts.File)
		transcodeFile(Opts.File)
	}
	fmt.Printf("\nTranscoding completed in %s\n", time.Since(startTime))
	fmt.Printf("Org Size: %d MB\n", total_org/1024)
	fmt.Printf("New Size: %d MB\n", total_new/1024)
	fmt.Printf("Saved:    %d%%\n",
		(total_org-total_new)*100/total_org)
}

////////////////////////////////////////////////////////////////////////////
// Function definitions

//==========================================================================
// Directory & files handling

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}

	appendVideo(Opts.Directory + string(os.PathSeparator) + path)
	return nil
}

// Append the video file to the list, unless it's encoded already
func appendVideo(fname string) {
	if fname[len(fname)-5:] == _encodedExt {
		return
	}

	fext := strings.ToUpper(fname[len(fname)-4:])
	if strings.Index(Opts.Exts, fext) < 0 {
		return
	}

	if Opts.NoClobber && fileExist(getOutputName(fname)) {
		return
	}

	videos = append(videos, fname)
}

//==========================================================================
// Transcode handling

// Transcode videos in the global videos array
func transcodeVideos(startTime time.Time) {
	videosTotal := len(videos)
	for i, inputName := range videos {
		videoNdx := i + 1
		fmt.Printf("\n== Transcoding [%d/%d]: '%s'\n   under %s\n",
			videoNdx, videosTotal, filepath.Base(inputName), filepath.Dir(inputName))
		transcodeFile(inputName)
		fmt.Printf("Time taken so far %s\n", time.Since(startTime))
		fmt.Printf("Finishing the remaining %d%% in %s\n",
			(videosTotal-videoNdx)*100/videosTotal,
			time.Duration(int(float32(time.Since(startTime))*
				float32(videosTotal-videoNdx)/float32(videoNdx))))
	}
}

func transcodeFile(inputName string) {
	startTime := time.Now()
	outputName := getOutputName(inputName)

	args := encodeParametersV(encodeParametersA(
		[]string{"-i", inputName}))
	if Opts.Force {
		args = append(args, "-y")
	}
	args = append(args, strings.Fields(Opts.OptExtra)...)
	args = append(args, flag.Args()...)
	args = append(args, outputName)
	debug(Opts.FFMpeg, 2)
	debug(strings.Join(args, " "), 1)

	if Opts.NoExec {
		fmt.Printf("%s: to execute -\n  %s %s\n",
			progname, Opts.FFMpeg, strings.Join(args, " "))
	} else {
		cmd := exec.Command(Opts.FFMpeg, args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Printf("%s: Exec error - %s", progname, err.Error())
		}
		fmt.Printf("%s\n", out.String())
		timeTake := time.Since(startTime)

		if err != nil {
			fmt.Println("Failed.")
		} else {
			originalSize := fileSize(inputName)
			transcodedSize := fileSize(outputName)
			sizeDifference := originalSize - transcodedSize

			total_org += originalSize
			total_new += transcodedSize

			fmt.Println("Done.")
			fmt.Printf("Org Size: %d KB\n", originalSize)
			fmt.Printf("New Size: %d KB\n", transcodedSize)
			fmt.Printf("Saved:    %d%% with %d KB\n",
				sizeDifference*100/originalSize, sizeDifference)
			fmt.Printf("Time: %v at %v\n\n", timeTake,
				time.Now().Format("2006-01-02 15:04:05"))
		}
	}

	return
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

//==========================================================================
// Dealing with Files

// Returns true if the file exist
func fileExist(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}

// Returns the file size
func fileSize(fname string) int64 {
	stat, err := os.Stat(fname)
	checkError(err)

	return stat.Size() / 1024
}

// Replaces the file extension from the input string with _.mkv, and optionally
// Opts.Suffix as well. If "-w" is defined, use it for output name.
func getOutputName(input string) string {
	index := strings.LastIndex(input, ".")
	if index > 0 {
		input = input[:index]
	}
	r := input + Opts.Suffix + encodedExt
	if Opts.WDirectory != "" {
		r = strings.Replace(r, Opts.Directory, Opts.WDirectory, 1)
	}
	return r
}

//==========================================================================
// Support functions

func debug(input string, threshold int) {
	if !(Opts.Debug >= threshold) {
		return
	}
	print("] ")
	print(input)
	print("\n")
}

func checkError(err error) {
	if err != nil {
		log.Printf("%s: Fatal error - %s", progname, err.Error())
		os.Exit(1)
	}
}
