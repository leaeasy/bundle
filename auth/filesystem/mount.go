package main

import (
    "sync"
    "os/exec"
)

//const (
//  MS_RDONLY = 1		/* Mount read-only.  */
//  MS_NOSUID = 2		/* Ignore suid and sgid bits.  */
//  MS_NODEV = 4		/* Disallow access to device special files.  */
//  MS_NOEXEC = 8		/* Disallow program execution.  */
//  MS_SYNCHRONOUS = 16		/* Writes are synced at once.  */
//  MS_REMOUNT = 32		/* Alter flags of a mounted FS.  */
//  MS_MANDLOCK = 64		/* Allow mandatory locks on an FS.  */
//  MS_DIRSYNC = 128		/* Directory modifications are synchronous.  */
//  MS_NOATIME = 1024		/* Do not update access times.  */
//  MS_NODIRATIME = 2048		/* Do not update directory access times.  */
//  MS_BIND = 4096		/* Bind directory at different place.  */
//  MS_MOVE = 8192
//  MS_REC = 16384
//  MS_SILENT = 32768
//  MS_POSIXACL = 1 << 16	/* VFS does not apply the umask.  */
//  MS_UNBINDABLE = 1 << 17	/* Change to unbindable.  */
//  MS_PRIVATE = 1 << 18		/* Change to private.  */
//  MS_SLAVE = 1 << 19		/* Change to slave.  */
//  MS_SHARED = 1 << 20		/* Change to shared.  */
//  MS_RELATIME = 1 << 21	/* Update atime relative to mtime/ctime.  */
//  MS_KERNMOUNT = 1 << 22	/* This is a kern_mount call.  */
//  MS_I_VERSION =  1 << 23	/* Update inode I_version field.  */
//  MS_STRICTATIME = 1 << 24	/* Always perform atime updates.  */
//  MS_ACTIVE = 1 << 30
//  MS_NOUSER = 1 << 31
//)

var mountLock sync.RWMutex
func Mount(fspath string, dest string, fstype string, options string)(*Volume, error){

    if isMounted(fspath) || isMounted(dest) {
        return MountInfoForPath(dest)
    }
    mountLock.Lock()
    defer mountLock.Unlock()

    args := []string{}
    args = append(args, fspath)
    args = append(args, "-t", fstype)
    if options != "" {
        args = append(args, "-o", options)
    }
    args = append(args, dest)
    command :=exec.Command("/bin/mount", args...)
    err :=  command.Run()
    if err != nil {
        return nil, err
    }
    //if err := syscall.Mount(fspath, dest, fstype, uintptr(flags), data); err != nil {
    //    return nil, err
    //}
    return MountInfoForPath(dest)
}


func isMounted(path string) bool{
    volume, _ := MountInfoForPath(path)
    if volume != nil {
        return true
    }
    return false
}

func Unmount(dest string, options []string) error {
    if !isMounted(dest){
        return nil
    }
    args := []string{}
    if len(options) != 0 {
        args = append(args, options...)
    }
    args = append(args, dest)
    command :=exec.Command("/bin/umount", args...)
    err :=  command.Run()
    return err
    //err := syscall.Unmount(fspath, flags)
}

func Bind(path string, dest string) error {
    if !isMounted(dest){
        return nil
    }
    command :=exec.Command("/bin/mount","--bind",path, dest)
    err := command.Run()
    return err
}

//map is [string]"ro" or [string]"rw" , may add options later
func Unionfs(unionPaths map[string]string,  mountPoint string)(*Volume, error) {
    if isMounted(mountPoint) {
        return MountInfoForPath(mountPoint)
    }
    mountLock.Lock()
    defer mountLock.Unlock()
    var deviceString string
    deviceString = "br="
    length := len(unionPaths)
    for path, info := range unionPaths {
        if length--;length != 0 {
            deviceString += (path +"="+ info+":")
        }else{
            deviceString += (path +"="+ info)
        }
    }

    args := []string{ "-t","aufs","-o"}
    args = append(args, deviceString , "-o","udba=none","none", mountPoint)
    command :=exec.Command("/bin/mount",  args...)
    err :=  command.Run()
    if err != nil {
        return nil, err
    }
    return MountInfoForPath(mountPoint)
}
