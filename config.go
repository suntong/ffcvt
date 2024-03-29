// !!! !!!
// WARNING: Code automatically generated. Editing discouraged.
// !!! !!!

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

const progname = "ffcvt" // os.Args[0]

// The Options struct defines the structure to hold the commandline values
type Options struct {
	Cfg        string        // cfg file to define your own targets: webm/wx/youtube etc
	Target     string        // target type: webm/x265-opus/x264-mp3/wx/youtube/copy, or empty
	Encoding                 // anonymous field to hold encoding values
	Directory  string        // directory that hold input files
	File       string        // input file name (either -d or -f must be specified)
	Links      bool          // symlinks will be processed as well
	Exts       string        // extension list for all the files to be queued
	Suffix     string        // suffix to the output file names
	WDirectory string        // work directory that hold output files
	AC         bool          // copy audio codec
	VC         bool          // copy video codec
	AN         bool          // no audio, output video only
	VN         bool          // no video, output audio only
	VSS        bool          // video: same size
	Cut        mFlags        // Cut segment(s) out to keep. Specify in the form of start-[end],\n\tstrictly in the format of hh:mm:ss, and may repeat
	Seg        string        // Split video into multiple segments (strictly in format: hh:mm:ss)
	Speed      string        // Speed up/down video playback speed (e.g. 1.28)
	Karaoke    bool          // Add a karaoke audio track to .mp4 MTV
	TranspFrom string        // Transpose song's key from (e.g. C/C#/Db/D etc)
	TranspTo   string        // Transpose song's key to (e.g. -tkf C -tkt Db)
	TranspBy   int           // Transpose song by (e.g. +2, -3, etc) chromatic scale
	Lang       string        // language selection for audio stream extraction
	SEL        mFlags        // subtitle encoding language (language picked for reencoded video)
	OptExtra   string        // more options that will pass to ffmpeg program
	A2Opus     bool          // audio encode to opus, using -abr
	V2X265     bool          // video video encode to x265, using -crf
	Par2C      bool          // par2create, create par2 files (in work directory)
	NoClobber  bool          // no clobber, do not queue those already been converted
	BreathTime time.Duration // breath time, interval between conversion to take a breath
	MaxC       int           // max conversion done each run (default no limit)
	NoExec     bool          // no exec, dry run
	Force      bool          // overwrite any existing none-empty file
	Debug      int           // debugging level
	FFMpeg     string        // ffmpeg program executable name
	FFProbe    string        // ffprobe program execution
	PrintV     bool          // print version then exit
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

// Opts holds the actual values from the command line parameters
var Opts Options

////////////////////////////////////////////////////////////////////////////
// Commandline definitions

func initVars() {

	// set default values for command line parameters
	flag.StringVar(&Opts.Cfg, "cfg", "",
		"cfg file to define your own targets: webm/wx/youtube etc")
	flag.StringVar(&Opts.Target, "t", "webm",
		"target type: webm/x265-opus/x264-mp3/wx/youtube/copy, or empty")
	flag.StringVar(&Opts.VES, "ves", "",
		"video encoding method set")
	flag.StringVar(&Opts.AES, "aes", "",
		"audio encoding method set")
	flag.StringVar(&Opts.SES, "ses", "",
		"subtitle encoding method set")
	flag.StringVar(&Opts.VEP, "vep", "",
		"video encoding method prepend")
	flag.StringVar(&Opts.AEP, "aep", "",
		"audio encoding method prepend")
	flag.StringVar(&Opts.SEP, "sep", "",
		"subtitle encoding method prepend")
	flag.StringVar(&Opts.VEA, "vea", "",
		"video encoding method append")
	flag.StringVar(&Opts.AEA, "aea", "",
		"audio encoding method append")
	flag.StringVar(&Opts.ABR, "abr", "",
		"audio bitrate (64k for opus, 256k for mp3)")
	flag.StringVar(&Opts.CRF, "crf", "",
		"the CRF value: 0-51. Higher CRF gives lower quality\n\t (28 for x265, ~ 23 for x264)")

	flag.StringVar(&Opts.Directory, "d", "",
		"directory that hold input files")
	flag.StringVar(&Opts.File, "f", "",
		"input file name (either -d or -f must be specified)")
	flag.BoolVar(&Opts.Links, "sym", false,
		"symlinks will be processed as well")
	flag.StringVar(&Opts.Exts, "exts", ".3GP.3G2.ASF.AVI.DAT.DIVX.FLV.M2TS.M4V.MKV.MOV.MPEG.MP4.MPG.RMVB.RM.TS.VOB.WEBM.WMV",
		"extension list for all the files to be queued")
	flag.StringVar(&Opts.Suffix, "suf", "",
		"suffix to the output file names")
	flag.StringVar(&Opts.Ext, "ext", "",
		"extension for the output file")
	flag.StringVar(&Opts.WDirectory, "w", "",
		"work directory that hold output files")

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
	flag.Var(&Opts.Cut, "C",
		"Cut segment(s) out to keep. Specify in the form of start-[end],\n\tstrictly in the format of hh:mm:ss, and may repeat")
	flag.Var(&Opts.Cut, "Cut",
		"Cut segment(s) out to keep. Specify in the form of start-[end],\n\tstrictly in the format of hh:mm:ss, and may repeat")
	flag.StringVar(&Opts.Seg, "S", "",
		"Split video into multiple segments (strictly in format: hh:mm:ss)")
	flag.StringVar(&Opts.Seg, "Seg", "",
		"Split video into multiple segments (strictly in format: hh:mm:ss)")
	flag.StringVar(&Opts.Speed, "Speed", "",
		"Speed up/down video playback speed (e.g. 1.28)")
	flag.BoolVar(&Opts.Karaoke, "K", false,
		"Add a karaoke audio track to .mp4 MTV")
	flag.BoolVar(&Opts.Karaoke, "karaoke", false,
		"Add a karaoke audio track to .mp4 MTV")
	flag.StringVar(&Opts.TranspFrom, "tkf", "",
		"Transpose song's key from (e.g. C/C#/Db/D etc)")
	flag.StringVar(&Opts.TranspTo, "tkt", "",
		"Transpose song's key to (e.g. -tkf C -tkt Db)")
	flag.IntVar(&Opts.TranspBy, "tkb", 0,
		"Transpose song by (e.g. +2, -3, etc) chromatic scale")
	flag.StringVar(&Opts.Lang, "lang", "eng",
		"language selection for audio stream extraction")
	flag.Var(&Opts.SEL, "sel",
		"subtitle encoding language (language picked for reencoded video)")
	flag.StringVar(&Opts.OptExtra, "o", "",
		"more options that will pass to ffmpeg program")
	flag.BoolVar(&Opts.A2Opus, "ato-opus", false,
		"audio encode to opus, using -abr")
	flag.BoolVar(&Opts.V2X265, "vto-x265", false,
		"video video encode to x265, using -crf")

	flag.BoolVar(&Opts.Par2C, "p", false,
		"par2create, create par2 files (in work directory)")
	flag.BoolVar(&Opts.NoClobber, "nc", false,
		"no clobber, do not queue those already been converted")
	flag.DurationVar(&Opts.BreathTime, "bt", 120*time.Second,
		"breath time, interval between conversion to take a breath")
	flag.IntVar(&Opts.MaxC, "maxc", 0,
		"max conversion done each run (default no limit)")
	flag.BoolVar(&Opts.NoExec, "n", false,
		"no exec, dry run")

	flag.BoolVar(&Opts.Force, "force", false,
		"overwrite any existing none-empty file")
	flag.IntVar(&Opts.Debug, "debug", 1,
		"debugging level")
	flag.StringVar(&Opts.FFMpeg, "ffmpeg", "ffmpeg",
		"ffmpeg program executable name")
	flag.StringVar(&Opts.FFProbe, "ffprobe", "ffprobe -print_format flat",
		"ffprobe program execution")
	flag.BoolVar(&Opts.PrintV, "version", false,
		"print version then exit")
}

func initVals() {
	exists := false
	// Now override those default values from environment variables
	if len(Opts.Cfg) == 0 &&
		len(os.Getenv("FFCVT_CFG")) != 0 {
		Opts.Cfg = os.Getenv("FFCVT_CFG")
	}
	if len(Opts.Target) == 0 &&
		len(os.Getenv("FFCVT_T")) != 0 {
		Opts.Target = os.Getenv("FFCVT_T")
	}
	if len(Opts.VES) == 0 &&
		len(os.Getenv("FFCVT_VES")) != 0 {
		Opts.VES = os.Getenv("FFCVT_VES")
	}
	if len(Opts.AES) == 0 &&
		len(os.Getenv("FFCVT_AES")) != 0 {
		Opts.AES = os.Getenv("FFCVT_AES")
	}
	if len(Opts.SES) == 0 &&
		len(os.Getenv("FFCVT_SES")) != 0 {
		Opts.SES = os.Getenv("FFCVT_SES")
	}
	if len(Opts.VEP) == 0 &&
		len(os.Getenv("FFCVT_VEP")) != 0 {
		Opts.VEP = os.Getenv("FFCVT_VEP")
	}
	if len(Opts.AEP) == 0 &&
		len(os.Getenv("FFCVT_AEP")) != 0 {
		Opts.AEP = os.Getenv("FFCVT_AEP")
	}
	if len(Opts.SEP) == 0 &&
		len(os.Getenv("FFCVT_SEP")) != 0 {
		Opts.SEP = os.Getenv("FFCVT_SEP")
	}
	if len(Opts.VEA) == 0 &&
		len(os.Getenv("FFCVT_VEA")) != 0 {
		Opts.VEA = os.Getenv("FFCVT_VEA")
	}
	if len(Opts.AEA) == 0 &&
		len(os.Getenv("FFCVT_AEA")) != 0 {
		Opts.AEA = os.Getenv("FFCVT_AEA")
	}
	if len(Opts.ABR) == 0 &&
		len(os.Getenv("FFCVT_ABR")) != 0 {
		Opts.ABR = os.Getenv("FFCVT_ABR")
	}
	if len(Opts.CRF) == 0 &&
		len(os.Getenv("FFCVT_CRF")) != 0 {
		Opts.CRF = os.Getenv("FFCVT_CRF")
	}

	if len(Opts.Directory) == 0 &&
		len(os.Getenv("FFCVT_D")) != 0 {
		Opts.Directory = os.Getenv("FFCVT_D")
	}
	if len(Opts.File) == 0 &&
		len(os.Getenv("FFCVT_F")) != 0 {
		Opts.File = os.Getenv("FFCVT_F")
	}
	if _, exists = os.LookupEnv("FFCVT_SYM"); Opts.Links || exists {
		Opts.Links = true
	}
	if len(Opts.Exts) == 0 &&
		len(os.Getenv("FFCVT_EXTS")) != 0 {
		Opts.Exts = os.Getenv("FFCVT_EXTS")
	}
	if len(Opts.Suffix) == 0 &&
		len(os.Getenv("FFCVT_SUF")) != 0 {
		Opts.Suffix = os.Getenv("FFCVT_SUF")
	}
	if len(Opts.Ext) == 0 &&
		len(os.Getenv("FFCVT_EXT")) != 0 {
		Opts.Ext = os.Getenv("FFCVT_EXT")
	}
	if len(Opts.WDirectory) == 0 &&
		len(os.Getenv("FFCVT_W")) != 0 {
		Opts.WDirectory = os.Getenv("FFCVT_W")
	}

	if _, exists = os.LookupEnv("FFCVT_AC"); Opts.AC || exists {
		Opts.AC = true
	}
	if _, exists = os.LookupEnv("FFCVT_VC"); Opts.VC || exists {
		Opts.VC = true
	}
	if _, exists = os.LookupEnv("FFCVT_AN"); Opts.AN || exists {
		Opts.AN = true
	}
	if _, exists = os.LookupEnv("FFCVT_VN"); Opts.VN || exists {
		Opts.VN = true
	}
	if _, exists = os.LookupEnv("FFCVT_VSS"); Opts.VSS || exists {
		Opts.VSS = true
	}
	if len(Opts.Seg) == 0 &&
		len(os.Getenv("FFCVT_S")) != 0 {
		Opts.Seg = os.Getenv("FFCVT_S")
	}
	if len(Opts.Seg) == 0 &&
		len(os.Getenv("FFCVT_SEG")) != 0 {
		Opts.Seg = os.Getenv("FFCVT_SEG")
	}
	if len(Opts.Speed) == 0 &&
		len(os.Getenv("FFCVT_SPEED")) != 0 {
		Opts.Speed = os.Getenv("FFCVT_SPEED")
	}
	if _, exists = os.LookupEnv("FFCVT_K,KARAOKE"); Opts.Karaoke || exists {
		Opts.Karaoke = true
	}
	if len(Opts.TranspFrom) == 0 &&
		len(os.Getenv("FFCVT_TKF")) != 0 {
		Opts.TranspFrom = os.Getenv("FFCVT_TKF")
	}
	if len(Opts.TranspTo) == 0 &&
		len(os.Getenv("FFCVT_TKT")) != 0 {
		Opts.TranspTo = os.Getenv("FFCVT_TKT")
	}
	if Opts.TranspBy == 0 &&
		len(os.Getenv("FFCVT_TKB")) != 0 {
		if i, err := strconv.Atoi(os.Getenv("FFCVT_TKB")); err == nil {
			Opts.TranspBy = i
		}
	}
	if len(Opts.Lang) == 0 &&
		len(os.Getenv("FFCVT_LANG")) != 0 {
		Opts.Lang = os.Getenv("FFCVT_LANG")
	}
	if len(Opts.OptExtra) == 0 &&
		len(os.Getenv("FFCVT_O")) != 0 {
		Opts.OptExtra = os.Getenv("FFCVT_O")
	}
	if _, exists = os.LookupEnv("FFCVT_ATO_OPUS"); Opts.A2Opus || exists {
		Opts.A2Opus = true
	}
	if _, exists = os.LookupEnv("FFCVT_VTO_X265"); Opts.V2X265 || exists {
		Opts.V2X265 = true
	}

	if _, exists = os.LookupEnv("FFCVT_P"); Opts.Par2C || exists {
		Opts.Par2C = true
	}
	if _, exists = os.LookupEnv("FFCVT_NC"); Opts.NoClobber || exists {
		Opts.NoClobber = true
	}
	if Opts.MaxC == 0 &&
		len(os.Getenv("FFCVT_MAXC")) != 0 {
		if i, err := strconv.Atoi(os.Getenv("FFCVT_MAXC")); err == nil {
			Opts.MaxC = i
		}
	}
	if _, exists = os.LookupEnv("FFCVT_N"); Opts.NoExec || exists {
		Opts.NoExec = true
	}

	if _, exists = os.LookupEnv("FFCVT_FORCE"); Opts.Force || exists {
		Opts.Force = true
	}
	if Opts.Debug == 0 &&
		len(os.Getenv("FFCVT_DEBUG")) != 0 {
		if i, err := strconv.Atoi(os.Getenv("FFCVT_DEBUG")); err == nil {
			Opts.Debug = i
		}
	}
	if len(Opts.FFMpeg) == 0 &&
		len(os.Getenv("FFCVT_FFMPEG")) != 0 {
		Opts.FFMpeg = os.Getenv("FFCVT_FFMPEG")
	}
	if len(Opts.FFProbe) == 0 &&
		len(os.Getenv("FFCVT_FFPROBE")) != 0 {
		Opts.FFProbe = os.Getenv("FFCVT_FFPROBE")
	}
	if _, exists = os.LookupEnv("FFCVT_VERSION"); Opts.PrintV || exists {
		Opts.PrintV = true
	}

}

const usageSummary = "  -cfg\tcfg file to define your own targets: webm/wx/youtube etc (FFCVT_CFG)\n  -t\ttarget type: webm/x265-opus/x264-mp3/wx/youtube/copy, or empty (FFCVT_T)\n  -ves\tvideo encoding method set (FFCVT_VES)\n  -aes\taudio encoding method set (FFCVT_AES)\n  -ses\tsubtitle encoding method set (FFCVT_SES)\n  -vep\tvideo encoding method prepend (FFCVT_VEP)\n  -aep\taudio encoding method prepend (FFCVT_AEP)\n  -sep\tsubtitle encoding method prepend (FFCVT_SEP)\n  -vea\tvideo encoding method append (FFCVT_VEA)\n  -aea\taudio encoding method append (FFCVT_AEA)\n  -abr\taudio bitrate (64k for opus, 256k for mp3) (FFCVT_ABR)\n  -crf\tthe CRF value: 0-51. Higher CRF gives lower quality\n\t (28 for x265, ~ 23 for x264) (FFCVT_CRF)\n\n  -d\tdirectory that hold input files (FFCVT_D)\n  -f\tinput file name (either -d or -f must be specified) (FFCVT_F)\n  -sym\tsymlinks will be processed as well (FFCVT_SYM)\n  -exts\textension list for all the files to be queued (FFCVT_EXTS)\n  -suf\tsuffix to the output file names (FFCVT_SUF)\n  -ext\textension for the output file (FFCVT_EXT)\n  -w\twork directory that hold output files (FFCVT_W)\n\n  -ac\tcopy audio codec (FFCVT_AC)\n  -vc\tcopy video codec (FFCVT_VC)\n  -an\tno audio, output video only (FFCVT_AN)\n  -vn\tno video, output audio only (FFCVT_VN)\n  -vss\tvideo: same size (FFCVT_VSS)\n  -C,Cut\tCut segment(s) out to keep. Specify in the form of start-[end],\n\tstrictly in the format of hh:mm:ss, and may repeat (FFCVT_C,CUT)\n  -S,Seg\tSplit video into multiple segments (strictly in format: hh:mm:ss) (FFCVT_S,SEG)\n  -Speed\tSpeed up/down video playback speed (e.g. 1.28) (FFCVT_SPEED)\n  -K,karaoke\tAdd a karaoke audio track to .mp4 MTV (FFCVT_K,KARAOKE)\n  -tkf\tTranspose song's key from (e.g. C/C#/Db/D etc) (FFCVT_TKF)\n  -tkt\tTranspose song's key to (e.g. -tkf C -tkt Db) (FFCVT_TKT)\n  -tkb\tTranspose song by (e.g. +2, -3, etc) chromatic scale (FFCVT_TKB)\n  -lang\tlanguage selection for audio stream extraction (FFCVT_LANG)\n  -sel\tsubtitle encoding language (language picked for reencoded video) (FFCVT_SEL)\n  -o\tmore options that will pass to ffmpeg program (FFCVT_O)\n  -ato-opus\taudio encode to opus, using -abr (FFCVT_ATO_OPUS)\n  -vto-x265\tvideo video encode to x265, using -crf (FFCVT_VTO_X265)\n\n  -p\tpar2create, create par2 files (in work directory) (FFCVT_P)\n  -nc\tno clobber, do not queue those already been converted (FFCVT_NC)\n  -bt\tbreath time, interval between conversion to take a breath (FFCVT_BT)\n  -maxc\tmax conversion done each run (default no limit) (FFCVT_MAXC)\n  -n\tno exec, dry run (FFCVT_N)\n\n  -force\toverwrite any existing none-empty file (FFCVT_FORCE)\n  -debug\tdebugging level (FFCVT_DEBUG)\n  -ffmpeg\tffmpeg program executable name (FFCVT_FFMPEG)\n  -ffprobe\tffprobe program execution (FFCVT_FFPROBE)\n  -version\tprint version then exit (FFCVT_VERSION)\n\nDetails:\n\n"

// Usage function shows help on commandline usage
func Usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage:\n %s [flags] \n\nFlags:\n\n",
		progname)
	fmt.Fprintf(os.Stderr, usageSummary)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		`
To reduce output, use '-debug 0', e.g., 'ffcvt -force -debug 0 -f testf.mp4 ...'
`)
	os.Exit(0)
}
