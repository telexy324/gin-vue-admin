package common

import (
	"github.com/jlaffaye/ftp"
	"io"
	"strconv"
	"time"
)

type FtpClient struct {
	Conn *ftp.ServerConn
}

func NewFtpClient(host string, port int, username, password string) (client FtpClient, err error) {
	addr := host + strconv.Itoa(port)
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return
	}
	err = c.Login(username, password)
	if err != nil {
		return
	}
	client.Conn = c
	return
}

func (f *FtpClient) Upload(path string, file io.Reader) error {
	return f.Conn.Stor(path, file)
}

func (f *FtpClient) Download(path string) (fileBytes []byte, err error) {
	response, err := f.Conn.Retr(path)
	if err != nil {
		return
	}
	defer response.Close()
	return io.ReadAll(response)
}
