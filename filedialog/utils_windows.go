// utils_windows.go

// +build windows

package filedialog

import(
	"log"
	filepath_go"path/filepath"
	"syscall"
)

const dotCharacter = 46
// isHidden checks if a file is hidden on Windows.
func isHidden(filepath string, filename string) (bool, error) {
	// dotfiles also count as hidden (if you want)
	if filename[0] == dotCharacter {
		return true, nil
	}

	absPath, err := filepath_go.Abs(filepath + filename)
	if err != nil {
		log.Println(err)
		return false, err
	}
	//log.Println("AbsPathb :", absPath)
	// Appending `\\?\` to the absolute path helps with
	// preventing 'Path Not Specified Error' when accessing
	// long paths and filenames
	// https://docs.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation?tabs=cmd
	pointer, err := syscall.UTF16PtrFromString(`\\?\` + absPath)
	if err != nil {
		return false, err
	}

	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false, err
	}

	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}