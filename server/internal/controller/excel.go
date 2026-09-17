package controller

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// columnLetter 将索引转换为 Excel 列字母（0->A, 26->AA）
func columnLetter(idx int) string {
	s := ""
	for idx >= 0 {
		s = string(rune('A'+idx%26)) + s
		idx = idx/26 - 1
		if idx < 0 {
			break
		}
	}
	return s
}

// exportExcel 通用 Excel 导出
func exportExcel(c *gin.Context, filename string, headers []string, rows [][]interface{}) {
	f := excelize.NewFile()
	const sheet = "Sheet1"
	for i, h := range headers {
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s1", columnLetter(i)), h)
	}
	for r, row := range rows {
		for ci, val := range row {
			_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", columnLetter(ci), r+2), val)
		}
	}
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		returnToClientError(c, "导出失败: "+err.Error())
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func returnToClientError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{"code": 400, "msg": msg, "data": nil})
}
