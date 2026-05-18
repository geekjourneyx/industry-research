package trace

type TracePlan struct {
	Question string  `json:"question"`
	Domain   string  `json:"domain"`
	Claims   []Claim `json:"claims"`
}

type Claim struct {
	ClaimID             string          `json:"claim_id"`
	Claim               string          `json:"claim"`
	Mechanism           string          `json:"mechanism"`
	ExpectedTraces      []ExpectedTrace `json:"expected_traces"`
	SourceFamilies      []string        `json:"source_families"`
	DisconfirmingTraces []string        `json:"disconfirming_traces"`
}

type ExpectedTrace struct {
	TraceType   string `json:"trace_type"`
	Trace       string `json:"trace"`
	WhyExpected string `json:"why_expected"`
}

func BuildChainBrandTracePlan(question string) TracePlan {
	return TracePlan{
		Question: question,
		Domain:   "chain-brand",
		Claims: []Claim{
			{
				ClaimID:   "claim_store_count_growth",
				Claim:     "门店数增长或扩张目标具备经营支撑",
				Mechanism: "门店增长需要选址、招聘、供应链、数字入口和用户需求共同支撑。",
				ExpectedTraces: []ExpectedTrace{
					{TraceType: "people_org", Trace: "目标城市出现店长、店员、区域运营、拓展岗位", WhyExpected: "门店扩张前后必须补充门店和区域运营人员。"},
					{TraceType: "digital_frontend", Trace: "小程序、App、外卖平台或门店列表出现可服务门店", WhyExpected: "真实门店必须被用户发现、选择或下单。"},
					{TraceType: "physical_fulfillment", Trace: "地图 POI、点评、外卖门店页或本地开业信息出现", WhyExpected: "真实运营门店会留下可定位和可评价的经营痕迹。"},
					{TraceType: "capital_legal", Trace: "直营网点、加盟主体、分支机构或许可信息出现", WhyExpected: "经营主体和合规经营通常会留下工商或许可痕迹。"},
					{TraceType: "management_narrative", Trace: "财报、公告、管理层访谈、公众号或权威媒体披露扩张计划", WhyExpected: "上市或准上市连锁品牌通常会公开解释扩张节奏和经营口径。"},
				},
				SourceFamilies: []string{"recruiting", "map_poi", "platform_frontend", "company_registry", "company_disclosure", "media_interview", "ugc"},
				DisconfirmingTraces: []string{
					"声称覆盖城市但无门店 POI",
					"无招聘或仅总部招聘",
					"小程序或外卖平台不可下单",
					"只有通稿转载，没有独立经营痕迹",
					"用户评价长期停滞或集中反映闭店",
				},
			},
		},
	}
}
