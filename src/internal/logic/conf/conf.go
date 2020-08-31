package conf

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strings"
	"time"
)

var (
	// Conf config
	Conf     *Config
	host     string
	confPath string
)

func init() {
	host, _ = os.Hostname()
	flag.StringVar(&confPath, "conf", "comet.conf", "comet config path")
}

// Init init config.
func Init() (err error) {
	Conf = Default()
	_, err = toml.DecodeFile(confPath, &Conf)
	fmt.Println(Conf.Kafka)
	return
}

// Default new a config with specified defualt value.
func Default() *Config {
	word := "阿扁推翻，阿宾，阿賓，挨了一炮，爱液横流，安街逆，安局办公楼，安局豪华，安门事，安眠藥，案的准确，八九民，八九学，八九政治，把病人整，把邓小平，把学生整，罢工门，白黄牙签，败培训，办本科，办理本科，办理各种，办理票据，办理文凭，办理真实，办理证书，办理资格，办文凭，办怔，办证，半刺刀，辦毕业，辦證，谤罪获刑，磅解码器，磅遥控器，宝在甘肃修，保过答案，报复执法，爆发骚，北省委门，被打死，被指抄袭，被中共，本公司担，本无码，毕业證，变牌绝，辩词与梦，冰毒，冰火毒，冰火佳，冰火九重，冰火漫，冰淫传，冰在火上，波推龙，博彩娱，博会暂停，博园区伪，不查都，不查全，不思四化，布卖淫女，部忙组阁，部是这样，才知道只生，财众科技，采花堂，踩踏事，苍山兰，苍蝇水，藏春阁，藏獨，操了嫂，操嫂子，策没有不，插屁屁，察象蚂，拆迁灭，车牌隐，成人电，成人卡通，成人聊，成人片，成人视，成人图，成人文，成人小，城管灭，惩公安，惩贪难，充气娃，冲凉死，抽着大中，抽着芙蓉，出成绩付，出售发票，出售军，穿透仪器，春水横溢，纯度白，纯度黄，次通过考，催眠水，催情粉，催情药，催情藥，挫仑，达毕业证，答案包，答案提供，打标语，打错门，打飞机专，打死经过，打死人，打砸办公，大鸡巴，大雞巴，大纪元，大揭露，大奶子，大批贪官，大肉棒，大嘴歌，代办发票，代办各，代办文，代办学，代办制，代辦，代表烦，代開，代考，代理发票，代理票据，代您考，代您考，代写毕，代写论，代孕，贷办，贷借款，贷开，戴海静，当代七整，当官要精，当官在于，党的官，党后萎，党前干劲，刀架保安，导的情人，导叫失，导人的最，导人最，导小商，到花心，得财兼，的同修，灯草和，等级證，等屁民，等人老百，等人是老，等人手术，邓爷爷转，邓玉娇，地产之歌，地下先烈，地震哥，帝国之梦，递纸死，点数优惠，电狗，电话监，电鸡，甸果敢，蝶舞按，丁香社，丁子霖，顶花心，东北独立，东复活，东京热，東京熱，洞小口紧，都当警，都当小姐，都进中央，毒蛇钻，独立台湾，赌球网，短信截，对日强硬，多美康，躲猫猫，俄羅斯，恶势力操，恶势力插，恩氟烷，儿园惨，儿园砍，儿园杀，儿园凶，二奶大，发牌绝，发票出，发票代，发票销，發票，法车仑，法伦功，法轮，法轮佛，法维权，法一轮，法院给废，法正乾，反测速雷，反雷达测，反屏蔽，范燕琼，方迷香，防电子眼，防身药水，房贷给废，仿真枪，仿真证，诽谤罪，费私服，封锁消，佛同修，夫妻交换，福尔马林，福娃的預，福娃頭上，福香巴，府包庇，府集中领，妇销魂，附送枪，复印件生，复印件制，富民穷，富婆给废，改号软件，感扑克，冈本真，肛交，肛门是邻，岡本真，钢针狗，钢珠枪，港澳博球，港馬會，港鑫華，高就在政，高考黑，高莺莺，搞媛交，告长期，告洋状，格证考试，各类考试，各类文凭，跟踪器，工程吞得，工力人，公安错打，公安网监，公开小姐，攻官小姐，共狗，共王储，狗粮，狗屁专家，鼓动一些，乖乖粉，官商勾，官也不容，官因发帖，光学真题，跪真相，滚圆大乳，国际投注，国家妓，国家软弱，国家吞得，国库折，国一九五七，國內美，哈药直销，海访民，豪圈钱，号屏蔽器，和狗交，和狗性，和狗做，黑火药的，红色恐怖，红外透视，紅色恐，胡江内斗，胡紧套，胡錦濤，胡适眼，胡耀邦，湖淫娘，虎头猎，华国锋，华门开，化学扫盲，划老公，还会吹萧，还看锦涛，环球证件，换妻，皇冠投注，黄冰，浑圆豪乳，活不起，火车也疯，机定位器，机号定，机号卫，机卡密，机屏蔽器，基本靠吼，绩过后付，激情电，激情短，激情妹，激情炮，级办理，级答案，急需嫖，集体打砸，集体腐，挤乳汁，擠乳汁，佳静安定，家一样饱，家属被打，甲虫跳，甲流了，奸成瘾，兼职上门，监听器，监听王，简易炸，江胡内斗，江太上，江系人，江贼民，疆獨，蒋彦永，叫自慰，揭贪难，姐包夜，姐服务，姐兼职，姐上门，金扎金，金钟气，津大地震，津地震，进来的罪，京地震，京要地震，经典谎言，精子射在，警察被，警察的幌，警察殴打，警察说保，警车雷达，警方包庇，警用品，径步枪，敬请忍，究生答案，九龙论坛，九评共，酒象喝汤，酒像喝汤，就爱插，就要色，举国体，巨乳，据说全民，绝食声，军长发威，军刺，军品特，军用手，开邓选，开锁工具，開碼，開票，砍杀幼，砍伤儿，康没有不，康跳楼，考答案，考后付款，考机构，考考邓，考联盟，考前答，考前答案，考前付，考设备，考试包过，考试保，考试答案，考试机构，考试联盟，考试枪，考研考中，考中答案，磕彰，克分析，克千术，克透视，空和雅典，孔摄像，控诉世博，控制媒，口手枪，骷髅死，快速办，矿难不公，拉登说，拉开水晶，来福猎，拦截器，狼全部跪，浪穴，老虎机，雷人女官，类准确答，黎阳平，李洪志，李咏曰，理各种证，理是影帝，理证件，理做帐报，力骗中央，力月西，丽媛离，利他林，连发手，聯繫電，炼大法，两岸才子，两会代，两会又三，聊视频，聊斋艳，了件渔袍，猎好帮手，猎枪销，猎槍，獵槍，领土拿，流血事，六合彩，六死，六四事，六月联盟，龙湾事件，隆手指，陆封锁，陆同修，氯胺酮，乱奸，乱伦类，乱伦小，亂倫，伦理大，伦理电影，伦理毛，伦理片，轮功，轮手枪，论文代，罗斯小姐，裸聊网，裸舞视，落霞缀，麻古，麻果配，麻果丸，麻将透，麻醉狗，麻醉枪，麻醉槍，麻醉藥，蟆叫专家，卖地财政，卖发票，卖银行卡，卖自考，漫步丝，忙爱国，猫眼工具，毛一鲜，媒体封锁，每周一死，美艳少妇，妹按摩，妹上门，门按摩，门保健，門服務，氓培训，蒙汗药，迷幻型，迷幻药，迷幻藥，迷昏口，迷昏药，迷昏藥，迷魂香，迷魂药，迷魂藥，迷奸药，迷情水，迷情药，迷藥，谜奸药，蜜穴，灭绝罪，民储害，民九亿商，民抗议，明慧网，铭记印尼，摩小姐，母乳家，木齐针，幕没有不，幕前戲，内射，南充针，嫩穴，嫩阴，泥马之歌，你的西域，拟涛哥，娘两腿之间，妞上门，浓精，怒的志愿，女被人家搞，女激情，女技师，女人和狗，女任职名，女上门，女優，鸥之歌，拍肩神药，拍肩型，牌分析，牌技网，炮的小蜜，陪考枪，配有消，喷尿，嫖俄罗，嫖鸡，平惨案，平叫到床，仆不怕饮，普通嘌，期货配，奇迹的黄，奇淫散，骑单车出，气狗，气枪，汽狗，汽枪，氣槍，铅弹，钱三字经，枪出售，枪的参，枪的分，枪的结，枪的制，枪货到，枪决女犯，枪决现场，枪模，枪手队，枪手网，枪销售，枪械制，枪子弹，强权政府，强硬发言，抢其火炬，切听器，窃听器，禽流感了，勤捞致，氢弹手，清除负面，清純壆，情聊天室，情妹妹，情视频，情自拍，氰化钾，氰化钠，请集会，请示威，请愿，琼花问，区的雷人，娶韩国，全真证，群奸暴，群起抗暴，群体性事，绕过封锁，惹的国，人权律，人体艺，人游行，人在云上，人真钱，认牌绝，任于斯国，柔胸粉，肉洞，肉棍，如厕死，乳交，软弱的国，赛后骚，三挫，三级片，三秒倒，三网友，三唑，骚妇，骚浪，骚穴，骚嘴，扫了爷爷，色电影，色妹妹，色视频，色小说，杀指南，山涉黑，煽动不明，煽动群众，上门激，烧公安局，烧瓶的，韶关斗，韶关玩，韶关旭，射网枪，涉嫌抄袭，深喉冰，神七假，神韵艺术，生被砍，生踩踏，生肖中特，圣战不息，盛行在舞，尸博，失身水，失意药，狮子旗，十八等，十大谎，十大禁，十个预言，十类人不，十七大幕，实毕业证，实体娃，实学历文，士康事件，式粉推，视解密，是躲猫，手变牌，手答案，手狗，手机跟，手机监，手机窃，手机追，手拉鸡，手木仓，手槍，守所死法，兽交，售步枪，售纯度，售单管，售弹簧刀，售防身，售狗子，售虎头，售火药，售假币，售健卫，售军用，售猎枪，售氯胺，售麻醉，售冒名，售枪支，售热武，售三棱，售手枪，售五四，售信用，售一元硬，售子弹，售左轮，书办理，熟妇，术牌具，双管立，双管平，水阎王，丝护士，丝情侣，丝袜保，丝袜恋，丝袜美，丝袜妹，丝袜网，丝足按，司长期有，司法黑，私房写真，死法分布，死要见毛，四博会，四大扯，个四小码，苏家屯集，诉讼集团，素女心，速代办，速取证，酸羟亚胺，蹋纳税，太王四神，泰兴幼，泰兴镇中，泰州幼，贪官也辛，探测狗，涛共产，涛一样胡，特工资，特码，特上门，体透视镜，替考，替人体，天朝特，天鹅之旅，天推广歌，田罢工，田田桑，田停工，庭保养，庭审直播，通钢总经，偷電器，偷肃贪，偷听器，偷偷贪，头双管，透视功能，透视镜，透视扑，透视器，透视眼镜，透视药，透视仪，秃鹰汽，突破封锁，突破网路，推油按，脱衣艳，瓦斯手，袜按摩，外透视镜，外围赌球，湾版假，万能钥匙，万人骚动，王立军，王益案，网民案，网民获刑，网民诬，微型摄像，围攻警，围攻上海，维汉员，维权基，维权人，维权谈，委坐船，谓的和谐，温家堡，温切斯特，温影帝，溫家寶，瘟加饱，瘟假饱，文凭证，文强，纹了毛，闻被控制，闻封锁，瓮安，我的西域，我搞台独，乌蝇水，无耻语录，无码专，五套功，五月天，午夜电，午夜极，武警暴，武警殴，武警已增，务员答案，务员考试，雾型迷，西藏限，西服进去，希脏，习进平，习晋平，席复活，席临终前，席指着护，洗澡死，喜贪赃，先烈纷纷，现大地震，现金投注，线透视镜，限制言，陷害案，陷害罪，相自首，香港论坛，香港马会，香港一类，香港总彩，硝化甘，小穴，校骚乱，协晃悠，写两会，泄漏的内，新建户，新疆叛，新疆限，新金瓶，新唐人，信访专班，信接收器，兴中心幼，星上门，行长王益，形透视镜，型手枪，姓忽悠，幸运码，性爱日，性福情，性感少，性推广歌，胸主席，徐玉元，学骚乱，学位證，學生妹，丫与王益，烟感器，严晓玲，言被劳教，言论罪，盐酸曲，颜射，恙虫病，姚明进去，要人权，要射精了，要射了，要泄了，夜激情，液体炸，一小撮别，遗情书，蚁力神，益关注组，益受贿，阴间来电，陰唇，陰道，陰戶，淫魔舞，淫情女，淫肉，淫騷妹，淫兽，淫兽学，淫水，淫穴，隐形耳，隐形喷剂，应子弹，婴儿命，咏妓，用手枪，幽谷三，游精佑，有奶不一，右转是政，幼齿类，娱乐透视，愚民同，愚民政，与狗性，玉蒲团，育部女官，冤民大，鸳鸯洗，园惨案，园发生砍，园砍杀，园凶杀，园血案，原一九五七，原装弹，袁腾飞，晕倒型，韵徐娘，遭便衣，遭到警，遭警察，遭武警，择油录，曾道人，炸弹教，炸弹遥控，炸广州，炸立交，炸药的制，炸药配，炸药制，张春桥，找枪手，找援交，找政法委副，赵紫阳，针刺案，针刺伤，针刺事，针刺死，侦探设备，真钱斗地，真钱投注真善忍，真实文凭，真实资格，震惊一个民，震其国土，证到付款，证件办，证件集团，证生成器，证书办，证一次性，政府操，政论区，證件，植物冰，殖器护，指纹考勤，指纹膜，指纹套，至国家高，志不愿跟，制服诱，制手枪，制证定金，制作证件，中的班禅，中共黑，中国不强，种公务员，种学历证，众像羔，州惨案，州大批贪，州三箭，宙最高法，昼将近，主席忏，住英国房，助考，助考网，专业办理，专业代，专业代写，专业助，转是政府，赚钱资料，装弹甲，装枪套，装消音，着护士的胸，着涛哥，姿不对死，资格證，资料泄，梓健特药，字牌汽，自己找枪，自慰用，自由圣，自由亚，总会美女，足球玩法，最牛公安，醉钢枪，醉迷药，醉乙醚，尊爵粉，左转是政，作弊器，作各种证，作硝化甘，唑仑，做爱小，做原子弹，做证件，傻逼，傻B，弱智，脑残"
	wordList := strings.Split(word, "，")
	return &Config{
		HTTPServer: &HTTPServer{
			Network:      "tcp",
			Addr:         "localhost:2333",
			ReadTimeout:  time.Duration(time.Second),
			WriteTimeout: time.Duration(time.Second),
		},
		Kafka: &Kafka{
			Topic:   "test",
			Brokers: []string{"localhost:9092"},
		},
		WordList: wordList,
	}
}

// Config config.
type Config struct {
	HTTPServer *HTTPServer
	Kafka      *Kafka
	WordList   []string
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IsLimit      bool
	RedisAddr    string
	Count        int64
	Dur          time.Duration
}

type Kafka struct {
	Topic   string
	Brokers []string
}
