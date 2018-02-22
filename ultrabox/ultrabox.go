package ultrabox

import (
    "../tvg"
    "fmt"
    "strings"
    "sort"
)

var (
    include_groups = map[string]bool{
        "EPIC EVENT":           true,
        "UK":		            true,
        "Sports Pass":		    true,
        "Sports":		        true,
        "International Sports": true,
        "Kids":                 true,
        "Movie Networks":       true,
        "News Networks":        true,
//        "Music Pass":           true,
        "English":              true,
        "USA":              true,
    }

    remap_groups = map[string]string{
        "International Sports": "Sports",
        "Sports": "Sports",
        "Sports Pass": "Sports",
        "EPIC EVENT": "Sports",
        "English": "USA",
    }

    choffset = map[string]int{
//        "EPIC EVENT":           true,
        "UK":		            2000,
        "Sports Pass":		    9000,
        "Sports":		        8000,
        "International Sports": 7000,
        "Kids":                 6000,
        "Movie Networks":       4000,
        "News Networks":        5000,
//        "Music Pass":           true,
        "USA":              3000,
    }

    group_prefix = map[string]string{
        "English": "USA",
        "USA": "USA",
        "News Networks": "USA",
        "Movie Networks": "USA",
        "Kids": "USA",
        "Sports": "USA",
        "Sports Pass": "USA",
        "International Sports": "USA",
        "EPIC EVENT": "USA",
    }

    prefix_prio = map[string]string{
        "UK": "a",
        "UK HD": "a",
        "UK FHD": "a",
        "USA": "b",
    }

)

func groupById(m3u *tvg.M3UData) (groups map[string][]*tvg.EXTINF){

    groups = make(map[string][]*tvg.EXTINF)

    for _, obj := range m3u.List{
        _, include := include_groups[obj.Group]; if !include{
            continue
        }

        // Remove SD/HD/FHD duplicates
        if obj.Id == ""{
            obj.Id = strings.Replace(obj.Title, " ", "", -1)
            //fmt.Println("New ID for", obj.Title, obj.Id)
        }
        _, exists := groups[obj.Id]; if !exists{
            groups[obj.Id] = make([]*tvg.EXTINF, 0)
        }

        groups[obj.Id] = append(groups[obj.Id], obj)
    }

    return
}

func chooseBestQuality(groups map[string][]*tvg.EXTINF) (m3u *tvg.M3UData){
    m3u = &tvg.M3UData{make([]*tvg.EXTINF, 0)}

    // Priority FHD > HD > ""
    for _, l := range groups{
        idx := -1
        idxq := ""
        for k, v := range l{
            if v.FHD{
                idx = k
                idxq = "FHD"
                break
            }else if v.HD{
                idx = k
                idxq = "HD"
            }else if v.SD{
                if idxq == "HD"{ continue }
                idx = k
                idxq = "SD"
            }
        }

        m3u.List = append(m3u.List, l[idx])
    }

    return m3u
}

func sortByPrefix(m3u *tvg.M3UData){

    sort.Slice(m3u.List, func(i, j int) bool{

        xi := prefix_prio[m3u.List[i].Prefix]
        xj := prefix_prio[m3u.List[j].Prefix]

        return xi + m3u.List[i].NewName < xj + m3u.List[j].NewName
    })
}

func setupCustomPrefix(inf *tvg.EXTINF){
    if inf.Prefix == ""{
        inf.Prefix = group_prefix[inf.Group]
    }else{ //FIXME needed?
        inf.Prefix = strings.Replace(inf.Prefix, "FHD", "", -1)
        inf.Prefix = strings.Replace(inf.Prefix, "HD", "", -1)
        inf.Prefix = strings.Trim(inf.Prefix, " ")
    }
    inf.Name = fmt.Sprintf("%s: %s", inf.Prefix, inf.Name)
    inf.Title = fmt.Sprintf("%s: %s", inf.Prefix, inf.Title)
}

func setupCustomGroup(inf *tvg.EXTINF){
    _, remap := remap_groups[inf.Group]; if remap{
        inf.Group = remap_groups[inf.Group]
    }
}

//func setupCustomName(inf *tvg.EXTINF){
//    replacer := strings.Replace
//}

func Filter(m3u *tvg.M3UData) *tvg.M3UData{
    groups := groupById(m3u)
    newm3u := chooseBestQuality(groups)
    sortByPrefix(newm3u)

    for _, obj := range newm3u.List{

        setupCustomGroup(obj)

        setupCustomPrefix(obj)

        choffset[obj.Group]++
        obj.Number = choffset[obj.Group]
    }

    return newm3u
}

