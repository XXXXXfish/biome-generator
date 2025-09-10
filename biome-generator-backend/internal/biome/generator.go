package biome

import (
	"math/rand"
	"time"
	"log"
)

// GenerationParameters 结构体用于接收前端传入的气候参数
type GenerationParameters struct {
	MoistureSpread    int `json:"moisture_spread"`    // 西侧相邻区块的气候影响强度
	TemperatureSpread int `json:"temperature_spread"` // 北侧相邻区块的气候影响强度
	ClimateStability  int `json:"climate_stability"`  // 气候稳定区的影响强度
}

// GenerateWorld 根据给定的参数生成一个 10x10 的生物群系网格
func GenerateWorld(params GenerationParameters) Grid {
	// 创建一个新的随机数生成器实例，使用当前时间作为种子，确保每次生成结果不同
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	var grid Grid // 声明一个 10x10 的网格

	// 1. 生成起始区块 (0,0)
	grid[0][0].X = 0
	grid[0][0].Y = 0
	// 随机选择一个生物群系作为起始区块的类型
	grid[0][0].Type = BiomeType(rng.Intn(len(BiomeMap)))
	grid[0][0].Info = BiomeMap[grid[0][0].Type] // 获取对应的生物群系信息
	log.Printf("Starting biome at (0,0): %s", BiomeMap[grid[0][0].Type].Name)
	// 2. 依次生成其他区块
	// 遍历顺序：先行再列，确保在生成 (x,y) 时，(x-1,y) 和 (x,y-1) 已经生成
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if x == 0 && y == 0 {
				continue // (0,0) 区块已在上面处理，跳过
			}

			// 为当前区块的生物群系类型初始化分数
			scores := make(map[BiomeType]float64)
			for i := BiomeType(0); i < BiomeType(len(BiomeMap)); i++ {
				scores[i] = 20.0 // 基础生成概率：每种生物群系自然出现率为 20%，我们用分数 20 来代表
			}

			var westBiome, northBiome BiomeType // 存储西侧和北侧区块的生物群系类型
			hasWest := x > 0                    // 判断是否有西侧相邻区块
			hasNorth := y > 0                   // 判断是否有北侧相邻区块

			// 应用西侧相邻区块的气候影响
			if hasWest {
				westBiome = grid[x-1][y].Type
				// 将 MoistureSpread 的值加到西侧生物群系的分数上
				scores[westBiome] += float64(params.MoistureSpread)
			}
			// 应用北侧相邻区块的气候影响
			if hasNorth {
				northBiome = grid[x][y-1].Type
				// 将 TemperatureSpread 的值加到北侧生物群系的分数上
				scores[northBiome] += float64(params.TemperatureSpread)
			}

			// 判断是否形成“气候稳定区”并应用其影响
			if hasWest && hasNorth && westBiome == northBiome {
				// 如果西侧和北侧生物群系相同，额外增加 ClimateStability 的值
				scores[westBiome] += float64(params.ClimateStability)
			}

			// --- 新增的调试日志开始 ---
            // log.Printf("Block (%d,%d) raw scores: %+v", x, y, scores) // 如果输出太多，可以注释掉这行
            
			totalScore := 0.0
			for _, score := range scores {
				totalScore += score
			}
            // 确保 totalScore 不为 0，避免除以零
            if totalScore == 0 {
                // 这种情况理论上不会发生，因为基础分至少是 20 * 5 = 100
                // 但作为防御性编程，如果发生，可以设置为默认某种生物群系或报错
                log.Printf("ERROR: totalScore is 0 for block (%d,%d)", x, y)
                grid[x][y].Type = Forest // fallback
                grid[x][y].Info = BiomeMap[Forest]
                continue
            }


			cumulativeProbs := make([]float64, len(BiomeMap))
			currentCumulativeProb := 0.0
			
            // 记录一下每个生物群系的归一化概率
            biomeProbabilities := make(map[BiomeType]float64) // 用于临时存储归一化概率
			for i := BiomeType(0); i < BiomeType(len(BiomeMap)); i++ {
                prob := scores[i] / totalScore
                biomeProbabilities[i] = prob // 存储归一化概率
				currentCumulativeProb += prob
				cumulativeProbs[i] = currentCumulativeProb
			}

            // 打印每个区块的概率分布，这样可以直观看到参数的影响
            log.Printf("Block (%d,%d) probabilities (M:%d,T:%d,C:%d):",
                x, y, params.MoistureSpread, params.TemperatureSpread, params.ClimateStability)
            for i := BiomeType(0); i < BiomeType(len(BiomeMap)); i++ {
                log.Printf("  %s: %.2f%%", BiomeMap[i].Name, biomeProbabilities[i]*100)
            }
            // --- 新增的调试日志结束 ---

			// 使用随机数选择最终的生物群系
			r := rng.Float64() // 生成一个 0 到 1 之间的随机浮点数
			chosenBiome := BiomeType(0) // 默认选择第一个生物群系
			for i := BiomeType(0); i < BiomeType(len(BiomeMap)); i++ {
				if r < cumulativeProbs[i] { // 如果随机数小于当前生物群系的累积概率，则选择它
					chosenBiome = i
					break
				}
			}

			// 设置当前区块的属性
			grid[x][y].X = x
			grid[x][y].Y = y
			grid[x][y].Type = chosenBiome
			grid[x][y].Info = BiomeMap[chosenBiome] // 获取并设置生物群系信息
		}
	}
	return grid // 返回生成的完整世界网格
}

