package main

import (
    //"reflect"
    //"log"
    "os"
    "io/ioutil"
    //"errors"
    "./tvg"
    "./ultrabox"
    "fmt"
)

func main(){
    data, _ := ioutil.ReadFile(os.Args[1])
    m3u, err := tvg.Parse(data); if err != nil{
        panic(err)
    }

    m3u = ultrabox.Filter(m3u)

    for _, inf := range m3u.List{
        suffix := ""
        if inf.HD{ suffix = " HD"}
        if inf.FHD{ suffix = " FHD"}

        name := fmt.Sprintf("%s: %s%s", inf.Prefix, inf.NewName, suffix)

        fmt.Printf("#EXTINF:-1 tvg-chno=\"%d\" tvg-id=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\", %s \n%s\n",
            inf.Number, inf.Id, name, inf.Logo, inf.Group, name, inf.Url)
    }
}

