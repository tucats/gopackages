package compiler

// Block compiles a statement block. The leading { has already
// been parse.
func (c *Compiler) Block() error {

	parsing := true
	for parsing {

		if c.t.IsNext("}") {
			break
		}

		err := c.Statement()
		if err != nil {
			return err
		}

		if c.t.IsNext(";") {
			// No action needed
		}

		if c.t.AtEnd() {
			return c.NewError("unclosed statement block")
		}
	}
	return nil
}
