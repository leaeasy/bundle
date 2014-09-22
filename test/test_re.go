package main

import (
    "fmt"
    "./bundle"
)

func main() {
    bundleString := "chrome-dev-27.0.1438.7-2-x86_64.cb"
    v,e := bundle.GetBundleInfoBase(bundleString)
    fmt.Println(v)
    fmt.Println(e)

    //r, e := bundle.GetBundleInfoPrivate("/tmp/b")
    //if e != nil {
    //    fmt.Println(e)
    //    return
    //}
    //for key, value := range r {
    //    fmt.Println(key, "------->", value)
    //}
    fmt.Println(bundle.GetMountPoint("tt", "firefox"))
}
