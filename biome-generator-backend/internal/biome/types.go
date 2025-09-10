package biome

// BiomeType 定义生物群系的枚举类型
type BiomeType int

// 定义所有生物群系的常量
const (
	Forest   BiomeType = iota // 森林
	Desert                    // 沙漠
	Ocean                     // 海洋
	Mountain                  // 山地
	Plains                    // 平原
)

// BiomeInfo 结构体存储生物群系的名称和对应的颜色
type BiomeInfo struct {
	Name  string `json:"name"`  // 生物群系名称
	Color string `json:"color"` // 用于可视化的颜色代码 (例如, CSS 颜色字符串)
}

// BiomeMap 全局映射，将 BiomeType 映射到 BiomeInfo
// 方便根据类型获取生物群系详细信息，并用于前端显示图例
var BiomeMap = map[BiomeType]BiomeInfo{
	Forest:   {Name: "森林", Color: "#228B22"}, // 绿色
	Desert:   {Name: "沙漠", Color: "#F0E68C"}, // 黄色
	Ocean:    {Name: "海洋", Color: "#4682B4"}, // 蓝色
	Mountain: {Name: "山地", Color: "#A9A9A9"}, // 灰色
	Plains:   {Name: "平原", Color: "#90EE90"}, // 浅绿色
}

// Block 结构体表示世界网格中的一个区块
type Block struct {
	X    int       `json:"x"`    // 区块的 X 坐标
	Y    int       `json:"y"`    // 区块的 Y 坐标
	Type BiomeType `json:"type"` // 该区块的生物群系类型
	Info BiomeInfo `json:"info"` // 该区块生物群系的详细信息 (名称和颜色)，方便前端直接使用
}

// Grid 类型定义 10x10 的二维 Block 数组，代表整个世界网格
type Grid [10][10]Block
