/*

Transcodes all videos in the given directory and all of it's subdirectories
using ffmpeg.

*/

package main

////////////////////////////////////////////////////////////////////////////
// Porgram: FfCvt
// Purpose: ffmpeg convert wrapper tool
// Authors: Tong Sun (c) 2015-2023, All rights reserved
////////////////////////////////////////////////////////////////////////////

//go:generate sh -x ffcvt_cli.sh

////////////////////////////////////////////////////////////////////////////
// Program start

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

// video encompass data around a single video
type video struct {
	name      string
	size, sum int64
	pct       int
}

// videoCol is the collection for the given set of videos
type videoCol struct {
	videos []video
	sum    int64
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	version = "1.9.0"
	date    = "2023-03-16"

	encodedExt string = _encodedExt
	totalOrg   int64  = 1
	totalNew   int64  = 1
	vidCol     videoCol
	workDirs   []string
	cutOps     string = ""
)

////////////////////////////////////////////////////////////////////////////
// Main

func main() {
	flag.Usage = Usage
	flag.Parse()

	if Opts.PrintV {
		fmt.Printf("%s version %s built on %s\n\n", progname, version, date)

		ver, err := getVersion(Opts.FFMpeg)
		if err != nil {
			log.Fatalf("%s: Version check error - %s\n", progname, err.Error())
		}
		fmt.Printf("%s\n\n", ver)

		// ffprobe, only use the first word
		ver, err = getVersion((strings.Fields(Opts.FFProbe))[0])
		if err != nil {
			log.Fatalf("%s: Version check error - %s\n", progname, err.Error())
		}
		fmt.Printf("%s\n", ver)

		os.Exit(0)
	}

	if len(Opts.Seg) > 0 {
		// sanity check
		_, err := time.Parse("15:04:05", Opts.Seg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Seg format error: '%s'\n", Opts.Seg)
			os.Exit(1)
		}
	}

	if len(Opts.Cut) > 0 {
		var b, vc strings.Builder
		var ci int
		var cv string
		d0, _ := time.Parse("15:04:05", "00:00:00")
		for ii, val := range Opts.Cut {
			ci = ii
			cv = val
			cRange := regexp.MustCompile(`\s*(.*?)\s*-\s*(\S*)\s*$`).
				FindStringSubmatch(val)
			if len(cRange) != 3 {
				fmt.Fprintf(os.Stderr, "pair - ")
				goto range_error
			}
			//fmt.Println("Cut range:", cRange[1], cRange[2])
			timeBgn, err := time.Parse("15:04:05", cRange[1])
			if err != nil {
				//fmt.Fprintf(os.Stderr, err.Error())
				fmt.Fprintf(os.Stderr, "start - ")
				goto range_error
			}
			endStr := ""
			if len(cRange[2]) > 0 {
				timeEnd, err := time.Parse("15:04:05", cRange[2])
				if err != nil {
					//fmt.Fprintf(os.Stderr, err.Error())
					fmt.Fprintf(os.Stderr, "end - ")
					goto range_error
				}
				endStr = fmt.Sprintf(":end=%d", int(timeEnd.Sub(d0).Seconds()))
			}
			secBgn := int(timeBgn.Sub(d0).Seconds())
			fmt.Fprintf(&b, "[0:v]trim=start=%d%s,setpts=PTS-STARTPTS[v%d];"+
				"[0:a]atrim=start=%d%s,asetpts=PTS-STARTPTS[a%d];",
				secBgn, endStr, ii, secBgn, endStr, ii)
			fmt.Fprintf(&vc, "[v%d][a%d]", ii, ii)
		}

		//fmt.Println("Cut(s):", len(Opts.Cut), Opts.Cut)
		//fmt.Println(vc.String())
		fmt.Fprintf(&b, "%sconcat=n=%d:v=1:a=1[vo][ao]", vc.String(), len(Opts.Cut))
		//fmt.Println(b.String())
		cutOps = b.String()
		goto cut_ok
	range_error:
		fmt.Fprintf(os.Stderr, "Cut range %d format error for '%s'\n", ci, cv)
		os.Exit(1)
	}
cut_ok:

	// One mandatory arguments, either -d or -f
	if len(Opts.Directory)+len(Opts.File) < 1 {
		Usage()
	}
	getDefault()
	//fmt.Fprintf(os.Stderr, "Defaults: '%+v'\n", Defaults)

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
		// calculate (finished) percentage for each video
		n := len(vidCol.videos)
		if n > 0 {
			totalSize := vidCol.videos[n-1].sum
			for i, _ := range vidCol.videos {
				vidCol.videos[i].pct = int(vidCol.videos[i].sum * 100 / totalSize)
			}
		}
		//fmt.Printf("%v", vidCol)
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
// Methods definitions

// Append the video file to the list, unless it's encoded already
// and duplicate none-video files to dest dir as well
func (vidCol *videoCol) append(path string, size int64) {
	fname := Opts.Directory + string(os.PathSeparator) + path

	if Opts.WDirectory == "" && fname[len(fname)-5:] == encodedExt {
		debug("Already-encoded file ignored: "+fname, 1)
		return
	}

	fext := strings.ToUpper(fname[len(fname)-4:])
	if strings.Index(Opts.Exts, fext) < 0 {
		// None-video files, dup to dest, hardlink 1st else copy
		if Opts.WDirectory != "" {
			src := Opts.Directory + "/" + fname
			dst := Opts.WDirectory + "/" + fname
			err := linkFile(src, dst)
			if err != nil {
				copyFile(src, dst)
			}
		}
		debug("None-video file '"+fname+"' duplicated to dest dir.", 1)
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

	vidCol.sum += size
	vidCol.videos = append(vidCol.videos,
		video{name: fname, size: size, sum: vidCol.sum})

}

////////////////////////////////////////////////////////////////////////////
// Function definitions

//==========================================================================
// Directory & files handling

// visit will visit the source directory and queue videos found there to vidCol
func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}

	vidCol.append(path, f.Size())
	return nil
}

