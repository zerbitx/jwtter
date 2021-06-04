package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Verify() *cobra.Command {
	var jwtSigningKey string
	var verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "verify a JWT",
		Long:  `Uses signing key in config, or environment variable JWT_SIGNING_KEY`,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("You must pass a token")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if jwtSigningKey == "" {
				jwtSigningKey = viper.GetString("jwt_signing_key")
			}

			if jwtSigningKey == "" {
				log.Fatal("No signing key configured.")
			}

			tkn, err := jwt.Parse(args[0], func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSigningKey), nil
			})

			if err != nil {
				log.Fatal("Failed to verify token ", err)
			}

			headerBytes, _ := json.MarshalIndent(tkn.Header, "", " ")
			fmt.Println(string(headerBytes))

			claimsBytes, _ := json.MarshalIndent(tkn.Claims, "", " ")
			fmt.Println(string(claimsBytes))
		},
	}

	verifyCmd.Flags().StringVarP(&jwtSigningKey, "key", "k", "", "JWT signing key")

	return verifyCmd
}
