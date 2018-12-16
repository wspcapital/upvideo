package jobworker

import (
	"bitbucket.org/marketingx/upvideo/app/accounts"
	"bitbucket.org/marketingx/upvideo/app/jobs"
	"bitbucket.org/marketingx/upvideo/app/videos"
	"bitbucket.org/marketingx/upvideo/app/videos/titles"
	"bitbucket.org/marketingx/upvideo/config"
	"bitbucket.org/marketingx/upvideo/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var (
	uploadSuccessfulRegexp = regexp.MustCompile("Upload successful! Video ID: ([a-zA-Z0-9]+)")
	uploadErrorRegexp      = regexp.MustCompile("Error making YouTube API call:([a-zA-Z0-9,.!_\\-\\s]+)")
)

type Service struct {
	Config         *config.Config
	VideoService   *videos.Service
	TitleService   *titles.Service
	JobService     *jobs.Service
	AccountService *accounts.Service
}

func (this *Service) Start() {
	go func() {
		gocron.Every(1).Seconds().Do(this.work)

		<-gocron.Start()
	}()
}

func (this *Service) work() {
	fmt.Printf("Job Worker has started work... \n")
	// check _jobs

	_jobs, err := this.JobService.FindAll(jobs.Params{ProcessId: "", Limit: 5})
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("Try to find idle _jobs err: \n%v\n", err)
		return
	}

	if _jobs == nil || err == sql.ErrNoRows {
		fmt.Println("Nothing to do")
		return
	}

	for _, _job := range _jobs {
		fmt.Printf("Process job id: %d\n", _job.Id)
		var _title *titles.Title

		switch _job.Type {
		case "convert-title":
			_title, err = this.processConvertTitleJob(_job)
			break
		case "upload-title":
			_title, err = this.processUploadTitleJob(_job)
			break
		default:
			fmt.Printf("Unknown job type: %s\n", _job.Type)
			break
		}

		if err != nil {
			fmt.Printf("Job id:'%d' failed.\n\tError:\n%v\n", _job.Id, err)
			err = this.JobService.JobFailed(_job, err.Error())
			if err != nil {
				fmt.Printf("JobService.JobFailed err: \n%v\n", err)
			}
		} else {
			// set pending if success to show failed status by title request
			_title.Pending = false
			err = this.TitleService.Update(_title)
			if err != nil {
				fmt.Printf("SQL error, title update. \n\tError:\n%v\n", err)
			}

			err = this.JobService.JobCompleted(_job)
			if err != nil {
				fmt.Printf("SQL error Job id:'%d', can not be deleted.\n\tError:\n%v\n", _job.Id, err)
			}
		}
	}
}

