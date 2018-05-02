package db


import (
	"database/sql"
	"fmt"
	"github.com/golang/glog"
	_ "github.com/go-sql-driver/mysql"

	"amulet/web/config"
	"amulet/web/msg"
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

func (this *MysqlDB) GetLatestList(page int) []*msg.LatestResp {
	tpl := "select fid,name, avater, title, description, linkid, source, pub_date from (" +
			"(" +
				"select tb_blog.fid as fid, tb_blog.title as title, tb_blog.description as description, tb_blog.linkid as linkid, tb_blog.source as source, tb_blog.pub_date as pub_date from tb_blog" +
			")" +
			"union all (" +
				"select tb_zhihu.fid as fid, tb_zhihu.title as title, tb_zhihu.description as description, tb_zhihu.linkid as linkid, tb_zhihu.source as source, tb_zhihu.pub_date as pub_date from tb_zhihu where verb = 'MEMBER_CREATE_ARTICLE'" +
			") " +
			"union all (" +
				"select tb_jianshu.fid as fid, tb_jianshu.title as title, tb_jianshu.description as description, tb_jianshu.linkid as linkid, tb_jianshu.source as source, tb_jianshu.pub_date as pub_date from tb_jianshu" +
			")" +
		") as tb left join tb_famous on tb.fid = tb_famous.id order by `pub_date` desc limit %d,%d"

	sql := fmt.Sprintf(tpl, (page - 1)*20, page*20)
	rows,err := this.DB.Query(sql)
	defer rows.Close()
	if err!=nil{
		panic(err.Error())
	}
	if rows.Err()!=nil{
		panic(rows.Err().Error())
	}
	ret := []*msg.LatestResp{}

	for rows.Next() {
		item := msg.LatestResp{}
		err = rows.Scan(&item.Fid,&item.Name,&item.Avater,&item.Title,&item.Description,&item.Linkid,&item.Source,&item.PubDate)
		ret = append(ret,&item)
	}
	return ret
}


func (this *MysqlDB) GetLink(tb string,linkid string) (int,string) {
	tpl := "select fid,link from `%s` where `linkid`='%s'"
	sql := fmt.Sprintf(tpl, tb, linkid)
	glog.Info("sql:",sql)
	rows, err := this.DB.Query(sql)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	if rows.Err() != nil {
		panic(rows.Err().Error())
	}
	var link = ""
	var fid = 0
	for rows.Next() {
		rows.Scan(&fid,&link)
	}
	return fid,link
}