package processor

var (
	defaultMatcher = NewHandlerMatcher()
)

func DefaultMatcher() HandlerMatcher {
	return defaultMatcher
}

type HandlerInfo struct {
	CommandKey         string
	CommandName        string
	CommandDescription string
	CommandExample     string
}

type HandlerMatcher interface {
	MatchHandlers(command string, content string) (mustHandlers []MustHandleFunction, optionalHandlers []OptionalHandleFunction)
	SetDefaultHandlerFunctions(mustHandlers []MustHandleFunction, optionalHandlers []OptionalHandleFunction)
	RegisterCommand(commandKey, commandName, commandDescription, commandExample string, contentFilter func(string) bool, mustHandlers []MustHandleFunction, optionalHandlers []OptionalHandleFunction)
	ListCommands(commandKey string) (commandInfoMap map[string]HandlerInfo)
	ListAllCommands() (commands []string)
	GetHelpInfo(commandKey, commandContent string) (commandInfo []HandlerInfo)
}

type handlerMatcherFilter struct {
	contentFilter    func(string) bool
	mustHandlers     []MustHandleFunction
	optionalHandlers []OptionalHandleFunction
	description      string
	example          string
}
type defaultHandlerMatcher struct {
	filters                 map[string]map[string]handlerMatcherFilter
	defaultMustHandlers     []MustHandleFunction
	defaultOptionalHandlers []OptionalHandleFunction
}

func (m *defaultHandlerMatcher) MatchHandlers(command string, content string) (mustHandlers []MustHandleFunction, optionalHandlers []OptionalHandleFunction) {
	if commandMap, existCommand := m.filters[command]; !existCommand {
		// 不存在对应的命令
		return m.defaultMustHandlers, m.defaultOptionalHandlers
	} else {
		for _, filter := range commandMap {
			if filter.contentFilter(content) {
				mustHandlers = append(mustHandlers, filter.mustHandlers...)
				optionalHandlers = append(optionalHandlers, filter.optionalHandlers...)
			}
		}

		// 没有找到对应的handler
		if len(mustHandlers) == 0 && len(optionalHandlers) == 0 {
			return m.defaultMustHandlers, m.defaultOptionalHandlers
		}

		// 找到了对应的handler
		return mustHandlers, optionalHandlers
	}
}

func (m *defaultHandlerMatcher) SetDefaultHandlerFunctions(mustHandlers []MustHandleFunction, optionalHandlers []OptionalHandleFunction) {
	if mustHandlers == nil {
		mustHandlers = []MustHandleFunction{}
	}
	for _, handler := range mustHandlers {
		if handler != nil {
			m.defaultMustHandlers = append(m.defaultMustHandlers, handler)
		}
	}

	if optionalHandlers == nil {
		optionalHandlers = []OptionalHandleFunction{}
	}
	for _, handler := range optionalHandlers {
		if handler != nil {
			m.defaultOptionalHandlers = append(m.defaultOptionalHandlers, handler)
		}
	}
}

func (m *defaultHandlerMatcher) RegisterCommand(commandKey, commandName, commandDescription, commandExample string, contentFilter func(string) bool, mustHandlers []MustHandleFunction, optionalHandlers []OptionalHandleFunction) {
	if m.filters == nil {
		m.filters = map[string]map[string]handlerMatcherFilter{}
	}
	if m.filters[commandKey] == nil {
		m.filters[commandKey] = map[string]handlerMatcherFilter{}
	}

	if mustHandlers == nil {
		mustHandlers = []MustHandleFunction{}
	}
	if optionalHandlers == nil {
		optionalHandlers = []OptionalHandleFunction{}
	}

	m.filters[commandKey][commandName] = handlerMatcherFilter{
		contentFilter:    contentFilter,
		mustHandlers:     mustHandlers,
		optionalHandlers: optionalHandlers,
		description:      commandDescription,
		example:          commandExample,
	}
}

func (m *defaultHandlerMatcher) ListCommands(commandKey string) (commandInfoMap map[string]HandlerInfo) {
	if m.filters == nil {
		return map[string]HandlerInfo{}
	} else if filter, existCommand := m.filters[commandKey]; !existCommand {
		return map[string]HandlerInfo{}
	} else {
		infoMap := map[string]HandlerInfo{}
		for commandName, filter := range filter {
			infoMap[commandName] = HandlerInfo{
				CommandKey:         commandKey,
				CommandName:        commandName,
				CommandDescription: filter.description,
				CommandExample:     filter.example,
			}
		}

		return infoMap
	}
}

func (m *defaultHandlerMatcher) GetHelpInfo(commandKey, commandContent string) (commandInfo []HandlerInfo) {
	if m.filters == nil {
		return []HandlerInfo{}
	} else if filter, existCommand := m.filters[commandKey]; !existCommand {
		return []HandlerInfo{}
	} else {
		for commandName, matcherFilter := range filter {
			if matcherFilter.contentFilter(commandContent) {
				commandInfo = append(commandInfo, HandlerInfo{
					CommandKey:         commandKey,
					CommandName:        commandName,
					CommandDescription: matcherFilter.description,
					CommandExample:     matcherFilter.example,
				})
			}
		}

		return commandInfo
	}
}

func (m *defaultHandlerMatcher) ListAllCommands() (commands []string) {
	if m.filters == nil {
		return []string{}
	}

	for command := range m.filters {
		commands = append(commands, command)
	}

	return commands
}

func NewHandlerMatcher() HandlerMatcher {
	return &defaultHandlerMatcher{
		filters:                 map[string]map[string]handlerMatcherFilter{},
		defaultMustHandlers:     []MustHandleFunction{},
		defaultOptionalHandlers: []OptionalHandleFunction{},
	}
}
