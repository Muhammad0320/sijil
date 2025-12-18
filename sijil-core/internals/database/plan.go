package database

type PlanLimits struct {
	RetentionDays int
	MaxMemebers   int
	MaxProject    int
}

var Plans = map[string]PlanLimits{
	"free": {
		RetentionDays: 3,
		MaxMemebers:   1,
		MaxProject:    3,
	},

	"pro": {
		RetentionDays: 30,
		MaxMemebers:   10,
		MaxProject:    15,
	},
	"ultra": {
		RetentionDays: 365,
		MaxMemebers:   100,
		MaxProject:    100,
	},
}

func GetPlanLimits(planName string) PlanLimits {
	if limit, ok := Plans[planName]; ok {
		return limit
	}
	return Plans["free"]
}