// Code generated by "stringer -type=TokenType TokenType.go"; DO NOT EDIT.

package lexer

import "strconv"

const _TokenType_name = "TokErrorTokNoEmitTokWhitespaceTokCharTokStringTokNumberTokDotTokElipsisTokOperTokNamespaceAccessTokOperatorStartTokStarTokPlusTokMinusTokDivTokExpTokLTTokLTETokGTTokGTETokOperatorEndTokSemiColonTokDefereferenceTokReferenceTokAssignmentTokEqualityTokRightParenTokLeftParenTokRightCurlyTokLeftCurlyTokRightBraceTokLeftBraceTokRightArrowTokLeftArrowTokCompoundAssignmentTokForTokWhileTokIfTokElseTokReturnTokFuncDefnTokClassDefnTokNamespaceTokDependencyTokTypeTokCommaTokIdentTokComment"

var _TokenType_index = [...]uint16{0, 8, 17, 30, 37, 46, 55, 61, 71, 78, 96, 112, 119, 126, 134, 140, 146, 151, 157, 162, 168, 182, 194, 210, 222, 235, 246, 259, 271, 284, 296, 309, 321, 334, 346, 367, 373, 381, 386, 393, 402, 413, 425, 437, 450, 457, 465, 473, 483}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
