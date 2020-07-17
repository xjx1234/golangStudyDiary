/**
* @Author: XJX
* @Description: Mysql数据库操作示例代码
* @File: mysql.go
* @Date: 2020/7/3 14:41
 */

package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"strings"
	"time"
)

//数据库配置结构体
type DbConfig struct {
	Host     string
	Network  string
	Port     string
	Username string
	Password string
	Charset  string
	Database string
}

//读取配置ENV文件信息，此处可以用 github.com/timest/env 替代
func getEnvConfig(fileName string) (dbconfig *DbConfig, err error) {
	config := &DbConfig{}
	dir, _ := os.Getwd()
	fileName = dir + "\\" + fileName
	if !checkConfigFileIsExist(fileName) {
		return config, errors.New(fileName + " is Not Exist!!!")
	} else {
		f, err := os.Open(fileName)
		if err != nil {
			return config, err
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		vConfigData := reflect.ValueOf(config)
		for {
			linestr, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			linestr = strings.TrimSpace(linestr)
			if linestr == "" {
				continue
			}
			pair := strings.Split(linestr, "=")
			key := strings.TrimSpace(pair[0])
			value := strings.TrimSpace(pair[1])
			vConfigData.Elem().FieldByName(key).SetString(value)
		}
		return config, nil
	}
}

//文件是否存在
func checkConfigFileIsExist(fileName string) bool {
	flag := true
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		flag = false
	}
	return flag
}

type Student struct {
	Id              int
	Name            string
	Age             int
	addtime, uptime int64
}

type SqlType struct {
	DB *sql.DB
	Tx *sql.Tx
}

//组合条件
func assemblyCondition(condition [][]string) (sql string) {
	sql = " WHERE 1 "
	if len(condition) > 0 {
		for _, v := range condition {
			if len(v) > 1 {
				for _, v1 := range v {
					sql += " " + v1
				}
			}
		}
	}
	return
}

//数据库操作接口
type MyStudent interface {
	insertOneData(dd *SqlType, table string, kv map[string]interface{}) (inserId int64)
	queryOne(dd *SqlType, table string, condition [][]string) (row Student)
	queryAll(dd *SqlType, table string, condition [][]string) (rows []map[string]interface{})
	deleteAll(dd *SqlType, table string, condition [][]string) (num int64)
	updateAll(dd *SqlType, table string, setSql string, condition [][]string) (affNums int64)
}


//插入数据
func (s *Student) insertOneData(dd *SqlType, table string, kv map[string]interface{}) (inserId int64) {
	inserId = 0
	field := ""
	values := ""
	if len(kv) > 0 {
		startFlag := true
		for k, v := range kv {
			if startFlag {
				values = values + "'" + fmt.Sprintf("%v", v) + "'"
				field = field + k
				startFlag = false
			} else {
				values = values + ",'" + fmt.Sprintf("%v", v) + "'"
				field += "," + k
			}
		}
		field += ",addtime,uptime"
		t := time.Now().Unix()
		values += ",'" + fmt.Sprintf("%v", t) + "','" + fmt.Sprintf("%v", t) + "'"
		sqlString := "INSERT INTO " + table + "(" + field + ")VALUES(" + values + ")"
		var ret sql.Result
		var err error
		if dd.DB != nil {
			ret, err = dd.DB.Exec(sqlString)
			fmt.Println("DB Exec")
		} else {
			ret, err = dd.Tx.Exec(sqlString)
			fmt.Println("TX Exec")
			fmt.Println(dd.Tx)
		}
		if err == nil {
			iId, _ := ret.LastInsertId()
			inserId = iId
		}
	}
	return
}

func (s *Student) updateAll(dd *SqlType, table string, setSql string, condition [][]string) (affNums int64) {
	sqlString := assemblyCondition(condition)
	var ret sql.Result
	if dd.DB != nil {
		ret, _ = dd.DB.Exec("UPDATE " + table + " SET " + setSql + sqlString)
	} else {
		ret, _ = dd.Tx.Exec("UPDATE " + table + " SET " + setSql + sqlString)
	}
	updateNums, _ := ret.RowsAffected()
	affNums = updateNums
	return
}

