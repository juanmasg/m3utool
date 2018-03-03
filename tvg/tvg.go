package tvg

import (
    "net/url"
    "regexp"
    "strings"
    "bytes"
    "strconv"
    "fmt"
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

func (m3u *M3UData) AsMapByNumber() map[int]*EXTINF{
    m := make(map[int]*EXTINF)

    for _, inf := range m3u.List{
        m[inf.Number] = inf
    }

    return m
}

func (m3u *M3UData) AsMapByName() map[string]*EXTINF{
    m := make(map[string]*EXTINF)

    for _, inf := range m3u.List{
        m[inf.NewName] = inf
    }

    return m
}

func (m3u *M3UData) Print(){
    for _, inf := range m3u.List{
        if inf == nil{
            continue
        }
        suffix, prefix := "", ""
        if inf.HD{ suffix = " HD"}
        if inf.FHD{ suffix = " FHD"}
        if inf.Prefix != ""{ prefix = inf.Prefix + ": " }

        name := prefix + inf.NewName + suffix

        //fmt.Printf("#EXTINF:-1 tvg-chno=\"%d\" tvg-id=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\", %s \n%s\n",
        //    inf.Number, inf.Id, name, inf.Logo, inf.Group, name, inf.Url)
        fmt.Printf("#EXTINF:-1 tvg-chno=\"%d\" tvg-id=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\", %s \n%s\n",
            inf.Number, inf.Id, name, "", inf.Group, name, inf.Url)
    }
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
            }else if key == "tvg-chno"{
                obj.Number, _ = strconv.Atoi(value)
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

    replacer := strings.NewReplacer("FHD", "", "HD", "")
    newname = name

    if strings.Contains(name, ":"){
        elems := strings.Split(name, ":")
        prefix = elems[0]
        newname = elems[1]
    }

    newname = replacer.Replace(newname)
    newname = strings.Title(strings.ToLower(newname))
    return strings.Trim(prefix, " "), strings.Trim(newname, " ")
}


