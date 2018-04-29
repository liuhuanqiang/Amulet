package db


import (
	"database/sql"
	"fmt"
	"github.com/golang/glog"
	_ "github.com/go-sql-driver/mysql"

	"amulet/config"
)
type MysqlDB struct {
	DB *sql.DB
}

type Blog struct {
	Id  			int  		`json:"id"`
	Fid 			int	 	`json:"fid"`
	Title 			string  	`json:"title"`
	Description 		string 		`json:"description"`
	Link  			string 		`json:"link"`
	LinkId 			string 		`json:"link_id"`
	PubDate			int		`json:"pub_date"`
	CreateTime  		string 		`json:"create_time"`
	Source 			int 		`json:"source"`
	Tag   			string 		`json:"tag"`
}

type Zhihu struct {
	Id              	int     	`json:"id"`
	Fid      	        int 		`json:"fid"`
	Title                   string		`json:"title"`
	Description		string		`json:"description"`
	Link  			string		`json:"link"`
	LinkId			string 		`json:"linkid"`
	Verb                    string		`json:"verb"`
	PubDate			int		`json:"pub_date"`
	CreateTime		int		`json:"create_time"`
	Source   		int		`json:"source"`
	Tag                     string		`json:"tag"`
}

type Jianshu struct {
	Id  			int  		`json:"id"`
	Fid 			int	 	`json:"fid"`
	Title 			string  	`json:"title"`
	Description 		string 		`json:"description"`
	Link  			string 		`json:"link"`
	LinkId 			string 		`json:"link_id"`
	PubDate			int		`json:"pub_date"`
	CreateTime  		string 		`json:"create_time"`
	Source 			int 		`json:"source"`
	Tag   			string 		`json:"tag"`
}


var db *MysqlDB

const (
	Table_ZhiHu = "tb_zhihu"
	Table_JianShu = "tb_jianshu"
	Table_Blog = "tb_blog"
	Table_JueJing = "tb_JueJing"
)
func init() {
	db = &MysqlDB{}
}

func GetDB() *MysqlDB {
	return db
}

func (this *MysqlDB)Init(){
	path:=fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s?charset=utf8mb4",config.Config().UserName,config.Config().Password,config.Config().Port,config.Config().Table)
	glog.Info("wait init db...",path)
	defer glog.Info("init db ok!")
	db,err := sql.Open("mysql",path)
	if err!=nil{
		glog.Info("db init error! ",err.Error())
	}else{
		this.DB = db
	}

}


func (this *MysqlDB) GetBlog() []Blog {
	sql := "SELECT * FROM `tb_blog`"
	rows,err := this.DB.Query(sql)
	defer rows.Close()
	if err!=nil{
		panic(err.Error())
	}
	if rows.Err()!=nil{
		panic(rows.Err().Error())
	}
	ret := []Blog{}

	for rows.Next() {
		item := Blog{}
		err = rows.Scan(&item.Id,&item.Fid,&item.Title,&item.Description,&item.Link,&item.LinkId,&item.PubDate,&item.CreateTime,&item.Source,&item.Tag)
		ret = append(ret,item)
	}
	return ret
}

func (this *MysqlDB) GetZhiHu() []Zhihu {
	sql := "SELECT * FROM `tb_zhihu` WHERE `verb` = 'MEMBER_CREATE_ARTICLE'"
	rows,err := this.DB.Query(sql)
	defer rows.Close()
	if err!=nil{
		panic(err.Error())
	}
	if rows.Err()!=nil{
		panic(rows.Err().Error())
	}
	ret := []Zhihu{}

	for rows.Next() {
		item := Zhihu{}
		err = rows.Scan(&item.Id,&item.Fid,&item.Title,&item.Description,&item.Link,&item.LinkId,&item.Verb,&item.PubDate,&item.CreateTime,&item.Source,&item.Tag)
		ret = append(ret,item)
	}
	return ret
}

func (this *MysqlDB) GetJianShu() []Jianshu {
	sql := "SELECT * FROM `tb_jianshu`"
	rows,err := this.DB.Query(sql)
	defer rows.Close()
	if err!=nil{
		panic(err.Error())
	}
	if rows.Err()!=nil{
		panic(rows.Err().Error())
	}
	ret := []Jianshu{}

	for rows.Next() {
		item := Jianshu{}
		err = rows.Scan(&item.Id,&item.Fid,&item.Title,&item.Description,&item.Link,&item.LinkId,&item.PubDate,&item.CreateTime,&item.Source,&item.Tag)
		ret = append(ret,item)
	}
	return ret
}
