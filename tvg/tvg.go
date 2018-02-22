package tvg

import (
    "net/url"
    "regexp"
    "strings"
    "bytes"
)

type URL struct{
    url.URL
    Raw     string
}
func (u *URL) String() string{
    return u.Raw
}
func (u *URL) Set(s string) (err error){
    parsed, err := url.Parse(s)
    u.URL = *parsed
    u.Raw = s
    return
}

type M3UData struct{
    List []*EXTINF
}

type EXTINF struct{
    Id          string  `extinf:"tvg-id"`
    Name        string  `extinf:"tvg-name"`
    Logo        string  `extinf:"tvg-logo"`
    Group       string  `extinf:"group-title"`
    Number      int     `extinf:"tvg-chno"`
    Title       string
    Url         string
    SD          bool
    HD          bool
    FHD         bool
    Prefix      string
    NewName     string
}

func Parse(b []byte) (*M3UData, error){

    list := make([]*EXTINF, 0)

    obj := &EXTINF{}

    r := regexp.MustCompile(`([a-z-]+)=+\"([^\"]+)\"`)
    for _, line := range bytes.Split(b, []byte{10}){

        line = bytes.Replace(line, []byte{0x0d}, []byte(""), -1) // \r

        if len(line) < 8{
            continue
        }

        s := string(line)
        if strings.Compare(s[:8], "#EXTINF:") != 0{
            obj.Url = strings.Replace(s, "\r", "", 0)
            list = append(list, obj)

            obj = &EXTINF{}
            continue
        }

        titles := strings.Split(s, ",")
        title := titles[len(titles)-1]
        title = strings.Trim(title, " ")
        title = strings.Replace(title, " :", ":", -1)
        obj.Title = title

        tags := r.FindAllStringSubmatch(s[8:], -1)

        for _, tag := range tags{

            key := tag[1]
            value := strings.Trim(tag[2], " ")
            value = strings.Replace(tag[2], " :", ":", -1)

            if key == "tvg-id"{
                obj.Id = value
            }else if key == "tvg-name"{
                obj.Name = value
            }else if key == "tvg-logo"{
                obj.Logo = value
            }else if key == "group-title"{
                obj.Group = value
            }

            if strings.Contains(obj.Name, "FHD"){
                obj.FHD = true
            }else if strings.Contains(obj.Name, "HD"){
                obj.HD = true
            }else{
                obj.SD = true
            }
        }

        obj.Prefix, obj.NewName = cleanName(obj.Name)
    }

    return &M3UData{list}, nil
}

func cleanName(name string) (prefix, newname string){
    if strings.Contains(name, ":"){
        elems := strings.Split(name, ":")
        prefix = elems[0]
        newname = elems[1]
        //noprefix = strings.Join(strings.Split(name, ":")[1:], " ")
    }else{
        newname = name
    }
    newname = strings.Replace(newname, "FHD", "", -1)
    newname = strings.Replace(newname, "HD", "", -1)
    //fmt.Println("name", name, "prefix", newname, "newname", newname)
    return strings.Trim(prefix, " "), strings.Trim(newname, " ")
}


