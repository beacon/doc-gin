package gin

func ExampleEngine() {
	g := NewEngine(true)
	g.GET("/doc", nil, func(c *Context) {

	})
}
