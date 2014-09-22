package main
/*
#include <stdio.h>
#include <stdlib.h>
#include <mntent.h>
*/
import "C"
import (
	"errors"
	"unsafe"
    "strings"
    "sync"
)


var mountInfoLock sync.RWMutex

// Represents a mounted volume on the host system
type Volume struct {
	Fspath string // The mount point of the volume
    Dest string // The mount point of dest
	Type string // The filesystem type
    Opts string // The filesyste opts
    Loop bool // The loop opts
}

// Gets a slice of all volumes that are currently
// mounted on the host system.
func GetMountedVolumes() ([]Volume, error) {
    mountInfoLock.Lock()
    defer mountInfoLock.Unlock()
	result := make([]Volume, 0)

	cpath := C.CString("/proc/mounts")
	defer C.free(unsafe.Pointer(cpath))
	cmode := C.CString("r")
	defer C.free(unsafe.Pointer(cmode))
	var file *C.FILE = C.setmntent(cpath, cmode)
	if file == nil {
		return nil, errors.New("Unable to open /proc/mounts")
	}
	defer C.endmntent(file)
	var ent *C.struct_mntent
    var loop bool

	for ent = C.getmntent(file); ent != nil; ent = C.getmntent(file) {
        fspath := C.GoString(ent.mnt_fsname)
        if strings.HasPrefix(fspath, "/dev/loop") {
            loop = true
        } else {
            loop = false
        }
		result = append(result, Volume{fspath, C.GoString(ent.mnt_dir), C.GoString(ent.mnt_type), C.GoString(ent.mnt_opts), loop})
	}

	return result, nil
}

// Get path mountinfo
func MountInfoForPath(path string) (*Volume, error){
    volumes, err := GetMountedVolumes()
    if err != nil {
        return nil, err
    }
    for _, volume := range volumes {
        if path == volume.Dest {
            return &volume, nil
        }
        if path == volume.Fspath {
            return &volume, nil
        }
    }
    return nil, errors.New("Fspath is not mounted")
}
