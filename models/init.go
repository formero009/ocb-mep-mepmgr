package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"mepmgr/models/certmd"
	"mepmgr/models/mepmd"
)

var PostgresDB *gorm.DB

func init() {
	err := ensurePqDatabase()
	if err != nil {
		panic(err)
	}
}

func ensurePqDatabase() error {
	dbName := beego.AppConfig.String("DBName")
	dbUser := beego.AppConfig.String("DBUser")
	dbPwd := beego.AppConfig.String("DBPasswd")
	dbHost := beego.AppConfig.String("DBHost")
	dbPort := beego.AppConfig.String("DBPort")
	dbTimeZone := beego.AppConfig.String("DBTimeZone")
	dataSourceName := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s timezone=%s sslmode=disable",
		dbUser, dbPwd, dbName, dbHost, dbPort, dbTimeZone)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		PrepareStmt:    true,
	})
	if err != nil {
		return err
	}

	fmt.Println("Initialize database connection:",
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s timezone=%s sslmode=disable",
			dbUser, "***", dbName, dbHost, dbPort, dbTimeZone))

	//autoMigrate apigw app models
	err = db.AutoMigrate(&mepmd.MepMeta{})
	if err != nil {
		panic("failed to autoMigrate table of mep_meta")
	}
	err = db.AutoMigrate(&mepmd.MepGroup{})
	if err != nil {
		panic("failed to autoMigrate table of mep_group")
	}
	err = db.AutoMigrate(&mepmd.MepGroupRelation{})
	if err != nil {
		panic("failed to autoMigrate table of mep_group_relation")
	}
	err = db.AutoMigrate(&mepmd.AlertInfo{})
	if err != nil {
		panic("failed to autoMigrate table of alert_info")
	}
	err = db.AutoMigrate(&mepmd.Preference{})
	if err != nil {
		panic("failed to autoMigrate table of preference")
	}
	err = db.AutoMigrate(&mepmd.ConfigParameter{})
	if err != nil {
		panic("failed to autoMigrate table of config_param")
	}

	err = db.AutoMigrate(&mepmd.MepProcessLog{})
	if err != nil {
		panic("failed to autoMigrate table of mep_process_log")
	}

	err = db.AutoMigrate(&certmd.Cert{})
	if err != nil {
		panic("failed to autoMigrate table of cert")
	}

	showSql := beego.AppConfig.DefaultBool("ShowSql", false)

	if showSql {
		PostgresDB = db.Debug()
	} else {
		PostgresDB = db
	}
	return nil
}
