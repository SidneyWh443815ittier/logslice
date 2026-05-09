//go:build !windows

package rotate

import (
	"os"
	"syscall"
)

// inode returns the inode number of the file described by info.
// On Unix systems this is taken from the underlying Sys() stat.
func inode(info os.FileInfo) uint64 {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
}
