package main

import (
    "fmt"
    "os"
    "path"
    "path/filepath"
    "pkg.linuxdeepin.com/lib/dbus"
    "strings"
    "./bundle"
)

var conn *dbus.Conn = nil

func getBus() *dbus.Conn {
    if conn == nil {
        var err error
        conn, err = dbus.SystemBus()
        if err != nil {
            panic(err)
        }
    }
    return conn
}

func parseArgsPrefix(args []string, prefix string) (result string) {
    for _, v := range args {
        if strings.HasPrefix(v, prefix) {
            result = strings.TrimPrefix(v,prefix)
            return result
        }
    }
    return ""
}

func removeArgsPrefix(args []string, prefix string) []string {
    s := -1
    for k, v := range args {
        if strings.HasPrefix(v, prefix) {
            s = k
            break
        }
    }
    if s == -1 {
        return args
    }else{
        copy(args[s:], args[s+1:])
        args = args[:len(args)-1]
        return args
    }

}

func main() {
    if length := len(os.Args); length == 1 || length > 3  {
        fmt.Println("Usage: ...")
        return
    }
    bundle.CreateWorkingDirectories()
    var bundleFile string
    switch os.Args[1] {
        case "-m":
            bundleFile = path.Base(os.Args[2])
            fmt.Println("Mount: ...", bundleFile)
            return
        case "-b":
            bundleFile = path.Base(os.Args[2])
            fmt.Println("Binary: ...", bundleFile)
            return
        case "-i":
            bundleFile = path.Base(os.Args[2])
            bundleInfo, err := bundle.GetBundleInfoBase(bundleFile)
            fmt.Println(bundleInfo)
            if err != nil {
                fmt.Println("Bundle Binary: ...")
                return
            }
            bundleInfo.Inactive()
            return

        default:
            bundleFile = os.Args[1]
            fmt.Println("Launcher: ...", bundleFile)
    }
    exec := parseArgsPrefix(os.Args, "-app=")
    fpath,_ := filepath.Abs(bundleFile)
    bundlePath := bundle.CheckLocation(fpath)
    //bundleInfoBasename := path.Base(bundlePath)
    bundleInfo, err := bundle.GetBundleInfoBase(path.Base(bundlePath))
    if err != nil {
        panic(err)
    }
    conn = getBus()
    //bundleInfoCompletename := bundlePath
    //fmt.Println(bundleInfoBasename, bundleInfoCompletename)
    fmt.Println(bundleInfo)

    fmt.Println("[INFO] prepare runtime environment")
    _, ifChroot := bundleInfo.PrepareEnvironment()
    fmt.Println("[INFO] prepare process environment")
    bundleInfo.ProcessEnvironment()
    ///unmount("chrome")
    var arguments []string
    arguments = os.Args[2:]
    arguments = removeArgsPrefix(arguments,"-app=")
    fmt.Println("[INFO] Run bundle, enjoy it!")
    bundleInfo.Run(ifChroot, exec, arguments)

}