func (s *Student) deleteAll(dd *SqlType, table string, condition [][]string) (num int64) {
	sqlString := assemblyCondition(condition)
	var ret sql.Result
	if dd.DB != nil {
		ret, _ = dd.DB.Exec("DELETE  FROM " + table + sqlString)
	} else {
		ret, _ = dd.Tx.Exec("DELETE  FROM " + table + sqlString)
	}
	delNums, _ := ret.RowsAffected()
	num = delNums
	return
}

//查询单笔记录
func (s *Student) queryOne(dd *SqlType, table string, condition [][]string) (row Student) {
	sqlString := assemblyCondition(condition)
	if dd.DB != nil {
		dd.DB.QueryRow("SELECT * FROM "+table+sqlString).Scan(&s.Id, &s.Name, &s.Age, &s.addtime, &s.uptime)
	} else {
		dd.Tx.QueryRow("SELECT * FROM "+table+sqlString).Scan(&s.Id, &s.Name, &s.Age, &s.addtime, &s.uptime)
	}
	return *s
}

//查询多笔记录
func (s *Student) queryAll(dd *SqlType, table string, condition [][]string) (rows []map[string]interface{}) {
	var rowsList *sql.Rows
	var err error
	sqlString := assemblyCondition(condition)
	if dd.DB != nil {
		rowsList, err = dd.DB.Query("SELECT * FROM " + table + sqlString)
	} else {
		rowsList, err = dd.Tx.Query("SELECT * FROM " + table + sqlString)
	}
	defer rowsList.Close()
	if err == nil {
		for rowsList.Next() {
			rowsList.Scan(&s.Id, &s.Name, &s.Age, &s.addtime, &s.uptime)
			rows = append(rows, map[string]interface{}{
				"id":      s.Id,
				"name":    s.Name,
				"age":     s.Age,
				"addtime": s.addtime,
				"uptime":  s.uptime,
			})
		}
	}
	return
}

var DB *sql.DB

func main() {
	config, err := getEnvConfig("mysql.env")
	if err != nil {
		fmt.Println(err)
	}
	dns := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s", config.Username, config.Password, config.Network, config.Host, config.Port, config.Database, config.Charset)
	db, err := sql.Open("mysql", dns)
	DB = db
	var myStudent MyStudent
	myStudent = &Student{}
	/*doDbType := &SqlType{DB: DB}
	insertStudentData := map[string]interface{}{"name": "qq", "age": 11}
	id := myStudent.insertOneData(doDbType, "student", insertStudentData)
	fmt.Println("Insert data succ!! id:" + fmt.Sprintf("%v", id))
	/*queryOneCondition := [][]string{{"AND", "id>1"}, {"AND", "name='qq'"}}
	oneData := myStudent.queryOne(doDbType, "student", queryOneCondition)
	fmt.Println(oneData.Id)
	queryCondition := [][]string{{"AND", "name LIKE '%xjx%'"}}
	datas := myStudent.queryAll(doDbType, "student", queryCondition)
	fmt.Println(datas)
	deleteCondition := [][]string{{"AND", "name='xjx3'"}}
	num := myStudent.deleteAll(doDbType, "student", deleteCondition)
	fmt.Println(num)
	upCondition := [][]string{{"AND", "name='xjx4'"}}
	num := myStudent.updateAll(doDbType, "student", " age=36,name='zzzz'", upCondition)
	fmt.Println(num)*/

	tx, err := DB.Begin()
	defer tx.Rollback()
	fmt.Println(err)
	doTypeTx := &SqlType{Tx: tx, DB: nil}
	insertStudentData := map[string]interface{}{"name": "fffff123", "age": 22}
	myStudent.insertOneData(doTypeTx, "student", insertStudentData)
	fmt.Println(doTypeTx.Tx)
	rr := doTypeTx.Tx.Commit()

	fmt.Println(rr)

}
