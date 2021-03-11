package userbusiness

import (
	"context"
	"fooddlv/appctx"
	"fooddlv/appctx/tokenprovider"
	"fooddlv/common"
	"fooddlv/module/user/usermodel"
	"go.opencensus.io/trace"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

type TokenConfig interface {
	GetAtExp() int
	GetRtExp() int
}

type loginBusiness struct {
	appCtx        appctx.AppContext
	storeUser     LoginStorage
	tokenProvider tokenprovider.Provider
	hasher        Hasher
	tkCfg         TokenConfig
}

func NewLoginBusiness(storeUser LoginStorage, tokenProvider tokenprovider.Provider,
	hasher Hasher, tkCfg TokenConfig) *loginBusiness {
	return &loginBusiness{
		storeUser:     storeUser,
		tokenProvider: tokenProvider,
		hasher:        hasher,
		tkCfg:         tkCfg,
	}
}

func (business *loginBusiness) Login(ctx context.Context, data *usermodel.UserLogin) (*usermodel.Account, error) {
	ctx2, span1 := trace.StartSpan(ctx, "user.biz.login")
	span1.AddAttributes(
		trace.StringAttribute("email", data.Email),
	)
	user, err := business.storeUser.FindUser(ctx2, map[string]interface{}{"email": data.Email})
	span1.End()

	if err != nil {
		return nil, common.ErrCannotGetEntity(usermodel.EntityName, err)
	}

	passHashed := business.hasher.Hash(data.Password + user.Salt)

	if user.Password != passHashed {
		return nil, usermodel.ErrUsernameOrPasswordInvalid
	}

	payload := tokenprovider.TokenPayload{
		UserId: user.Id,
		Role:   user.Role,
	}

	_, span2 := trace.StartSpan(ctx, "user.biz.gen_access_token")
	//defer span2.End() // Do not do this
	accessToken, err := business.tokenProvider.Generate(payload, business.tkCfg.GetAtExp())
	if err != nil {
		span2.End()
		return nil, common.ErrInternal(err)
	}
	span2.End()

	_, span3 := trace.StartSpan(ctx, "user.biz.gen_refresh_token")
	refreshToken, err := business.tokenProvider.Generate(payload, business.tkCfg.GetRtExp())
	if err != nil {
		span3.End()
		return nil, common.ErrInternal(err)
	}
	span3.End()

	account := usermodel.NewAccount(accessToken, refreshToken)

	return account, nil
}
