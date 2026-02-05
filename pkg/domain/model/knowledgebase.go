package model

type KnowledgeBaseCreateReq struct {
	Name                  string `json:"name" validate:"required,min=1,max=100"`
	ModelProvider         string `json:"model_provider" validate:"required,oneof=openai anthropic google azure cohere huggingface local"`
	Model                 string `json:"model" validate:"required,min=1,max=50"`
	VectorDimension       int    `json:"vector_dimension" validate:"required,min=1,max=10000"`
	MaxBatchDocumentCount int    `json:"max_batch_document_count" validate:"required,min=1,max=1000"`
}

type KnowledgeBaseUpdateReq struct {
	ID                    int    `json:"id" validate:"required"`
	Name                  string `json:"name" validate:"required,min=1,max=100"`
	ModelProvider         string `json:"model_provider" validate:"required,oneof=openai anthropic google azure cohere huggingface local"`
	Model                 string `json:"model" validate:"required,min=1,max=50"`
	VectorDimension       int    `json:"vector_dimension" validate:"required,min=1,max=10000"`
	MaxBatchDocumentCount int    `json:"max_batch_document_count" validate:"required,min=1,max=1000"`
}

type KnowledgeBaseQueryReq struct {
	Name          string `json:"name"`
	ModelProvider string `json:"model_provider" validate:"omitempty,oneof=openai anthropic google azure cohere huggingface local"`
	PageQuery
}

type KnowledgeBaseResp struct {
	ID                    int        `json:"id"`
	Name                  string     `json:"name"`
	ModelProvider         string     `json:"model_provider"`
	Model                 string     `json:"model"`
	VectorDimension       int        `json:"vector_dimension"`
	MaxBatchDocumentCount int        `json:"max_batch_document_count"`
	CreatedAt             *LocalTime `json:"created_at"`
	UpdatedAt             *LocalTime `json:"updated_at"`
}
