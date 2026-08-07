package icons

// FooterStatIcon returns compact footer count markers (open/ready/blocked/closed).
// Emoji mode keeps single-width ○◉◈● for alignment; nerd mode uses registry glyphs.
func FooterStatIcon(kind string) string {
	if ActiveSet() == SetNerd {
		switch kind {
		case "open":
			return Get(StatusGraphOpen)
		case "ready":
			return Get(CheckCircle)
		case "blocked":
			return Get(StatusBlocked)
		case "closed":
			return Get(StatusClosed)
		default:
			return "•"
		}
	}
	switch kind {
	case "open":
		return "○"
	case "ready":
		return "◉"
	case "blocked":
		return "◈"
	case "closed":
		return "●"
	default:
		return "•"
	}
}
