package auth

import (
	"io/ioutil"
	"os"

	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/auth/contexts"

	"github.com/Jeffail/gabs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleAuthService() *GoogleAuthService {
	return &GoogleAuthService{
		userRepository: repositories.UserRepo(),
	}
}

type GoogleAuthService struct {
	userRepository *repositories.UserRepository
}

func (gas *GoogleAuthService) Handle(context *contexts.GoogleAuthContext) (string, error) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			os.Getenv("GOOGLE_AUTH_SCOPE"),
		},
		Endpoint: google.Endpoint,
	}

	token, err := conf.Exchange(oauth2.NoContext, context.Code())

	if err != nil {
		return "", err
	}

	response, err := conf.Client(oauth2.NoContext, token).Get(GOOGLE_OAUTH2_URL)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	profile, err := gabs.ParseJSON(data)

	if err != nil {
		return "", err
	}

	users := []*models.User{}

	if err := gas.userRepository.FindWhere(map[string]interface{}{
		"email = ?": profile.S("email").Data().(string),
	}, &users); err != nil {
		return "", err
	}

	var user *models.User

	if len(users) == 0 {
		user = &models.User{
			Email:         profile.S("email").Data().(string),
			EmailVerified: true,
			Name:          profile.S("name").Data().(string),
			LastName:      profile.S("last_name").Data().(string),
			Picture:       profile.S("picture").Data().(string),
		}

		if err := gas.userRepository.Create(user); err != nil {
			return "", err
		}
	} else {
		user = users[0]
	}

	return issueJwt(user.ID, user.Name, user.Email, user.Picture)
}
