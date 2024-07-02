package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/captcha/cache"
	"gitee.com/go-nianxi/go-common/pkg/captcha/util"
	"log"
)

const (
	TEXT = "明有奇巧人曰王叔远能以径寸之木为宫室器皿人物以至鸟兽木石罔不因势象形各具情态尝贻余核舟一盖大苏泛赤壁云舟首尾长约八分有奇高可二黍许中轩敞者为舱箬篷覆之旁开小窗左右各四共八扇启窗而观雕栏相望焉闭之则右刻山高月小水落石出左刻清风徐来水波不兴”石青糁之船头坐三人中峨冠而多髯者为东坡佛印居右鲁直居左苏黄共阅一手卷东坡右手执卷端左手抚鲁直背鲁直左手执卷末右手指卷如有所语东坡现右足鲁直现左足各微侧其两膝相比者各隐卷底衣褶中佛印绝类弥勒袒胸露乳矫首昂视神情与苏黄不属卧右膝诎右臂豫章故郡洪都新府星分翼轸地接衡庐襟三江而带五湖控蛮荆而引瓯越物华天宝龙光射牛斗之墟人杰地灵徐孺下陈蕃之榻雄州雾列俊采星驰台隍枕夷夏之交宾主尽东南之美都督阎公之雅望棨戟遥临宇文新州之懿范襜帷暂驻十旬休假胜友如云千里逢迎高朋满座腾蛟起凤孟学士之词宗紫电青霜王将军之武库家君作宰路出名区童子何知躬逢胜饯(豫章故郡时维九月序属三秋潦水尽而寒潭清烟光凝而暮山紫俨骖騑于上路访风景于崇阿临帝子之长洲得天人之旧馆层峦耸翠上出重霄飞阁流丹下临无地鹤汀凫渚穷岛屿之萦回桂殿兰宫即冈峦之体势"
)

type ClickWordCaptchaService struct {
	cache  cache.CacheInterface
	config *CaptchaConfig
	//factory      *config.CaptchaServiceFactory
	fontPath string
}

func NewClickWordCaptchaService(resourcePath string, config *CaptchaConfig, cache cache.CacheInterface) CaptchaInterface {
	util.SetUp(resourcePath, DefaultBackgroundImageDirectory, DefaultClickBackgroundImageDirectory, DefaultTemplateImageDirectory, DefaultFontFile)
	return &ClickWordCaptchaService{cache: cache, config: config}
}

func (c *ClickWordCaptchaService) Get() (map[string]interface{}, error) {
	// 初始化背景图片
	backgroundImage := util.NewImageByClickBackgroundImage(util.FontPathDirectory)

	pointList, wordList, err := c.getImageData(backgroundImage)
	if err != nil {
		return nil, err
	}

	originalImageBase64, err := backgroundImage.Base64()

	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["originalImageBase64"] = originalImageBase64
	data["wordList"] = wordList
	data["secretKey"] = pointList[0].SecretKey
	data["token"] = util.GetUuid()

	//data := make(map[string]interface{})
	//data["originalImageBase64"] = originalImageBase64
	//data["jigsawImageBase64"] = jigsawImageBase64
	//data["secretKey"] = b.point.SecretKey
	//data["token"] = util.GetUuid()

	codeKey := fmt.Sprintf(CodeKeyPrefix, data["token"])
	jsonPoint, err := json.Marshal(pointList)
	if err != nil {
		log.Printf("point json Marshal err: %v", err)
		return nil, err
	}

	c.cache.Set(codeKey, string(jsonPoint), c.config.CacheExpireSec)
	return data, nil
}

func (c *ClickWordCaptchaService) Check(token string, pointJson string) error {
	cache := c.cache
	codeKey := fmt.Sprintf(CodeKeyPrefix, token)

	cachePointInfo := cache.Get(codeKey)

	if cachePointInfo == "" {
		return errors.New("验证码已失效")
	}

	// 解析结构体
	var cachePoint []util.PointVO

	var userPoint []util.PointVO

	err := json.Unmarshal([]byte(cachePointInfo), &cachePoint)

	if err != nil {
		return err
	}

	// 解密前端传递过来的数据
	userPointJson := util.AesDecrypt(pointJson, cachePoint[0].SecretKey)

	err = json.Unmarshal([]byte(userPointJson), &userPoint)

	if err != nil {
		return err
	}

	for i, pointVO := range cachePoint {
		targetPoint := userPoint[i]
		fontsize := c.config.ClickWord.FontSize
		if targetPoint.X-fontsize > pointVO.X || targetPoint.X > pointVO.X+fontsize || targetPoint.Y-fontsize > pointVO.Y || targetPoint.Y > pointVO.Y+fontsize {
			return errors.New("验证失败")
		}
	}

	return nil
}

func (c *ClickWordCaptchaService) Verification(token string, pointJson string) error {
	err := c.Check(token, pointJson)
	if err != nil {
		return err
	}
	codeKey := fmt.Sprintf(CodeKeyPrefix, token)
	c.cache.Delete(codeKey)
	return nil
}

func (c *ClickWordCaptchaService) getImageData(image *util.ImageUtil) ([]util.PointVO, []string, error) {
	wordCount := c.config.ClickWord.FontNum

	// 某个字不参与校验
	num := util.RandomInt(1, wordCount)
	currentWord := c.getRandomWords(wordCount)

	var pointList []util.PointVO
	var wordList []string

	i := 0

	// 构建本次的 secret
	key := util.RandString(16)

	for _, s := range currentWord {
		point := c.randomWordPoint(image.Width, image.Height, i, wordCount)
		point.SetSecretKey(key)
		// 随机设置文字 TODO 角度未设置
		err := image.SetArtText(s, c.config.ClickWord.FontSize, point)

		if err != nil {
			return nil, nil, err
		}

		if (num - 1) != i {
			pointList = append(pointList, point)
			wordList = append(wordList, s)
		}
		i++
	}
	return pointList, wordList, nil
}

// getRandomWords 获取随机文件
func (c *ClickWordCaptchaService) getRandomWords(count int) []string {
	runesArray := []rune(TEXT)
	size := len(runesArray)

	set := make(map[string]bool)
	var wordList []string

	for {
		word := runesArray[util.RandomInt(0, size-1)]
		set[string(word)] = true
		if len(set) >= count {
			for str, _ := range set {
				wordList = append(wordList, str)
			}
			break
		}
	}
	return wordList
}

func (c *ClickWordCaptchaService) randomWordPoint(width int, height int, i int, count int) util.PointVO {
	avgWidth := width / (count + 1)
	fontSizeHalf := c.config.ClickWord.FontSize / 2

	var x, y int
	if avgWidth < fontSizeHalf {
		x = util.RandomInt(1+fontSizeHalf, width)
	} else {
		if i == 0 {
			x = util.RandomInt(1+fontSizeHalf, avgWidth*(i+1)-fontSizeHalf)
		} else {
			x = util.RandomInt(avgWidth*i+fontSizeHalf, avgWidth*(i+1)-fontSizeHalf)
		}
	}
	y = util.RandomInt(c.config.ClickWord.FontSize, height-fontSizeHalf)
	return util.PointVO{X: x, Y: y}
}
