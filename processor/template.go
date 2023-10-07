package processor

import "fmt"

var (
	NotMatchedCommandMessageTemplate        = "很抱歉，没有找到这个命令，这个命令可能不存在或者仍在测试中，命令：%s，参数：%s"
	IncorrectCommandArgumentMessageTemplate = "很抱歉，您的命令输入有误，您使用的命令可能是：\n%s"
	CommandHelpArgumentMessageTemplate      = "功能%s，参数%s的命令使用提示如下：\n%s"
	CommandArgumentInfoTemplate             = "%d. %s: %s\n描述：%s\n\n使用示例：\n%s\n"
	ListHelpInfoMessageTemplate             = "功能%s的命令如下：\n%s"
	HelpInfoTemplate                        = "%d. %s: %s"
	ListAllServiceInfoTemplate              = "所有功能如下：\n%s"
	ServiceInfoTemplate                     = "%d. %s"
)

func FormatNotMatchedCommandMessage(command, content string) string {
	return fmt.Sprintf(NotMatchedCommandMessageTemplate, command, content)
}

func FormatCommandArgumentInfo(index int, info HandlerInfo) string {
	return fmt.Sprintf(CommandArgumentInfoTemplate, index, info.CommandKey, info.CommandName, info.CommandDescription, info.CommandExample)
}

func FormatIncorrectCommandArgumentMessage(commandInfos []HandlerInfo) string {
	if commandInfos == nil {
		return fmt.Sprintf(IncorrectCommandArgumentMessageTemplate, "没有找到结果")
	} else if len(commandInfos) == 0 {
		return fmt.Sprintf(IncorrectCommandArgumentMessageTemplate, "没有找到结果")
	}

	infoString := ""
	for i, info := range commandInfos {
		infoString = fmt.Sprintf("%s\n%s", infoString, FormatCommandArgumentInfo(i+1, info))
	}

	return fmt.Sprintf(IncorrectCommandArgumentMessageTemplate, infoString)
}

func FormatCommandHelpArgumentMessage(command, content string, commandInfos []HandlerInfo) string {
	if commandInfos == nil {
		return fmt.Sprintf(CommandHelpArgumentMessageTemplate, command, content, "没有找到结果")
	} else if len(commandInfos) == 0 {
		return fmt.Sprintf(CommandHelpArgumentMessageTemplate, command, content, "没有找到结果")
	}

	infoString := ""
	for i, info := range commandInfos {
		infoString = fmt.Sprintf("%s\n%s", infoString, FormatCommandArgumentInfo(i+1, info))
	}

	return fmt.Sprintf(CommandHelpArgumentMessageTemplate, command, content, infoString)
}

func FormatHelpInfo(index int, info HandlerInfo) string {
	return fmt.Sprintf(HelpInfoTemplate, index, info.CommandName, info.CommandDescription)
}

func FormatListHelpInfoMessage(command string, commandInfos []HandlerInfo) string {
	if commandInfos == nil {
		return fmt.Sprintf(ListHelpInfoMessageTemplate, command, "该功能下没有命令")
	} else if len(commandInfos) == 0 {
		return fmt.Sprintf(ListHelpInfoMessageTemplate, command, "该功能下没有命令")
	}

	infoString := ""
	for i, info := range commandInfos {
		infoString = fmt.Sprintf("%s\n%s", infoString, FormatHelpInfo(i+1, info))
	}

	return fmt.Sprintf(ListHelpInfoMessageTemplate, command, infoString)
}

func FormatServiceInfo(index int, service string) string {
	return fmt.Sprintf(ServiceInfoTemplate, index, service)
}

func FormatListAllServiceInfoMessage() string {
	if services := defaultMatcher.ListAllCommands(); services == nil {
		return fmt.Sprintf(ListAllServiceInfoTemplate, "没有找到结果")
	}

	infoString := ""
	for i, service := range defaultMatcher.ListAllCommands() {
		infoString = fmt.Sprintf("%s\n%s", infoString, FormatServiceInfo(i+1, service))
	}

	return fmt.Sprintf(ListAllServiceInfoTemplate, infoString)
}
