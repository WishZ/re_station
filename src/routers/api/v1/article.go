package v1

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"go-blog/dao"
	"go-blog/dto/inputDto"
	"go-blog/pkg/e"
	"go-blog/pkg/logging"
	"go-blog/pkg/setting"
	"go-blog/pkg/util"
	"net/http"
)

//获取单个文章
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.InvalidParams
	var data interface{}
	if ! valid.HasErrors() {
		if dao.ExistArticleByID(id) {
			data = dao.GetArticle(id)
			code = e.Success
		} else {
			code = e.ErrorNotExistArticle
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//获取多个文章
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.InvalidParams
	if ! valid.HasErrors() {
		code = e.Success

		data["lists"] = dao.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = dao.GetArticleTotal(maps)

	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//新增文章
func AddArticle(c *gin.Context) {
	var articleCreateDto inputDto.ArticleCreateInputDto
	err := c.BindJSON(&articleCreateDto)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": "Post Data Err"})
		return
	}
	valid := validation.Validation{}
	valid.Min(articleCreateDto.TagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(articleCreateDto.Title, "title").Message("标题不能为空")
	valid.Required(articleCreateDto.Desc, "desc").Message("简述不能为空")
	valid.Required(articleCreateDto.Content, "content").Message("内容不能为空")
	valid.Required(articleCreateDto.CreatedBy, "created_by").Message("创建人不能为空")
	valid.Range(articleCreateDto.State, 0, 1, "state").Message("状态只允许0或1")
	code := e.InvalidParams
	if ! valid.HasErrors() {
		if dao.ExistTagByID(articleCreateDto.TagId) {
			dao.AddArticle(articleCreateDto)
			code = e.Success
		} else {
			code = e.ErrorNotExistTag
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

//修改文章
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	var articleEditDto inputDto.ArticleEditInputDto
	err := c.BindJSON(&articleEditDto)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": "Post Data Err"})
		return
	}

	valid.Range(articleEditDto.State, 0, 1, "state").Message("状态只允许0或1")
	valid.Min(articleEditDto.Id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(articleEditDto.Title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(articleEditDto.Desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(articleEditDto.Content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(articleEditDto.Modifier, "modified_by").Message("修改人不能为空")
	valid.MaxSize(articleEditDto.Modifier, 100, "modified_by").Message("修改人最长为100字符")

	code := e.InvalidParams
	if ! valid.HasErrors() {
		if dao.ExistArticleByID(articleEditDto.Id) {
			if dao.ExistTagByID(articleEditDto.TagId) {
				data := make(map[string]interface{})
				if articleEditDto.TagId > 0 {
					data["tag_id"] = articleEditDto.TagId
				}
				if articleEditDto.Title != "" {
					data["title"] = articleEditDto.Title
				}
				if articleEditDto.Desc != "" {
					data["desc"] = articleEditDto.Desc
				}
				if articleEditDto.Content != "" {
					data["content"] = articleEditDto.Content
				}

				data["modifier"] = articleEditDto.Modifier

				dao.EditArticle(articleEditDto.Id, data)
				code = e.Success
			} else {
				code = e.ErrorNotExistTag
			}
		} else {
			code = e.ErrorNotExistArticle
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.InvalidParams
	if ! valid.HasErrors() {
		if dao.ExistArticleByID(id) {
			dao.DeleteArticle(id)
			code = e.Success
		} else {
			code = e.ErrorNotExistArticle
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}