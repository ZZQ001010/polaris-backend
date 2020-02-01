package po

var WhiteColor = "#FFFFFF"

var ProjectType = map[string]PpmPrsProjectType{
	"normal_task": ProjectTypeNormalTask,
	"agile":       ProjectTypeAgile,
}

var ProjectTypeNormalTask = PpmPrsProjectType{
	LangCode: "ProjectType.NormalTask",
	Name:     "普通任务项目",
	Sort:     1,
}

var ProjectTypeAgile = PpmPrsProjectType{
	LangCode: "ProjectType.Agile",
	Name:     "敏捷研发项目",
	Sort:     2,
}

var ProjectObjectType = map[string]PpmPrsProjectObjectType{
	"iteration": ProjectObjectTypeIteration,
	"feature":   ProjectObjectTypeFeature,
	"demand":    ProjectObjectTypeDemand,
	"task":      ProjectObjectTypeTask,
	"bug":       ProjectObjectTypeBug,
	"test_task": ProjectObjectTypeTestTask,
}

var ProjectObjectTypeIteration = PpmPrsProjectObjectType{
	LangCode:   "Project.ObjectType.Iteration",
	PreCode:    "I",
	Name:       "迭代",
	ObjectType: 1,
	Sort:       1,
	IsReadonly: 1,
}

var ProjectObjectTypeFeature = PpmPrsProjectObjectType{
	LangCode:   "Project.ObjectType.Feature",
	PreCode:    "F",
	Name:       "特性",
	ObjectType: 2,
	Sort:       2,
	IsReadonly: 1,
}

var ProjectObjectTypeDemand = PpmPrsProjectObjectType{
	LangCode:   "Project.ObjectType.Demand",
	PreCode:    "D",
	Name:       "需求",
	ObjectType: 2,
	Sort:       3,
	IsReadonly: 1,
}

var ProjectObjectTypeTask = PpmPrsProjectObjectType{
	LangCode:   "Project.ObjectType.Task",
	PreCode:    "T",
	Name:       "任务",
	ObjectType: 2,
	Sort:       4,
	IsReadonly: 1,
}

var ProjectObjectTypeBug = PpmPrsProjectObjectType{
	LangCode:   "Project.ObjectType.Bug",
	PreCode:    "B",
	Name:       "缺陷",
	ObjectType: 2,
	Sort:       5,
	IsReadonly: 1,
}

var ProjectObjectTypeTestTask = PpmPrsProjectObjectType{
	LangCode:   "Project.ObjectType.TestTask",
	PreCode:    "TT",
	Name:       "测试任务",
	ObjectType: 2,
	Sort:       6,
	IsReadonly: 1,
}
