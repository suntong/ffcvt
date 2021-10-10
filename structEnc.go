package main

import "log"

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The Encoding struct defines the structure to hold encoding values
type Encoding struct {
	VES string // video encoding method set
	AES string // audio encoding method set
	SES string // subtitle encoding method set
	VEP string // video encoding method prepend
	AEP string // audio encoding method prepend
	SEP string // subtitle encoding method prepend
	VEA string // video encoding method append
	AEA string // audio encoding method append
	ABR string // audio bitrate
	CRF string // the CRF value: 0-51. Higher CRF gives lower quality
	Ext string // extension for the output file
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var Defaults map[string]Encoding

func init() {

	Defaults = map[string]Encoding{
		"webm": {
			AES: "libopus",
			VES: "libvpx-vp9",
			SES: "-c:s copy",
			ABR: "64k",
			CRF: "42",
			Ext: _encodedExt,
		},
		"x265-opus": {
			AES: "libopus",
			VES: "libx265",
			ABR: "64k",
			CRF: "28",
			Ext: _encodedExt,
		},
		"wx": {
			AES: "aac",
			AEA: "-q:a 3",
			VES: "libx264",
			ABR: "48k",
			CRF: "33",
			Ext: "_.m4v",
		},
		"x264-mp3": {
			AES: "libmp3lame",
			AEA: "-q:a 3",
			VES: "libx264",
			ABR: "256k",
			CRF: "23",
			Ext: "_.mp4",
		},
		"youtube": {
			// https://trac.ffmpeg.org/wiki/Encode/YouTube
			// https://trac.ffmpeg.org/wiki/Encode/HighQualityAudio
			AES: "libvorbis",
			AEA: "-q:a 5",
			VES: "libx264",
			VEA: "-pix_fmt yuv420p",
			ABR: "",
			CRF: "20",
			Ext: "_.avi",
		},
	}

}

func getDefault() {
	// preserve command line values
	abr, crf := Opts.ABR, Opts.CRF

	if encDefault, ok := Defaults[Opts.Target]; ok {
		// debug(encDefault.Ext, 2)
		Opts.Encoding = encDefault
		// debug(Opts.Encoding.Ext, 2)
		// debug(Opts.Ext, 2)
	} else {
		log.Fatal(progname + " Error: Wrong target option passed to -t.")
	}

	// restore command line values
	if abr != "" {
		Opts.ABR = abr
	}
	if crf != "" {
		Opts.CRF = crf
	}

	if Opts.Suffix != "" {
		Opts.Suffix = "_" + Opts.Suffix
	}

}
