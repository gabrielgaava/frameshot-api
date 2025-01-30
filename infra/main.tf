provider "aws" {
  region = var.aws_region
}

resource "aws_ecs_cluster" "frameshot_cluster" {
  name = "frameshot-cluster"
}

resource "aws_ecs_task_definition" "frameshot-api-task" {
  family                   = "frameshot-api-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn            = var.ecs_task_execution_role_arn

  container_definitions = jsonencode([
    {
      name      = "frameshot-api"
      image     = var.ecr_image_url
      essential = true
      environment = [
        { name = "APP_NAME", value = "frameshot-api" },
        { name = "APP_ENV", value = "production" },
        { name = "HTTP_URL", value = "127.0.0.1" },
        { name = "HTTP_PORT", value = "8080" },
        { name = "HTTP_ALLOWED_ORIGINS", value = "*" },
        { name = "DB_CONNECTION", value = var.db_connection },
        { name = "DB_HOST", value = var.db_host },
        { name = "DB_PORT", value = var.db_port },
        { name = "DB_NAME", value = var.db_name },
        { name = "DB_USER", value = var.db_user },
        { name = "DB_PASSWORD", value = var.db_password },
        { name = "AWS_REGION", value = var.aws_region },
        { name = "AWS_ACCESS_KEY_ID", value = var.aws_access_key },
        { name = "AWS_SECRET_ACCESS_KEY", value = var.aws_secret_key },
        { name = "AWS_SESSION_TOKEN", value = var.aws_session_token },
        { name = "AWS_BUCKET_NAME", value = var.aws_bucket_name },
        { name = "AWS_COGNITO_JWKS_URL", value = var.cognito_jwks_url },
        { name = "AWS_S3_QUEUE_URL", value = var.s3_queue_url },
        { name = "AWS_VIDEO_INPUT_QUEUE_URL", value = var.video_input_queue_url },
        { name = "AWS_VIDEO_OUTPUT_QUEUE_URL", value = var.video_output_queue_url },
        { name = "SENDGRID_API_KEY", value = var.sendgrid_api_key },
        { name = "SENDGRID_TEMPLATE_ID", value = var.sendgrid_template_id }
      ]
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = "/ecs/frameshot-api"
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}

resource "aws_ecs_service" "frameshot_service" {
  name            = "frameshot-api"
  cluster         = aws_ecs_cluster.frameshot_cluster.id
  task_definition = aws_ecs_task_definition.frameshot-api-task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = var.subnets
    security_groups = var.security_groups
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.frameshot_tg.arn
    container_name   = "frameshot-api"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.frameshot_listener]
}

resource "aws_lb_target_group" "frameshot_tg" {
  name     = "frameshot-api-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = var.vpc_id
  target_type = "ip"

  health_check {
    path                = "/healthcheck"
    port                = 8080
    protocol            = "HTTP"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    matcher             = "200"
  }
}

resource "aws_lb" "frameshot_lb" {
  name               = "frameshot-api-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = var.security_groups
  subnets            = var.subnets
}

resource "aws_lb_listener" "frameshot_listener" {
  load_balancer_arn = aws_lb.frameshot_lb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.frameshot_tg.arn
  }
}