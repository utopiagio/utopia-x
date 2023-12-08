// utils_unix.go

// +build !windows

package filedialog

const dotCharacter = 46

func isHidden(filepath string, filename string) (bool, error) {
	if filename[0] == dotCharacter {
		return true, nil
	}

	return false, nil
}


/*func modeToType(mode os.FileMode) byte {
	switch {
	case mode.IsDir():
		return 'd'
	case mode.IsRegular():
		return 'f'
	case mode&os.ModeSymlink != 0:
		return 'l'
	case mode&os.ModeSocket != 0:
		return 's'
	case mode&os.ModeNamedPipe != 0:
		return 'p'
	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0:
		return 'C'
	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice == 0:
		return 'D'
	default:
		return 'u'
	}
}*/
