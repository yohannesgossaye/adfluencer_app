package helper

import (
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"math"
)

func ComputeContractFees(amountETB float64) contractDto.FeeBreakdown {
	businessFee := roundETB(amountETB * contractDto.BusinessFeePercent / 100)
	influencerFee := roundETB(amountETB * contractDto.InfluencerFeePercent / 100)

	return contractDto.FeeBreakdown{
		AmountETB:               amountETB,
		BusinessFeePercent:      contractDto.BusinessFeePercent,
		InfluencerFeePercent:    contractDto.InfluencerFeePercent,
		BusinessFeeAmountETB:    businessFee,
		InfluencerFeeAmountETB:  influencerFee,
		TotalBusinessPayableETB: roundETB(amountETB + businessFee),
		InfluencerReceivableETB: roundETB(amountETB - influencerFee),
		PlatformTotalFeeETB:     roundETB(businessFee + influencerFee),
	}
}

func roundETB(v float64) float64 {
	return math.Round(v*100) / 100
}
