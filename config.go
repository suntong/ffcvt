// !!! !!!
// WARNING: Code automatically generated. Editing discouraged.
// !!! !!!

package main

import (
	"flag"
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

const progname = "ffcvt" // os.Args[0]

// The Options struct defines the structure to hold the commandline values
type Options struct {
	Encoding         // anonymous field to hold encoding values
	Target    string // target type: x265-opus/x264-mp3/youtube
	Directory string // directory that hold input files
	File      string // input file name (either -d or -f must be specified)
	Suffix    string // suffix to the output file names
	AC        bool   // copy audio codec
	VC        bool   // copy video codec
	AN        bool   // no audio, output video only
	VN        bool   // no video, output audio only
	VSS       bool   // video: same size
	OptExtra  string // more options that will pass to ffmpeg program
	A2Opus    bool   // audio encode to opus, using -abr
	V2X265    bool   // video video encode to x265, using -crf
	Force     bool   // overwrite any existing none-empty file
	Debug     int    // debugging level
	FFMpeg    string // ffmpeg program executable name
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

// Opts holds the actual values from the command line paramters
var Opts Options

////////////////////////////////////////////////////////////////////////////
// Commandline definitions

func init() {

	// set default values for command line paramters
	flag.StringVar(&Opts.AES, "aes", "",
		"audio encoding method set")
	flag.StringVar(&Opts.VES, "ves", "",
		"video encoding method set")
	flag.StringVar(&Opts.AEA, "aea", "",
		"audio encoding method append")
	flag.StringVar(&Opts.VEA, "vea", "",
		"video encoding method append")
	flag.StringVar(&Opts.ABR, "abr", "",
		"audio bitrate (64k for opus, 256k for mp3)")
	flag.StringVar(&Opts.CRF, "crf", "",
		"the CRF value: 0-51. Higher CRF gives lower quality\n\t (28 for x265, ~ 23 for x264)")

	flag.StringVar(&Opts.Target, "t", "x265-opus",
		"target type: x265-opus/x264-mp3/youtube")
	flag.StringVar(&Opts.Directory, "d", "",
		"directory that hold input files")
	flag.StringVar(&Opts.File, "f", "",
		"input file name (either -d or -f must be specified)")
	flag.StringVar(&Opts.Suffix, "suf", "",
		"suffix to the output file names")

	flag.BoolVar(&Opts.AC, "ac", false,
		"copy audio codec")
	flag.BoolVar(&Opts.VC, "vc", false,
		"copy video codec")
	flag.BoolVar(&Opts.AN, "an", false,
		"no audio, output video only")
	flag.BoolVar(&Opts.VN, "vn", false,
		"no video, output audio only")
	flag.BoolVar(&Opts.VSS, "vss", true,
		"video: same size")
	flag.StringVar(&Opts.OptExtra, "o", "",
		"more options that will pass to ffmpeg program")
	flag.BoolVar(&Opts.A2Opus, "ato-opus", false,
		"audio encode to opus, using -abr")
	flag.BoolVar(&Opts.V2X265, "vto-x265", false,
		"video video encode to x265, using -crf")

	flag.BoolVar(&Opts.Force, "force", false,
		"overwrite any existing none-empty file")
	flag.IntVar(&Opts.Debug, "debug", 0,
		"debugging level")
	flag.StringVar(&Opts.FFMpeg, "ffmpeg", "ffmpeg",
		"ffmpeg program executable name")

	// Now override those default values from environment variables
	if len(Opts.AES) == 0 ||
		len(os.Getenv("FFCVT_AES")) != 0 {
		Opts.AES = os.Getenv("FFCVT_AES")
	}
	if len(Opts.VES) == 0 ||
		len(os.Getenv("FFCVT_VES")) != 0 {
		Opts.VES = os.Getenv("FFCVT_VES")
	}
	if len(Opts.AEA) == 0 ||
		len(os.Getenv("FFCVT_AEA")) != 0 {
		Opts.AEA = os.Getenv("FFCVT_AEA")
	}
	if len(Opts.VEA) == 0 ||
		len(os.Getenv("FFCVT_VEA")) != 0 {
		Opts.VEA = os.Getenv("FFCVT_VEA")
	}
	if len(Opts.ABR) == 0 ||
		len(os.Getenv("FFCVT_ABR")) != 0 {
		Opts.ABR = os.Getenv("FFCVT_ABR")
	}
	if len(Opts.CRF) == 0 ||
		len(os.Getenv("FFCVT_CRF")) != 0 {
		Opts.CRF = os.Getenv("FFCVT_CRF")
	}

	if len(Opts.Target) == 0 ||
		len(os.Getenv("FFCVT_T")) != 0 {
		Opts.Target = os.Getenv("FFCVT_T")
	}
	if len(Opts.Directory) == 0 ||
		len(os.Getenv("FFCVT_D")) != 0 {
		Opts.Directory = os.Getenv("FFCVT_D")
	}
	if len(Opts.File) == 0 ||
		len(os.Getenv("FFCVT_F")) != 0 {
		Opts.File = os.Getenv("FFCVT_F")
	}
	if len(Opts.Suffix) == 0 ||
		len(os.Getenv("FFCVT_SUF")) != 0 {
		Opts.Suffix = os.Getenv("FFCVT_SUF")
	}

	if len(Opts.OptExtra) == 0 ||
		len(os.Getenv("FFCVT_O")) != 0 {
		Opts.OptExtra = os.Getenv("FFCVT_O")
	}

	if len(Opts.FFMpeg) == 0 ||
		len(os.Getenv("FFCVT_FFMPEG")) != 0 {
		Opts.FFMpeg = os.Getenv("FFCVT_FFMPEG")
	}

}

const USAGE_SUMMARY = "  -aes\taudio encoding method set (FFCVT_AES)\n  -ves\tvideo encoding method set (FFCVT_VES)\n  -aea\taudio encoding method append (FFCVT_AEA)\n  -vea\tvideo encoding method append (FFCVT_VEA)\n  -abr\taudio bitrate (64k for opus, 256k for mp3) (FFCVT_ABR)\n  -crf\tthe CRF value: 0-51. Higher CRF gives lower quality\n\t (28 for x265, ~ 23 for x264) (FFCVT_CRF)\n\n  -t\ttarget type: x265-opus/x264-mp3/youtube (FFCVT_T)\n  -d\tdirectory that hold input files (FFCVT_D)\n  -f\tinput file name (either -d or -f must be specified) (FFCVT_F)\n  -suf\tsuffix to the output file names (FFCVT_SUF)\n\n  -ac\tcopy audio codec (FFCVT_AC)\n  -vc\tcopy video codec (FFCVT_VC)\n  -an\tno audio, output video only (FFCVT_AN)\n  -vn\tno video, output audio only (FFCVT_VN)\n  -vss\tvideo: same size (FFCVT_VSS)\n  -o\tmore options that will pass to ffmpeg program (FFCVT_O)\n  -ato-opus\taudio encode to opus, using -abr (FFCVT_ATO_OPUS)\n  -vto-x265\tvideo video encode to x265, using -crf (FFCVT_VTO_X265)\n\n  -force\toverwrite any existing none-empty file (FFCVT_FORCE)\n  -debug\tdebugging level (FFCVT_DEBUG)\n  -ffmpeg\tffmpeg program executable name (FFCVT_FFMPEG)\n\nDetails:\n\n"

// The Usage function shows help on commandline usage
func Usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage:\n %s [flags] \n\nFlags:\n\n",
		progname)
	fmt.Fprintf(os.Stderr, USAGE_SUMMARY)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		"\nThe `ffcvt -f testf.mp4 -debug 1 -force` will invoke\n\n  ffmpeg -i testf.mp4 -c:a libopus -b:a 64k -c:v libx265 -x265-params crf=28 -y testf_.mkv\n\nTo use `preset`, do the following or set it in env var FFCVT_O\n\n  cm=medium\n  ffcvt -f testf.mp4 -debug 1 -force -suf $cm -- -preset $cm\n\nWhich will invoke\n\n  ffmpeg -i testf.mp4 -c:a libopus -b:a 64k -c:v libx265 -x265-params crf=28 -y -preset medium testf_medium_.mkv\n\nHere are the final sizes and the conversion time (in seconds):\n\n  2916841  testf.mp4\n  1807513  testf_.mkv\n  1743701  testf_veryfast_.mkv   41\n  2111667  testf_faster_.mkv     44\n  1793216  testf_fast_.mkv       85\n  1807513  testf_medium_.mkv    120\n  1628502  testf_slow_.mkv      366\n  1521889  testf_slower_.mkv    964\n  1531154  testf_veryslow_.mkv 1413\n\nI.e., if `preset` is not used, the default is `medium`.\n\nHere is another set of results, sizes and the conversion time (in minutes):\n\n 171019470  testf.avi\n  55114663  testf_veryfast_.mkv  39.2\n  57287586  testf_faster_.mkv    51.07\n  52950504  testf_fast_.mkv     147.11\n  55641838  testf_medium_.mkv   174.25\n\nSame source file, using the fixed `-preset fast`, altering the crf:\n\n  52950504  testf_28_.mkv       147.11\n  43480573  testf_30_.mkv       146.5\n  36609186  testf_32_.mkv       144.5\n  31427912  testf_34_.mkv       143.9\n  27397348  testf_36_.mkv       139.33\n\nSo it confirms that `-preset` determines the conversion time,\nwhile crf controls the final file size, not conversion time.\n")
	os.Exit(0)
}
