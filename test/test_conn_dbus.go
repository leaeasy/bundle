package main

import (
    "fmt"
    "pkg.linuxdeepin.com/lib/dbus"
)


func main() {
    conn, err := dbus.SystemBus()
    if err != nil {
        panic(err)
    }
    obj := conn.Object("com.linuxdeepin.bundle.filesystem", "/")
    var result string
    err2 := obj.Call("com.linuxdeepin.bundle.filesystem.Rmount", 0, "/dev/sda1", "/mnt", "ext2", int32(0)).Store(&result) 
    fmt.Println(err2)
    fmt.Println(result)

}
