package po

var WhiteColor = "#FFFFFF"

var Process = map[string]PpmPrsProcess{
	"default_project":    ProcessDefaultProject,
	"default_iteration":  ProcessDefaultIteration,
	"default_feature":    ProcessDefaultFeature,
	"default_demand":     ProcessDefaultDemand,
	"default_task":       ProcessDefaultTask,
	"default_agile_task": ProcessDefaultAgileTask,
	"default_bug":        ProcessDefaultBug,
	"default_test_task":  ProcessDefaultTestTask,
}

var ProcessDefaultProject = PpmPrsProcess{
	LangCode:  "Process.DefaultProject",
	Name:      "默认项目流程",
	IsDefault: 1,
	Type:      1,
	Sort:      1,
}

var ProcessDefaultIteration = PpmPrsProcess{
	LangCode:  "Process.DefaultIteration",
	Name:      "默认迭代流程",
	IsDefault: 1,
	Type:      2,
	Sort:      1,
}

var ProcessDefaultFeature = PpmPrsProcess{
	LangCode:  "Process.Issue.DefaultFeature",
	Name:      "默认特性流程",
	IsDefault: 0,
	Type:      3,
	Sort:      1,
}

var ProcessDefaultDemand = PpmPrsProcess{
	LangCode:  "Process.Issue.DefaultDemand",
	Name:      "默认需求流程",
	IsDefault: 0,
	Type:      3,
	Sort:      2,
}

var ProcessDefaultTask = PpmPrsProcess{
	LangCode:  "Process.Issue.DefaultTask",
	Name:      "默认任务流程",
	IsDefault: 1,
	Type:      3,
	Sort:      3,
}

var ProcessDefaultAgileTask = PpmPrsProcess{
	LangCode:  "Process.Issue.DefaultAgileTask",
	Name:      "默认敏捷项目任务流程",
	IsDefault: 1,
	Type:      3,
	Sort:      4,
}

var ProcessDefaultBug = PpmPrsProcess{
	LangCode:  "Process.Issue.DefaultBug",
	Name:      "默认缺陷流程",
	IsDefault: 0,
	Type:      3,
	Sort:      5,
}

var ProcessDefaultTestTask = PpmPrsProcess{
	LangCode:  "Process.Issue.DefaultTestTask",
	Name:      "默认测试任务流程",
	IsDefault: 0,
	Type:      3,
	Sort:      6,
}
