package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/todo"
	"github.com/todo/pkg/handler"
	"github.com/todo/pkg/queue"
	"github.com/todo/pkg/repository"
	"github.com/todo/pkg/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title TodoItem App API
// @version 1.0
// @description API Server for TodoItem Application

// @host localhost:3333
// @BasePath /

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	appMode := os.Getenv("APP_MODE")

	if appMode == "" {
		appMode = "app"
	}

	switch appMode {
	case "consumer":
		initConsumer()
	case "app":
	default:
		initApp()
	}

	initApp()
}

func initApp() {
	db, err := repository.NewMySQLDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	q, err := queue.NewSQSQueue(queue.Config{
		Options: session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("AWS_ACCESS_KEY_ID"),
					os.Getenv("AWS_SECRET_ACCESS_KEY"),
					"",
				),
				Endpoint: aws.String(viper.GetString("queue.host")),
				Region:   aws.String(viper.GetString("queue.region")),
			},
		},
		QueueName: viper.GetString("queue.name"),
		ConsumerConfig: queue.ConsumerConfig{
			MaxNumberOfMessages:      viper.GetInt64("queue.consumer.max_messages"),
			MessageVisibilityTimeout: viper.GetInt64("queue.consumer.visibility_timeout"),
			PollDelayInMilliseconds:  viper.GetInt("queue.consumer.poll_delay"),
			Receivers:                viper.GetInt("queue.consumer.receivers"),
		},
	})

	if err != nil {
		log.Fatalf("failed to initialize queue: %s", err.Error())
	}

	queues := queue.NewQueue(q)
	repositories := repository.NewRepository(db)
	services := service.NewService(repositories, queues)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Print("Application started and running")

	scheduleStopper := initSchedule(services)

	log.Print("Schedule started and running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Application shutting down")

	scheduleStopper <- true

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConsumer() {
	q, err := queue.NewSQSQueue(queue.Config{
		Options: session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("AWS_ACCESS_KEY_ID"),
					os.Getenv("AWS_SECRET_ACCESS_KEY"),
					"",
				),
				Endpoint: aws.String(viper.GetString("queue.host")),
				Region:   aws.String(viper.GetString("queue.region")),
			},
		},
		QueueName: viper.GetString("queue.name"),
		ConsumerConfig: queue.ConsumerConfig{
			MaxNumberOfMessages:      viper.GetInt64("queue.consumer.max_messages"),
			MessageVisibilityTimeout: viper.GetInt64("queue.consumer.visibility_timeout"),
			PollDelayInMilliseconds:  viper.GetInt("queue.consumer.poll_delay"),
			Receivers:                viper.GetInt("queue.consumer.receivers"),
		},
	})

	if err != nil {
		log.Fatalf("failed to initialize queue: %s", err.Error())
	}

	queues := queue.NewQueue(q)

	queues.Consumer.Consume(queues.Consumer.DefaultHandler)

	log.Print("Application consumer started and running")
}

func initConfig() error {
	viper.AddConfigPath("conf")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func initSchedule(services *service.Service) chan bool {
	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				f := faker.New()
				todoItem := todo.TodoItem{
					Description: f.Lorem().Sentence(5),
					DueDate:     time.Now(),
				}

				if _, err := services.TodoItem.Create(todoItem); err != nil {
					log.Errorf("error occured on creating todoitem: %s", err.Error())
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}