// visitWDir will visit the work dir and add all directories there to workDirs
func visitWDir(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		return nil
	}

	debug(path, 2)
	workDirs = append(workDirs, path)
	return nil
}

// createPar2s will create par2s files for each dir in workDirs
func createPar2s(workDirs []string) {
	fmt.Printf("\n== Creating par2 files\n\n")
	for ii, dir := range workDirs {
		if ii == 0 && len(workDirs) > 1 {
			// skip the root folder, if there are sub folders
			continue
		}
		os.Chdir(dir)
		dirName := filepath.Base(dir)

		// create par2s for all files within dest dir
		cmd := []string{"par2create", "-u", "zz_" + dirName + ".par2", "*"}
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
	videosTotal := len(vidCol.videos)
	for i, v := range vidCol.videos {
		inputName := v.name
		videoNdx := i + 1

		if Opts.NoClobber && fileExist(getOutputName(inputName)) {
			debug("Encoded file exist for: "+inputName, 1)
			continue
		}

		fmt.Printf("\n== Transcoding [%d/%d] (%d%%): '%s'\n   under %s\n",
			videoNdx, videosTotal, v.pct, filepath.Base(inputName), filepath.Dir(inputName))
		transcodeFile(inputName)
		fmt.Printf("Time taken so far %s\n", time.Since(startTime))
		fmt.Printf("Finishing the remaining %d%% in %s\n",
			100-v.pct,
			time.Duration(int64(float32(time.Since(startTime))*
				float32(100-v.pct)/float32(v.pct))))
	}
}

func transcodeFile(inputName string) {
	startTime := time.Now()
	outputName, outputGrpName := getOutputName(inputName), ""
	if len(Opts.Seg) > 0 {
		outputName, outputGrpName = getOutputNameSeg(inputName)
	}
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
		/*

		   Cases when `-map`s are necessary

		   - more than one audio stream, and we pick eng stream only
		   - more than one subtitle stream, and
		     * output all subtitle streams (default, no -sel), or
		     * pick specific subtitle stream(s) via -sel

		*/
		allAudioStreams := regexp.MustCompile(`Stream #0:.+: Audio: (.+)`).
			FindAllStringSubmatch(fsinfo, -1)
		if len(allAudioStreams) > 1 ||
			len(regexp.MustCompile(`Stream #0:.+: Subtitle: (.+)`).
				FindAllStringSubmatch(fsinfo, -1)) > 1 {
			// then use the designated audio stream language
			// via *temporarily* using the AEP/VEP/SEP setting
			oldAEP, oldVEP, oldSEP = Opts.AEP, Opts.VEP, Opts.SEP
			oldEPUsed = true
			Opts.VEP += " -map 0:v"
			Opts.AEP += " -map 0:a:m:language:" + Opts.Lang
			//log.Printf("%s: Opts.SEL - %#v", progname, Opts.SEL)
			if len(Opts.SEL) == 0 {
				Opts.SEP += " -map 0:s"
			} else {
				for _, val := range Opts.SEL {
					Opts.SEP += " -map 0:s:m:language:" + val
				}
			}
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
	if len(Opts.Seg) > 0 {
		args = append(args, "-f")
		args = append(args, "segment")
		args = append(args, "-segment_time")
		args = append(args, Opts.Seg)
		args = append(args, "-reset_timestamps")
		args = append(args, "1")
	}
	if len(Opts.Speed) > 0 {
		args = append(args, "-filter_complex")
		args = append(args, fmt.Sprintf("[0:v]setpts=PTS/%s[v];[0:a]atempo=%s[a]",
			Opts.Speed, Opts.Speed))
		args = append(args, "-map")
		args = append(args, "[v]")
		args = append(args, "-map")
		args = append(args, "[a]")
	}
	if len(cutOps) != 0 {
		args = append(args, "-filter_complex")
		args = append(args, cutOps)
		args = append(args, "-map")
		args = append(args, "[vo]")
		args = append(args, "-map")
		args = append(args, "[ao]")
	}
	args = append(args, flag.Args()...)
	if len(Opts.Seg) > 0 {
		args = append(args, outputGrpName)
	} else {
		args = append(args, outputName)
	}
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

func getVersion(execName string) (string, error) {
	out := &bytes.Buffer{}

	cmd := exec.Command(execName, "-version")
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	// OK, then only return the first two lines
	lines := strings.Split(string(out.Bytes()), "\n")
	return strings.Join(lines[:2], "\n"), nil
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
		if len(Opts.VES) > 6 && Opts.VES[:6] == "libx26" {
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

// getOutputNameSeg will do getOutputName() but tailored toward video segmenting.
// The first return will be the first video segment file name, 00_.mkv, while
// the second return will be the segment group file name, %02d_.mkv
func getOutputNameSeg(input string) (string, string) {
	r := getOutputName(input)
	fileFirst := strings.Replace(r, encodedExt, "00"+encodedExt, 1)
	fileGroup := strings.Replace(r, encodedExt, "%02d"+encodedExt, 1)
	return fileFirst, fileGroup
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
