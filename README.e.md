## {{toc 5}}
- [Install Debian/Ubuntu package](#install-debianubuntu-package)
- [Download/install binaries](#downloadinstall-binaries)
  - [The binary executables](#the-binary-executables)
  - [Distro package](#distro-package)
  - [Debian package](#debian-package)
- [Install Source](#install-source)
- [Author](#author)
- [Contributors](#contributors-)

## {{.Name}} - ffmpeg convert wrapper tool

### Latest Update(s)


#### Latest Releases

- Release v1.9.0
  * `ffcvt -version` now checks/outputs dependent program versions too
  * now finished percentage are calculated from file size
- Release v1.8.1, enable parallel execution

#### Release v1.8.0

* Now able to define your own defaults. Just make a copy of [ffcvt.json](ffcvt.json) and customize it to your heart's content, then use the `-cfg` option to point to it. Better yet, set `FFCVT_CFG` environment variable and forget all about it afterwards.
  * This means that `ffcvt` is now not only limited to its own predefined transcoding sets, but you can also define your own transcoding rules and names and then fully enjoy its advanced addon assistants.
  * BTW, If you have a good set, don't forget to send in a PR so that everybody can also benefit from it.
* Now the subtitles, nfo, html or any files in the source directory will be duplicated into the output (work) directory, first by hard-link and if it fails due to cross storage devices, a copy will be used instead.
* And when creating `par2` checksum/repair files, all files in the output (work) directory will be covered.

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

### $ ffcvt -version

```sh
$ {{shell "ffcvt -version"}}
```

### $ {{exec "ffcvt" | color "sh"}}


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

    apt install {{.Name}}
