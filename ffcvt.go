////////////////////////////////////////////////////////////////////////////
// Porgram: FfCvt
// Purpose: ffmpeg convert wrapper tool
// Authors: Tong Sun (c) 2015-2019, All rights reserved
////////////////////////////////////////////////////////////////////////////

/*

Transcodes all videos in the given directory and all of it's subdirectories
using ffmpeg.

*/

//go:generate sh -x ffcvt_cli.sh

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
	"regexp"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

const _encodedExt = "_.mkv"

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	version = "1.6.1"
	date    = "2020-01-08"

	sprintf           = fmt.Sprintf
	encodedExt string = _encodedExt
	totalOrg   int64  = 1
	totalNew   int64  = 1
	videos     []string
	workDirs   []string
)

////////////////////////////////////////////////////////////////////////////
// Main

func main() {
	flag.Usage = Usage
	flag.Parse()

	if Opts.PrintV {
		fmt.Fprintf(os.Stderr, "%s\nVersion %s built on %s\n", progname, version, date)
		os.Exit(0)
	}

	// One mandatory arguments, either -d or -f
	if len(Opts.Directory)+len(Opts.File) < 1 {
		Usage()
	}
	getDefault()

	encodedExt = Opts.Ext
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

		// The basename of the source directory will be created under the work
		//directory, which will become the new work directory
		if Opts.File == "" {
			Opts.WDirectory += string(os.PathSeparator) + filepath.Base(absd)
		}
		absw, _ := filepath.Abs(Opts.WDirectory)
		if absd == absw {
			log.Fatalf("[%s] Error: work directory\n\t\t  (%s)\n\t\t is the same as the source directory\n\t\t  (%s).", progname, absw, absd)
		}

		debug("Transcoding to "+Opts.WDirectory, 2)
		err := os.MkdirAll(Opts.WDirectory, os.ModePerm)
		checkError(err)
	} else {
		Opts.Par2C = false
	}

	startTime := time.Now()
	// transcoding
	if Opts.File != "" {
		fmt.Printf("\n== Transcoding: %s\n", Opts.File)
		transcodeFile(Opts.File)
	} else if Opts.Directory != "" {
		filepath.Walk(Opts.Directory, visit)
		transcodeVideos(startTime)
	}
	// par2 creating
	if Opts.Par2C {
		filepath.Walk(Opts.WDirectory, visitWDir)
		createPar2s(workDirs)
	}
	// reporting
	fmt.Printf("\nTranscoding completed in %s\n", time.Since(startTime))
	fmt.Printf("Org Size: %d MB\n", totalOrg/1024)
	fmt.Printf("New Size: %d MB\n", totalNew/1024)
	fmt.Printf("Saved:    %d%%\n",
		(totalOrg-totalNew)*100/totalOrg)
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
	if Opts.WDirectory == "" && fname[len(fname)-5:] == encodedExt {
		debug("Already-encoded file ignored: "+fname, 1)
		return
	}

	fext := strings.ToUpper(fname[len(fname)-4:])
	if strings.Index(Opts.Exts, fext) < 0 {
		debug("None-video file ignored: "+fname, 3)
		return
	}

	if !Opts.Links && isSymlink(fname) {
		debug("Skip symlink file: "+fname, 1)
		return
	}

	if Opts.NoClobber && fileExist(getOutputName(fname)) {
		debug("Encoded file exist for: "+fname, 1)
		return
	}

	videos = append(videos, fname)
}

func visitWDir(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		return nil
	}

	debug(path, 2)
	workDirs = append(workDirs, path)
	return nil
}

func createPar2s(workDirs []string) {
	fmt.Printf("\n== Creating par2 files\n\n")
	for ii, dir := range workDirs {
		if ii == 0 && len(workDirs) > 1 {
			// skip the root folder, if there are sub folders
			continue
		}
		os.Chdir(dir)
		dirName := filepath.Base(dir)

		cmd := []string{"par2create", "-u", "zz_" + dirName + ".par2", "*" + encodedExt}
		debug(strings.Join(cmd, " "), 1)

		out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		if err != nil {
			log.Printf("%s: Exec error - %s", progname, err.Error())
		}
		fmt.Printf("%s\n", out)
	}
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
			time.Duration(int64(float32(time.Since(startTime))*
				float32(videosTotal-videoNdx)/float32(videoNdx))))
	}
}

