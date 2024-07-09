/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/hertz-contrib/thrift-gen-mongo/consts"

	gv "github.com/hashicorp/go-version"
)

// PathExist is used to judge whether the path exists in file system.
func PathExist(path string) (bool, error) {
	abPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(abPath)
	if err != nil {
		return os.IsExist(err), nil
	}
	return true, nil
}

func ReadFileContent(filePath string) (content []byte, err error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

// CamelString converts the string 's' to a camel string
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return Bytes2Str(data[:])
}

func Bytes2Str(in []byte) (out string) {
	op := (*reflect.StringHeader)(unsafe.Pointer(&out))
	ip := (*reflect.SliceHeader)(unsafe.Pointer(&in))
	op.Data = ip.Data
	op.Len = ip.Len
	return
}

func UnpackArgs(args []string, c interface{}) error {
	m, err := MapForm(args)
	if err != nil {
		return fmt.Errorf("unmarshal args failed, err: %v", err.Error())
	}

	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("passed c must be struct or pointer of struct")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		x := v.Field(i)
		n := f.Name
		values, ok := m[n]
		if !ok || len(values) == 0 || values[0] == "" {
			continue
		}
		switch x.Kind() {
		case reflect.Bool:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			x.SetBool(values[0] == "true")
		case reflect.String:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			x.SetString(values[0])
		case reflect.Slice:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			ss := strings.Split(values[0], ";")
			if x.Type().Elem().Kind() == reflect.Int {
				n := reflect.MakeSlice(x.Type(), len(ss), len(ss))
				for i, s := range ss {
					val, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return err
					}
					n.Index(i).SetInt(val)
				}
				x.Set(n)
			} else {
				for _, s := range ss {
					val := reflect.Append(x, reflect.ValueOf(s))
					x.Set(val)
				}
			}
		case reflect.Map:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			ss := strings.Split(values[0], ";")
			out := make(map[string]string, len(ss))
			for _, s := range ss {
				sk := strings.SplitN(s, "=", 2)
				if len(sk) != 2 {
					return fmt.Errorf("map filed %v invalid key-value pair '%v'", n, s)
				}
				out[sk[0]] = sk[1]
			}
			x.Set(reflect.ValueOf(out))
		default:
			return fmt.Errorf("field %s has unsupported type %+v", n, f.Type)
		}
	}
	return nil
}

func MapForm(input []string) (map[string][]string, error) {
	out := make(map[string][]string, len(input))

	for _, str := range input {
		parts := strings.SplitN(str, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid argument: '%s'", str)
		}
		key, val := parts[0], parts[1]
		out[key] = append(out[key], val)
	}

	return out, nil
}

func CreateFile(path, content string) (err error) {
	return os.WriteFile(path, []byte(content), os.FileMode(0o644))
}

const ThriftgoMiniVersion = "v0.2.0"

// QueryVersion will query the version of the corresponding executable.
func QueryVersion(exe string) (version string, err error) {
	var buf strings.Builder
	cmd := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe, "--version",
		},
		Stdin:  os.Stdin,
		Stdout: &buf,
		Stderr: &buf,
	}
	err = cmd.Run()
	if err == nil {
		version = strings.Split(buf.String(), " ")[1]
		if strings.HasSuffix(version, "\n") {
			version = version[:len(version)-1]
		}
	}
	return
}

func GetBuildGoPaths() []string {
	var all []string
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		if p == "" || p == build.Default.GOROOT {
			continue
		}
		if strings.HasPrefix(p, "~") {
			continue
		}
		all = append(all, p)
	}
	for k, v := range all {
		if strings.HasSuffix(v, "/") || strings.HasSuffix(v, string(os.PathSeparator)) {
			v = v[:len(v)-1]
		}
		all[k] = v
	}
	return all
}

func GetGOPATH() (gopath string, err error) {
	ps := filepath.SplitList(os.Getenv("GOPATH"))
	if len(ps) > 0 {
		gopath = ps[0]
	}
	if gopath == "" {
		cmd := exec.Command("go", "env", "GOPATH")
		var out bytes.Buffer
		cmd.Stderr = &out
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			gopath = strings.Trim(out.String(), " \t\n\r")
		}
	}
	if gopath == "" {
		ps := GetBuildGoPaths()
		if len(ps) > 0 {
			gopath = ps[0]
		}
	}
	isExist, err := PathExist(gopath)
	if !isExist {
		return "", err
	}
	return strings.Replace(gopath, "/", string(os.PathSeparator), -1), nil
}

