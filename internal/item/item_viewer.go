package item

type Viewer func(prop string, c Configuration)

func ItemView(c Configuration, viewer Viewer) {
	ins := c.Reference
	for k, v := range ins {
		//fmt.Printf("Key %s\n", k)
		ref, ok := v.([]Configuration)
		if ok {
			for i := range ref {
				viewer(k, ref[i])
				ItemView(ref[i], viewer)
			}
		}
	}
}
