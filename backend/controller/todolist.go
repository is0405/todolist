package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"regexp"
	"errors"
	// "fmt"
	"strconv"

	"github.com/is0405/model"
	"github.com/is0405/repository"
	"github.com/is0405/service"
	"github.com/is0405/httputil"
	"github.com/jmoiron/sqlx"
)

type ToDO struct {
	db *sqlx.DB
}

func NewToDO(db *sqlx.DB) *ToDO {
	return &ToDO{db: db}
}

func (a *ToDO) Create(_ http.ResponseWriter, r *http.Request) (int, interface{}, error) {

	// accessTokenから取り出した情報を取得する
	getc, err := httputil.GetClaimsFromContext(r.Context())
	if err != nil {
		//fmt.Println("accessToken")
		return http.StatusInternalServerError, nil, err
	}
	
	model := &model.ToDO{}
	
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		return http.StatusBadRequest, nil, err
	}

	//Title
	if model.Title == "" {
		return http.StatusUnprocessableEntity, nil, errors.New("Not enough parameters")
	}

	if !StringCheck(model.Title) {
		return http.StatusBadRequest, nil, errors.New("Not string")
	}

	model.OK = false

	model.AccountID = getc.ID

	Service := service.NewToDO(a.db)
	_, err = Service.Create(model)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	
	return http.StatusOK, nil, nil
}

func (a *ToDO) Get(_ http.ResponseWriter, r *http.Request) (int, interface{}, error) {

	// accessTokenから取り出した情報を取得する
	getc, err := httputil.GetClaimsFromContext(r.Context())
	if err != nil {
		//fmt.Println("accessToken")
		return http.StatusInternalServerError, nil, err
	}
	
	res, err := repository.GetToDOList(a.db, getc.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, res, nil
}

func (a *ToDO) Update(_ http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// accessTokenから取り出した情報を取得する
	getc, err := httputil.GetClaimsFromContext(r.Context())
	if err != nil {
		//fmt.Println("accessToken")
		return http.StatusInternalServerError, nil, err
	}

	list_id, err := URLToID(r)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	
	model := &model.ToDO{}
	
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		return http.StatusBadRequest, nil, err
	}

	model.ID = list_id
	//Title
	if model.Title == "" {
		return http.StatusUnprocessableEntity, nil, errors.New("Not enough parameters")
	}

	if !StringCheck(model.Title) {
		return http.StatusBadRequest, nil, errors.New("Not string")
	}

	Service := service.NewToDO(a.db)
	_, err = Service.Update(model)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	res, err := repository.GetToDOList(a.db, getc.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	
	return http.StatusOK, res, nil
}

func (a *ToDO) Delete(_ http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	
	// accessTokenから取り出した情報を取得する
	getc, err := httputil.GetClaimsFromContext(r.Context())
	if err != nil {
		//fmt.Println("accessToken")
		return http.StatusInternalServerError, nil, err
	}

	list_id, err := URLToID(r)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	
	Service := service.NewToDO(a.db)
	_, err = Service.Delete(list_id, getc.ID)
	
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	
	return http.StatusOK, nil, nil
}

func StringCheck(str string) bool {
	chars := []string{"?", "!", "\\*","\\_", "\\#", "<", ">", "\\\\", "(", ")", "\\$", "\"", "%", "=", "~", "|", "[", "]", ";", "\\+", ":", "{", "}", "@", "\\`", "/", "；", "＠", "＋", "：", "＊", "｀", "「", "」", "｛", "｝", "＿", "？", "。", "、", "＞", "＜"}
    r := strings.Join(chars, "")
	symbol := regexp.MustCompile("[" + r + "]+")
	if symbol.Match([]byte(str)) {
		//上記が含まれている
		return false
	}
	return true
}

func URLToID(r *http.Request) (int, error) {
	url := r.URL.Path
	strID := url[strings.LastIndex(url, "/")+1:]
	id, err := strconv.Atoi( strID )
	if err != nil {
		return 0, err
	}
	
	return id, nil
}
