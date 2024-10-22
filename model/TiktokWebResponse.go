package model

type Item struct {
	ID                string      `json:"id"`
	Desc              string      `json:"desc"`
	CreateTime        int64       `json:"createTime"`
	Video             Video       `json:"video"`
	Author            Author      `json:"author"`
	Music             Music       `json:"music"`
	Challenges        []Challenge `json:"challenges"`
	Stats             Stats       `json:"stats"`
	DuetInfo          DuetInfo    `json:"duetInfo"`
	OriginalItem      bool        `json:"originalItem"`
	OfficialItem      bool        `json:"officalItem"`
	TextExtra         []TextExtra `json:"textExtra"`
	Secret            bool        `json:"secret"`
	ForFriend         bool        `json:"forFriend"`
	Digged            bool        `json:"digged"`
	ItemCommentStatus int         `json:"itemCommentStatus"`
	ShowNotPass       bool        `json:"showNotPass"`
	Vl1               bool        `json:"vl1"`
	ItemMute          bool        `json:"itemMute"`
	AuthorStats       AuthorStats `json:"authorStats"`
	PrivateItem       bool        `json:"privateItem"`
	DuetEnabled       bool        `json:"duetEnabled"`
	StitchEnabled     bool        `json:"stitchEnabled"`
	ShareEnabled      bool        `json:"shareEnabled"`
	IsAd              bool        `json:"isAd"`
	Collected         bool        `json:"collected"`
}

type Video struct {
	ID            string   `json:"id"`
	Height        int      `json:"height"`
	Width         int      `json:"width"`
	Duration      int      `json:"duration"`
	Ratio         string   `json:"ratio"`
	Cover         string   `json:"cover"`
	OriginCover   string   `json:"originCover"`
	DynamicCover  string   `json:"dynamicCover"`
	PlayAddr      string   `json:"playAddr"`
	ShareCover    []string `json:"shareCover"`
	ReflowCover   string   `json:"reflowCover"`
	Bitrate       int      `json:"bitrate"`
	EncodedType   string   `json:"encodedType"`
	Format        string   `json:"format"`
	VideoQuality  string   `json:"videoQuality"`
	EncodeUserTag string   `json:"encodeUserTag"`
}

type Author struct {
	ID              string `json:"id"`
	UniqueId        string `json:"uniqueId"`
	Nickname        string `json:"nickname"`
	AvatarThumb     string `json:"avatarThumb"`
	AvatarMedium    string `json:"avatarMedium"`
	AvatarLarger    string `json:"avatarLarger"`
	Signature       string `json:"signature"`
	Verified        bool   `json:"verified"`
	SecUid          string `json:"secUid"`
	Secret          bool   `json:"secret"`
	Ftc             bool   `json:"ftc"`
	Relation        int    `json:"relation"`
	OpenFavorite    bool   `json:"openFavorite"`
	CommentSetting  int    `json:"commentSetting"`
	DuetSetting     int    `json:"duetSetting"`
	StitchSetting   int    `json:"stitchSetting"`
	PrivateAccount  bool   `json:"privateAccount"`
	DownloadSetting int    `json:"downloadSetting"`
}

type Music struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	PlayUrl     string `json:"playUrl"`
	CoverThumb  string `json:"coverThumb"`
	CoverMedium string `json:"coverMedium"`
	CoverLarge  string `json:"coverLarge"`
	AuthorName  string `json:"authorName"`
	Original    bool   `json:"original"`
	Duration    int    `json:"duration"`
	Album       string `json:"album"`
}

type Challenge struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	ProfileThumb  string `json:"profileThumb"`
	ProfileMedium string `json:"profileMedium"`
	ProfileLarger string `json:"profileLarger"`
	CoverThumb    string `json:"coverThumb"`
	CoverMedium   string `json:"coverMedium"`
	CoverLarger   string `json:"coverLarger"`
	IsCommerce    bool   `json:"isCommerce"`
}

type Stats struct {
	DiggCount    int `json:"diggCount"`
	ShareCount   int `json:"shareCount"`
	CommentCount int `json:"commentCount"`
	PlayCount    int `json:"playCount"`
	CollectCount int `json:"collectCount"`
}

type DuetInfo struct {
	DuetFromId string `json:"duetFromId"`
}

type TextExtra struct {
	AwemeId      string `json:"awemeId"`
	Start        int    `json:"start"`
	End          int    `json:"end"`
	HashtagName  string `json:"hashtagName"`
	HashtagId    string `json:"hashtagId"`
	Type         int    `json:"type"`
	UserId       string `json:"userId"`
	IsCommerce   bool   `json:"isCommerce"`
	UserUniqueId string `json:"userUniqueId"`
	SecUid       string `json:"secUid"`
	SubType      int    `json:"subType"`
}

type AuthorStats struct {
	FollowingCount int `json:"followingCount"`
	FollowerCount  int `json:"followerCount"`
	HeartCount     int `json:"heartCount"`
	VideoCount     int `json:"videoCount"`
	DiggCount      int `json:"diggCount"`
	Heart          int `json:"heart"`
}

type TiktokWebResponse struct {
	Code   int32            `json:"code"`
	Router string           `json:"router"`
	Params TiktokWebParamsR `json:"params"`
	Data   TiktokDataDetail `json:"data"`
}

type TiktokWebParamsSend struct {
	Keyword  string `json:"keyword"`
	Count    int    `json:"count"`
	Offset   int    `json:"offset"`
	SearchId string `json:"search_id"`
	Cookie   string `json:"cookie"`
}

type TiktokWebParamsR struct {
	Keyword  string `json:"keyword"`
	Count    string `json:"count"`
	Offset   string `json:"offset"`
	SearchId string `json:"search_id"`
	Cookie   string `json:"cookie"`
}

type TiktokDataDetail struct {
	StatusCode int32  `json:"status_code"`
	ItemList   []Item `json:"item_list"`
	HasMore    int32  `json:"has_more"`
	Cursor     int32  `json:"cursor"`
	extra      `json:"extra"`
	LogPb      struct {
		ImprId string `json:"impr_id"`
	} `json:"log_pb"`
	Backtrace string `json:"backtrace"`
}

type extra struct {
	Now             int64         `json:"now"`
	Logid           string        `json:"logid"`
	FatalItemIds    []interface{} `json:"fatal_item_ids"`
	SearchRequestId string        `json:"search_request_id"`
	ApiDebugInfo    string        `json:"api_debug_info"`
}
