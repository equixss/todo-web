package users_service

import auth_core "github.com/equixss/todo-web/internal/core/auth"

type AuthConfigAdapter struct {
	config auth_core.Config
}

func NewAuthConfigAdapter(config auth_core.Config) *AuthConfigAdapter {
	return &AuthConfigAdapter{config: config}
}

func (a *AuthConfigAdapter) GetSecret() string {
	return a.config.Secret
}

func (a *AuthConfigAdapter) GetExpiry() string {
	return a.config.Expiry.String()
}

func (a *AuthConfigAdapter) GetIssuer() string {
	return a.config.Issuer
}
