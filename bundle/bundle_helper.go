package bundle

import (
    "fmt"
    "os"
    "regexp"
    "bufio"
    "strings"
    "io"
    "errors"
    "io/ioutil"
)

type BundleInfoPrivate struct {
    Path string
    options map[string]string
}

type BundleInfoBase struct {
    name string
    version string
    reversion string
    architecture string
}


func GetBundleInfoBase(path string)(*BundleInfoBase, error){
    re := regexp.MustCompile("(.+)-([.\\d\\w]+)-(\\d+)-(x86_64|i686)")
    result := re.FindStringSubmatch(path)
    fmt.Println(result)
    if result == nil {
        return nil, errors.New("Cann't regexp string")
    }
    info := new(BundleInfoBase)
    info.name = result[1]
    info.version = result[2]
    info.reversion = result[3]
    info.architecture = result[4]
    return info, nil
}

// mountPoint should contain filename PKGINFO
//func GetBundleInfoPrivate(mountPoint string)(*BundleInfoPrivate, error){
func GetBundleInfoPrivate(mountPoint string)(*BundleInfoPrivate, error){
    fp, err := os.Open(mountPoint+"/PKGINFO")
    if err != nil {
        return nil, errors.New("PKGINFO cannot found or cannot Open")
    }
    result := make(map[string]string)
    buf := bufio.NewReader(fp)
    for {
        line, err := buf.ReadString('\n')
        line = strings.TrimSpace(line)
        if err != nil {
            if err != io.EOF {
                panic(err)
            }
            if len(line) == 0 {
                break
            }
        }
        switch {
            //detect blank line or code comments
            case len(line) == 0:
            case line[0] == '#' :
            default:
                i := strings.IndexAny(line,"=")
                result[strings.TrimSpace(line[0:i])] = strings.TrimSpace(line[i+1:])
            }
    }
    info := new(BundleInfoPrivate)
    info.Path = mountPoint
    info.options = result
    return info, nil
}

func (b *BundleInfoBase) GetBundlePath() string{
    path := fmt.Sprintf("%s/%s-%s-%s-x86_64.cb", GetRepositoryDirectory(), b.name, b.version, b.reversion)
    return path
}
func (b *BundleInfoPrivate) GetBundlePath() string{
    path := fmt.Sprintf("%s/%s-%s-%s-x86_64.cb", GetRepositoryDirectory(), b.options["name"], b.options["version"], b.options["release"])
    return path
}

func (b *BundleInfoBase) Inactive(){
    mountPoint := b.GetMountPoint()
    appPath := b.GetApplicationPath()
    unionPath := b.GetUnionPath()
    if _, err := os.Stat(unionPath); err == nil {
        Unmount(unionPath+"/home")
        Unmount(unionPath+"/proc")
        Unmount(unionPath+"/var")
        Unmount(unionPath+"/sys")
        Unmount(unionPath+"/dev/pts")
        Unmount(unionPath+"/dev/shm")
        Unmount(unionPath+"/dev")
        Unmount(unionPath+"/run")
        Unmount(unionPath+"/tmp")
        Unmount(unionPath)
    }
    if _, err := os.Stat(appPath); err == nil {
        Unmount(appPath)
    }
    Unmount(mountPoint)

}

func (b *BundleInfoBase) Active() {
    b.PrepareEnvironment()
}

func (b *BundleInfoBase) PrepareEnvironment() (workingPath string, ifChroot bool) {
    mountPoint := b.GetMountPoint()
    appPath := b.GetApplicationPath()
    absPath := b.GetBundlePath()
    Mount(absPath, mountPoint, "iso9660")
    bundleInfoPrivate, err := GetBundleInfoPrivate(mountPoint)
    if err != nil {
        panic(err)
    }
    Mount(mountPoint+"/app.sqfs", appPath, "squashfs")
    chroot, ok := bundleInfoPrivate.options["chroot"]
    if !ok {
        chroot = "true"
    }
    if chroot == "true" {
        unionPath := b.GetUnionPath()
        unionPaths := make(map[string]string,2)
        unionPaths[appPath]="ro"
        unionPaths["/"]="ro"
        Unionfs(unionPaths, unionPath)
        return unionPath, true
    }
    return appPath, false
}

