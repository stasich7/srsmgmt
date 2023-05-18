package playlist

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gopkg.in/vansante/go-ffprobe.v2"
)

const (
	AllPlaylists  = ""
	LiveNumChunks = 6
	PL_WaitTime   = 12 * time.Second
	PL_RetryTime  = 1 * time.Second
)

var (
	ErrOperationTimedout = errors.New("TIMEOUT")
	ErrInternalError     = errors.New("INTERNAL_ERROR")
	ErrIvalidStartTime   = errors.New("INVALID_START_TIME")
	playlistTypes        = []string{"low.m3u8", "mid.m3u8", "high.m3u8"}
	once                 sync.Once
)

type mediaInfo struct {
	Path       string
	Name       string
	Bitrate    string
	Resolution string
	Codec      string
	Audio      string
}

type Playlist struct {
	logger log.Logger
}

type PLline struct {
	Num  int
	Line string
	TS   bool
}
type PLstruct struct {
	Header []PLline
	TS     []PLline
}

func New() *Playlist {
	var newlogger log.Logger
	once.Do(func() {
		newlogger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		newlogger = log.With(newlogger, "ts", log.TimestampFormat(time.Now, time.RFC3339))
		newlogger = log.With(newlogger, "caller", log.Caller(5))
		newlogger = level.NewFilter(newlogger, level.AllowDebug())
	})
	return &Playlist{
		logger: newlogger,
	}
}

