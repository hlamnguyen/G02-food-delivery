package userstorage

import (
	"context"
	"fooddlv/common"
	"fooddlv/module/user/usermodel"
	"go.opencensus.io/trace"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error) {
	_, span := trace.StartSpan(ctx, "user.mysql_store.find")
	defer span.End()

	db := s.db.Table(usermodel.User{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var user usermodel.User

	if err := db.Where(conditions).First(&user).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return &user, nil
}
