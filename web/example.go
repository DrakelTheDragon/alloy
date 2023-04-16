package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/drakelthedragon/alloy/web/request"
	"github.com/drakelthedragon/alloy/web/response"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalServer     = errors.New("internal server error")
)

type (
	LoginInput            struct{ EmailAddress string }
	LoginOutput           struct{ AccessToken string }
	LoginProcessor        struct{ users map[string]string }
	LoginRequestProcessor = requestProcessor[LoginInput, LoginOutput]
)

func NewLoginProcessor() *LoginProcessor {
	return &LoginProcessor{users: map[string]string{"user@example.com": "supersecrettoken"}}
}

func NewLoginRequestMapper() RequestMapperFunc[LoginInput] {
	return func(r *http.Request, input *LoginInput) error {
		var req struct {
			EmailAddress string `json:"email_address"`
		}

		if err := request.UnJSONify(r, &req); err != nil {
			return err
		}

		input.EmailAddress = req.EmailAddress

		return nil
	}
}

func NewLoginResponseMapper() ResponseMapperFunc[LoginOutput] {
	return func(result Result[LoginOutput]) Response {
		output, err := result.Unwrap()
		if err != nil {
			switch {
			case errors.Is(err, ErrInvalidCredentials):
				return response.Unauthorized()
			default:
				return response.InternalServerError(err)
			}
		}

		return response.OK(struct {
			AccessToken string `json:"access_token"`
		}{AccessToken: output.AccessToken})
	}
}

func NewLoginRequestProcessor() *LoginRequestProcessor {
	return NewRequestProcessor[LoginInput, LoginOutput](
		NewLoginRequestMapper(),
		NewLoginResponseMapper(),
		NewLoginProcessor().Process,
	)
}

func (p *LoginProcessor) Process(ctx context.Context, input LoginInput) (LoginOutput, error) {
	token, ok := p.users[input.EmailAddress]
	if !ok {
		return LoginOutput{}, errors.New("invalid credentials")
	}

	return LoginOutput{AccessToken: token}, nil
}