func (p *Playlist) Refresh(livePath, plType string, startTs *time.Time, outputPlaylistPrefix string) error {
	re, _ := regexp.Compile(plType)
	plsToParse := []string{}
	for _, pl := range playlistTypes {
		if re.MatchString(pl) {
			plsToParse = append(plsToParse, pl)
		}
	}
	wg := sync.WaitGroup{}
	errC := make(chan error, len(plsToParse))

	re_h1, _ := regexp.Compile(`EXT-X-MEDIA-SEQUENCE.+?(\d+)`)
	re_h2, _ := regexp.Compile(`#EXT-X-DISCONTINUITY`)

	for _, pl := range plsToParse {
		wg.Add(1)
		go func(wg *sync.WaitGroup, pl string) {
			defer wg.Done()
			start := time.Now()
			defer level.Info(p.logger).Log(fmt.Sprintf("RefreshPlaylistFrom_for_[%s/%s]_took", livePath, pl), time.Since(start))
			// level.Info(p.logger).Log("Looking for " + livePath + pl)
			plLines, err := getTSName(path.Join(livePath, pl))
			if err != nil {
				errC <- ErrIvalidStartTime
				return
			}

			var f_live, f_dvr *os.File
			f_live, err = os.OpenFile(path.Join(livePath, outputPlaylistPrefix+pl), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
			if err != nil {
				errC <- err
				return
			}
			level.Info(p.logger).Log("openF_live", fmt.Sprintf("[%s]-[%s]-[%s]", livePath, outputPlaylistPrefix, pl))
			if startTs != nil && !startTs.IsZero() {
				f_dvr, err = os.OpenFile(path.Join(livePath, fmt.Sprintf("%s%d-%s", outputPlaylistPrefix, startTs.Unix(), pl)), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
				if err != nil {
					errC <- err
					return
				}
				defer f_dvr.Close()
				level.Info(p.logger).Log("openF_dvr", fmt.Sprintf("[%s]-[%s]-[%d]-[%s]", livePath, outputPlaylistPrefix, startTs.Unix(), pl))
			}
			defer f_live.Close()

			tsCounter := 0
			for _, line := range plLines.Header {
				if re_h2.MatchString(line.Line) {
					continue
				}
				seqMatch := re_h1.FindStringSubmatch(line.Line)
				if len(seqMatch) > 0 {
					tsCounter, _ = strconv.Atoi(seqMatch[1])
					continue
				}
				_, err = fmt.Fprintln(f_live, line.Line)
				if err != nil {
					errC <- err
					return
				}
				if startTs != nil && !startTs.IsZero() {
					_, err = fmt.Fprintln(f_dvr, line.Line)
					if err != nil {
						errC <- err
						return
					}
				}
			}

			for _, v := range plLines.TS {
				if v.TS {
					tsCounter++
				}
			}
			if tsCounter-LiveNumChunks > 0 {
				tsCounter = tsCounter - LiveNumChunks
				fmt.Fprintf(f_live, "%s%d\n", "#EXT-X-MEDIA-SEQUENCE:", tsCounter)
			}
			if startTs == nil || startTs.IsZero() {
				fmt.Fprintf(f_live, "%s%s\n", "#EXT-X-PROGRAM-DATE-TIME:", time.Now().Format("2006-01-02T15:04:05.000Z"))
				fmt.Fprintf(f_live, "%s\n", "#EXT-X-PLAYLIST-TYPE:EVENT")
			} else {
				fmt.Fprintf(f_dvr, "%s%s\n", "#EXT-X-PROGRAM-DATE-TIME:", startTs.Format("2006-01-02T15:04:05.000Z"))
				fmt.Fprintf(f_dvr, "%s\n", "#EXT-X-PLAYLIST-TYPE:EVENT")
			}

			startPlaylist := false
			prevLine := ""
			for i, tss := range plLines.TS {
				if tss.TS && !startPlaylist {
					tsTime, err := p.timeFromTS(tss.Line)
					if err != nil {
						errC <- err
						return
					}
					if (startTs == nil) || (startTs != nil && tsTime.After(*startTs)) {
						startPlaylist = true
						if startTs != nil && !startTs.IsZero() {
							_, err = fmt.Fprintln(f_dvr, prevLine)
							if err != nil {
								errC <- err
								return
							}
						}
						if len(plLines.TS)-i <= LiveNumChunks+2 {
							_, err = fmt.Fprintln(f_live, prevLine)
							if err != nil {
								errC <- err
								return
							}
						}
					}
				}
				prevLine = tss.Line
				if !startPlaylist {
					continue
				}
				if startTs != nil && !startTs.IsZero() {
					_, err = fmt.Fprintln(f_dvr, tss.Line)
					if err != nil {
						errC <- err
						return
					}
				}

				if len(plLines.TS)-i <= LiveNumChunks*2 {
					_, err = fmt.Fprintln(f_live, tss.Line)
					if err != nil {
						errC <- err
						return
					}
				}

			}
		}(&wg, pl)
	}

	wg.Wait()
	close(errC)

	select {
	case err, ok := <-errC:
		if ok {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (p *Playlist) timeFromTS(TSfilename string) (*time.Time, error) {
	re := regexp.MustCompile(`\w+-(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-`)
	loc, _ := time.LoadLocation("Europe/Moscow")

	result := re.FindAllStringSubmatch(TSfilename, -1)
	if result == nil {
		return nil, ErrInternalError
	}

	respTime, err := time.ParseInLocation("2006-01-02T15:04:05", fmt.Sprintf("%s-%s-%sT%s:%s:%s", result[0][1], result[0][2], result[0][3], result[0][4], result[0][5], result[0][6]), loc)
	if err != nil {
		return nil, ErrInternalError
	}

	return &respTime, nil
}

func (p *Playlist) Stop(livePath string, outputPlaylistPrefix string, startTs *time.Time) error {
	for _, pl := range playlistTypes {
		appendStr := []byte("#EXT-X-ENDLIST\n")
		f, err := os.OpenFile(path.Join(livePath, outputPlaylistPrefix+pl), os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(appendStr)
		if err != nil {
			return err
		}

		var f_ts *os.File
		if startTs != nil && !startTs.IsZero() {
			f_ts, err = os.OpenFile(path.Join(livePath, fmt.Sprintf("%s%d-%s", outputPlaylistPrefix, startTs.Unix(), pl)), os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				return err
			}
			defer f_ts.Close()
			_, err = f_ts.Write(appendStr)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (p *Playlist) Create(livePath string, outputPlaylistPrefix string, startTs *time.Time) error {
	plInfoC := make(chan map[string]*mediaInfo)
	plInfo := make(map[string]*mediaInfo)

	tcanel := time.After(PL_WaitTime)
	done := make(chan bool)

	if startTs != nil && !startTs.IsZero() {
		outputPlaylistPrefix = fmt.Sprintf("%s%d-", outputPlaylistPrefix, startTs.Unix())
	}

	for _, item := range playlistTypes {
		go func(item string) {
			t := time.NewTicker(PL_RetryTime)
			defer func() {
				t.Stop()
			}()

			plInfo := make(map[string]*mediaInfo)
			for {
				plLines, err := getTSName(path.Join(livePath, item))
				if err == nil {
					for _, ts := range plLines.TS {
						if !ts.TS {
							continue
						}
						// if _, err := os.Stat(path.Join(livePath, ts.Line)); err == nil {
						plInfo[outputPlaylistPrefix+item] = &mediaInfo{} // данные временно не нужны
						plInfoC <- plInfo
						return
						// }
					}
				}
				select {
				case _, ok := <-done:
					if !ok {
						return
					}
				case <-t.C:
					continue
				}
			}
		}(item)
	}

	for _, item := range playlistTypes {
		select {
		case <-tcanel:
			level.Debug(p.logger).Log("getmediaInfo", item)
			return ErrOperationTimedout
		case inf := <-plInfoC:
			for k, v := range inf {
				plInfo[k] = v
			}
		}
	}

	plsToParse := []string{}
	for _, item := range playlistTypes {
		plsToParse = append(plsToParse, outputPlaylistPrefix+item)
	}
	level.Debug(p.logger).Log("genMasterPlaylistfor", fmt.Sprintf("%+v", plsToParse))
	if err := p.genMasterPlaylist(outputPlaylistPrefix+"index", plsToParse, livePath); err != nil {
		level.Error(p.logger).Log("ERROR gen master playlist %s", err)
		return err
	}

	return nil
}

func getMediaInfo(filePath string) (*mediaInfo, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, filePath)
	if err != nil {
		return nil, fmt.Errorf("file %s %v", filePath, err)
	}

	mi := &mediaInfo{
		Path:       filePath,
		Name:       data.Format.Filename,
		Bitrate:    data.Format.BitRate,
		Resolution: fmt.Sprintf("%dx%d", data.FirstVideoStream().Width, data.FirstVideoStream().Height),
		Codec:      data.FirstVideoStream().CodecName,
		Audio:      data.FirstAudioStream().CodecName,
	}

	return mi, err
}

func getTSName(filePath string) (pl *PLstruct, err error) {
	reTs := regexp.MustCompile(`.+\.ts`)
	reHd := regexp.MustCompile(`EXTINF`)

	readFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fScan := bufio.NewScanner(readFile)
	fScan.Split(bufio.ScanLines)

	lineNum := 1
	header := true
	pl = &PLstruct{}

	for fScan.Scan() {
		line := fScan.Text()
		if header && reHd.Match([]byte(line)) {
			header = false
		}

		if header {
			pl.Header = append(pl.Header, PLline{Num: lineNum, Line: fScan.Text()})
		} else {
			header = false
			tsFlag := false
			if reTs.Match([]byte(line)) {
				tsFlag = true
			}
			pl.TS = append(pl.TS, PLline{
				Num:  lineNum,
				Line: fScan.Text(),
				TS:   tsFlag,
			})
		}
		lineNum++
	}

	return
}

func (p *Playlist) genMasterPlaylist(key string, info []string, pl string) error {
	buf := bytes.NewBufferString("#EXTM3U\n") //#EXT-X-VERSION:3\n")
	for _, v := range info {
		if strings.Contains(v, "low") {
			buf.WriteString(`#EXT-X-STREAM-INF:BANDWIDTH=1100000,RESOLUTION=852x480,FRAME-RATE=25.000,CODECS="avc1.4d0028,mp4a.40.2",CLOSED-CAPTIONS=NONE` + "\n")
		} else if strings.Contains(v, "mid") {
			buf.WriteString(`#EXT-X-STREAM-INF:BANDWIDTH=2200000,RESOLUTION=1280x720,FRAME-RATE=25.000,CODECS="avc1.4d0028,mp4a.40.2",CLOSED-CAPTIONS=NONE` + "\n")
		} else if strings.Contains(v, "high") {
			buf.WriteString(`#EXT-X-STREAM-INF:BANDWIDTH=6500000,RESOLUTION=1920x1080,FRAME-RATE=25.000,CODECS="avc1.4d0028,mp4a.40.2",CLOSED-CAPTIONS=NONE` + "\n")
		}
		buf.WriteString(fmt.Sprintf("%s\n", v))
	}

	masterPlaylistName := fmt.Sprintf("%s.m3u8", key)

	f, err := os.OpenFile(path.Join(pl, masterPlaylistName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return err
}

func (p *Playlist) Delete(filePath string) error {
	re := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
	if !re.Match([]byte(filePath)) {
		return ErrInternalError
	}

	if err := os.RemoveAll(filePath); err != nil {
		return err
	}
	return nil
}
