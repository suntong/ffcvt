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
	Target    string // target type: x265-opus/x264-mp3
	Directory string // directory that hold input files
	File      string // input file name (either -d or -f must be specified)
	Base      string // used as basename for output files
	AC        bool   // copy audio codec
	VC        bool   // copy video codec
	VSS       bool   // video: same size
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
		"target type: x265-opus/x264-mp3")
	flag.StringVar(&Opts.Directory, "d", "",
		"directory that hold input files")
	flag.StringVar(&Opts.File, "f", "",
		"input file name (either -d or -f must be specified)")
	flag.StringVar(&Opts.Base, "base", "",
		"used as basename for output files")

	flag.BoolVar(&Opts.AC, "ac", false,
		"copy audio codec")
	flag.BoolVar(&Opts.VC, "vc", false,
		"copy video codec")
	flag.BoolVar(&Opts.VSS, "vss", true,
		"video: same size")
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
	if len(Opts.Base) == 0 ||
		len(os.Getenv("FFCVT_BASE")) != 0 {
		Opts.Base = os.Getenv("FFCVT_BASE")
	}

	if len(Opts.FFMpeg) == 0 ||
		len(os.Getenv("FFCVT_FFMPEG")) != 0 {
		Opts.FFMpeg = os.Getenv("FFCVT_FFMPEG")
	}

}

const USAGE_SUMMARY = "  -aes\taudio encoding method set\n  -ves\tvideo encoding method set\n  -aea\taudio encoding method append\n  -vea\tvideo encoding method append\n  -abr\taudio bitrate (64k for opus, 256k for mp3)\n  -crf\tthe CRF value: 0-51. Higher CRF gives lower quality\n\t (28 for x265, ~ 23 for x264)\n\n  -t\ttarget type: x265-opus/x264-mp3\n  -d\tdirectory that hold input files\n  -f\tinput file name (either -d or -f must be specified)\n  -base\tused as basename for output files\n\n  -ac\tcopy audio codec\n  -vc\tcopy video codec\n  -vss\tvideo: same size\n  -ato-opus\taudio encode to opus, using -abr\n  -vto-x265\tvideo video encode to x265, using -crf\n\n  -force\toverwrite any existing none-empty file\n  -debug\tdebugging level\n  -ffmpeg\tffmpeg program executable name\n\nDetails:\n\n"

// The Usage function shows help on commandline usage
func Usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage:\n %s [flags] \n\nFlags:\n\n",
		progname)
	fmt.Fprintf(os.Stderr, USAGE_SUMMARY)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		"\n \n\t \n")
	os.Exit(0)
}
