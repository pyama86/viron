package controller

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"

	"github.com/cam-inc/dmc/example-go/common"
	"github.com/cam-inc/dmc/example-go/gen/app"
	"github.com/cam-inc/dmc/example-go/models"
	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"

	"golang.org/x/crypto/scrypt"
)

// AuthController implements the auth resource.
type AuthController struct {
	*goa.Controller
	privateKey *rsa.PrivateKey
}

// NewAuthController creates a auth controller.
func NewAuthController(service *goa.Service) *AuthController {
	b := common.GetPrivateKey()
	privateKey, err := jwtgo.ParseRSAPrivateKeyFromPEM([]byte(b))
	if err != nil {
		panic(err)
	}
	return &AuthController{
		Controller: service.NewController("AuthController"),
		privateKey: privateKey,
	}
}

// Signin runs the signin action.
func (c *AuthController) Signin(ctx *app.SigninAuthContext) error {
	// AuthController_Signin: start_implement

	// Put your logic here
	// Authorize
	adminUserTable := models.NewAdminUserDB(common.DB)
	m, err := adminUserTable.GetByLoginID(ctx.Context, *ctx.Payload.LoginID)
	if err == gorm.ErrRecordNotFound {
		return ctx.NotFound()
	} else if err != nil {
		panic(err)
	}

	hash, err := scrypt.Key([]byte(*ctx.Payload.Password), []byte(m.Salt), 16384, 8, 1, 64)
	if m.Password != base64.StdEncoding.EncodeToString(hash) {
		return ctx.Unauthorized()
	}

	// Generate JWT
	token := jwtgo.New(jwtgo.SigningMethodRS512)
	in1day := time.Now().Add(time.Duration(24) * time.Hour).Unix()
	token.Claims = jwtgo.MapClaims{
		"iss":    "DMC",                 // Token発行者
		"aud":    "dmc.local",           // このTokenを利用する対象の識別子
		"exp":    in1day,                // Tokenの有効期限
		"jti":    uuid.NewV4().String(), // Tokenを一意に識別するためのID
		"iat":    time.Now().Unix(),     // Tokenを発行した日時(now)
		"nbf":    0,                     // Tokenが有効になるのが何分後か
		"sub":    *ctx.Payload.LoginID,  // ユーザー識別子
		"scopes": "api:access",          // このTokenが有効なSCOPE - not a standard claim
		// TODO: roleも入れる
	}
	signedToken, err := token.SignedString(c.privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %s", err)
	}

	// Set auth header for client retrieval
	ctx.ResponseData.Header().Set("Authorization", fmt.Sprintf("Bearer %s", signedToken))

	// AuthController_Signin: end_implement
	return nil
}

// Signout runs the signout action.
func (c *AuthController) Signout(ctx *app.SignoutAuthContext) error {
	// AuthController_Signout: start_implement

	// Put your logic here
	ctx.ResponseData.Header().Del("Authorization")

	// AuthController_Signout: end_implement
	return nil
}