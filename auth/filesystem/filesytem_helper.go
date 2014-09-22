package main

import (
    //"./filehelper"
    "pkg.linuxdeepin.com/lib/dbus"
)


type mountHelper struct {
}

func (f *mountHelper) GetDBusInfo() dbus.DBusInfo {
    return dbus.DBusInfo{
        "com.linuxdeepin.bundle.filesystem",
        "/",
        "com.linuxdeepin.bundle.filesystem",
    }
}

func (f *mountHelper) Mount(fspath string, dest string, fstype string, options string)( volume *Volume, err error){
    return Mount(fspath, dest, fstype, options)
}

func (f *mountHelper) Rmount(fspath string, dest string, fstype string) string {
    result := fspath + dest + fstype
    return result
}


func (f *mountHelper) Unmount(fspath string, options []string)( err error){
    return Unmount(fspath, options)
}

func (f *mountHelper) Bind(fspath string, dest string)( err error){
    return Bind(fspath, dest)
}

func (f *mountHelper) Unionfs(unionPaths map[string]string, mountPoint string)( volume *Volume, err error){
    return Unionfs(unionPaths, mountPoint)
}

func main() {
    f := &mountHelper{}
    if err:= dbus.InstallOnSystem(f); err != nil {
        panic(err)
    }

    dbus.DealWithUnhandledMessage()

    dbus.Wait()
}
