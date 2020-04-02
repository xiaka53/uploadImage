#### 基于github.com/nfnt/resize这个包做的一个图片上传服务端的接收以及压缩处理


可配置项：

    type UploadFile struct {
        Param       string     //前端传文件字段
        StandByType []filetype //接收图片类型
        MaxSize     int64      //文件最大限制
        Width       float64    //宽度px
        Height      float64    //高度px
        Narrow      narrow     //缩小类型
        FileName    string     //生成文件名
    }

上传图片支持类型：JPEG｜PNG｜GIF

文件名：目前仅支持png类型保存

示例：

    package main

    import (
	    "github.com/gin-gonic/gin"
	    "github.com/xiaka53/uploadImage"
	    "strconv"
	    "time"
    )

    func main() {
	    r := gin.Default()
	    r.POST("/upload", uploadFile)
	    r.Run()
    }

    func uploadFile(c *gin.Context) {
	    now := int(time.Now().Unix())
	    path := "img/" + strconv.Itoa(now) + ".png"
	    upload := uploadImage.NewUploadFile("img", path, 240, 240, 10*1024*1024, uploadImage.Fixed, uploadImage.PNG, uploadImage.JPEG)
	    str, err := upload.UploadImage(c)
	    if err != nil {
		    c.JSON(500, err)
		    return
	    }
	    c.JSON(200, str)
    }