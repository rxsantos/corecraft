package services

import "math"

// CalculateMempoolStats calcula estatísticas detalhadas da mempool do Bitcoin.
// Recebe os dados do getmempoolinfo e getrawmempool(verbose=true).
func CalculateMempoolStats(info map[string]interface{}, raw map[string]interface{}) map[string]interface{} {

	// Tratamento inicial: se não houver transações
	if raw == nil || len(raw) == 0 {
		return map[string]interface{}{
			"result": "Mempool vazia",
		}
	}

	var totalVsize float64 		// soma de todos os vsizes (tamanho virtual)
	var totalFeeRate float64	// soma de todas as fee rates

	minFee := 999999.0
	maxFee := 0.0

	low := 0					// fee rate < 10 sat/vB
	medium := 0					// 10 <= fee rate <= 50 sat/vB
	high := 0					// fee rate > 50 sat/vB

	count := 0.0

	// Itera sobre cada transação na mempool
	// raw tem o formato: map[txid] -> dados da transação
	for _, tx := range raw {

		txData, ok := tx.(map[string]interface{})
		if !ok {
			continue	// ignora entradas inválidas
		}

		// Extrai vsize (tamanho virtual em vbytes)
		vsize, vok := txData["vsize"].(float64)

		// Extrai as taxas (fees)
		feesMap, fok := txData["fees"].(map[string]interface{})
		if !vok || !fok || vsize == 0 {
			continue
		}

		fee, ok := feesMap["base"].(float64)
		if !ok || fee == 0 {
			continue
		}
		// Calcula fee rate em satoshis por vbyte (sat/vB)
		feeRate := (fee / vsize) * 1e8	//100 000 0000 
		

		// Acumula totais
		totalVsize += vsize
		totalFeeRate += feeRate

		// Atualiza mínimo e máximo
		if feeRate < minFee {
			minFee = feeRate
		}
		if feeRate > maxFee {
			maxFee = feeRate
		}

		// Classifica por faixa de fee rate
		switch {
		case feeRate < 10:
			low++
		case feeRate <= 50:
			medium++
		default:
			high++
		}

		count++
	}

	// Se nenhuma transação válida foi processada
	if count == 0 {
		return map[string]interface{}{
			"error": "nenhuma transação válida",
		}
	}

	// Formatação para ficar mais legível
	avgFeeRate := totalFeeRate / float64(count)


	// Retorna estatísticas calculadas
	return map[string]interface{}{
		"tx_count":     count,
		"total_vsize":  totalVsize,
		// "avg_fee_rate": totalFeeRate / count,
		// "min_fee_rate": minFee,
		// "max_fee_rate": maxFee,
	"avg_fee_rate": roundTo(avgFeeRate, 0),   // 0 casas decimais
		"min_fee_rate": roundTo(minFee, 0),
		"max_fee_rate": roundTo(maxFee, 0),		
		// Quantidades por faixa de fee rate
		"fee_distribution": map[string]int{
			"low":    low,			// < 10 sat/vB
			"medium": medium,		// 10 a 50 sat/vB
			"high":   high,			// > 50 sat/vB
		},
	}
}

// roundTo arredonda um float64 para N casas decimais
func roundTo(value float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Round(value*shift) / shift
}