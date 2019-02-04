package protocol3

//验证部分的一些常量
const (
	severity         = 'S' //严重性：该域的内容是ERROR、FATAL或者PANIC（在一个错误消息里），或者WARNING、NOTICE、DEBUG、INFO或者LOG（在一条通知消息里），或者是这些形式的某种本地化翻译。总是会出现。
	severityV        = 'V' //严重性：该域的内容是ERROR、FATAL或者PANIC（在一个错误消息里），或者WARNING、NOTICE、DEBUG、INFO或者LOG（在一条通知消息里）。这和S域相同，不过其内容没有被本地化。这只存在于PostgreSQL版本 9.6 及其后版本产生的消息中。
	code             = 'C' //代码：错误的SQLSTATE代码（参阅附录 A）。非本地化。总是出现。
	message          = 'M' //消息：人类可读的错误消息的主体。这些信息应该准确并且简洁（通常是一行）。总是出现。
	detail           = 'D' //细节：一个可选的二级错误消息，携带了有关问题的更多错误消息。可以是多行。
	hint             = 'H' //提示：一个可选的有关如何处理问题的建议。它和细节不同的地方是它提出了建议（可能并不合适）而不仅仅是事实。可以是多行。
	position         = 'P' //位置：这个域值是一个十进制ASCII整数，表示一个错误游标的位置，它是一个指向原始查询字符串的索引。第一个字符的索引是 1，位置是以字符计算而非字节计算的。
	internalPosition = 'p' //内部位置：这个域和P域定义相同，但是它被用于当游标位置指向一个内部生成的命令的情况， 而不是用于客户端提交的命令。这个域出现的时候，总是会出现q域。
	internalQuery    = 'q' //内部查询：失败的内部生成的命令的文本。比如，它可能是一个PL/pgSQL函数发出的 SQL 查询。
	where            = 'W' //哪里：一个指示错误发生的环境的指示器。目前，它包含一个活跃的过程语言函数的调用堆栈的路径和内部生成的查询。 这个路径每个项记录一行，最新的在最前面。
	schemaName       = 's' //模式名：如果错误与一个指定数据库对象相关，这里是包含该对象的模式名（如果有）。
	tableName        = 't' //表名：如果错误与一个指定表相关，这里是表的名字（引用该表模式的名字的模式名域）。
	columnName       = 'c' //列名：如果错误与一个指定表列相关，这里是该列的名字（引用该模式和表的名字来标识该表）。
	dataTypeName     = 'd' //数据类型名：如果错误与一个指定数据类型相关，这里是该数据类型的名字（引用该数据类型模式的名字的模式名域）。
	constraintName   = 'n' //约束名：如果错误是和一个指定约束相关，这里是该约束的名字。引用至上面列出的相关表或域的域（为了这个目的，索引被视作约束，即使它们并不是按照约束语法创建的）。
	file             = 'F' //文件：报告的错误在源代码中的文件名。
	line             = 'L' //行：报告的错误所在的源代码的位置的行号。
	routine          = 'R' //例程：报告错误的例程在源代码中的名字。

)

var codePrefix = []byte{0, 0x43}
var bodyPrefix = []byte{0, 0x4D}

type ErrorResponse struct {
	backMessage
	Key  byte
	Code string
	Body string
}

func (e *ErrorResponse) Decode(data []byte) {
	e.Key = data[5]
	var ind int

	switch e.Key {
	case 'S':
		e.Code, ind = takeContent(data, codePrefix)
		e.Body, ind = takeContent(data[ind:], bodyPrefix)
	default:
		e.Code, ind = takeContent(data, codePrefix)
		e.Body, ind = takeContent(data[ind:], bodyPrefix)
		log.Error().Msg("未处理的验证方式。未操作。直接返回")
	}

}
