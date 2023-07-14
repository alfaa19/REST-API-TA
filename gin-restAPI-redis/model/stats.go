package model

type Stats struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	Tag              string `json:"tag"`
	Region           string `json:"region"`
	Rating           string `json:"rating"`
	Damage_round     string `json:"damage_round"`
	Headshot_percent string `json:"headshot_percent"`
	First_bloods     string `json:"first_bloods"`
	Kills            string `json:"kills"`
	Deaths           string `json:"deaths"`
	Assists          string `json:"assists"`
	Kd_ratio         string `json:"kd_ratio"`
}