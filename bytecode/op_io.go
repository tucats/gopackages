package bytecode

import (
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*           B A S I C   I / O             *
*                                         *
\******************************************/

// PrintOpcode implementation. If the operand
// is given, it represents the number of items
// to remove from the stack.
func PrintOpcode(c *Context, i interface{}) error {

	count := 1
	if i != nil {
		count = util.GetInt(i)
	}

	for n := 0; n < count; n = n + 1 {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		s := util.FormatUnquoted(v)
		if c.output == nil {
			fmt.Printf("%s", s)
		} else {
			c.output.WriteString(s)
		}
	}

	// If we are instruction tracing, print out a newline anyway so the trace
	// display isn't made illegible.
	if c.output == nil && c.Tracing {
		fmt.Println()
	}

	return nil
}

// LogOpcode imeplements the Log option.
func LogOpcode(c *Context, i interface{}) error {
	logger := util.GetString(i)
	msg, err := c.Pop()
	if err == nil {
		ui.Debug(logger, "%v", msg)
	}
	return err
}

// SayOpcode implementation. This can be used in place
// of NewLine to end buffered output, but the output is
// only displayed if we are not in --quiet mode.
func SayOpcode(c *Context, i interface{}) error {
	ui.Say("%s\n", c.output.String())
	c.output = nil
	return nil
}

// NewlineOpcode implementation.
func NewlineOpcode(c *Context, i interface{}) error {

	if c.output == nil {
		fmt.Printf("\n")
	} else {
		c.output.WriteString("\n")
	}
	return nil
}

/******************************************\
*                                         *
*           T E M P L A T E S             *
*                                         *
\******************************************/

// TemplateOpcode compiles a template string from the
// stack and stores it in the template manager for the
// context.
func TemplateOpcode(c *Context, i interface{}) error {

	name := util.GetString(i)
	t, err := c.Pop()
	if err == nil {
		t, err = template.New(name).Parse(util.GetString(t))
		if err == nil {
			err = c.Push(t)
		}
	}
	return err
}

/******************************************\
*                                         *
*            R E S T   I / O              *
*                                         *
\******************************************/

func AuthOpcode(c *Context, i interface{}) error {

	if _, ok := c.Get("_authenticated"); !ok {
		return c.NewError(NotAServiceError)
	}
	kind := util.GetString(i)
	var user, pass string
	if v, ok := c.Get("_user"); ok {
		user = util.GetString(v)
	}
	if v, ok := c.Get("_password"); ok {
		user = util.GetString(v)
	}
	tokenValid := false
	if v, ok := c.Get("_token_valid"); ok {
		tokenValid = util.GetBool(v)
	}

	if (kind == "token" || kind == "tokenadmin") && !tokenValid {
		_ = c.SetAlways("_rest_status", 403)
		if c.output != nil {
			c.output.WriteString("403 Forbidden")
		}
		c.running = false
		ui.Debug(ui.ServerLogger, "@authenticated token: no valid token")
		return nil
	}

	if kind == "user" && user == "" && pass == "" {
		_ = c.SetAlways("_rest_status", 401)
		if c.output != nil {
			c.output.WriteString("401 Not authorized")
		}
		c.running = false
		ui.Debug(ui.ServerLogger, "@authenticated user: no credentials")
		return nil
	} else {
		kind = "any"
	}

	if kind == "any" {
		isAuth := false
		if v, ok := c.Get("_authenticated"); ok {
			isAuth = util.GetBool(v)
		}
		if !isAuth {
			_ = c.SetAlways("_rest_status", 403)
			if c.output != nil {
				c.output.WriteString("403 Forbidden")
			}
			c.running = false
			ui.Debug(ui.ServerLogger, "@authenticated any: not authenticated")
			return nil
		}
	}

	if kind == "admin" || kind == "admintoken" {
		isAuth := false
		if v, ok := c.Get("_superuser"); ok {
			isAuth = util.GetBool(v)
		}
		if !isAuth {
			_ = c.SetAlways("_rest_status", 403)
			if c.output != nil {
				c.output.WriteString("403 Forbidden")
			}
			c.running = false
			ui.Debug(ui.ServerLogger, fmt.Sprintf("@authenticated %s: not admin", kind))
		}
	}

	return nil
}

func ResponseOpcode(c *Context, i interface{}) error {

	// See if we have a media type specified.
	isJSON := false
	if v, found := c.Get("_json"); found {
		isJSON = util.GetBool(v)
	}

	var output string
	v, err := c.Pop()
	if err != nil {
		return err
	}

	if isJSON {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		output = string(b)
	} else {
		output = util.FormatUnquoted(v)
	}

	if c.output == nil {
		fmt.Println(output)
	} else {
		c.output.WriteString(output)
		c.output.WriteRune('\n')
	}
	return nil
}
