package pkgep

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

const fileRootDir = "assets"

type FileType int

const (
	ImageType FileType = iota
	VideoType
	DocType
)

type FsErrorLevel int

const (
	Fatal FsErrorLevel = 1
	Chill FsErrorLevel = 2
)

type FsError struct {
	Code    int
	Level   FsErrorLevel
	Message string
}

func (e FsError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// 💡move file from root dir to destination path, return success file name.
func MoveFile(f string, dstDir string) (string, error) {

	src, err := os.Open(filepath.Join(fileRootDir, f))
	// err case file does not exist in root dir.
	if err != nil {
		return "", FsError{Code: 2, Level: Chill, Message: err.Error()}
	}
	defer src.Close()

	// Check and Create Destination Directory
	if _, err := os.Stat(filepath.Join(fileRootDir, dstDir)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(fileRootDir, dstDir), os.ModePerm); err != nil {
			return "", FsError{Code: 1, Level: Fatal, Message: err.Error()}
		}
	}

	// Create New File in Destination Directory.
	dstSrc, err := os.Create(filepath.Join(fileRootDir, dstDir, f))
	// err case file does not exist in root dir.
	if err != nil {
		return "", FsError{Code: 1, Level: Fatal, Message: err.Error()}
	}
	defer dstSrc.Close()

	// copy the content from source to destination.
	_, err = io.Copy(dstSrc, src)
	if err != nil {
		return "", FsError{Code: 1, Level: Fatal, Message: err.Error()}
	}

	src.Close() // for Windows, close before trying to remove file.

	// Remove Source File.
	err = os.Remove(filepath.Join(fileRootDir, f))
	if err != nil {
		//return "", fmt.Errorf("couldn't remove source file: %v", err)
		return "", FsError{Code: 2, Level: Chill, Message: fmt.Sprintf("couldn't remove source file: %v", err)}
	}

	return f, nil
}

func RemoveFile(fileName string, dirName string) error {

	sourcePath := filepath.Join(fileRootDir, dirName, fileName)

	err := os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't remove source file: %v", err)
	}

	return nil
}
func RemoveOneFile(fileName string) error {

	sourcePath := filepath.Join(fileRootDir, fileName)

	err := os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't remove source file: %v", err)
	}

	return nil
}

// ? if it cannot find a file ?
func FileUpdateManager(destDir string, files string) (string, error) {

	claimList := strings.Split(files, "|")
	if len(claimList) == 0 {
		return "", nil
	}

	// Map to keep track of what to be deleted and to be move in.
	shouldMoveFileInMap := make(map[string]bool)
	shouldDeleteFileMap := make(map[string]bool)

	for i := range claimList {
		// Split Out Paths And Get The Last Index File.
		fpath := strings.Split(claimList[i], "/")

		if len(fpath) > 0 {
			idx := fpath[len(fpath)-1]
			claimList[i] = idx
			// Set Default Of Should Move Map to true.
			shouldMoveFileInMap[idx] = true
		}
	}

	// Get Access To Target Directory.
	entries, err := os.ReadDir(filepath.Join(fileRootDir, destDir))
	// Dir Not found In Case This Specific Object Has Not Been Created File in its Creation.
	if err != nil {
		logrus.Printf("directory: %s not found, creating %s/%s/", destDir, fileRootDir, destDir)

		// Create new directory.
		newDirPath := filepath.Join(fileRootDir, destDir)

		// Move files in, no delete for new create dir.
		err := os.MkdirAll(newDirPath, os.ModePerm)
		if err != nil {
			logrus.Error(err)
			// the only error case to be return in this function.
			return "", err
		}

		// move new uploaded files into new directory, since the directory does not exist, all files are meant to be new uploaded.
		for _, f := range claimList {
			_, err := MoveFile(f, destDir)
			if err != nil {
				log.Println(err)
			}
		}

		// get stringify file names
		var ret []string
		entries, _ := os.ReadDir(filepath.Join(fileRootDir, destDir))

		// iterate over target directory and get its file name stringified.
		for _, entry := range entries {
			if !entry.IsDir() {
				n := fmt.Sprintf("%s/%s", destDir, entry.Name())
				ret = append(ret, n)
			}
		}
		return strings.Join(ret, "|"), nil
	}

	// in case directory already exist.
	var existList []string

	// map file that does not exist in request parameter but exist in file system.

	for _, entry := range entries {
		if !entry.IsDir() {
			existList = append(existList, entry.Name())
			// Set Default Should Delete Map to true.
			shouldDeleteFileMap[entry.Name()] = true
		}
	}

	// Should not delete if file exist in both places, claimList and existList.
	for _, i := range claimList {
		shouldDeleteFileMap[i] = false
	}

	// if fs file exists, set request file to false
	for _, i := range existList {
		shouldMoveFileInMap[i] = false
	}

	// Move new file into target dir
	for fname, shouldMove := range shouldMoveFileInMap {
		if shouldMove {
			fmt.Printf("moving file: %s\n", fname)
			_, err := MoveFile(fname, destDir)
			if err != nil {
				logrus.Printf("error moving file %s err: %v\n", fname, err)
			} else {
				logrus.Printf("success moved %s\n", fname)
			}
		}
	}

	// Delete Exlcluded Files
	for fname, shouldDelete := range shouldDeleteFileMap {
		if shouldDelete {
			fmt.Printf("deleting file: %s\n", fname)
			err := RemoveFile(fname, destDir)
			if err != nil {
				logrus.Printf("error deleting file %s err: %v\n", fname, err)
			} else {
				logrus.Printf("success deleted %s\n", fname)
			}
		}
	}

	// Get to dir and get all the files name to stringify it for return
	var ret []string
	ent, _ := os.ReadDir(filepath.Join(fileRootDir, destDir))
	for _, n := range ent {
		if !n.IsDir() {
			ret = append(ret, fmt.Sprintf("%s/%s", destDir, n.Name()))
		}
	}

	return strings.Join(ret, "|"), nil
}

