package util

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

var addrRegex = regexp.MustCompile("^(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5]):([0-9]|[1-9]\\d|[1-9]\\d{2}|[1-9]\\d{3}|[1-5]\\d{4}|6[0-4]\\d{3}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])$")

// Empty checks if a string referenced by s or s itself is empty.
func Empty(s *string) bool {
	return s == nil || *s == ""
}

func TransIntToCstTime(timeMs int64) string {
	localTime := time.Unix(timeMs/1000, 0)
	//获取本地时区

	cstLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logs.Error("get cst location err, %v", err)
		return ""
	}

	localTime = localTime.In(cstLocation)
	return localTime.Format("2006-01-02 15:04:05")
}

func UUID() string {
	id := uuid.NewV4()
	ids := id.String()
	return ids
}

func IsBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func SaltEncode(password string, salt ...interface{}) string {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		password = fmt.Sprintf(password+strings.Join(slice, "%v"), salt...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func WriteFile(content string, filename string) (err error) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		logs.Error(err.Error())
		return err
	} else {
		_, err := f.Write([]byte(content))
		if err != nil {
			return err
		}
	}
	return nil
}

func ReloadNgninx(dockerSock, command, container string) (err error) {
	cli, err := client.NewClientWithOpts(client.WithHost(dockerSock), client.WithAPIVersionNegotiation())
	if err != nil {
		logs.Error("error: could not create docker client handle")
		return err
	}

	ctx := context.Background()

	config := types.ExecConfig{
		Cmd:          strings.Split(command, " "),
		AttachStdout: true,
		AttachStderr: true,
	}

	response, err := cli.ContainerExecCreate(ctx, container, config)
	if err != nil {
		return err
	}
	execID := response.ID

	resp, err := cli.ContainerExecAttach(ctx, execID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	defer resp.Close()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	_, err = stdcopy.StdCopy(stdout, stderr, resp.Reader)
	if err != nil {
		return err
	}
	logs.Info(stdout.String())
	logs.Error(stderr.String())
	return nil
}