func (this *Service) processConvertTitleJob(job *jobs.Job) (_title *titles.Title, err error) {
	_title, err = this.TitleService.FindOne(titles.Params{Id: job.RelatedId})
	if err != nil {
		if err == sql.ErrNoRows {
			return _title, errors.New(fmt.Sprintf("Related title id:'%d' not found", job.RelatedId))
		}
		return _title, errors.New(fmt.Sprintf("Sql err: %v", err))
	}

	_video, err := this.VideoService.FindOne(videos.Params{Id: _title.VideoId})
	if err != nil {
		if err == sql.ErrNoRows {
			return _title, errors.New(fmt.Sprintf("Related video id:'%d' not found", _title.VideoId))
		}
		return _title, errors.New(fmt.Sprintf("Sql err: %v", err))
	}

	if _video.File == "" {
		return _title, errors.New(fmt.Sprintf("Video id:'%d' file is empty", _video.Id))
	}

	// download from url
	if _video.TmpFile == "" {
		_video.TmpFile, _ = utils.FindUniqueFileNameInPath(_video.File, this.Config.YoutubeUploaderDirs.TempVideosDir)

		out, err := os.Create(_video.TmpFile)
		if err != nil {
			return _title, errors.New(fmt.Sprintf("Create _video.TmpFile err: %v", err))
		}
		defer out.Close()

		resp, err := http.Get(_video.File)
		if err != nil {
			return _title, errors.New(fmt.Sprintf("Download _video.File url: %s, err: %v", err, _video.File))
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return _title, errors.New(fmt.Sprintf("Save _video.TmpFile path: %s, err: %v", err, _video.TmpFile))
		}

		err = this.VideoService.Update(_video)
		if err != nil {
			return _title, errors.New(fmt.Sprintf("VideoService.Update Sql err: %v", err))
		}
	}

	_title.TmpFile, _ = utils.FindUniqueFileNameInPath(_title.GetPreparedFilename(), this.Config.YoutubeUploaderDirs.TitlesVideosDir)
	// force overwrite file
	cmd := exec.Command("ffmpeg", "-i", _video.TmpFile, "-vf", fmt.Sprintf("scale=%d:trunc(ow/a/2)*2", _title.Resolution), "-r", fmt.Sprintf("%d", _title.FrameRate), _title.TmpFile, "-y")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		fmt.Println("ffmpeg execute Error: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return _title, errors.New(fmt.Sprintf("ffmpeg err: %v", err))
	}

	job.ProcessId = cmd.Process.Pid
	err = this.JobService.Update(job)
	if err != nil {
		err2 := cmd.Process.Kill()
		if err2 != nil {
			fmt.Printf("Can not kill ffmpeg process. \n\tError:\n%v\n", err)
		}
		return _title, errors.New(fmt.Sprintf("Sql JobService.Update err: %v", err))
	}

	// stuff to update progress
	//var stdoutBuf, stderrBuf bytes.Buffer
	//stdoutIn, _ := cmd.StdoutPipe()
	//stderrIn, _ := cmd.StderrPipe()
	//
	//var errStdout, errStderr error
	//stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	//stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	//err := cmd.Start()
	//if err != nil {
	//	log.Fatalf("cmd.Start() failed with '%s'\n", err)
	//}
	//
	//go func() {
	//	_, errStdout = io.Copy(stdout, stdoutIn)
	//}()
	//
	//go func() {
	//	_, errStderr = io.Copy(stderr, stderrIn)
	//}()
	//
	//err = cmd.Wait()
	//if err != nil {
	//	log.Fatalf("cmd.Start() failed with %s\n", err)
	//}
	//if errStdout != nil || errStderr != nil {
	//	log.Fatal("failed to capture stdout or stderr\n")
	//}
	//outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())

	err = cmd.Wait()
	if err != nil {
		fmt.Println("ffmpeg Error: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return _title, errors.New(fmt.Sprintf("ffmpeg output: %s, err: %v, ssterr: %s", out.String(), err, stderr.String()))
	}

	// no error - all is ok, check file exists
	if _, err := os.Stat(_title.TmpFile); os.IsNotExist(err) {
		return _title, errors.New(fmt.Sprintf("ffmpeg finished ok, but file not exists. output: %v, err: %v, ssterr: %s", out.String(), err, stderr.String()))
	}

	_title.Converted = true

	return
}

func (this *Service) processUploadTitleJob(job *jobs.Job) (_title *titles.Title, err error) {
	_title, err = this.TitleService.FindOne(titles.Params{Id: job.RelatedId})
	if err != nil {
		if err == sql.ErrNoRows {
			return _title, errors.New(fmt.Sprintf("Related title id:'%d' not found", job.RelatedId))
		}
		return _title, errors.New(fmt.Sprintf("Sql err: %v", err))
	}

	_video, err := this.VideoService.FindOne(videos.Params{Id: _title.VideoId})
	if err != nil {
		if err == sql.ErrNoRows {
			return _title, errors.New(fmt.Sprintf("Related video id:'%d' not found", _title.VideoId))
		}
		return _title, errors.New(fmt.Sprintf("Sql err: %v", err))
	}

	_account, err := this.AccountService.FindOne(accounts.Params{Id: _video.AccountId})
	if err != nil {
		if err == sql.ErrNoRows {
			return _title, errors.New(fmt.Sprintf("Related account id:'%d' not found", _video.AccountId))
		}
		return _title, errors.New(fmt.Sprintf("Sql err: %v", err))
	}

	// prepare meta
	tmpDescription := strings.Replace(_video.Description, "{title}", _title.Title, -1)
	tmpDescription = strings.Replace(tmpDescription, "{tags}", _title.Tags, -1)

	metadata := map[string]interface{}{
		"title":               _title.Title,
		"description":         tmpDescription,
		"tags":                _title.Tags,
		"privacyStatus":       "public",
		"embeddable":          true,
		"license":             "creativeCommon",
		"publicStatsViewable": true,
		"categoryId":          "28",
		"language":            "EN",
	}

	metaFilePath := path.Join(this.Config.YoutubeUploaderDirs.MetasDir, fmt.Sprintf("job_upload_titleid_%d,json", _title.Id))
	f, err := os.OpenFile(metaFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return _title, errors.New(fmt.Sprintf("Can not open file: %s, err: %v", metaFilePath, err))
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return _title, errors.New(fmt.Sprintf("json.Marshal err: %v", err))
	}

	_, err = f.Write(b)
	err = f.Close()
	if err != nil {
		return _title, errors.New(fmt.Sprintf("Can not write file: %s, err: %v", metaFilePath, err))
	}

	cmd := exec.Command(this.Config.YoutubeUploaderCmd, "-headlessAuth", "-secrets", _account.ClientSecrets, "-cache", _account.RequestToken, "-filename", _title.TmpFile, "-metaJSON", metaFilePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		fmt.Println("Youtube uploader Error: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return _title, errors.New(fmt.Sprintf("Youtube uploader Error: %v", err))
	}

	job.ProcessId = cmd.Process.Pid
	err = this.JobService.Update(job)
	if err != nil {
		err2 := cmd.Process.Kill()
		if err2 != nil {
			fmt.Printf("Can not kill YoutubeUploader process. \n\tError:\n%v\n", err)
		}
		return _title, errors.New(fmt.Sprintf("Sql JobService.Update err: %v", err))
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Youtube uploader Error: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return _title, errors.New(fmt.Sprintf("Youtube uploader output: %v, err: %v, ssterr: %s", out.String(), err, stderr.String()))
	}

	if !uploadSuccessfulRegexp.Match(out.Bytes()) {
		fmt.Println("Youtube uploader is not successful. Result: " + out.String())
		return _title, errors.New(fmt.Sprintf("Youtube uploader is not successful, output: %v, err: %v, ssterr: %s", out.String(), err, stderr.String()))
	}

	matches := uploadSuccessfulRegexp.FindStringSubmatch(out.String())
	_title.YoutubeId = matches[1]
	_title.Posted = true

	_ = os.Remove(metaFilePath)
	_ = os.Remove(_title.TmpFile)

	return
}