func FileUpdateManagerX(destDir string, claimList []string) ([]string, error) {

	// Map to keep track of what to be deleted and to be move in.
	shouldMoveFileInMap := make(map[string]bool)
	shouldDeleteFileMap := make(map[string]bool)

	for i := range claimList {
		// Split Out Paths And Get The Last Index File.
		fpath := strings.Split(claimList[i], "/")

		if len(fpath) > 0 {
			idx := fpath[len(fpath)-1]
			claimList[i] = idx
			// Set Default Of Should Move Map to true.
			shouldMoveFileInMap[idx] = true
		}
	}

	// Get Access To Target Directory.
	entries, err := os.ReadDir(filepath.Join(fileRootDir, destDir))
	// Dir Not found In Case This Specific Object Has Not Been Created File in its Creation.
	if err != nil {
		logrus.Printf("directory: %s not found, creating %s/%s/", destDir, fileRootDir, destDir)

		// Create new directory.
		newDirPath := filepath.Join(fileRootDir, destDir)

		// Move files in, no delete for new create dir.
		err := os.MkdirAll(newDirPath, os.ModePerm)
		if err != nil {
			logrus.Error(err)
			// the only error case to be return in this function.
			return []string{}, err
		}

		// move new uploaded files into new directory, since the directory does not exist, all files are meant to be new uploaded.
		for _, f := range claimList {
			_, err := MoveFile(f, destDir)
			if err != nil {
				log.Println(err)
			}
		}

		// get stringify file names
		var ret []string
		entries, _ := os.ReadDir(filepath.Join(fileRootDir, destDir))

		// iterate over target directory and get its file name stringified.
		for _, entry := range entries {
			if !entry.IsDir() {
				n := fmt.Sprintf("%s/%s", destDir, entry.Name())
				ret = append(ret, n)
			}
		}
		return ret, nil
	}

	// in case directory already exist.
	var existList []string

	// map file that does not exist in request parameter but exist in file system.

	for _, entry := range entries {
		if !entry.IsDir() {
			existList = append(existList, entry.Name())
			// Set Default Should Delete Map to true.
			shouldDeleteFileMap[entry.Name()] = true
		}
	}

	// Should not delete if file exist in both places, claimList and existList.
	for _, i := range claimList {
		shouldDeleteFileMap[i] = false
	}

	// if fs file exists, set request file to false
	for _, i := range existList {
		shouldMoveFileInMap[i] = false
	}

	// Move new file into target dir
	for fname, shouldMove := range shouldMoveFileInMap {
		if shouldMove {
			fmt.Printf("moving file: %s\n", fname)
			_, err := MoveFile(fname, destDir)
			if err != nil {
				logrus.Printf("error moving file %s err: %v\n", fname, err)
			} else {
				logrus.Printf("success moved %s\n", fname)
			}
		}
	}

	// Delete Exlcluded Files
	for fname, shouldDelete := range shouldDeleteFileMap {
		if shouldDelete {
			fmt.Printf("deleting file: %s\n", fname)
			err := RemoveFile(fname, destDir)
			if err != nil {
				logrus.Printf("error deleting file %s err: %v\n", fname, err)
			} else {
				logrus.Printf("success deleted %s\n", fname)
			}
		}
	}

	// Get to dir and get all the files name to stringify it for return
	var ret []string
	ent, _ := os.ReadDir(filepath.Join(fileRootDir, destDir))
	for _, n := range ent {
		if !n.IsDir() {
			ret = append(ret, fmt.Sprintf("%s/%s", destDir, n.Name()))
		}
	}

	return ret, nil
}

