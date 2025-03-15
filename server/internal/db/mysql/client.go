package mysql

/**
 * @Description
 * @Date 2025/3/15 19:31
 **/
import "gorm.io/gorm"

type Client struct {
	*gorm.DB //可定制
}
