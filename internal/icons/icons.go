package icons

import "komorebit/internal/icons/generated"

var icons map[string][]byte

func init() {
	icons = make(map[string][]byte)
	loadAllIcons()
}

func loadAllIcons() {
	icons["one"] = generated.One
	icons["two"] = generated.Two
	icons["three"] = generated.Three
	icons["four"] = generated.Four
	icons["five"] = generated.Five
	icons["six"] = generated.Six
	icons["seven"] = generated.Seven
	icons["eight"] = generated.Eight
	icons["nine"] = generated.Nine
	icons["ten"] = generated.Ten
	icons["pause"] = generated.Pause
	icons["sad"] = generated.Sad
	icons["tilde"] = generated.Tilde
}

func WorkspaceIcon(workspaceIndex int) []byte {
	switch int(workspaceIndex) {
	case 0:
		return icons["one"]
	case 1:
		return icons["two"]
	case 2:
		return icons["three"]
	case 3:
		return icons["four"]
	case 4:
		return icons["five"]
	case 5:
		return icons["six"]
	case 6:
		return icons["seven"]
	case 7:
		return icons["eight"]
	case 8:
		return icons["nine"]
	case 9:
		return icons["ten"]
	default:
		return icons["tilde"]
	}
}

func PauseIcon() []byte {
	return icons["pause"]
}

func SadIcon() []byte {
	return icons["sad"]
}

func TildeIcon() []byte {
	return icons["tilde"]
}
