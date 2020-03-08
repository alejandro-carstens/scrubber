package options

var availableThresholdTypes []string = []string{"count", "average_count", "stats"}

var availableAlertChannels []string = []string{"slack", "email", "log", "pager_duty"}

var availableIntervalUnits []string = []string{"seconds", "minutes", "hours", "days", "months", "years"}

var availableOperators []string = []string{"=", "<>", ">", "<", "<=", ">="}

var availableMetrics []string = []string{
	"min",
	"max",
	"avg",
	"sum_of_squares",
	"variance",
	"std_deviation",
	"upper_std_deviation_bound",
	"lower_std_deviation_bound",
}

var nonMatchingClauses []string = []string{"limit", "order_by", "order_by_nested"}

var availableInClauses []string = []string{
	"where_in",
	"where_not_in",
	"where_nested_in",
	"where_nested_not_in",
	"filter_in",
	"filter_in_nested",
	"match_in",
	"match_in_nested",
	"match_not_in",
	"match_not_in_nested",
	"match_phrase_in",
	"match_phrase_in_nested",
	"match_phrase_not_in",
	"match_phrase_not_in_nested",
}

var availableClauses []string = append([]string{
	"where",
	"where_nested",
	"filter",
	"filter_nested",
	"match",
	"match_nested",
	"match_phrase",
	"match_phrase_nested",
}, append(availableInClauses, nonMatchingClauses...)...)
