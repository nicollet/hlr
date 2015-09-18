package parse

type Variables map[string]string

func BuildVariables(input string) *Variables {
	l := lex(input)

	v := make(Variables)
	var varName string
	for item := range l.items {
		switch item.typ {
		case itemVarName:
			varName = item.val
		case itemVarValue:
			// Does the value already exists in our variables ?
			content := item.val
			newVal, replace := v[item.val]
			if replace {
				content = newVal
			}

			// Add or create new variable
			_, ok := v[varName]
			if ok {
				v[varName] += "," + content
			} else {
				v[varName] = content
			}
		}
	}
	return &v
}

// vim: set sw=2 ts=2 list:
