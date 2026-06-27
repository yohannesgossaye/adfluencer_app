package influencer

import contractDto "adinfluence_service_app/internal/constant/model/dto/contract"

type ListContractsResponse struct {
	Items      []contractDto.ListItem         `json:"items"`
	Pagination contractDto.PaginationResponse `json:"pagination"`
}

type RespondContractResponse struct {
	Contract contractDto.DetailResponse `json:"contract"`
}
