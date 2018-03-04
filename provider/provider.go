package provider

import (
    "../tvg"
    "fmt"
    "strings"
    "sort"
)

func groupById(include_groups map[string]bool, m3u *tvg.M3UData) (groups map[string][]*tvg.EXTINF){

    groups = make(map[string][]*tvg.EXTINF)

    for _, obj := range m3u.List{
        _, include := include_groups[obj.Group]; if !include{
            continue
        }

        // Set Id if empty
        if obj.Id == ""{
            obj.Id = strings.Replace(obj.Title, " ", "", -1)
        }
        _, exists := groups[obj.Id]; if !exists{
            groups[obj.Id] = make([]*tvg.EXTINF, 0)
        }

        groups[obj.Id] = append(groups[obj.Id], obj)
    }

    return
}

func groupByMatchName(include_groups map[string]bool, m3u *tvg.M3UData) (groups map[string][]*tvg.EXTINF){
    groups = make(map[string][]*tvg.EXTINF)

    for _, obj := range m3u.List{
        _, include := include_groups[obj.Group]; if !include{
            continue
        }

        _, exists := groups[obj.MatchName]; if !exists{
            groups[obj.MatchName] = make([]*tvg.EXTINF, 0)
        }

        groups[obj.MatchName] = append(groups[obj.MatchName], obj)
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

func sortByPrefix(prefix_prio map[string]string, m3u *tvg.M3UData){

    sort.Slice(m3u.List, func(i, j int) bool{

        xi := prefix_prio[m3u.List[i].Prefix]
        xj := prefix_prio[m3u.List[j].Prefix]

        return xi + m3u.List[i].NewName < xj + m3u.List[j].NewName
    })
}

func setupCustomPrefix(group_prefix map[string]string, inf *tvg.EXTINF){
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

func setupCustomGroup(remap_groups map[string]string, inf *tvg.EXTINF){
    _, remap := remap_groups[inf.Group]; if remap{
        inf.Group = remap_groups[inf.Group]
    }
}

func Filter(m3u *tvg.M3UData,
  include_groups map[string]bool,
  remap_groups map[string]string,
  choffset map[string]int,
  group_prefix map[string]string,
  prefix_prio map[string]string) *tvg.M3UData{

    groups := groupByMatchName(include_groups, m3u)
    newm3u := chooseBestQuality(groups)
    sortByPrefix(prefix_prio, newm3u)

    for _, obj := range newm3u.List{

        setupCustomGroup(remap_groups, obj)

        setupCustomPrefix(group_prefix, obj)

        choffset[obj.Group]++
        obj.Number = choffset[obj.Group]
    }

    return newm3u
}

