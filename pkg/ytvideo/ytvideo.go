package ytvideo

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/89z/mech/youtube"
	"github.com/SultanKs4/ytarchiver/config"
	"github.com/SultanKs4/ytarchiver/pkg/fileutils"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

const chunk int64 = 10_000_000

type YtVideo struct {
	id       string
	address  string
	basePath string
	config   *config.Config
	player   *youtube.Player
}

func NewYtVideo(id, address, basePath string, config *config.Config) (*YtVideo, error) {
	if id == "" {
		err := youtube.Video_ID(address, &id)
		if err != nil {
			return nil, err
		}
	}
	player, err := youtube.Android().Player(id)
	if err != nil {
		return nil, err
	}
	return &YtVideo{id: id, address: address, basePath: basePath, config: config, player: player}, nil
}

// generate output file name
func (myt *YtVideo) getOutputFile() (outputFilename string, err error) {
	outputFilename = fileutils.SanitizeFilename(myt.player.VideoDetails.Title) + "." + myt.config.Youtube.Mimetype

	// if base path empty then save to resource > author
	// else custom base path > author
	if myt.basePath == "" {
		myt.basePath = filepath.Join(fileutils.ResourcePath(), fileutils.SanitizeFilename(myt.player.VideoDetails.Author))
	} else {
		// if base path contain folder name same with author skip create path
		author := fileutils.SanitizeFilename(myt.player.VideoDetails.Author)
		if !strings.Contains(myt.basePath, author) {
			myt.basePath = filepath.Join(myt.basePath, author)
		}
	}

	if err := os.MkdirAll(myt.basePath, os.ModePerm); err != nil {
		return "", err
	}

	outputFilename = filepath.Join(myt.basePath, outputFilename)
	return
}

// get seperate video and audio file
func (myt *YtVideo) getVideoAudioFormat() (videoFormat *youtube.Format, audioFormat *youtube.Format, err error) {
	formats := myt.player.StreamingData.AdaptiveFormats
	// sort ascending based from height property
	sort.SliceStable(formats, func(i, j int) bool {
		return formats[i].Height < formats[j].Height
	})
	for i := range formats {
		if strings.Contains(formats[i].MimeType, "video/"+myt.config.Youtube.Mimetype) && formats[i].Height >= myt.config.Youtube.Height && formats[i].ContentLength > 0 {
			if videoFormat != nil {
				continue
			}
			videoFormat = &formats[i]
		} else if strings.Contains(formats[i].MimeType, "audio/"+myt.config.Youtube.Mimetype) && myt.config.Youtube.Audio == formats[i].AudioQuality {
			if audioFormat != nil {
				continue
			}
			audioFormat = &formats[i]
		}
	}

	if videoFormat == nil || videoFormat.ContentLength == 0 {
		err = fmt.Errorf("no video format found after filtering, try to change height")
		return
	}

	if audioFormat == nil || audioFormat.ContentLength == 0 {
		err = fmt.Errorf("no audio format found after filtering, try to change quality")
		return
	}

	return
}

// custom encode using progress bar mpb
//
// TODO: Search when downloaad big size video potentially OOM (Out of Memory)
func (myt *YtVideo) Encode(f *youtube.Format, w io.Writer) error {
	req, err := http.NewRequest("GET", f.URL, nil)
	if err != nil {
		return err
	}

	// create progress bar
	progress := mpb.New(mpb.WithWidth(64))
	bar := progress.AddBar(
		int64(f.ContentLength),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" speed: "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)

	// download chunked for faster download
	var pos int64
	for pos < f.ContentLength {
		b := []byte("bytes=")
		b = strconv.AppendInt(b, pos, 10)
		b = append(b, '-')
		b = strconv.AppendInt(b, pos+chunk-1, 10)
		req.Header.Set("Range", string(b))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		reader := bar.ProxyReader(res.Body)
		defer reader.Close()

		if _, err = io.Copy(w, reader); err != nil {
			return err
		}

		pos += chunk
	}

	progress.Wait()
	return nil
}

// log every process
func (myt *YtVideo) logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// main process
func (myt *YtVideo) DownloadComposite() error {
	videoFormat, audioFormat, err := myt.getVideoAudioFormat()
	if err != nil {
		return err
	}

	myt.logf("Video '%s' - Quality '%s' - Video Codec '%s' - Audio Codec '%s'", myt.player.VideoDetails.Title, videoFormat.QualityLabel, videoFormat.MimeType, audioFormat.MimeType)
	destFile, err := myt.getOutputFile()
	if err != nil {
		return err
	}

	if _, err := os.Stat(destFile); !os.IsNotExist(err) {
		return fmt.Errorf("file exist")
	}
	outputDir := filepath.Dir(destFile)

	log.Printf("save video to: %s", outputDir)

	// Create temporary video file
	myt.logf("Downloading video file...")
	videoFile, err := os.CreateTemp(outputDir, "youtube_*.m4v")
	if err != nil {
		return err
	}
	defer videoFile.Close()
	defer os.Remove(videoFile.Name())

	err = myt.Encode(videoFormat, videoFile)
	if err != nil {
		return err
	}

	// Create temporary audio file
	myt.logf("Downloading audio file...")
	audioFile, err := os.CreateTemp(outputDir, "youtube_*.m4a")
	if err != nil {
		return err
	}
	defer audioFile.Close()
	defer os.Remove(audioFile.Name())

	err = myt.Encode(audioFormat, audioFile)
	if err != nil {
		return err
	}

	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", videoFile.Name(),
		"-i", audioFile.Name(),
		"-c", "copy", // Just copy without re-encoding
		"-shortest", // Finish encoding when the shortest input stream ends
		destFile,
		"-loglevel", "warning",
	)
	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout
	myt.logf("merging video and audio to %s", destFile)

	err = ffmpegVersionCmd.Run()
	if err != nil {
		return err
	}
	myt.logf("created: %s", destFile)
	return nil
}
