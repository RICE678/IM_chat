package contact

import (
	"IM_chat/dao/sql"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"go.uber.org/zap"
)

func SearchContactList(userID int64) ([]models.ListContact, string) {
	var users []sql.User
	var err error
	if users, err = sql.SearchContact(userID); err != nil {
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

func SearchFriendAll(friendID, userID int64) (friend models.ContactMain, err1 string) {
	if friendID <= 0 || userID <= 0 {
		err1 = errcode.Msg(errcode.InvalidParams)
		return
	}
	var err error
	if ok := sql.IsFriend(userID, friendID); !ok {
		err1 = errcode.Msg(errcode.DontFriend)
		return
	}
	friend.FriendID = friendID
	if friend.FriendPicture, err = sql.SearchPicture(friendID); err != nil {
		zap.L().Error("search friend picture failed", zap.Int64("friend_id", friendID), zap.Error(err))
		err1 = errcode.Msg(errcode.ERROR)
		return
	}
	u, err := sql.SearchUserMain(friendID)
	if err != nil {
		zap.L().Error("search friend main failed", zap.Int64("friend_id", friendID), zap.Error(err))
		err1 = errcode.Msg(errcode.ERROR)
		return
	}
	friend.FriendName = u.Name
	friend.FriendEmail = u.Email
	friend.Signature = u.Signature
	friend.Gender = u.Gender
	friend.Remark, err = sql.SearchRemark(userID, friendID)
	if err != nil {
		err1 = errcode.Msg(errcode.ERROR)
		return
	}
	err1 = errcode.Msg(errcode.SUCCESS)
	return
}

func RemarkChange(userID, friendID int64, remark string) (err string) {
	if friendID <= 0 || userID <= 0 {
		err = errcode.Msg(errcode.InvalidParams)
		return
	}
	if ok := sql.IsFriend(userID, friendID); !ok {
		err = errcode.Msg(errcode.DontFriend)
		return
	}
	if err1 := sql.ChangeRemark(userID, friendID, remark); err1 != nil {
		err = errcode.Msg(errcode.ERROR)
		return
	}
	err = errcode.Msg(errcode.SUCCESS)
	return
}
