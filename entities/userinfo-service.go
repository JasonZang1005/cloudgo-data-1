package entities

//UserInfoAtomicService .
type UserInfoAtomicService struct{}

//UserInfoService .
var UserInfoService = UserInfoAtomicService{}

var userInfoInsertStmt = "INSERT userinfo SET username=?,departname=?,created=?"
var userInfoQueryAll = "SELECT * FROM userinfo"
var userInfoQueryByID = "SELECT * FROM userinfo where uid = ?"

// Save .
func (*UserInfoAtomicService) Save(u *UserInfo) error {
	_, err := engine.Insert(&u)
	checkErr(err)
	return err
}

// FindAll .
func (*UserInfoAtomicService) FindAll() []UserInfo {
	allUserInfo := make([]UserInfo, 0)
	err := engine.Find(&allUserInfo)
	checkErr(err)

	return allUserInfo
}

// FindByID .
func (*UserInfoAtomicService) FindByID(id int) *UserInfo {
	user := &UserInfo{UID: id}
	has, err := engine.Get(user)
	checkErr(err)
	if has {
		return user
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
