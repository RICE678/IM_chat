package ws

import (
	"IM_chat/dao/redisdao"
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

func DelFriendMain(user models.DelResponse) string {
	if err := sql.DeleteContact(user.UserID, user.FriendID); err != nil {
		return errcode.Msg(errcode.ErrDelFriend)
	}
	if err := sql.DelUnReadMessage(user.UserID, user.FriendID); err != nil {
		return errcode.Msg(errcode.ErrDelFriend)
	}
	if err := redisdao.DelUnreadCount(user.FriendID, user.UserID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err := redisdao.InitUnreadCount(user.FriendID, user.UserID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}
