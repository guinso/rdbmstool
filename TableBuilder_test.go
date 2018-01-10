package rdbmstool

import "testing"
import "strings"

func TestTableBuilder_SQL(t *testing.T) {
	builder := NewTableBuilder()

	expectedSQL := "CREATE TABLE `member`(\n" +
		"`id` int(10) NOT NULL,\n" +
		"`name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`score` float NULL,\n" +
		"`range` decimal(10,2) NOT NULL,\n" +
		"`join_on` date NOT NULL,\n" +
		"`last_acccess` datetime NULL,\n" +
		"`remark` text COLLATE utf8mb4_unicode_ci NULL,\n" +
		"`is_vip` tinyint(1) NOT NULL,\n" +
		"`class` char(3) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"PRIMARY KEY(`id`)\n" +
		") ENGINE=innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

	SQLStr, SQLErr := builder.TableName("member").
		AddColumnInt("id", 10, false).
		AddColumnVarchar("name", 100, false).
		AddColumnFloat("score", true).
		AddColumnDecimal("range", 10, 2, false).
		AddColumnDate("join_on", false).
		AddColumnDateTime("last_acccess", true).
		AddColumnText("remark", true).
		AddColumnBoolean("is_vip", false).
		AddColumnChar("class", 3, false).
		AddPrimaryKey("id").
		SQL()

	if SQLErr != nil {
		t.Error(SQLErr)
	}

	if strings.Compare(expectedSQL, SQLStr) != 0 {
		t.Errorf("Expect:\n%s\n\nbut get:\n\n%s", expectedSQL, SQLStr)
	}
}
