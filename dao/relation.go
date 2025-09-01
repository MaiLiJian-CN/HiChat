package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"

	"go.uber.org/zap"
)

// FirendList
func FriendList(userId uint) (*[]models.UserBasic, error) {
	// 查询当前用户的好友关系
	relation := make([]models.Relation, 0)
	tx := global.DB.Where("owner_id=? and type=1", userId).Find(&relation)
	if tx.RowsAffected == 0 {
		zap.S().Info("Not Found Relation data")
		return nil, errors.New("Not Found Firend relation")
	}

	// 提取好友ID
	userID := make([]uint, 0)
	for _, v := range relation {
		userID = append(userID, v.TargetID)
	}
	// 查询好友信息
	user := make([]models.UserBasic, 0)
	tx = global.DB.Where("id in ?", userID).Find(&user)
	if tx.RowsAffected == 0 {
		zap.S().Info("Not Found Relation friend relation")
		return nil, errors.New("Not Found Friend")
	}
	return &user, nil
}

// AddFriend 加好友
func AddFriend(userID, targetId uint) (int, error) {

	if userID == targetId {
		return -2, errors.New("userID和TargetId相等")
	}
	//通过id查询用户
	targetUser, err := FindUserByID(targetId)
	if err != nil {
		return -1, errors.New("Not Found User")
	}
	if targetUser.ID == 0 {
		zap.S().Info("Not Found User")
		return -1, errors.New("Not Found User")
	}

	relation := models.Relation{}
	tx := global.DB.Where("owner_id=? and target_id=? and type=1", userID, targetId).First(&relation)
	if tx.RowsAffected == 1 {
		zap.S().Info("该好友存在")
		return 0, errors.New("好友已经存在")
	}
	if tx := global.DB.Where("owner_id = ? and target_id = ?  and type = 1", targetId, userID).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("该好友存在")
		return 0, errors.New("好友已经存在")
	}

	//start
	tx = global.DB.Begin()

	relation.OwnerId = userID
	relation.TargetID = targetUser.ID
	relation.Type = 1

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("Create Error")
		//事务回滚
		tx.Rollback()
		return -1, errors.New("创建好友记录失败")
	}

	relation = models.Relation{}
	relation.OwnerId = targetUser.ID
	relation.TargetID = userID
	relation.Type = 1

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("Create Error")
		//事务回滚
		tx.Rollback()
		return -1, errors.New("创建好友记录失败")
	}

	//submit
	tx.Commit()
	return 1, nil
}

//AddFriendByName 昵称加好友
func AddFriendByName(userId uint,targetName string)(int,error){
	// user:=models.UserBasic{}

	targetUser,err:=FIndUserByName(targetName)
	if err!=nil{
		return -1,errors.New("该用户不存在")
	}
	if targetUser.ID==0{
		zap.S().Info("Not Found User")
		return -1,errors.New("该用户不存在")
	}
	return AddFriend(userId,targetUser.ID)
}