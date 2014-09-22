package bundle

import (
    "github.com/alyu/configparser"
    "errors"
    "strings"
    "log"
    "fmt"
    "os"
)


//Get value from PKGINFO
// --- Exec=/usr/bin/gimp;/usr/bin/google-chrome -> ["/usr/bin/gimp" "/usr/bin/google-chrome"]
func (p *BundleInfoPrivate) GetValue(key string)([]string,error){
    value, ok := p.options[key]
    if !ok {
        return nil, errors.New("No section found")
    }
    result := strings.Split(value, ";")
    return result, nil
}

func (p *BundleInfoPrivate) CreateWorkFile() error {
    //check if multi program
    //TODO: please reWrite it first
    var multi bool
    desktops, _ := p.GetValue("desktop")
    execs, _ := p.GetValue("exec")
    icons, _ := p.GetValue("icon")
    if len(execs) != len(desktops) || len(icons) != len(desktops) {
        panic(errors.New("PKGINFO file error"))
    }
    if len(desktops) > 1 {
        multi = true
    } else {
        multi  = false
    }
    pkgname := p.options["name"]
    version := fmt.Sprintf("%s-%s", p.options["version"],p.options["release"])
    for n, desktop := range desktops {
        exec := execs[n]
        icon := icons[n]
        icon_prefix, icon_suffix, _:= ParseLastDot(icon)
        desktop_prefix, _,_ := ParseLastDot(desktop)
        if desktop_prefix == "" {
            continue
        }
        fmt.Println(desktop_prefix, "-----------------------")
        var process,iconname string
        if multi {
            process = fmt.Sprintf("%s-%s", pkgname, strings.TrimSuffix(desktop, ".desktop"))
            iconname = fmt.Sprintf("%s-%s", pkgname, icon_prefix)
            desktop_prefix = fmt.Sprintf("%s-%s", pkgname, desktop_prefix)
        } else {
            process = pkgname
            iconname = fmt.Sprintf("%s", icon_prefix)
        }
        fixExec := fmt.Sprintf("%s/%s-%s",GetLaunchersDirectory(), process, version)
        fixIcon := fmt.Sprintf("%s/%s-%s.%s",GetIconsDirectory(), iconname, version, icon_suffix)
        fixdesktopPath := fmt.Sprintf("%s/%s-%s.%s", GetCacheDirectory(), desktop_prefix, version, "desktop")
        desktopAbsPath := fmt.Sprintf("%s/%s", p.Path, desktop)
        info := make(map[string]string)
        info["Exec"] = fixExec
        info["Icon"] = fixIcon

        err := createDesktop(desktopAbsPath, info, version, fixdesktopPath)
        if err != nil {
            return errors.New("Write desktopfile failed")
        }
        os.Chmod(fixdesktopPath, 0755)
        p.writeScript(exec, fixExec)
        iconAbsPath := fmt.Sprintf("%s/icons/%s", p.Path, icon)
        fmt.Println(iconAbsPath , "--->", fixIcon)
        CopyFile(iconAbsPath, fixIcon)

    }
    return nil

}

func ParseLastDot(source string) (prefix string, suffix string, err error) {
    p := strings.Split(source, ".")
    suffix = p[len(p)-1]
    if suffix == "" {
        return "","",errors.New("suffix should not be nil")
    }
    prefix = strings.TrimSuffix(source, "."+suffix)
    return prefix,suffix, nil
}

func createDesktop(desktop string, info map[string]string, version string, dest string) error {
    content, err := configparser.Read(desktop)
    if err != nil {
        log.Fatal(err)
        return  err
    }
    section, _ := content.Section("Desktop Entry")
    options := section.Options()
    for k, v := range info {
        if strings.Title(k) == "Exec" {
            exec := options["Exec"]
            if strings.Contains(exec, "%"){
                p := strings.Split(exec, "%")
                suffix := p[len(p)-1]
                v += (" "+"%"+suffix)
            }
        }
        section.Add(strings.Title(k),v)
    }
    new_content := configparser.NewConfiguration()
    activeSection := new_content.NewSection("Desktop Entry")

    for k, v := range options {
        if strings.HasPrefix(k,"Name") {
            activeSection.Add(k, v+" "+version)
        }else if k == "TryExec"{
            log.Print("Strip TryExec option")
        }else {
            activeSection.Add(k, v)
        }
    }
    err = configparser.Save(new_content, dest)
    return err
}

//func (p *BundleInfoPrivate) getFileName(exec string) string {
//    var f string
//    if exec != ""{
//        f = fmt.Sprintf("%s/%s-%s-%s", GetLaunchersDirectory(), p.options["name"], p.options["version"], p.options["release"])
//    }else{
//        f = fmt.Sprintf("%s/%s-%s-%s-%s", GetLaunchersDirectory(), p.options["name"], exec, p.options["version"], p.options["release"])
//    }
//    return f
//}

func (p *BundleInfoPrivate) writeScript(exec string, scriptPath string) error {
    f, err := os.OpenFile(scriptPath, os.O_WRONLY |os.O_TRUNC| os.O_CREATE, 0x755)
    if err != nil{
        return err
    }
    defer f.Close()
    f.WriteString("#!/bin/bash\n")
    var script string
    if exec == "" {
        script = "bundler-helper"+" "+p.GetBundlePath()+" \"${@}\""+"\n"
    }else{
        script = "bundler-helper"+" "+p.GetBundlePath()+" "+"-app="+exec+" \"${@}\""+"\n"
    }
    fmt.Println(script)
    f.WriteString(script)
    return nil
}
