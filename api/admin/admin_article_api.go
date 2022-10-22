package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"r0Website-server/models/po"
	"r0Website-server/models/vo"
	"r0Website-server/service"
	"r0Website-server/utils/msg"
)

type ArticleController struct {
	ArticleService *service.ArticleService `R0Ioc:"true"`
}

// ArticleList 文章列表
func (articleCon *ArticleController) ArticleList(c *gin.Context) {

}

// ArticleFileWay 增加文章通过上传文件
func (articleCon *ArticleController) ArticleFileWay(c *gin.Context) {
	articleID := c.Param("id")
	var articleFile vo.AdminArticleAddFileVo
	if err := c.ShouldBind(&articleFile); err != nil {
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed("参数异常"))
		return
	}
	if articleFile.Author == "" {
		// 从上下文中获取作者
		if curUser, exists := c.Get("userInfo"); exists {
			// 其实此key不可能不存在
			if val, ok := curUser.(po.User); ok {
				// 其实类型推导也不可能不ok
				articleFile.Author = val.Username
			}
		}
	}
	ans, err := articleCon.ArticleService.ArticleADDFile(articleFile, articleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, msg.NewMsg().Success(ans))
}

// ArticleFormWay 增加文章通过编辑方式 提交一个**表单**
func (articleCon *ArticleController) ArticleFormWay(c *gin.Context) {
	articleID := c.Param("id")
	var articleForm vo.AdminArticleAddFormVo
	if err := c.ShouldBind(&articleForm); err != nil {
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed("参数异常"))
		return
	}
	if articleForm.Author == "" {
		if curUser, exists := c.Get("userInfo"); exists {
			// 其实此key不可能不存在
			if val, ok := curUser.(po.User); ok {
				// 其实类型推导也不可能不ok
				articleForm.Author = val.Username
			}
		}
	}
	ans, err := articleCon.ArticleService.ArticleADDForm(articleForm, articleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
		return
	}
	c.JSON(http.StatusOK, msg.NewMsg().Success(ans))
}

// ArticleDelete 文章的删除（支持批量删除）
func (articleCon *ArticleController) ArticleDelete(c *gin.Context) {

}

// ArticleOverhead 文章的顶置
func (articleCon *ArticleController) ArticleOverhead(c *gin.Context) {

}

// ArticleContent 文章内容
func (articleCon *ArticleController) ArticleContent(c *gin.Context) {

}

// ArticleUpdate 文章的更新
func (articleCon *ArticleController) ArticleUpdate(c *gin.Context) {

}

// ArticleBackup 文章的备份
func (articleCon *ArticleController) ArticleBackup(c *gin.Context) {

}
