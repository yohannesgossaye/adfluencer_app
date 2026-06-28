package influencer

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
)

type ListContractsResponse struct {
	Items      []contractDto.ListItem `json:"items"`
	Pagination commonDto.ListMeta     `json:"pagination"`
}

type RespondContractResponse struct {
	Contract contractDto.DetailResponse `json:"contract"`
}