func transcodeFile(inputName string) {
	startTime := time.Now()
	outputName := getOutputName(inputName)
	debug(outputName, 4)
	os.MkdirAll(filepath.Dir(outputName), os.ModePerm)
	var oldAEP, oldVEP, oldSEP string
	oldEPUsed := false

	if !Opts.NoExec {
		// probe the file stream info first, only when not using -n
		fsinfo, err := probeFile(inputName)
		if err != nil {
			log.Printf("%s: Probe error - %s", progname, err.Error())
			return
		}
		debug(fsinfo, 4)
		// if there are more than one audio stream
		allAudioStreams := regexp.MustCompile(`Stream #0:.+: Audio: (.+)`).
			FindAllStringSubmatch(fsinfo, -1)
		if len(allAudioStreams) > 1 {
			// then find the designated audio stream language
			audioStreams := regexp.
				MustCompile(`Stream #(.+)\(` + Opts.Lang + `\): Audio: (.+)`).
				FindStringSubmatch(fsinfo)
			if len(audioStreams) >= 1 {
				// and use the 1st audio stream of the designated language
				// via *temporarily* tweaking AEP/VEP/SEP setting
				oldAEP, oldVEP, oldSEP = Opts.AEP, Opts.VEP, Opts.SEP
				oldEPUsed = true
				debug(audioStreams[1], 3)

				Opts.AEP += " -map " + audioStreams[1]
				dealSurroundSound(audioStreams[2])
				// when `-map` is used (for audio), then all else need mapping as well
				videoStreams := regexp.MustCompile(`Stream #(.+?)(\(.+\))*: Video: `).
					FindStringSubmatch(fsinfo)
				Opts.VEP += " -map " + videoStreams[1]
				subtitleStreams := regexp.MustCompile(`Stream #(.+)\(.+\): Subtitle: `).
					FindAllStringSubmatch(fsinfo, -1)
				// keep all subtitle streams
				for _, subtitleStream := range subtitleStreams {
					Opts.SEP += " -map " + subtitleStream[1]
				}
			}
			// else: designated audio language not found, use `default` instead
		} else {
			debug(inputName+" has single audio stream", 2)
			dealSurroundSound(allAudioStreams[0][1])
		}
	}

	args := []string{"-i", inputName}
	args = append(args, strings.Fields(Opts.OptExtra)...)
	args = encodeParametersS(encodeParametersA(encodeParametersV(args)))
	if Opts.Force {
		args = append(args, "-y")
	}
	args = append(args, flag.Args()...)
	args = append(args, outputName)
	debug(Opts.FFMpeg+" "+strings.Join(args, " "), 1)

	if Opts.NoExec {
		fmt.Printf("%s: to execute -\n  %s %s\n",
			progname, Opts.FFMpeg, strings.Join(args, " "))
	} else {
		//fmt.Printf("] %#v\n", args)
		cmd := exec.Command(Opts.FFMpeg, args...)
		var out, errOut bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errOut
		err := cmd.Run()
		if err != nil {
			log.Printf("%s: Exec error - %s\n\n%s", progname, err.Error(),
				string(errOut.Bytes()))
		}
		fmt.Printf("%s\n", out.String())
		timeTake := time.Since(startTime)

		if err != nil {
			fmt.Println("Failed.")

			// == remove zero-sized output file
			file, err := os.Open(outputName)
			checkError(err)

			// get the file size
			stat, err := file.Stat()
			file.Close()
			checkError(err)
			// fmt.Printf("Size of file '%s' is %d\n", outputName, stat.Size())
			if stat.Size() <= 500 {
				err := os.Remove(outputName)
				checkError(err)
			}
			debug("Failed output file '"+outputName+"' removed.", 1)

		} else {
			originalSize := fileSize(inputName)
			transcodedSize := fileSize(outputName)
			sizeDifference := originalSize - transcodedSize

			totalOrg += originalSize
			totalNew += transcodedSize

			fmt.Println("Done.")
			fmt.Printf("Org Size: %d KB\n", originalSize)
			fmt.Printf("New Size: %d KB\n", transcodedSize)
			fmt.Printf("Saved:    %d%% with %d KB\n",
				sizeDifference*100/originalSize, sizeDifference)
			fmt.Printf("Time: %v at %v\n\n", timeTake,
				time.Now().Format("2006-01-02 15:04:05"))
		}

		if oldEPUsed {
			// restored *temporarily* tweaked AEP/VEP/SEP setting
			Opts.AEP, Opts.VEP, Opts.SEP = oldAEP, oldVEP, oldSEP
		}
	}

	return
}

// dealSurroundSound will append to Opts.AEP proper setting to encode
// 5.1 surround sound channels
func dealSurroundSound(channelFeatures string) {
	if regexp.MustCompile(`, 5.1\(side\), `).MatchString(channelFeatures) {
		Opts.AEP += " -ac 2"
	}
}

func probeFile(inputName string) (string, error) {
	out := &bytes.Buffer{}

	cmdFFProbe := Opts.FFProbe + " " + Quote(inputName) + " 2>&1 | grep 'Stream #'"
	debug("Probing with "+cmdFFProbe, 2)
	cmd := exec.Command("sh", "-c", cmdFFProbe)
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	return string(out.Bytes()), err
}

// Returns the encode parameters for Subtitle
func encodeParametersS(args []string) []string {
	if Opts.SEP != "" {
		args = append(args, strings.Fields(Opts.SEP)...)
	}
	if Opts.SES != "" {
		args = append(args, strings.Fields(Opts.SES)...)
	}
	return args
}

// Returns the encode parameters for Audio
func encodeParametersA(args []string) []string {
	if Opts.AEP != "" {
		args = append(args, strings.Fields(Opts.AEP)...)
	}
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
	if Opts.VEP != "" {
		args = append(args, strings.Fields(Opts.VEP)...)
	}
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
		if Opts.VES == "libvpx-vp9" {
			// -b:v 0 -crf 37
			args = append(args, "-b:v", "0", "-crf", Opts.CRF)
		}
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

// Returns true if the file is symbolic link
func isSymlink(fname string) bool {
	fi, err := os.Lstat(fname)
	checkError(err)
	return fi.Mode()&os.ModeSymlink != 0
}

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
	//fmt.Printf("] (r, od, owd) %+v, %+v, %+v\n", r, Opts.Directory, Opts.WDirectory)
	if Opts.WDirectory != "" {
		// transcoding single file
		if Opts.File != "" {
			r = Opts.WDirectory + "/" + filepath.Base(input) + Opts.Suffix + encodedExt
		} else {
			r = strings.Replace(r, Opts.Directory, Opts.WDirectory, 1)
		}
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
