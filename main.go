package main

import (
    "io/ioutil"
    "./tvg"
    "./provider"
    "./provider/ultrabox"
    "./provider/happytv"
    "os"
    "strconv"
    "strings"
    "fmt"
    "flag"
    "regexp"
    "github.com/juanmasg/xmltvtool/xmltv"
)

type RemapperInt interface{
    Map(m3u *tvg.M3UData, spec string) bool
}



type BasicRemapper struct{}
func (r BasicRemapper) Map(m3u *tvg.M3UData, from, to int) bool{
    var fromid, toid int
    var fromfound, tofound bool
    for i, inf := range m3u.List{
        if inf == nil{ continue }
        if inf.Number == from{
            fromid = i
            fromfound = true
            if tofound{ break }
        }

        if inf.Number == to{
            toid = i
            tofound = true
            if fromfound{ break }
        }
    }

    if fromfound{

        if ! tofound{
            toid = len(m3u.List)
            m3u.List = append(m3u.List, &tvg.EXTINF{})
        }

        m3u.List[toid] = m3u.List[fromid]
        m3u.List[fromid] = nil

        m3u.List[toid].Number = to
        return true
    }

    return false
}

type BasicSwapper struct{}

func (s *BasicSwapper) Map(m3u *tvg.M3UData, from, to int) bool{
    return false
}

