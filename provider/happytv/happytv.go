package happytv

var (
    Include_groups = map[string]bool{
        "USA": true,
        "VIP Sports UK": true,
        "VIP Sports Canada": true,
        "VIP Sports USA": true,
        "UK General": true,
        "UK Entertainment": true,
        "Canada": true,
        "UK Documentaries": true,
        "UK Movies": true,
        "UK Kids": true,
        "UK News": true,
    }

    Remap_groups = map[string]string{
        "Canada": "CAN",
        "VIP Sports UK": "Sports",
        "VIP Sports Canada": "Sports",
        "VIP Sports USA": "Sports",
        "UK General": "UK",
        "UK Entertainment": "UK",
        "UK Documentaries": "UK",
        "UK News": "UK",
        "UK Kids": "Kids",
        "UK Movies": "Movies",
    }

    Choffset = map[string]int{
        "UK":           2000,
        "Sports":       8000,
        "Kids":         6000,
        "Movies":       5000,
        "USA":          3000,
        "CAN":          4000,
    }

    Group_prefix = map[string]string{
        "UK": "UK",
        "Sports": "Sports",
        "Kids": "Kids",
        "Movies": "Movies",
        "USA": "USA",
        "CAN": "CAN",
    }

    Prefix_prio = map[string]string{
        "UK": "a",
        "USA": "b",
        "CAN": "c",
    }
)
