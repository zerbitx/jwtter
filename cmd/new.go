package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	// Claims wraps a generic map to allow for arbitrary input claims.
	Claims map[string]interface{}
)

// Valid implements jwt.Claims
func (c Claims) Valid() error {
	return nil
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   `new '{"json":"of", "your": "claims"}'`,
	Short: "Create a new JWT",
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

		claims := Claims{
			"iat": issuedAt,
		}
		err := json.Unmarshal([]byte(args[0]), &claims)

		if err != nil {
			log.Fatal("Failed to map claims ", err)
		}

		if duration != 0 {
			claims["exp"] = time.Now().Add(duration).Unix()
		}

		if issuer != "" {
			claims["iss"] = issuer
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(mySigningKey)

		if err != nil {
			log.Fatal("Failed to encode JWT ", err)
		}

		fmt.Println(ss)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	flags := newCmd.Flags()

	flags.DurationVarP(&duration, "duration", "d", time.Duration(0), `-d 120h (Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".)`)
	flags.StringVarP(&issuer, "issuer", "i", "", "-i veridian-dynamics")
	flags.Int64VarP(&issuedAt, "issuedAt", "a", int64(time.Now().Unix()), "-a 1590524782")
	flags.StringVarP(&jwtSigningKey, "key", "k", "", "JWT signing key")
}
