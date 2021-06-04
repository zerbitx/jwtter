package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

func New() *cobra.Command {
	var signingMethods = map[string]jwt.SigningMethod{
		"HS256": jwt.SigningMethodHS256,
		"HS512": jwt.SigningMethodHS512,
		"HS384": jwt.SigningMethodHS384,
		"ES256": jwt.SigningMethodES256,
		"ES384": jwt.SigningMethodES384,
		"ES512": jwt.SigningMethodES512,
		"PS256": jwt.SigningMethodPS256,
		"PS384": jwt.SigningMethodPS384,
		"PS512": jwt.SigningMethodPS512,
	}
	// Used for help messages
	var signingMethodList []string

	for sm := range signingMethods {
		signingMethodList = append(signingMethodList, sm)
	}

	// flags
	var signingMethod, jwtSigningKey, issuer string
	var addIat bool
	var issuedAt int64
	var duration time.Duration
	var headersStr string

	var newCmd = &cobra.Command{
		Use:     `new '{"json":"of", "your": "claims"}'`,
		Short:   "Create a new JWT",
		Long:    `Uses signing key in config, or environment variable JWT_SIGNING_KEY`,
		Example: `jwtter new -k somekey -s HS256 '<claims json>'` + "\n valid signing methods include: " + strings.Join(signingMethodList, ","),
		Run: func(cmd *cobra.Command, args []string) {

			if jwtSigningKey == "" {
				jwtSigningKey = viper.GetString("jwt_signing_key")
			}

			if jwtSigningKey == "" {
				log.Fatal("No signing key configured.")
			}

			if len(args) == 0 {
				log.Fatal(`You must provide a set of claims. example {"env":"dev","iss":"veridian-dynamics"}`)
			}

			claims := Claims{}
			if addIat {
				claims["iat"] = issuedAt
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

			method, ok := signingMethods[signingMethod]

			if !ok {
				log.Fatalf("%s is not a valid signing method\n Valid methods include: %v", signingMethod, signingMethodList)
			}

			token := makeToken(headersStr, method, claims)

			ss, err := token.SignedString([]byte(jwtSigningKey))

			if err != nil {
				log.Fatal("Failed to encode JWT ", err)
			}

			fmt.Println(ss)
		},
	}

	flags := newCmd.Flags()

	flags.DurationVarP(&duration, "duration", "d", time.Duration(0), `-d 120h (Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".)`)
	flags.StringVarP(&issuer, "issuer", "i", "", "-i veridian-dynamics")
	flags.Int64VarP(&issuedAt, "issuedAt", "a", int64(time.Now().Unix()), "-a 1590524782")
	flags.BoolVarP(&addIat, "addIat", "n", true, "don't created the default iat (issued at) claim")
	flags.StringVarP(&jwtSigningKey, "key", "k", "", "JWT signing key")
	flags.StringVarP(&signingMethod, "signingMethod", "s", "HS256", "JWT signing method")
	flags.StringVarP(&headersStr, "headers", "", "", `Defaults to {"typ":"JWT","alg":"<signing method>"} if unset`)

	return newCmd
}

func makeToken(headersStr string, method jwt.SigningMethod, claims Claims) *jwt.Token {
	var token *jwt.Token
	// if headers were specified by the user, use them, but set alg properly
	if headersStr != "" {
		headers := map[string]interface{}{}

		if err := json.Unmarshal([]byte(headersStr), &headers); err != nil {
			log.Fatalf("failed to map headers: %s", err)
		}

		headers["alg"] = method.Alg()

		token = &jwt.Token{
			Header: headers,
			Claims: claims,
			Method: method,
		}
	} else {
		token = jwt.NewWithClaims(method, claims)
	}
	return token
}
