package setup

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) GetSetup(c *fiber.Ctx) error {
	// HTML Create an HRMIS Administrator form template
	htmlTemplate := `
    <html>
    <head>
        <title>Create an HRMIS Administrator</title>
    </head>
    <body>
        <h2>Create an HRMIS Administrator</h2>
        <form method="post" action="/setup">
            <label for="username">Username:</label>
            <input type="text" id="username" name="username"><br><br>
            <label for="email">Email:</label>
            <input type="text" id="email" name="email"><br><br>
            <label for="password">Password:</label>
            <input type="text" id="password" name="password"><br><br>
            <label for="repassword">Re-Type Password:</label>
            <input type="text" id="repassword" name="repassword"><br><br>
            <input type="submit" value="Submit">
        </form>
    </body>
    </html>
    `

	// Parse the HTML template
	tmpl, err := template.New("setup").Parse(htmlTemplate)
	if err != nil {
		// Handle parsing error
		return err
	}

	// Execute the template and send it as the response
	c.Set("Content-Type", "text/html") // Set the content type to HTML
	if err := tmpl.Execute(c.Response().BodyWriter(), nil); err != nil {
		// Handle execution error
		return err
	}

	return nil
}
