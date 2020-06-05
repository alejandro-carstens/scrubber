package auth

import (
	"io/ioutil"
	"os"

	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/auth/contexts"

	"github.com/Jeffail/gabs"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleAuthService() *GoogleAuthService {
	return &GoogleAuthService{
		userRepository:          repositories.NewUserRepository(),
		accessControlRepository: repositories.NewAcessControlRepository(),
	}
}

type GoogleAuthService struct {
	userRepository          *repositories.UserRepository
	accessControlRepository *repositories.AccessControlRepository
}

func (gas *GoogleAuthService) Handle(context *contexts.GoogleAuthContext) (string, error) {
	profile, err := gas.exchangeCode(context.Code())

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
		if err := gas.userRepository.FindWhere(map[string]interface{}{
			"email != ?": profile.S("email").Data().(string),
		}, &users); err != nil {
			return "", err
		}

		if err := gas.userRepository.DB().Transaction(func(tx *gorm.DB) error {
			userRepository := gas.userRepository.FromTx(tx)

			user = &models.User{
				Email:         profile.S("email").Data().(string),
				EmailVerified: true,
				Name:          profile.S("given_name").Data().(string),
				LastName:      profile.S("family_name").Data().(string),
				Picture:       profile.S("picture").Data().(string),
			}

			if err := userRepository.Create(user); err != nil {
				return err
			}

			if len(users) > 0 {
				return nil
			}

			return gas.accessControlRepository.FromTx(tx).Create(&models.AccessControl{
				UserID: user.ID,
				Action: repositories.ACCESS_CONTROL_ALL_ACTIONS,
				Scope:  repositories.ACCESS_CONTROL_WRITE_SCOPE,
			})
		}); err != nil {
			return "", err
		}
	} else {
		user = users[0]
	}

	return issueJwt(user.ID, user.Name, user.Email, user.Picture)
}

func (gas *GoogleAuthService) exchangeCode(code string) (*gabs.Container, error) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			os.Getenv("GOOGLE_AUTH_SCOPE"),
		},
		Endpoint: google.Endpoint,
	}

	token, err := conf.Exchange(oauth2.NoContext, code)

	if err != nil {
		return nil, err
	}

	response, err := conf.Client(oauth2.NoContext, token).Get(GOOGLE_OAUTH2_URL)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return gabs.ParseJSON(data)
}
