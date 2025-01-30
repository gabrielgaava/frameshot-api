variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "aws_access_key" {
  description = "AWS aws_access_key"
  type        = string
}

variable "aws_secret_key" {
  description = "AWS aws_secret_key"
  type        = string
}

variable "aws_session_token" {
  description = "AWS aws_session_token"
  type        = string
}

variable "cognito_jwks_url" {
    description = "Coginito JWKS Url"
    type = string
}

variable "ecr_image_url" {
  description = "URL da imagem no ECR"
  type        = string
}

variable "ecs_task_execution_role_arn" {
  description = "ARN da role de execução do ECS"
  type        = string
}

variable "vpc_id" {
  description = "ID da VPC"
  type        = string
}

variable "public_subnets" {
  description = "Lista de public subnets"
  type        = list(string)
}

variable "private_subnets" {
  description = "Lista de private subnets"
  type        = list(string)
}

variable "security_groups" {
  description = "Lista de security groups"
  type        = list(string)
}

variable "db_connection" {
  description = "String de conexão do banco de dados"
  type        = string
}

variable "db_host" {
  description = "Host do banco de dados"
  type        = string
}

variable "db_port" {
  description = "Porta do banco de dados"
  type        = string
}

variable "db_name" {
  description = "Nome do banco de dados"
  type        = string
}

variable "db_user" {
  description = "Usuário do banco de dados"
  type        = string
}

variable "db_password" {
  description = "Senha do banco de dados"
  type        = string
}

variable "aws_bucket_name" {
  description = "Nome do bucket S3"
  type        = string
}

variable "sendgrid_api_key" {
  description = "Chave da API do SendGrid"
  type        = string
}

variable "sendgrid_template_id" {
    description = "Template utilizado para email"
    type = string
}

variable "s3_queue_url" {
    description = "URL da fila do S3"
    type = string
}

variable "video_input_queue_url" {
    description = "URL da fila para processamento do video"
    type = string
}

variable "video_output_queue_url" {
    description = "URL da fila de retorno do processamento do video"
    type = string
}