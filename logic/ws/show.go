package ws

import (
	"IM_chat/dao/sql"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
)

func SearchShowList(userID int64) ([]models.ListContact, string) {
	var users []sql.User
	var err error
	if users, err = sql.SearchShow(userID); err != nil {
		return nil, errcode.Msg(errcode.ERROR)
	}
	if len(users) == 0 {
		return nil, errcode.Msg(errcode.NotFriend)
	}
	var contactList []models.ListContact
	for _, user := range users {
		pic, _ := sql.SearchPicture(user.ID)
		contactList = append(contactList, models.ListContact{
			UserID:        userID,
			FriendID:      user.ID,
			FriendName:    user.Name,
			FriendPicture: pic,
			Remark:        user.Remark,
		})
	}
	return contactList, errcode.Msg(errcode.SUCCESS)
}
