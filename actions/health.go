package actions

import "github.com/gobuffalo/buffalo"

func Health(c buffalo.Context) error {
	return c.Render(200, r.JSON("OK"))
}