// Validate File Extension
func ValidFileStr(str string, t FileType) string {
	if str != "" {

		switch t {
		case ImageType:
			imgExts := []string{"png", "jpeg", "jpg", "ico"}
			validPattern := regexp.MustCompile(fmt.Sprintf(`\.(%s)$`, strings.Join(imgExts, "|")))
			files := strings.Split(str, "|")
			var validList []string
			for _, f := range files {
				if validPattern.MatchString(f) {
					validList = append(validList, f)
				}
			}
			return strings.Join(validList, "|")

		case VideoType:
			vidExts := []string{"mp4"}
			validPattern := regexp.MustCompile(fmt.Sprintf(`\.(%s)$`, strings.Join(vidExts, "|")))
			files := strings.Split(str, "|")
			var validList []string
			for _, f := range files {
				if validPattern.MatchString(f) {
					validList = append(validList, f)
				}
			}
			return strings.Join(validList, "|")

		case DocType:
			return str

		default:
			return str
		}
	}
	return ""
}

func ValidOneFile(s string) string {
	parts := strings.Split(s, "|")
	if len(parts) > 1 {
		// Use the first one.
		return parts[0]
	}
	return s
}

func ExtractFileExt(s string) string {
	parts := strings.Split(filepath.Base(s), ".")
	if len(parts) > 1 {
		return "." + strings.Join(parts[1:], ".")
	}
	return ""
}

func SaveFile(f *multipart.FileHeader, dstDir string) (string, error) {

	src, err := f.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileName := RandStr() + ExtractFileExt(f.Filename)

	dstPath := filepath.Join(dstDir, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// copy the content from source to destination
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// 💡Validate list of file name match file type
func ValidFileExist(ls []string) []string {
	ret := []string{}
	for _, f := range ls {
		_, err := os.Stat(filepath.Join(fileRootDir, f))
		if err == nil {
			ret = append(ret, f)
		}
	}
	return ret
}

func ValidFileName(lsf []string, t FileType) []string {
	if len(lsf) != 0 {
		switch t {
		case ImageType:
			imgExts := []string{"png", "jpeg", "jpg", "ico"}
			pattern := regexp.MustCompile(fmt.Sprintf(`\.(%s)$`, strings.Join(imgExts, "|")))
			validLs := []string{}
			for _, f := range lsf {
				if pattern.MatchString(f) {
					validLs = append(validLs, f)
				}
			}
			return validLs

		case VideoType:
			vidExts := []string{"mp4"}
			pattern := regexp.MustCompile(fmt.Sprintf(`\.(%s)$`, strings.Join(vidExts, "|")))
			validLs := []string{}
			for _, f := range lsf {
				if pattern.MatchString(f) {
					validLs = append(validLs, f)
				}
			}
			return validLs

		case DocType:
			vidExts := []string{"pdf", "csv", "html", "json"}
			pattern := regexp.MustCompile(fmt.Sprintf(`\.(%s)$`, strings.Join(vidExts, "|")))
			validLs := []string{}
			for _, f := range lsf {
				if pattern.MatchString(f) {
					validLs = append(validLs, f)
				}
			}
			return validLs

		default:
			return lsf
		}
	}
	return lsf
}

func FilterExistingFilesV1(files []string, t FileType) []string {
	ret := []string{}
	for _, file := range files {
		_, err := os.Stat(filepath.Join(fileRootDir, file))
		if !os.IsNotExist(err) && ValidOneFileName(file, t) {
			ret = append(ret, file)
		}
	}
	return ret
}

func FilterExistingFiles(arrs [][]string, t FileType) [][]string {
	ret := [][]string{}

	subRetCount := 0

	for _, arr := range arrs {
		var subRet []string
		for _, file := range arr {
			_, err := os.Stat(filepath.Join(fileRootDir, file))
			if !os.IsNotExist(err) && ValidOneFileName(file, t) {
				subRet = append(subRet, file)
			}
		}

		// check if sub is empty
		if len(subRet) == 0 {
			subRetCount++
		}
		ret = append(ret, subRet)
	}
	if len(ret) == subRetCount {
		return [][]string{}
	}
	return ret
}

func ValidOneFileName(file string, t FileType) bool {
	switch t {
	case ImageType:
		imgExts := []string{"png", "jpeg", "jpg", "svg", "ico", "webp", "gif"}
		pattern := regexp.MustCompile(fmt.Sprintf(`\.(%s)$`, strings.Join(imgExts, "|")))
		if pattern.MatchString(file) {
			return true
		}
		return false
	case VideoType:
		vidExts := []string{"mp4", "mkv", "avi", "mov", "webm", "m4v"}
		pattern := regexp.MustCompile(fmt.Sprintf(`\.(%s)$`, strings.Join(vidExts, "|")))
		if pattern.MatchString(file) {
			return true
		}
		return false
	case DocType:
		vidExts := []string{"pdf", "csv", "html", "json"}
		pattern := regexp.MustCompile(fmt.Sprintf(`\.(%s)$`, strings.Join(vidExts, "|")))
		if pattern.MatchString(file) {
			return true
		}
		return false

	default:
		return false
	}
}
