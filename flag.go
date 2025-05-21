package command

func NewFlag(name, short, description, defaultValue string, required bool) Flag {
	return Flag{
		Name:         name,
		Shorthand:    short,
		Description:  description,
		DefaultValue: defaultValue,
		Required:     required,
	}
}