func LookupTool(idlType string) (string, error) {
	tool := "thriftgo"
	if idlType == "proto" {
		tool = "protoc"
	}

	path, err := exec.LookPath(tool)
	if err != nil {
		goPath, err := GetGOPATH()
		if err != nil {
			return "", fmt.Errorf("get 'GOPATH' failed for find %s : %v", tool, path)
		}
		path = filepath.Join(goPath, "bin", tool)
	}

	isExist, err := PathExist(path)
	if err != nil {
		return "", fmt.Errorf("check '%s' path error: %v", path, err)
	}

	if !isExist {
		if tool == "thriftgo" {
			// If thriftgo does not exist, the latest version will be installed automatically.
			err := InstallAndCheckThriftgo()
			if err != nil {
				return "", fmt.Errorf("can't install '%s' automatically, please install it manually for https://github.com/cloudwego/thriftgo, err : %v", tool, err)
			}
		} else {
			return "", fmt.Errorf("%s is not installed, please install it first", tool)
		}
	}

	if tool == "thriftgo" {
		// If thriftgo exists, the version is detected; if the version is lower than v0.2.0 then the latest version of thriftgo is automatically installed.
		err := CheckAndUpdateThriftgo()
		if err != nil {
			return "", fmt.Errorf("update thriftgo version failed, please install it manually for https://github.com/cloudwego/thriftgo, err: %v", err)
		}
	}

	return path, nil
}

// CheckAndUpdateThriftgo checks the version of thriftgo and updates the tool to the latest version if its version is less than v0.2.0.
func CheckAndUpdateThriftgo() error {
	path, err := exec.LookPath("thriftgo")
	if err != nil {
		return fmt.Errorf("can not find %s", "thriftgo")
	}
	curVersion, err := QueryVersion(path)
	log.Printf("current thriftgo version is %s", curVersion)
	if ShouldUpdate(curVersion, ThriftgoMiniVersion) {
		log.Println(" current thriftgo version is less than v0.2.0, so update thriftgo version")
		err = InstallAndCheckThriftgo()
		if err != nil {
			return fmt.Errorf("update thriftgo version failed, err: %v", err)
		}
	}

	return nil
}

// InstallAndCheckThriftgo will automatically install thriftgo and judge whether it is installed successfully.
func InstallAndCheckThriftgo() error {
	exe, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("can not find tool 'go': %v", err)
	}
	var buf strings.Builder
	cmd := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe, "install", "github.com/cloudwego/thriftgo@latest",
		},
		Stdin:  os.Stdin,
		Stdout: &buf,
		Stderr: &buf,
	}

	done := make(chan error)
	log.Println("installing thriftgo automatically")
	go func() {
		done <- cmd.Run()
	}()
	select {
	case err = <-done:
		if err != nil {
			return fmt.Errorf("can not install thriftgo, err: %v. Please install it manual, and make sure the version of thriftgo is greater than v0.2.0", cmd.Stderr)
		}
	case <-time.After(time.Second * 30):
		return fmt.Errorf("install thriftgo time out.Please install it manual, and make sure the version of thriftgo is greater than v0.2.0")
	}

	exist, err := CheckCompiler("thriftgo")
	if err != nil {
		return fmt.Errorf("check %s exist failed, err: %v", "thriftgo", err)
	}
	if !exist {
		return fmt.Errorf("install thriftgo failed. Please install it manual, and make sure the version of thriftgo is greater than v0.2.0")
	}

	return nil
}

func InstallAndCheckMongoPlugin() error {
	exe, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("can not find tool 'go': %v", err)
	}
	var buf strings.Builder
	cmd := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe, "install", "github.com/hertz-contrib/thrift-gen-mongo@latest",
		},
		Stdin:  os.Stdin,
		Stdout: &buf,
		Stderr: &buf,
	}

	done := make(chan error)
	log.Println("installing thrift-gen-mongo automatically")
	go func() {
		done <- cmd.Run()
	}()
	select {
	case err = <-done:
		if err != nil {
			return fmt.Errorf("can not install thrift-gen-mongo, err: %v. Please install it manual, and make sure the version of thrift-gen-mongo is greater than v0.2.0", cmd.Stderr)
		}
	case <-time.After(time.Second * 30):
		return fmt.Errorf("install thrift-gen-mongo time out.Please install it manual, and make sure the version of thrift-gen-mongo is greater than v0.2.0")
	}

	exist, err := CheckCompiler("thrift-gen-mongo")
	if err != nil {
		return fmt.Errorf("check %s exist failed, err: %v", "thrift-gen-mongo", err)
	}
	if !exist {
		return fmt.Errorf("install thrift-gen-mongo failed. Please install it manual, and make sure the version of thrift-gen-mongo is greater than v0.2.0")
	}

	return nil
}

func ShouldUpdate(current, latest string) bool {
	cv, err := gv.NewVersion(current)
	if err != nil {
		return false
	}
	lv, err := gv.NewVersion(latest)
	if err != nil {
		return false
	}

	return cv.Compare(lv) < 0
}

func CheckCompiler(tool string) (bool, error) {
	path, err := exec.LookPath(tool)
	if err != nil {
		goPath, err := GetGOPATH()
		if err != nil {
			return false, fmt.Errorf("get 'GOPATH' failed for find %s : %v", tool, path)
		}
		path = filepath.Join(goPath, "bin", tool)
	}

	isExist, err := PathExist(path)
	if err != nil {
		return false, fmt.Errorf("can not check %s exist, err: %v", tool, err)
	}
	if !isExist {
		return false, nil
	}

	return true, nil
}

func commandAndNotice(cmd, notice string) {
	argv := strings.Split(cmd, consts.BlackSpace)
	err := exec.Command(argv[0], argv[1:]...).Run()

	res := "Done"
	if err != nil {
		res = err.Error()
	}
	log.Println(notice, res)
}

func ReplaceThriftVersion() {
	cmd := "go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0"
	notice := "Adding apache/thrift@v0.13.0 to go.mod for generated code .........."
	commandAndNotice(cmd, notice)
}

var goModReg = regexp.MustCompile(`^\s*module\s+(\S+)\s*`)

func SearchGoMod(cwd string, recurse bool) (moduleName, path string, found bool) {
	for {
		path = filepath.Join(cwd, "go.mod")
		data, err := os.ReadFile(path)
		if err == nil {
			for _, line := range strings.Split(string(data), consts.LineBreak) {
				m := goModReg.FindStringSubmatch(line)
				if m != nil {
					return m[1], cwd, true
				}
			}
			return fmt.Sprintf("<module name not found in '%s'>", path), path, true
		}

		if !os.IsNotExist(err) {
			return
		}
		if !recurse {
			break
		}
		cwd = filepath.Dir(cwd)
		// the root directory will return itself by using "filepath.Dir()"; to prevent dead loops, so jump out
		if cwd == filepath.Dir(cwd) {
			break
		}
	}
	return
}

// GetIdlType is used to return the idl type.
func GetIdlType(path string, pbName ...string) (string, error) {
	ext := filepath.Ext(path)
	if ext == "" || ext[0] != '.' {
		return "", fmt.Errorf("idl path %s is not a valid file", path)
	}
	ext = ext[1:]
	switch ext {
	case consts.IdlThrift:
		return consts.IdlThrift, nil
	case consts.IdlProto:
		if len(pbName) > 0 {
			return pbName[0], nil
		}
		return consts.IdlProto, nil
	default:
		return "", fmt.Errorf("IDL type %s is not supported", ext)
	}
}

func InitGoMod(module string) error {
	isExist, err := PathExist(consts.GoMod)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}
	gg, err := exec.LookPath(consts.Go)
	if err != nil {
		return err
	}
	cmd := &exec.Cmd{
		Path:   gg,
		Args:   []string{consts.Go, consts.Mod, consts.Init, module},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmd.Run()
}

// IsWindows determines whether the current operating system is Windows
func IsWindows() bool {
	return consts.SysType == consts.WindowsOS
}
