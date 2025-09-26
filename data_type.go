package daogext

import "strings"

// isMySQLTypeCompatibleWithGo 判断 MySQL 数据类型是否可映射到指定的 Go 类型
func isMySQLTypeCompatibleWithGo(mysqlType, goTypeName string) bool {
	// 统一转为小写便于比较
	mysqlType = strings.ToLower(mysqlType)

	// 去除无符号等修饰符，只保留基础类型部分（如：int(11) unsigned -> int）
	baseMySQLType := extractBaseMySQLType(mysqlType)

	// 定义映射关系：Go 类型 -> 兼容的 MySQL 类型列表
	compatible := map[string][]string{
		"int": {
			"tinyint", "smallint", "mediumint", "int", "integer",
		},
		"int64": {
			"bigint",
		},
		"uint": {
			"tinyint unsigned", "smallint unsigned", "mediumint unsigned", "int unsigned", "integer unsigned",
		},
		"uint64": {
			"bigint unsigned",
		},
		"float32": {
			"float",
		},
		"float64": {
			"double", "decimal", "dec", "numeric",
		},
		"decimal.Decimal": {
			"double", "decimal", "dec", "numeric",
		},
		"string": {
			"char", "varchar", "text", "tinytext", "mediumtext", "longtext",
			"enum", "set", "json",
		},
		"ttypes.NilableString": {
			"char", "varchar", "text", "tinytext", "mediumtext", "longtext",
			"enum", "set", "json",
		},
		"bool": {
			"tinyint", // 通常用 tinyint(1) 表示布尔
			"boolean", "bool",
		},
		"[]byte": {
			"blob", "tinyblob", "mediumblob", "longblob",
			"binary", "varbinary", "bit",
		},
		"time.time": {
			"datetime", "timestamp", "date", "time", "year",
		},
		"ttypes.NilableDatetime": {
			"datetime", "timestamp", "date", "time", "year",
		},
		"ttypes.NormalDatetime": {
			"datetime", "timestamp", "date", "time", "year",
		},
	}

	// 获取该 Go 类型支持的 MySQL 类型列表
	validTypes, exists := compatible[goTypeName]
	if !exists {
		return false // 不支持的 Go 类型
	}

	// 检查 baseMySQLType 是否在兼容列表中
	for _, t := range validTypes {
		if baseMySQLType == t {
			return true
		}
	}

	return false
}

// extractBaseMySQLType 提取基础 MySQL 类型（去掉长度、括号、unsigned 等）
// 例如: "int(11) unsigned" -> "int unsigned", "varchar(255)" -> "varchar"
func extractBaseMySQLType(t string) string {
	// 去掉括号及内部内容，例如 (11), (10,2)
	noParen := removeParentheses(t)

	// 分割并重新组合，保留 unsigned 等关键字
	parts := strings.Fields(noParen)
	if len(parts) == 0 {
		return ""
	}

	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		// 只保留有意义的词：类型名和 unsigned
		if trimmed == "unsigned" || trimmed == "zerofill" || trimmed == "auto_increment" {
			continue // 我们单独处理 unsigned，在主类型后添加
		}
		result = append(result, trimmed)
	}

	// 重新组合主类型
	base := strings.Join(result, " ")

	// 特殊处理：如果原字符串包含 unsigned，则附加
	hasUnsigned := strings.Contains(strings.ToLower(t), "unsigned")
	if hasUnsigned {
		if base == "tinyint" || base == "smallint" || base == "mediumint" || base == "int" || base == "integer" || base == "bigint" {
			return base + " unsigned"
		}
	}

	return base
}

// removeParentheses 去除字符串中的括号及其内容
func removeParentheses(s string) string {
	var result strings.Builder
	parenCount := 0
	for _, r := range s {
		if r == '(' {
			parenCount++
		} else if r == ')' {
			if parenCount > 0 {
				parenCount--
			}
		} else if parenCount == 0 {
			result.WriteRune(r)
		}
		// 在括号内则跳过
	}
	return result.String()
}
