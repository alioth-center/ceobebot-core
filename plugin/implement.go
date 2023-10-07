package plugin

type TextReplyMessageBaseImplementation struct{}

func (t TextReplyMessageBaseImplementation) Type() MessageCommandType {
	return TextReplyMessageCommandType
}

type ImageReplyMessageBaseImplementation struct{}

func (i ImageReplyMessageBaseImplementation) Type() MessageCommandType {
	return ImageReplyMessageCommandType
}
