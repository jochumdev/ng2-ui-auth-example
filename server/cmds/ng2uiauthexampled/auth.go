package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth"
)

// Command represents the auth command
var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Authentication stuff",
	Long:  ``,
}

var authCommandAddUser = &cobra.Command{
	Use:   "adduser <username> <password> <role> [role...]",
	Short: "Add a user",
	Long:  ``,
	Run: commandWrapper(func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("ERROR: Give an username and a password.")
			os.Exit(1)
		}

		// Set username and password
		u := &auth.User{
			Username: args[0],
		}
		u.PasswordSet(args[1])

		// Add Roles
		if len(args) > 2 {
			for _, role := range args[2:] {
				// Ignore duplicated roles here.
				u.RoleAdd(role)
			}
		}

		db, apiErr := auth.DBGet()
		if apiErr != nil {
			fmt.Printf("ERROR: %v\n", apiErr.Reason)
			os.Exit(1)
		}

		if apiErr = db.UserCreate(u); apiErr != nil {
			fmt.Printf("ERROR: While creating the user: %v\n", apiErr)
			os.Exit(1)
		}

		fmt.Printf("User '%s' successfully created.\n", args[0])
	}),
}

var authCommandRemoveUser = &cobra.Command{
	Use:   "removeuser <username>",
	Short: "Remove a user",
	Long:  ``,
	Run: commandWrapper(func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("ERROR: Give a username to remove.")
			os.Exit(1)
		}

		db, apiErr := auth.DBGet()
		if apiErr != nil {
			fmt.Printf("ERROR: %v\n", apiErr.Reason)
			os.Exit(1)
		}

		if apiErr = db.UserDelete(args[0]); apiErr != nil {
			fmt.Printf("ERROR: While removing the user: %v\n", apiErr.Reason)
			os.Exit(1)
		}

		fmt.Printf("User '%s' successfully removed.\n", args[0])
	}),
}

var authCommandLogin = &cobra.Command{
	Use:   "login <username>",
	Short: "Login to generate a JWT for testing use",
	Long:  ``,
	Run: commandWrapper(func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("ERROR: Give an username.")
			os.Exit(1)
		}

		db, apiErr := auth.DBGet()
		if apiErr != nil {
			fmt.Printf("ERROR: %v\n", apiErr.Reason)
			os.Exit(1)
		}

		user, apiErr := db.UserFindByUsername(args[0])
		if apiErr != nil {
			if apiErr.Reason == auth.ErrWrongUsernameOrPassword {
				fmt.Println("ERROR: User not found.")
				os.Exit(1)
			}

			fmt.Printf("ERROR: %v\n", apiErr.Reason)
			os.Exit(1)
		}

		token, err := user.TokenGenerate()
		if err != nil {
			fmt.Printf("ERROR: Failed to generate the token %v\n", err)
			os.Exit(1)
		}
		fmt.Println(token)
	}),
}

func init() {
	authCommand.AddCommand(authCommandAddUser)
	authCommand.AddCommand(authCommandRemoveUser)
	authCommand.AddCommand(authCommandLogin)
}
