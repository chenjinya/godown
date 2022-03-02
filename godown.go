package godown

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/chenjinya/loji"
)

type HTTPFileReader struct {
	io.Reader
	Filename string
	Total    int64
	Current  int64
	lo       *loji.LoadingEmoji
}

func BeautifulSize(n int64) string {
	if n > 1024*1024*1024 {
		return fmt.Sprintf("%.2fGB", float64(n)/1024/1024/1024)
	}
	if n > 1024*1024 {
		return fmt.Sprintf("%.2fMB", float64(n)/1024/1024)
	}
	if n > 1024 {
		return fmt.Sprintf("%.2fKB", float64(n)/1024)
	}
	if n >= 0 {
		return fmt.Sprintf("%dB", n)
	}
	return "0B"
}
func (r *HTTPFileReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.Current += int64(n)

	if r.Current > r.Total {
		r.Current = r.Total
	}

	r.lo.Loading(fmt.Sprintf("下载 %s 进度 %.2f%% %s/%s",
		r.Filename, float64(r.Current*10000/r.Total)/100,
		BeautifulSize(r.Current), BeautifulSize(r.Total),
	))
	return
}

func Download(url, filename string) (err error) {
	err = retry.Do(func() error {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%s", e)
			}
		}()
		r, err := http.Get(url)
		if err != nil {
			return err
		}
		defer func() { _ = r.Body.Close() }()

		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer func() { _ = f.Close() }()

		reader := &HTTPFileReader{
			Reader:   r.Body,
			Filename: filename,
			Total:    r.ContentLength,
			lo:       loji.New(),
		}
		defer func() {
			reader.lo.Stop()
		}()
		_, _ = io.Copy(f, reader)
		return nil
	}, retry.Attempts(3), retry.Delay(time.Second*2), retry.DelayType(retry.BackOffDelay))

	return err

}
