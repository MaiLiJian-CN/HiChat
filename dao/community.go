package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"
)

//create new group

func CreateCommunity(community models.Community) (int, error) {
	com := models.Community{}

	//if group haven
	if tx := global.DB.Where("name = ? ", community.Name).First(&com); tx.RowsAffected == 1 {
		return -1, errors.New("当前群记录已存在")
	}
	tx := global.DB.Begin()
	if t := tx.Create(&community); t.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("群记录创建失败")
	}
	relation := models.Relation{}
	relation.OwnerId = community.OwnerId //group owner Id
	relation.TargetID = community.ID     //group ID
	relation.Type = 2                    //group
	if t := tx.Create(&relation); t.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("群记录创建失败")
	}
	tx.Commit()
	return 0, nil
}

// GetCommunityList 获取群列表
func GetCommunityList(ownerId uint) (*[]models.Community, error) {

	//获取我加入的群
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type = 2", ownerId).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("不存在群记录")
	}

	communityIDs := make([]uint, 0)
	for _, v := range relation {
		communityIDs = append(communityIDs, v.TargetID)
	}
	community := make([]models.Community, 0)
	if tx := global.DB.Where("id in ?", communityIDs).Find(&community); tx.RowsAffected == 0 {
		return nil, errors.New("获取群数据失败")
	}
	return &community, nil
}

// JoinCommunity 根据群昵称搜索并加入群
func JoinCommunity(ownerId uint, cname string) (int, error) {
	community := models.Community{}
	if tx := global.DB.Where("name = ? ", cname).First(&community); tx.RowsAffected == 0 {
		return -1, errors.New("群记录不存在")
	}

	//重复加群
	relation := models.Relation{}
	if tx := global.DB.Where("owner_id = ? and target_id = ? and type = 2", ownerId, community.ID).First(&relation); tx.RowsAffected == 1 {
		return -1, errors.New("该群已经加入")
	}

	relation = models.Relation{}
	relation.OwnerId = ownerId
	relation.Type = 2
	relation.TargetID = community.ID

	if tx := global.DB.Create(&relation); tx.RowsAffected == 0 {
		return -1, errors.New("加入失败")
	}
	return 0, nil

}
