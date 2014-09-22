package main

// #cgo pkg-config: libkmod 
// #include "kmod.c"
import "C"
import "unsafe"
import "fmt"
import "pkg.linuxdeepin.com/lib/dbus"

type moduleHelper struct {
}

func (f *moduleHelper) GetDBusInfo() dbus.DBusInfo {
    return dbus.DBusInfo{
        "com.linuxdeepin.bundle.modulehelper",
        "/",
        "com.linuxdeepin.bundle.modulehelper",
    }
}

func (f *moduleHelper) Load(name string) int32{
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    return_code := int32(C.load(cname))
    fmt.Println("return code:",return_code)
    return return_code

    //if return_code == 0 {
     //   return "Insert module Successful"
    //}else{
     //   return "Insert module Failed"
    //}
}

func main() {
    m := &moduleHelper{}
    if err := dbus.InstallOnSystem(m); err != nil {
        panic(err)
    }
    dbus.DealWithUnhandledMessage()
    dbus.Wait()
}