func main(){

    flag_ultrabox_filter := flag.String("ultrabox_filter", "", "Filter ultrabox m3u")
    flag_happytv_filter := flag.String("happytv_filter", "", "Filter happytv m3u")
    flag_with_xmltv := flag.String("with_xmltv", "", "Filter also xmltv data")
    flag_ultrabox_movies := flag.String("ultrabox_movies", "", "Extract movies")
    flag_check_epg := flag.String("check_epg", "", "Check m3u<-xmltv")
    flag_gen_epg := flag.Bool("gen_epg", false, "Scan multiple XMLTV files and generate a complete EPG")
    flag_remap  := flag.String("remap", "", "Remap channel numbers (10:0,12-40:10,...)")

    flag.Parse()

    if *flag_remap != ""{

        data, _ := ioutil.ReadFile(flag.Args()[0])
        m3u, err := tvg.Parse(data); if err != nil{
            panic(err)
        }

        maps := strings.Split(*flag_remap, ",")
        for _, m := range maps{
            ft := strings.Split(m, ":")
            from, _ := strconv.Atoi(ft[0])
            to, _ := strconv.Atoi(ft[1])
            //if (BasicRemapper{}).Map(m3u.AsMapByNumber(), from, to){
            if (BasicRemapper{}).Map(m3u, from, to){
                fmt.Println("Map", from, to, "success")
            }else{
                fmt.Println("Map", from, to, "fail!")
            }
        }

        m3u.Print()
    }

    if *flag_ultrabox_filter != ""{
        data, _ := ioutil.ReadFile(*flag_ultrabox_filter)
        m3u, err := tvg.Parse(data); if err != nil{
            panic(err)
        }

        m3u = provider.Filter(m3u, ultrabox.Include_groups, ultrabox.Remap_groups, ultrabox.Choffset, ultrabox.Group_prefix, ultrabox.Prefix_prio)
        m3u.Print()

        if *flag_with_xmltv != ""{
        }
    }

    if *flag_happytv_filter != ""{
        data, _ := ioutil.ReadFile(*flag_happytv_filter)
        m3u, err := tvg.Parse(data); if err != nil{
            panic(err)
        }

        m3u = provider.Filter(m3u, happytv.Include_groups, happytv.Remap_groups, happytv.Choffset, happytv.Group_prefix, happytv.Prefix_prio)
        m3u.Print()

        if *flag_with_xmltv != ""{
        }
    }

    if *flag_ultrabox_movies != ""{
        //data, _ := ioutil.ReadFile(*flag_ultrabox_movies)
        //m3u, err := tvg.Parse(data); if err != nil{
        //    panic(err)
        //}
    }

    if *flag_check_epg != ""{
        data, _ := ioutil.ReadFile(os.Args[0])
        m3u, err := tvg.Parse(data); if err != nil{
            panic(err)
        }
        tv, err := xmltv.ReadFile(*flag_check_epg); if err != nil{
            panic(err)
        }

        for _, inf := range m3u.List{
            for _, prog := range tv.Programme{
                if inf.Id == prog.Channel{
                    fmt.Println(inf.Id, prog.Channel)
                }
            }
        }
    }

    if *flag_gen_epg{

        tvmaster := xmltv.NewXMLTVFile()

        tvsources := make([]*xmltv.Tv, 0)

        data, _ := ioutil.ReadFile(flag.Args()[0])
        m3u, err := tvg.Parse(data); if err != nil{
            panic(err)
        }
        for _, tvpath := range flag.Args()[1:]{
            tv, err := xmltv.ReadFile(tvpath); if err != nil{
                panic(err)
            }
            tvsources = append(tvsources, tv)
        }

        var chcount, matchcount int

        for _, ch := range m3u.List{
            chcount++
            found := false
//            fmt.Printf("\nmatch? % 20s" , ch.NewName)
            for _, tv := range tvsources{
                for _, tvc := range tv.Channel{
                    //fmt.Println("\n", normalize(ch.NewName), normalize(tvc.Name))
                    if ch.Id == tvc.Id{
//                        fmt.Printf(" -> OMATC! CH % 20s -- |TVC| % -20s |CHID| % -20s |TVCID| % -20s", ch.NewName, tvc.Name, ch.Id, tvc.Id)
                        found = true
                    } else if normalize(ch.NewName) == normalize(tvc.Name){
//                        fmt.Printf(" -> MATCH! CH % 20s -- |TVC| % -20s |CHID| % -20s |TVCID| % -20s", ch.NewName, tvc.Name, ch.Id, tvc.Id)
                        found = true
                    }
                    if found{
                        programme := getProgramme(tv, tvc.Id)
                        if len(programme) == 0{
//                            fmt.Println("\nChannel match", ch.NewName, tvc.Id, "but no programme found")
                            found = false
                            continue
                        }

                        tvmaster.Channel = append(tvmaster.Channel, &xmltv.Channel{ch.Id, ch.Name})
                        tvmaster.Programme = append(tvmaster.Programme, programme...)
//                        fmt.Println(" ", len(programme), "programme found")
                        break
                    }
                }
                if found{
                    matchcount++
                    break
                }
            }
//            if !found{
//                fmt.Println(" -> NO MATCH!", normalize(ch.NewName))
//            }
        }

//        fmt.Println("\nCHCOUNT", chcount, "MATCHCOUNT", matchcount)

        tvmasterdata, _ := xmltv.Marshal(tvmaster)
        os.Stdout.Write(tvmasterdata)
    }
}

func getProgramme(tv *xmltv.Tv, id string) []*xmltv.Programme{
    programme := make([]*xmltv.Programme, 0)
    for _, tvp := range tv.Programme{
        if tvp.Channel == id{
            tvpmod := tvp
            tvpmod.Channel = id
            programme = append(programme, tvpmod)
        }
    }

    return programme
}

func normalize(name string) (norm string){

    norm = name

    r := regexp.MustCompile(`[^:]:+(.*)$`)

    sub := r.FindStringSubmatch(norm)
    //fmt.Println("NORM", norm, "SUB", sub)

    if len(sub) > 1{
        norm = sub[1]
    }

    replacer := strings.NewReplacer("one", "1",
        "two", "2",
        "three", "3",
        "four", "4")

    norm = strings.Replace(norm, "(east)", "", -1)
    norm = strings.Replace(norm, "FHD", "", -1)
    norm = strings.Replace(norm, "HDTV", "", -1)
    norm = strings.Replace(norm, "HD", "", -1)
    norm = strings.ToLower(norm)
    replacer.Replace(norm)
    norm = strings.Replace(norm, " ", "", -1)

    norm = strings.Trim(norm, " ")

    return norm
}
