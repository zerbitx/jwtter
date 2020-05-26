package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify a JWT",
	Long:  `Uses signing key in config, or environment variable JWT_SIGNING_KEY`,
	Run: func(cmd *cobra.Command, args []string) {
		if jwtSigningKey == "" {
			jwtSigningKey = viper.GetString("jwt_signing_key")
		}

		if jwtSigningKey == "" {
			log.Fatal("No signing key configured.")
		}

		if len(args) == 0 {
			log.Fatal(`You must provide a set of claims even. example {"env":"dev","iss":"veridian-dynamics"}`)
		}

		mySigningKey := []byte(jwtSigningKey)

		if len(args) == 0 {
			log.Fatal("You must provide a JWT string")
		}

		tkn, err := jwt.Parse(args[0], func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		if err != nil {
			log.Fatal("Failed to verify token ", err)
		}

		jb, _ := json.MarshalIndent(tkn.Claims, "", " ")
		fmt.Println(string(jb))
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&jwtSigningKey, "key", "k", "", "JWT signing key")
}
