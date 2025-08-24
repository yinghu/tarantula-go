package item

type Validator func(prop string, c Configuration)

func ItemValidator(c Configuration, validator Validator) {
	ins := c.Reference
	for k, v := range ins {
		ref, ok := v.([]Configuration)
		if ok {
			for i := range ref {
				validator(k, ref[i])
				ItemValidator(ref[i], validator)
			}
		}
	}
}
