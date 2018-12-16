package jobworker

import (
	"bitbucket.org/marketingx/upvideo/app/jobs"
	"bitbucket.org/marketingx/upvideo/app/videos"
	"bitbucket.org/marketingx/upvideo/app/videos/titles"
	"bitbucket.org/marketingx/upvideo/config"
	"bitbucket.org/marketingx/upvideo/utils"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type Service struct {
	Config       *config.Config
	VideoService *videos.Service
	TitleService *titles.Service
	JobService   *jobs.Service
}

func (this *Service) Run() {
	gocron.Every(1).Seconds().Do(this.work)
}

func (this *Service) work() {
	fmt.Printf("Job Worker has started work... \n")
	// check _jobs

	_jobs, err := this.JobService.FindAll(jobs.Params{ProcessId: "", Limit: 5})
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("Try to find idle _jobs err: \n%v\n", err)
		return
	}

	if err == sql.ErrNoRows {
		fmt.Println("Nothing to do")
		return
	}

	for _, _job := range _jobs {
		fmt.Printf("Process job id: %d\n", _job.Id)
		switch _job.Type {
		case "convert-title":
			err = this.processConvertTitleJob(_job)
			break
		case "upload-title":
			err = this.processUploadTitleJob(_job)
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
		}
	}
}

func (this *Service) processConvertTitleJob(job *jobs.Job) (err error) {
	_title, err := this.TitleService.FindOne(titles.Params{Id: job.RelatedId})
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(fmt.Sprintf("Related title id:'%d' not found", job.RelatedId))
		}
		return errors.New(fmt.Sprintf("Sql err: %v", err))
	}

	_video, err := this.VideoService.FindOne(videos.Params{Id: _title.VideoId})
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(fmt.Sprintf("Related video id:'%d' not found", _title.VideoId))
		}
		return errors.New(fmt.Sprintf("Sql err: %v", err))
	}

	if _video.File == "" {
		return errors.New(fmt.Sprintf("Video id:'%d' file is empty", _video.Id))
	}

	// download from url
	if _video.TmpFile != "" {
		// look for unique filename
		var unic string
		for {
			unic = utils.RandomString(32)
			_video.TmpFile = path.Join(this.Config.YoutubeUploaderDirs.TempVideosDir, filepath.Base(_video.File)+"_"+unic+filepath.Ext(_video.File))

			if _, err := os.Stat(_video.TmpFile); os.IsNotExist(err) {
				break
			}
		}

		out, err := os.Create(_video.TmpFile)
		if err != nil {
			return errors.New(fmt.Sprintf("Create _video.TmpFile err: %v", err))
		}
		defer out.Close()

		resp, err := http.Get(_video.File)
		if err != nil {
			return errors.New(fmt.Sprintf("Download _video.File url: %s, err: %v", err, _video.File))
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("Save _video.TmpFile path: %s, err: %v", err, _video.TmpFile))
		}
	}

	//'ffmpeg -i {} -vf scale="{}:trunc(ow/a/2)*2" -r {} {}'.format(argv.video_file, resolution_size, frame_rate, file_output)
	cmd := exec.Command("ffmpeg", "-i", _video.TmpFile, "-vf", fmt.Sprintf(`scale="%d:trunc(ow/a/2)*2"`, _title.Resolution), "-r", string(_title.FrameRate), _title.File)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		fmt.Println("ffmpeg execute Error: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return errors.New(fmt.Sprintf("ffmpeg err: %v", err))
	}

	job.ProcessId = cmd.Process.Pid
	err = this.JobService.Update(job)
	if err != nil {
		err2 := cmd.Process.Kill()
		if err2 != nil {

		}
		return errors.New(fmt.Sprintf("Sql JobService.Update err: %v", err))
	}

	// this is important, otherwise the process becomes in S mode
	go func() {
		err = cmd.Wait()
		fmt.Printf("Command finished with error: %v", err)
	}()

	if err != nil {
		fmt.Println("Youtube uploader Error: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return errors.New(fmt.Sprintf("ffmpeg err: %v", err))
	}
	if !uploadSuccessfulRegexp.Match(out.Bytes()) {
		fmt.Println("Youtube uploader Result: " + out.String())
		c.Status(http.StatusInternalServerError)
		return
	}

	matches := uploadSuccessfulRegexp.FindStringSubmatch(out.String())
	url := "https://www.youtube.com/watch?v=" + matches[1]

}

func (this *Service) processUploadTitleJob(job *jobs.Job) (err error) {

	cmd := exec.Command(this.Config.YoutubeUploaderCmd, "-headlessAuth", "-secrets", clientSecretsPath, "-cache", tokenPath, "-filename", this.Config.TestVideoPath, "-metaJSON", this.Config.TestVideoMetaPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Youtube uploader Error: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		c.Status(http.StatusInternalServerError)
		return
	}
	if !uploadSuccessfulRegexp.Match(out.Bytes()) {
		fmt.Println("Youtube uploader Result: " + out.String())
		c.Status(http.StatusInternalServerError)
		return
	}

	matches := uploadSuccessfulRegexp.FindStringSubmatch(out.String())
	url := "https://www.youtube.com/watch?v=" + matches[1]
}
