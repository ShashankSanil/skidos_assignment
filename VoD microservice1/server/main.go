package main

import (
	_ "VoD_microservice1/chassisHandlers"
	"VoD_microservice1/common"
	"VoD_microservice1/database"
	authRepository "VoD_microservice1/repository"
	"VoD_microservice1/resource"
	service "VoD_microservice1/service"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2"
	"github.com/go-chassis/openlog"
)

func getServices(dbname string) *service.Service {
	Repo := authRepository.Repository{DBClient: database.GetClient(), DBName: dbname}
	return &service.Service{Rep: &Repo}
}

func duplicateInArray(arr []interface{}) []string {
	visited := make(map[string]bool)
	result := make(map[string]bool)
	for _, obj := range arr {
		errdetails := obj.(map[string]interface{})
		errcode := errdetails["errorcode"].(string)
		if visited[errcode] {
			result[errcode] = true
		} else {
			visited[errcode] = true
		}
	}
	res := []string{}
	for k := range result {
		res = append(res, k)
	}
	return res
}

func CheckandLoadErrors() error {
	errbytes, err := ioutil.ReadFile("./conf/errorcode.json")
	if err != nil {
		openlog.Error(err.Error())
		return err
	}
	errorslist := make([]interface{}, 0)

	err = json.Unmarshal(errbytes, &errorslist)
	if err != nil {
		openlog.Error(err.Error())
		return err
	}

	dups := duplicateInArray(errorslist)
	if len(dups) > 0 {
		fmt.Println("Duplicates exists in Errorcode: ", dups)
		return errors.New(" Duplicates exists in Errorcode")
	}

	for _, errs := range errorslist {
		ERROR := errs.(map[string]interface{})
		errcode := ERROR["errorcode"].(string)
		common.Errorcodes[errcode] = ERROR
	}

	return nil
}

func main() {
	tempresource := resource.Resources{ServiceProvider: getServices}
	chassis.RegisterSchema("rest", &tempresource)
	if err := chassis.Init(); err != nil {
		openlog.Fatal("Init failed." + err.Error())
		return
	}
	if err := archaius.AddFile("./conf/database.yaml"); err != nil {
		openlog.Error("add props configurations failed." + err.Error())
		return
	}

	if err := archaius.AddFile("./conf/watcher.yaml"); err != nil {
		openlog.Error("add props configurations failed." + err.Error())
		return
	}

	if err := archaius.AddFile("./conf/payloadSchemas.yaml"); err != nil {
		openlog.Error("add props configurations failed." + err.Error())
		return
	}

	common.DbName = archaius.GetString("database.mongodb.dbname", "")
	common.SECRET_KEY = archaius.Get("secretkey").(string)

	if err := database.Connect(); err != nil {
		openlog.Fatal("Error occured while connecting to database")
		return
	}

	err := CheckandLoadErrors()
	if err != nil {
		openlog.Fatal(err.Error())
		return
	}

	chassis.Run()
}
