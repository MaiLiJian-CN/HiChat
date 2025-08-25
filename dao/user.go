package dao

import (
	"HiChat/common"
	"HiChat/global"
	"HiChat/models"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func GetUserList()([]*models.UserBasic,error){
	var list []*models.UserBasic
	if tx:=global.DB.Find(&list);tx.RowsAffected==0{
		return nil,errors.New("Get User List Error")
	}
	return list,nil
}

func FindNuserByNameAndPwd(name string,password string) (*models.UserBasic,error){
	user:=models.UserBasic{}
	
	if tx:=global.DB.Where("name=? and pass_word=?",name,password).First(&user);tx.RowsAffected==0{
		return nil,errors.New("Not Found this message")
	}

	t:=strconv.Itoa(int(time.Now().Unix()))

	temp := common.Md5encoder(t)

	if tx:=global.DB.Model(&user).Where("id=?",user.ID).Update("identify",temp);tx.RowsAffected==0{
		return  nil,errors.New("write identity error")
	}
	return &user,nil
}


//login
func FIndUserByName(name string) (*models.UserBasic,error){
	user:=models.UserBasic{}
	if tx:=global.DB.Where("name=?",name).First(&user);tx.RowsAffected==0{
		return nil,errors.New("Not Found the User")
	}
	return  &user,nil
}

//register
func FindUser(name string)(*models.UserBasic,error){
	user:=models.UserBasic{}
	if tx:=global.DB.Where("name=?",name).First(&user);tx.RowsAffected==1{
		return nil,errors.New("User have registered")
	}
	return &user,nil

}

//select user by Id
func FindUserByID(ID uint) (*models.UserBasic,error){
	user:=models.UserBasic{}
	tx:=global.DB.Where(ID).First(&user);
	if tx.RowsAffected==0{
		return  nil,errors.New("Not Found User By Id")
	}
	return &user,nil
}

func FindUserByPhone(phone string)(*models.UserBasic,error){
	user:=models.UserBasic{}
	tx:=global.DB.Where("phone=?",phone).First(&user)
	if tx.RowsAffected==0 {
		return  nil,errors.New("Not Found User By Phone")
	}

	return &user,nil
}


func FindUserByEmail(email string) (*models.UserBasic,error){
	user:=models.UserBasic{}
	tx:=global.DB.Where("email=?",email).First(&user)
	if tx.RowsAffected==0{
		return nil,errors.New("Not Found User By Email")
	}
	return &user,nil
}

func CreateUser(user models.UserBasic) (*models.UserBasic,error){
	tx:=global.DB.Create(&user)
	if tx.RowsAffected==0{
		zap.S().Info("Create User Error")
		return nil,errors.New("Create User Error")
	}
	return &user,nil
}

func UpdateUser(user models.UserBasic) (*models.UserBasic, error) {
    tx := global.DB.Model(&user).Updates(models.UserBasic{
        Name:     user.Name,
        PassWord: user.PassWord,
        Gender:   user.Gender,
        Phone:    user.Phone,
        Email:    user.Email,
        Avatar:   user.Avatar,
        Salt:     user.Salt,
    })
    if tx.RowsAffected == 0 {
        zap.S().Info("更新用户失败")
        return nil, errors.New("更新用户失败")
    }
    return &user, nil
}
func DeleteUser(user models.UserBasic) error {
    if tx := global.DB.Delete(&user); tx.RowsAffected == 0 {
        zap.S().Info("删除失败")
        return errors.New("删除用户失败")
    }
    return nil
}