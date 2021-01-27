package mysql

import (
	"database/sql"
	"fmt"
	"usermana/config"
	"usermana/utils"

	_ "github.com/go-sql-driver/mysql" // mysqldrive
)

var (
	createAccountSt    *sql.Stmt
	loginAuthSt        *sql.Stmt
	createProfileSt    *sql.Stmt
	getProfileSt       *sql.Stmt
	updateProfileSt    *sql.Stmt
	updateNickNameSt   *sql.Stmt
	updateProfilePicSt *sql.Stmt
)

//init,  mysql的初始化函数.
func init() {
	//连接数据库
	db, err := sql.Open("mysql", config.MysqlDB)
	if err != nil {
		panic(err)
	}
	//配置数据库连接的限制.
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)

	//测试是否连接成功.
	if err = db.Ping(); err != nil {
		panic(err)
	}

	//预处理mysql语句
	createAccountSt = dbPrepare(db, "INSERT INTO tbl_login_info (user_name, password) values (?, ?)")
	createProfileSt = dbPrepare(db, "INSERT INTO tbl_user_info (user_name, nick_name) values (?, ?)")
	loginAuthSt = dbPrepare(db, "SELECT password FROM tbl_login_info WHERE user_name = ?")
	getProfileSt = dbPrepare(db, "SELECT nick_name, pic_name FROM tbl_user_info WHERE user_name = ?")
	updateProfileSt = dbPrepare(db, "UPDATE tbl_user_info SET nick_name = ?, pic_name = ? where user_name = ?")
	updateNickNameSt = dbPrepare(db, "UPDATE tbl_user_info SET nick_name = ? where user_name = ?")
	updateProfilePicSt = dbPrepare(db, "UPDATE tbl_user_info SET pic_name = ? where user_name = ?")

	fmt.Println("mysql init done.")
}

//dbPrepare 预处理sql语句.
func dbPrepare(db *sql.DB, query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	return stmt
}

// CreateAccount 创建账号.
func CreateAccount(userName string, password string) error {
	//先对密码进行sha256的编码再保存到数据库.
	pwd := utils.Sha256(password)
	_, err := createAccountSt.Exec(userName, pwd)
	if err != nil {
		return err
	}
	return nil
}

// CheckAccountExist  判断账号是否存在.
func CheckAccountExist(userName string) (bool, error) {
	rows, err := loginAuthSt.Query(userName)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

// LoginAuth 登录校验.
func LoginAuth(userName string, password string) (bool, error) {
	var pwd string
	//t := time.Now()
	rows, err := loginAuthSt.Query(userName)
	if err != nil {
		return false, err
	}
	//连接归还到连接池中
	defer rows.Close()
	//从数据库中过去用户密码.
	for rows.Next() {
		err = rows.Scan(&pwd)
	}

	if err != nil {
		return false, err
	}
	//进行校验.
	if pwd == utils.Sha256(password) {
		//log.Infof("%q", time.Since(t))
		return true, nil
	}
	return false, nil
}

// CreateProfile 创建用户信息.
func CreateProfile(userName string, nickName string) error {
	_, err := createProfileSt.Exec(userName, nickName)
	if err != nil {
		return err
	}
	return nil
}

// GetProfile 获取用户信息.
func GetProfile(userName string) (nickName string, picName string, hasData bool, err error) {
	rows, err := getProfileSt.Query(userName)
	if err != nil {
		return nickName, picName, hasData, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&nickName, &picName)
	}

	if err != nil {
		return nickName, picName, hasData, err
	}
	if nickName != "" {
		hasData = true
	}
	return nickName, picName, hasData, nil
}

// CheckProfileExist 判断用户信息是否存在.
func CheckProfileExist(userName string) (bool, error) {
	rows, err := getProfileSt.Query(userName)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

// UpdateProfile 更新用户信息.
func UpdateProfile(userName string, nickName string, picName string) (bool, error) {
	res, err := updateProfileSt.Exec(nickName, picName, userName)
	if err != nil {
		return false, err
	}
	if afrows, _ := res.RowsAffected(); afrows > 0 {
		return true, nil
	}
	return CheckProfileExist(userName)
}

// UpdateNikcName 更新用户昵称.
func UpdateNikcName(userName string, nickName string) (bool, error) {
	_, err := updateNickNameSt.Exec(nickName, userName)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateProfilePic 更新用户头像.
func UpdateProfilePic(userName string, picName string) (bool, error) {
	_, err := updateProfilePicSt.Exec(picName, userName)
	if err != nil {
		return false, err
	}
	return true, nil
}
