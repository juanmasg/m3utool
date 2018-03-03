package ultrabox

var (
    Include_groups = map[string]bool{
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

    Remap_groups = map[string]string{
        "International Sports": "Sports",
        "Sports": "Sports",
        "Sports Pass": "Sports",
        "EPIC EVENT": "Sports",
        "English": "USA",
    }

    Choffset = map[string]int{
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

    Group_prefix = map[string]string{
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

    Prefix_prio = map[string]string{
        "UK": "a",
        "UK HD": "a",
        "UK FHD": "a",
        "USA": "b",
    }

)

