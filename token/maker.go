package token

import "time"

type Maker interface {
	CreateToken(userName string, duration time.Duration)

	verifyToken(token string) (*Payload, error)
}
