package notemodel

import "fooddlv/common"

type NoteCreate struct {
	common.SQLModelCreate `json:",inline"`
	UserId                int            `json:"-" gorm:"column:user_id;"`
	Title                 string         `json:"title" gorm:"column:title;"`
	Content               string         `json:"content" gorm:"column:content;"`
	Cover                 *common.Image  `json:"cover" gorm:"column:cover;"`
	Photos                *common.Images `json:"photos" gorm:"column:photos;"`
	CoverImgId            int            `json:"cover_img_id" gorm:"-"`
}

func (NoteCreate) TableName() string {
	return Note{}.TableName()
}

func (n *NoteCreate) Mask(isAdmin bool) {
	n.GenUID(common.DBTypeNote, 1)
}

func (n *NoteCreate) GetImageIds() []int {
	return []int{n.CoverImgId}
}
