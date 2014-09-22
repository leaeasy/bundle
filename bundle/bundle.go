package bundle

import (
    "os"
    "fmt"
)

const clientName = "cinstall"

func GetClientDirectory() string {
    return fmt.Sprintf("%s/.%s",os.Getenv("HOME"), clientName)
}

func GetRepositoryDirectory() string {
    return fmt.Sprintf("%s/repo", GetClientDirectory())
}

func  GetLaunchersDirectory() string {
    return fmt.Sprintf("%s/launchers", GetClientDirectory())
}

func GetCacheDirectory() string {
    return fmt.Sprintf("%s/cache", GetClientDirectory())
}

func GetIconsDirectory() string {
    return fmt.Sprintf("%s/icons", GetClientDirectory())
}

func CreateWorkingDirectories(){
    var dirs []string
    dirs = append(dirs, GetClientDirectory(), GetRepositoryDirectory(), GetLaunchersDirectory(), GetCacheDirectory(), GetIconsDirectory())
    for _, dir := range dirs {
        if _, err := os.Stat(dir); os.IsNotExist(err) {
            fmt.Printf("No such directory: %s", dir)
            os.MkdirAll(dir, 0755)
        }
    }
}

func Mount(path string, dest string, fstype string) {
    conn := getBus()
    if _, err := os.Stat(dest); os.IsNotExist(err) {
        os.MkdirAll(dest, 0755)
    }
    obj := conn.Object("com.linuxdeepin.bundle.filesystem","/")
    err := obj.Call("com.linuxdeepin.bundle.filesystem.Mount", 0, path, dest, fstype, "")
    fmt.Println(err)
}

func Unmount(path string) {
    conn := getBus()
    obj := conn.Object("com.linuxdeepin.bundle.filesystem","/")
    options := make([]string,0)
    err := obj.Call("com.linuxdeepin.bundle.filesystem.Unmount", 0, path, options)
    fmt.Println(err)
}

func Bind(path string, dest string) {
    conn := getBus()
    if _, err := os.Stat(dest); os.IsNotExist(err) {
        os.MkdirAll(dest, 0755)
    }
    obj := conn.Object("com.linuxdeepin.bundle.filesystem","/")
    err := obj.Call("com.linuxdeepin.bundle.filesystem.Bind", 0, path, dest)
    fmt.Println(err)
}

func Unionfs (unionPaths map[string]string, dest string) {
    conn := getBus()
    if _, err := os.Stat(dest); os.IsNotExist(err) {
        os.MkdirAll(dest, 0755)
    }
    obj := conn.Object("com.linuxdeepin.bundle.filesystem","/")
    obj.Call("com.linuxdeepin.bundle.filesystem.Unionfs", 0, unionPaths, dest)
}
