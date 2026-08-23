package ntp

func StratumName(s uint8) string {
	switch {
	case s == 0:
		return "kiss"
	case s == 1:
		return "primary"
	case s < 16:
		return "secondary"
	default:
		return "unsynchronized"
	}
}
