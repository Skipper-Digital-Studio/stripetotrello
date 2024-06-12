package trello

import (
	"time"
)

const DEFAULT_CARDS_PATH = "static/trello/cards.json"

type (
	ListRes struct {
		Id         string  `json:"id"`
		Name       string  `json:"name"`
		Closed     bool    `json:"closed"`
		Pos        float64 `json:"pos"`
		SoftLimit  string  `json:"softLimit"`
		Color      string  `json:"color"`
		Subscribed bool    `json:"subscribed"`
		Limits     struct {
			Attachements struct {
				PerBoard struct {
					Status    string `json:"status"`
					DisableAt int64  `json:"disableAt"`
					WarnAt    int64  `json:"warnAt"`
				} `json:"perBoard"`
			} `json:"attachements"`
		} `json:"limits"`
	}

	BoardRes struct {
		Closed         bool        `json:"closed"`
		Desc           string      `json:"desc"`
		DescData       interface{} `json:"descData"`
		Id             string      `json:"id"`
		IdEnterprise   interface{} `json:"idEnterprise"`
		IdOrganization string      `json:"idOrganization"`
		LabelNames     struct {
			Black       string `json:"black"`
			BlackDark   string `json:"black_dark"`
			BlackLight  string `json:"black_light"`
			Blue        string `json:"blue"`
			BlueDark    string `json:"blue_dark"`
			BlueLight   string `json:"blue_light"`
			Green       string `json:"green"`
			GreenDark   string `json:"green_dark"`
			GreenLight  string `json:"green_light"`
			Lime        string `json:"lime"`
			LimeDark    string `json:"lime_dark"`
			LimeLight   string `json:"lime_light"`
			Orange      string `json:"orange"`
			OrangeDark  string `json:"orange_dark"`
			OrangeLight string `json:"orange_light"`
			Pink        string `json:"pink"`
			PinkDark    string `json:"pink_dark"`
			PinkLight   string `json:"pink_light"`
			Purple      string `json:"purple"`
			PurpleDark  string `json:"purple_dark"`
			PurpleLight string `json:"purple_light"`
			Red         string `json:"red"`
			RedDark     string `json:"red_dark"`
			RedLight    string `json:"red_light"`
			Sky         string `json:"sky"`
			SkyDark     string `json:"sky_dark"`
			SkyLight    string `json:"sky_light"`
			Yellow      string `json:"yellow"`
			YellowDark  string `json:"yellow_dark"`
			YellowLight string `json:"yellow_light"`
		} `json:"labelNames"`
		Limits struct {
		} `json:"limits"`
		Name   string `json:"name"`
		Pinned bool   `json:"pinned"`
		Prefs  struct {
			Background               string        `json:"background"`
			BackgroundBottomColor    string        `json:"backgroundBottomColor"`
			BackgroundBrightness     string        `json:"backgroundBrightness"`
			BackgroundColor          string        `json:"backgroundColor"`
			BackgroundImage          interface{}   `json:"backgroundImage"`
			BackgroundImageScaled    interface{}   `json:"backgroundImageScaled"`
			BackgroundTile           bool          `json:"backgroundTile"`
			BackgroundTopColor       string        `json:"backgroundTopColor"`
			CalendarFeedEnabled      bool          `json:"calendarFeedEnabled"`
			CanBeEnterprise          bool          `json:"canBeEnterprise"`
			CanBeOrg                 bool          `json:"canBeOrg"`
			CanBePrivate             bool          `json:"canBePrivate"`
			CanBePublic              bool          `json:"canBePublic"`
			CanInvite                bool          `json:"canInvite"`
			CardAging                string        `json:"cardAging"`
			CardCounts               bool          `json:"cardCounts"`
			CardCovers               bool          `json:"cardCovers"`
			Comments                 string        `json:"comments"`
			HiddenPluginBoardButtons []interface{} `json:"hiddenPluginBoardButtons"`
			HideVotes                bool          `json:"hideVotes"`
			Invitations              string        `json:"invitations"`
			IsTemplate               bool          `json:"isTemplate"`
			PermissionLevel          string        `json:"permissionLevel"`
			SelfJoin                 bool          `json:"selfJoin"`
			SharedSourceUrl          interface{}   `json:"sharedSourceUrl"`
			SwitcherViews            []struct {
				Enabled  bool   `json:"enabled"`
				ViewType string `json:"viewType"`
			} `json:"switcherViews"`
			Voting string `json:"voting"`
		} `json:"prefs"`
		ShortUrl string `json:"shortUrl"`
		Url      string `json:"url"`
	}
	NewCardRes struct {
		ID      string `json:"id"`
		Address string `json:"address"`
		Badges  struct {
			AttachmentsByType struct {
				Trello struct {
					Board int `json:"board"`
					Card  int `json:"card"`
				} `json:"trello"`
			} `json:"attachmentsByType"`
			Location           bool   `json:"location"`
			Votes              int    `json:"votes"`
			ViewingMemberVoted bool   `json:"viewingMemberVoted"`
			Subscribed         bool   `json:"subscribed"`
			Fogbugz            string `json:"fogbugz"`
			CheckItems         int    `json:"checkItems"`
			CheckItemsChecked  int    `json:"checkItemsChecked"`
			Comments           int    `json:"comments"`
			Attachments        int    `json:"attachments"`
			Description        bool   `json:"description"`
			Due                string `json:"due"`
			Start              string `json:"start"`
			DueComplete        bool   `json:"dueComplete"`
		} `json:"badges"`
		CheckItemStates  []string  `json:"checkItemStates"`
		Closed           bool      `json:"closed"`
		Coordinates      string    `json:"coordinates"`
		CreationMethod   string    `json:"creationMethod"`
		DateLastActivity time.Time `json:"dateLastActivity"`
		Desc             string    `json:"desc"`
		DescData         struct {
			Emoji struct {
			} `json:"emoji"`
		} `json:"descData"`
		Due          string `json:"due"`
		DueReminder  string `json:"dueReminder"`
		IDBoard      string `json:"idBoard"`
		IDChecklists []struct {
			ID string `json:"id"`
		} `json:"idChecklists"`
		IDLabels []struct {
			ID      string `json:"id"`
			IDBoard string `json:"idBoard"`
			Name    string `json:"name"`
			Color   string `json:"color"`
		} `json:"idLabels"`
		IDList         string   `json:"idList"`
		IDMembers      []string `json:"idMembers"`
		IDMembersVoted []string `json:"idMembersVoted"`
		IDShort        int      `json:"idShort"`
		Labels         []string `json:"labels"`
		Limits         struct {
			Attachments struct {
				PerBoard struct {
					Status    string `json:"status"`
					DisableAt int    `json:"disableAt"`
					WarnAt    int    `json:"warnAt"`
				} `json:"perBoard"`
			} `json:"attachments"`
		} `json:"limits"`
		LocationName          string `json:"locationName"`
		ManualCoverAttachment bool   `json:"manualCoverAttachment"`
		Name                  string `json:"name"`
		Pos                   int    `json:"pos"`
		ShortLink             string `json:"shortLink"`
		ShortURL              string `json:"shortUrl"`
		Subscribed            bool   `json:"subscribed"`
		URL                   string `json:"url"`
		Cover                 struct {
			Color                string `json:"color"`
			IDUploadedBackground bool   `json:"idUploadedBackground"`
			Size                 string `json:"size"`
			Brightness           string `json:"brightness"`
			IsTemplate           bool   `json:"isTemplate"`
		} `json:"cover"`
	}
)