// fuck my code
func (b *BundleInfoBase) Run(ifChroot bool, exec string, arguments []string) {
    mountPoint := b.GetMountPoint()
    unionPath := b.GetUnionPath()
    bundleInfoPrivate, _ := GetBundleInfoPrivate(mountPoint)
    fmt.Println("[Debug]", ifChroot, arguments, bundleInfoPrivate)
    if exec == "" {
        exec = bundleInfoPrivate.options["exec"]
    }
    if ifChroot {
        StartChrootedProcess(exec, unionPath, arguments)
    } else {
        StartNormalProcess(exec, arguments)
    }

}

func (b *BundleInfoBase) ProcessEnvironment() {
    envPrefix := "CBUNDLE_"
    //mountPoint := b.GetMountPoint()
    //bundleInfoPrivate, _ := GetBundleInfoPrivate(mountPoint)
    //chroot, ok := bundleInfoPrivate["chroot"]
    //if !ok {
    //    chroot = "true"
    //}
    //if chroot == "true" {
    //    mP := mountPoint
    //} else {
    //    mP := b.GetApplicationPath()
    //}
    //
    //UnsetEnv("LD_LIBRARAY_PATH")
    //os.Setenv(envPrefix+"LD_LIBRARAY_PATH", fmt.Sprintf("%s/usr/lib/:%s/lib",mP, mP))

    //oldPath := os.Getenv("PATH")
    //UnsetEnv("PATH")
    //newPath := "/bin:/usr/bin:/sbin:/usr/sbin:/opt/java/jre/bin:" + "/usr/bin/perlbin/site:/usr/bin/perlbin/vendor:/usr/bin/perlbin/core:" + "/opt/lib32/bin:/opt/lib32/usr/bin:/usr/local/bin:"
    //os.Setenv(envPrefix+"PATH", newPath + oldPath)

    oldDbusAddress := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
    UnsetEnv("DBUS_SESSION_BUS_ADDRESS")
    os.Setenv(envPrefix+"DBUS_SESSION_BUS_ADDRESS",oldDbusAddress)
    // enhance by PKGINFO --- TODO
}

//bundleFile is the abs path 
func CheckLocation(bundlePath string)(fixedBundlePath string){
    r, _ := regexp.MatchString(`.*/\.cinstall/repo/.*\.cb`, bundlePath)
    if r {
        //check if the bundleFile is valided
        return bundlePath
    } else {
        //TODO check file exists 
        tempdirPrefix := fmt.Sprintf("/tmp/cinstall-%d/tempdir", os.Getuid())
        if _, err := os.Stat(tempdirPrefix); os.IsNotExist(err) {
            os.MkdirAll(tempdirPrefix, 0755)
        }
        tempDir, _ := ioutil.TempDir(tempdirPrefix,"app")
        Mount(bundlePath, tempDir, "iso9660")
        result,_ := GetBundleInfoPrivate(tempDir)
        //need add architure to PKGINFO
        bundleName := fmt.Sprintf("%s-%s-%s-x86_64.cb", result.options["name"], result.options["version"], result.options["release"])
        fixedBundlePath := fmt.Sprintf("%s/%s", GetRepositoryDirectory(), bundleName)
        fmt.Println("Debug", fixedBundlePath)
        if _, err := os.Stat(fixedBundlePath); err == nil {
            //TODO: should repace it , now just skip
            fmt.Println("Same bundle exists, please fix it soon~")
            bundleInfo, _ := GetBundleInfoBase(bundleName)
            fmt.Println("Same bundle exists, Now Inactive it~")
            bundleInfo.Inactive()
        } else {
            fmt.Println("Copy bundle:", bundlePath, " --->  ",fixedBundlePath)
            if _, err := CopyFile(bundlePath, fixedBundlePath); err != nil {
                panic(err)
            }
            fmt.Println("[w] CreateWorkFile")
            if err := result.CreateWorkFile(); err != nil {
                panic(err)
            }
        }
        fmt.Println("[w] Before unmount", tempDir)
        Unmount(tempDir)
        fmt.Println("[w] After unmount", tempDir)
        defer os.Remove(tempDir)
        return fixedBundlePath
    }

}

func (b *BundleInfoBase) GetMountPoint() string {
    uid := os.Getuid()
    return fmt.Sprintf("/tmp/%s-%d/%s", clientName, uid, b.name+"-"+b.version+"-"+b.reversion)
}

func (b *BundleInfoBase) GetApplicationPath() string {
    return fmt.Sprintf("%s-app", b.GetMountPoint())
}

func (b *BundleInfoBase) GetUnionPath() string {
    return fmt.Sprintf("%s-union", b.GetMountPoint())
}


