package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/deis/deis/client-go/controller/client"
	"github.com/deis/deis/client-go/controller/models/apps"
	"github.com/deis/deis/client-go/controller/models/config"
)

// AppCreate creates an app with the application
func AppCreate(id string, buildpack string, remote string, noRemote bool) error {
	c, err := client.Load()

	fmt.Print("Creating Application... ")
	quit := progress()
	app, err := apps.Create(*c, id)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Printf("done, created %s\n", app.ID)

	if buildpack != "" {
		configValues := map[string]string{
			"BUILDPACK_URL": buildpack,
		}
		err = config.Set(*c, app.ID, configValues)

		if err != nil {
			return err
		}
	}

	if !noRemote {
		return c.CreateRemote(remote, app.ID)
	}

	fmt.Println("remote available at", c.RemoteURL(app.ID))

	return nil
}

// AppsList lists apps on the controller
func AppsList() error {
	c, err := client.Load()

	if err != nil {
		return err
	}

	users, err := apps.List(*c)

	if err != nil {
		return err
	}

	fmt.Println("=== Apps")

	for _, user := range users {
		fmt.Println(user.ID)
	}
	return nil
}

// AppInfo prints info about app
func AppInfo(appID string) error {
	c, err := client.Load()

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = c.DetectApp()

		if err != nil {
			return err
		}
	}

	app, err := apps.Get(*c, appID)

	if err != nil {
		return err
	}

	fmt.Printf("=== %s Application\n", app.ID)
	fmt.Println("updated: ", app.Updated)
	fmt.Println("uuid: ", app.UUID)
	fmt.Println("created: ", app.Created)
	fmt.Println("url: ", app.URL)
	fmt.Println("owner: ", app.Owner)
	fmt.Println("id: ", app.ID)

	return nil
}

// AppOpen opens an app in the default webbrowser
func AppOpen(appID string) error {
	c, err := client.Load()

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = c.DetectApp()

		if err != nil {
			return err
		}
	}

	app, err := apps.Get(*c, appID)

	if err != nil {
		return err
	}

	URL, err := url.Parse(app.URL)

	if err != nil {
		return err
	}

	URL, err = chooseScheme(*URL)

	return client.Webbrowser(URL.String())
}

// AppLogs return the logs from an app
func AppLogs(appID string, lines int) error {
	c, err := client.Load()

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = c.DetectApp()

		if err != nil {
			return err
		}
	}

	logs, err := apps.Logs(*c, appID, lines)

	if err != nil {
		return err
	}

	for _, log := range strings.Split(strings.Trim(logs, `\n`), `\n`) {
		catagory := strings.Split(strings.Split(log, ": ")[0], " ")[1]
		printColor(log, chooseColor(catagory))
	}

	return nil
}

// AppRun runs a one time command in the app
func AppRun(appID, command string) error {
	c, err := client.Load()

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = c.DetectApp()

		if err != nil {
			return err
		}
	}

	fmt.Printf("Running '%s'...\n", command)

	out, err := apps.Run(*c, appID, command)

	if err != nil {
		return err
	}

	fmt.Print(out.Output)
	os.Exit(out.ReturnCode)
	return nil
}

// AppDestroy destroys an app
func AppDestroy(appID, confirm string) error {
	gitSession := false

	c, err := client.Load()

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = c.DetectApp()

		if err != nil {
			return err
		}

		gitSession = true
	}

	if confirm == "" {
		fmt.Printf(` !    WARNING: Potentially Destructive Action
 !    This command will destroy the application: %s
 !    To proceed, type "%s" or re-run this command with --confirm=%s

>`, appID, appID, appID)

		fmt.Scanln(&confirm)
	}

	if confirm != appID {
		return fmt.Errorf("App %s does not match confirm %s, aborting.", appID, confirm)
	}

	fmt.Printf("Destroying %s...", appID)

	err = apps.Delete(*c, appID)

	if err != nil {
		return err
	}

	if gitSession {
		return c.DeleteRemote(appID)
	}

	return nil
}
