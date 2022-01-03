# ffcvt
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

[![MIT License](http://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/suntong/ffcvt?status.svg)](http://godoc.org/github.com/suntong/ffcvt)
[![Go Report Card](https://goreportcard.com/badge/github.com/suntong/ffcvt)](https://goreportcard.com/report/github.com/suntong/ffcvt)
[![Build Status](https://github.com/suntong/ffcvt/actions/workflows/go-release-build.yml/badge.svg?branch=master)](https://github.com/suntong/ffcvt/actions/workflows/go-release-build.yml)
[![PoweredBy WireFrame](https://github.com/go-easygen/wireframe/blob/master/PoweredBy-WireFrame-B.svg)](http://godoc.org/github.com/go-easygen/wireframe)


## TOC
- [ffcvt - ffmpeg convert wrapper tool](#ffcvt---ffmpeg-convert-wrapper-tool)
  - [Latest Update(s)](#latest-update(s))
    - [Release v1.7.5](#release-v175)
    - [Release v1.7.3](#release-v173)
    - [Release v1.7.2](#release-v172)
    - [Release v1.7.1](#release-v171)
    - [Release v1.7.0](#release-v170)
  - [Introduction](#introduction)
  - [Quick Usage](#quick-usage)
    - [$ ffcvt](#-ffcvt)
  - [Environment Variables](#environment-variables)
  - [Encoding Help](#encoding-help)
  - [Tools Choices](#tools-choices)
    - [Install Debian/Ubuntu package](#install-debianubuntu-package)
- [Install Debian/Ubuntu package](#install-debianubuntu-package)
- [Download/install binaries](#downloadinstall-binaries)
  - [The binary executables](#the-binary-executables)
  - [Distro package](#distro-package)
  - [Debian package](#debian-package)
- [Install Source](#install-source)
- [Author](#author)
- [Contributors](#contributors-)

## ffcvt - ffmpeg convert wrapper tool

### Latest Update(s)

#### Release v1.7.5

* Now able to speed up playback speed (`-Speed`). Details in [\#22](https://github.com/suntong/ffcvt/issues/22)
* Also have added a `copy` target type that can speed up the `Seg` (split video) operation (v1.7.4). Details in [\#21](https://github.com/suntong/ffcvt/issues/21)

#### Release v1.7.3

* Now able to split video into multiple segments (`-S,Seg`) by the given time. Details in [\#16](https://github.com/suntong/ffcvt/issues/16)

#### Release v1.7.2

* Able to [choose streams by language, instead of streams index. ](https://github.com/suntong/ffcvt/commit/f649609356ef06d22d17d6dbe3f89b945cf18643)Details in [\#9](https://github.com/suntong/ffcvt/issues/9)
* Fixed [\#8](https://github.com/suntong/ffcvt/issues/8). Now [force copy all subtitle streams. ](https://github.com/suntong/ffcvt/commit/46ce6725f9b036d373c6836d3bd66b429d5c4b2f)Details in [\#8](https://github.com/suntong/ffcvt/issues/8)
* [Added option -sel](https://github.com/suntong/ffcvt/commit/defc5df5168216e279b944590f1d92523ecadc60), so now able to pick subtitle language(s). Details in [\#12](https://github.com/suntong/ffcvt/issues/12)

#### Release v1.7.1

Added option `-C,Cut` which allows cutting multiple segments.

For further details, check out the wiki https://git.io/JuK0c,
in which the source file of

https://user-images.githubusercontent.com/422244/132961501-a2344db0-c48c-4a57-90fa-c3746bf3025f.mp4

is cut-short into

https://user-images.githubusercontent.com/422244/132961530-ea65cd03-19f8-4e7c-a871-40218f7289cc.mp4


#### Release v1.7.0

Added `wx` type for weixin.

Convert to video that is recognizable and playable by weixin/wechat, by using the `-t wx` option as the convertion type. Here is a converted sample:

https://user-images.githubusercontent.com/422244/132617136-e1371ef3-6a21-4f12-8324-6db003c12468.mp4

(credit [here](https://www.youtube.com/watch?v=2-UzBitLmf8))

For further details, check out the wiki https://git.io/JuK0q

## Introduction

- The next-generation codec like [High Efficiency Video codec (HEVC), H.265](https://goo.gl/IZrDH2) or [VP9](https://developers.google.com/media/vp9/) can produce videos visually comparable to H.264's result, but in [about half the file size](https://trac.ffmpeg.org/wiki/Encode/H.265).
- Meanwhile the [Opus](https://goo.gl/BPUkTf) [audio codec](https://goo.gl/IZrDH2) is becoming the best thing ever for compressing audio -- A 64K Opus audio stream is comparable to mp3 files of 128K to 256K bandwidth.
- Such fantastic high efficiency audio/video codec/encoding capability has long been available in `ffmpeg`, but fewer people know it or use it, partly because the `ffmpeg` command line is not that simple for every one.
- The `ffcvt` is designed to take the burden from normal Joe -- All you need to do to encode a video is to give one parameter to `ffcvt`, i.e., the path and file name of the video to be encoded, and `ffcvt` will take care of the rest, using the recommended values for both audio/video encoding to properly encode it for you.
- It can't be more simpler than that. However, beneath the simple surface, `ffcvt` is versatile and powerful enough to allow you to touch every corner of audio/video encoding. There is a huge list of environment variables (or command-line parameters) which will allow you tweak the encoding methods and parameters to exactly what you prefer instead.
- Moreover, to encode a directory full of video files, including under its sub-directories, you need just to give `ffcvt` one single parameter, the directory location, and `ffcvt` will go ahead and encode all video files under that directory, including all its sub-directories as well. 

## Quick Usage

There is a quick usage help that comes with `ffcvt`, produced when it is invoked without any parameters:

### $ ffcvt
```sh
Usage:
 ffcvt [flags] 

Flags:

  -t	target type: webm/x265-opus/x264-mp3/wx/youtube/copy (FFCVT_T)
  -ves	video encoding method set (FFCVT_VES)
  -aes	audio encoding method set (FFCVT_AES)
  -ses	subtitle encoding method set (FFCVT_SES)
  -vep	video encoding method prepend (FFCVT_VEP)
  -aep	audio encoding method prepend (FFCVT_AEP)
  -sep	subtitle encoding method prepend (FFCVT_SEP)
  -vea	video encoding method append (FFCVT_VEA)
  -aea	audio encoding method append (FFCVT_AEA)
  -abr	audio bitrate (64k for opus, 256k for mp3) (FFCVT_ABR)
  -crf	the CRF value: 0-51. Higher CRF gives lower quality
	 (28 for x265, ~ 23 for x264) (FFCVT_CRF)

  -d	directory that hold input files (FFCVT_D)
  -f	input file name (either -d or -f must be specified) (FFCVT_F)
  -sym	symlinks will be processed as well (FFCVT_SYM)
  -exts	extension list for all the files to be queued (FFCVT_EXTS)
  -suf	suffix to the output file names (FFCVT_SUF)
  -ext	extension for the output file (FFCVT_EXT)
  -w	work directory that hold output files (FFCVT_W)

  -ac	copy audio codec (FFCVT_AC)
  -vc	copy video codec (FFCVT_VC)
  -an	no audio, output video only (FFCVT_AN)
  -vn	no video, output audio only (FFCVT_VN)
  -vss	video: same size (FFCVT_VSS)
  -C,Cut	Cut segment(s) out to keep. Specify in the form of start-[end],
	strictly in the format of hh:mm:ss, and may repeat (FFCVT_C,CUT)
  -S,Seg	Split video into multiple segments (strictly in format: hh:mm:ss) (FFCVT_S,SEG)
  -Speed	Speed up/down video playback speed (e.g. 1.28) (FFCVT_SPEED)
  -lang	language selection for audio stream extraction (FFCVT_LANG)
  -sel	subtitle encoding language (language picked for reencoded video) (FFCVT_SEL)
  -o	more options that will pass to ffmpeg program (FFCVT_O)
  -ato-opus	audio encode to opus, using -abr (FFCVT_ATO_OPUS)
  -vto-x265	video video encode to x265, using -crf (FFCVT_VTO_X265)

  -p	par2create, create par2 files (in work directory) (FFCVT_P)
  -nc	no clobber, do not queue those already been converted (FFCVT_NC)
  -n	no exec, dry run (FFCVT_N)

  -force	overwrite any existing none-empty file (FFCVT_FORCE)
  -debug	debugging level (FFCVT_DEBUG)
  -ffmpeg	ffmpeg program executable name (FFCVT_FFMPEG)
  -ffprobe	ffprobe program execution (FFCVT_FFPROBE)
  -version	print version then exit (FFCVT_VERSION)

Details:

  -C value
    	Cut segment(s) out to keep. Specify in the form of start-[end],
    		strictly in the format of hh:mm:ss, and may repeat
  -Cut value
    	Cut segment(s) out to keep. Specify in the form of start-[end],
    		strictly in the format of hh:mm:ss, and may repeat
  -S string
    	Split video into multiple segments (strictly in format: hh:mm:ss)
  -Seg string
    	Split video into multiple segments (strictly in format: hh:mm:ss)
  -Speed string
    	Speed up/down video playback speed (e.g. 1.28)
  -abr string
    	audio bitrate (64k for opus, 256k for mp3)
  -ac
    	copy audio codec
  -aea string
    	audio encoding method append
  -aep string
    	audio encoding method prepend
  -aes string
    	audio encoding method set
  -an
    	no audio, output video only
  -ato-opus
    	audio encode to opus, using -abr
  -crf string
    	the CRF value: 0-51. Higher CRF gives lower quality
    		 (28 for x265, ~ 23 for x264)
  -d string
    	directory that hold input files
  -debug int
    	debugging level (default 1)
  -ext string
    	extension for the output file
  -exts string
    	extension list for all the files to be queued (default ".3GP.3G2.ASF.AVI.DAT.DIVX.FLV.M2TS.M4V.MKV.MOV.MPEG.MP4.MPG.RMVB.RM.TS.VOB.WEBM.WMV")
  -f string
    	input file name (either -d or -f must be specified)
  -ffmpeg string
    	ffmpeg program executable name (default "ffmpeg")
  -ffprobe string
    	ffprobe program execution (default "ffprobe -print_format flat")
  -force
    	overwrite any existing none-empty file
  -lang string
    	language selection for audio stream extraction (default "eng")
  -n	no exec, dry run
  -nc
    	no clobber, do not queue those already been converted
  -o string
    	more options that will pass to ffmpeg program
  -p	par2create, create par2 files (in work directory)
  -sel value
    	subtitle encoding language (language picked for reencoded video)
  -sep string
    	subtitle encoding method prepend
  -ses string
    	subtitle encoding method set
  -suf string
    	suffix to the output file names
  -sym
    	symlinks will be processed as well
  -t string
    	target type: webm/x265-opus/x264-mp3/wx/youtube/copy (default "webm")
  -vc
    	copy video codec
  -vea string
    	video encoding method append
  -vep string
    	video encoding method prepend
  -version
    	print version then exit
  -ves string
    	video encoding method set
  -vn
    	no video, output audio only
  -vss
    	video: same size (default true)
  -vto-x265
    	video video encode to x265, using -crf
  -w string
    	work directory that hold output files

To reduce output, use `-debug 0`, e.g., `ffcvt -force -debug 0 -f testf.mp4 ...`
```


## Environment Variables

For each `ffcvt` command line parameter, there is a environment variable corresponding to it. For example you can use `export FFCVT_FFMPEG=avconv` to use `avconv` instead of `ffmpeg` (Don't, I use it for my [CommandLineArgs](https://github.com/suntong001/lang/blob/master/lang/Go/src/sys/CommandLineArgs.go) to develop/test `ffcvt` without invoking `ffmpeg` each time). 

## Encoding Help

The detailed guide to choose/provide proper parameters to `ffcvt` have been moved to [wiki](https://github.com/suntong/ffcvt/wiki/). For example,

- [HEVC vs VP9](https://github.com/suntong/ffcvt/wiki/KB:-WebM-(VP9)-Encoding#hevc-vs-vp9)
- [HEVC Preset Method Comparison](https://github.com/suntong/ffcvt/wiki/KB:-HEVC-(x265)-Encoding#preset-method-comparison)
- [The HEVC CRF Comparison](https://github.com/suntong/ffcvt/wiki/KB:-HEVC-(x265)-Encoding#the-crf-comparison)
- [Example 1: YouTube Encoding](https://github.com/suntong/ffcvt/wiki/Example:-YouTube-Encoding)
- [Example 2: Talk Encoding](https://github.com/suntong/ffcvt/wiki/Example:-Talk-Encoding)

Please check them out in the [wiki](https://github.com/suntong/ffcvt/wiki/), and for other documents like "Most used ffmpeg options", "How to crop a video", etc.

## Tools Choices

As suggested before, don't use `avconv`, use `ffmpeg` instead (the `avconv` fork was more for political reasons. I personally believe `ffmpeg` is technically superior although might not be politically).

As for video/movie play back, use [mpv](http://mpv.io/). It is a fork of mplayer2 and MPlayer, and is a true *modern* *all-in-one* movie player that can play ANYTHING, and one of the few movie players being actively developed all the time. Download link is in [mpv.io](http://mpv.io/), from which Ubuntu repo I get my Ubuntu `ffmpeg` package as well. If you are unsatisfied with mpv's simple user interface, check out https://wiki.archlinux.org/index.php/Mpv#Front_ends.

### Install Debian/Ubuntu package

    apt install ffcvt

## Download/install binaries

- The latest binary executables are available 
as the result of the Continuous-Integration (CI) process.
- I.e., they are built automatically right from the source code at every git release by [GitHub Actions](https://docs.github.com/en/actions).
- There are two ways to get/install such binary executables
  * Using the **binary executables** directly, or
  * Using **packages** for your distro

### The binary executables

- The latest binary executables are directly available under  
https://github.com/suntong/ffcvt/releases/latest 
- Pick & choose the one that suits your OS and its architecture. E.g., for Linux, it would be the `ffcvt_verxx_linux_amd64.tar.gz` file. 
- Available OS for binary executables are
  * Linux
  * Mac OS (darwin)
  * Windows
- If your OS and its architecture is not available in the download list, please let me know and I'll add it.
- The manual installation is just to unpack it and move/copy the binary executable to somewhere in `PATH`. For example,

``` sh
tar -xvf ffcvt_*_linux_amd64.tar.gz
sudo mv -v ffcvt_*_linux_amd64/ffcvt /usr/local/bin/
rmdir -v ffcvt_*_linux_amd64
```


### Distro package

- [Packages available for Linux distros](https://cloudsmith.io/~suntong/repos/repo/packages/) are
  * [Alpine Linux](https://cloudsmith.io/~suntong/repos/repo/setup/#formats-alpine)
  * [Debian](https://cloudsmith.io/~suntong/repos/repo/setup/#formats-deb)
  * [RedHat](https://cloudsmith.io/~suntong/repos/repo/setup/#formats-rpm)

The repo setup instruction url has been given above.
For example, for [Debian](https://cloudsmith.io/~suntong/repos/repo/setup/#formats-deb) --

### Debian package


```sh
curl -1sLf \
  'https://dl.cloudsmith.io/public/suntong/repo/setup.deb.sh' \
  | sudo -E bash

# That's it. You then can do your normal operations, like

sudo apt-get update
apt-cache policy ffcvt

sudo apt-get install -y ffcvt
```

## Install Source

To install the source code instead:

```
go get -v -u github.com/suntong/ffcvt
```

## Author

Tong SUN  
![suntong from cpan.org](https://img.shields.io/badge/suntong-%40cpan.org-lightgrey.svg "suntong from cpan.org")

_Powered by_ [**WireFrame**](https://github.com/go-easygen/wireframe)  
[![PoweredBy WireFrame](https://github.com/go-easygen/wireframe/blob/master/PoweredBy-WireFrame-Y.svg)](http://godoc.org/github.com/go-easygen/wireframe)  
the _one-stop wire-framing solution_ for Go cli based projects, from _init_ to _deploy_.

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/suntong"><img src="https://avatars.githubusercontent.com/u/422244?v=4?s=100" width="100px;" alt=""/><br /><sub><b>suntong</b></sub></a><br /><a href="https://github.com/suntong/ffcvt/commits?author=suntong" title="Code">💻</a> <a href="#ideas-suntong" title="Ideas, Planning, & Feedback">🤔</a> <a href="#design-suntong" title="Design">🎨</a> <a href="#data-suntong" title="Data">🔣</a> <a href="https://github.com/suntong/ffcvt/commits?author=suntong" title="Tests">⚠️</a> <a href="https://github.com/suntong/ffcvt/issues?q=author%3Asuntong" title="Bug reports">🐛</a> <a href="https://github.com/suntong/ffcvt/commits?author=suntong" title="Documentation">📖</a> <a href="#blog-suntong" title="Blogposts">📝</a> <a href="#example-suntong" title="Examples">💡</a> <a href="#tutorial-suntong" title="Tutorials">✅</a> <a href="#tool-suntong" title="Tools">🔧</a> <a href="#platform-suntong" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/suntong/ffcvt/pulls?q=is%3Apr+reviewed-by%3Asuntong" title="Reviewed Pull Requests">👀</a> <a href="#question-suntong" title="Answering Questions">💬</a> <a href="#maintenance-suntong" title="Maintenance">🚧</a> <a href="#infra-suntong" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
    <td align="center"><a href="https://github.com/sanjaymsh"><img src="https://avatars.githubusercontent.com/u/66668807?v=4?s=100" width="100px;" alt=""/><br /><sub><b>sanjaymsh</b></sub></a><br /><a href="#platform-sanjaymsh" title="Packaging/porting to new platform">📦</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
