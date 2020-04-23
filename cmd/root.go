package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var AppName = "go-calendar"

var cfgFile string

// RootCmd корневая команда Cobra, родитель всех дочерних команд
var RootCmd = &cobra.Command{
	Use:   AppName,
	Short: "Another calendar service",
	Long:  `Another calendar service`,
}

// Execute выполняет главную команду приложения.
// В зависимости от параметров запуска буду выполнены подчиненные команды, либо выведена справка
func Execute() int {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

// init инициализирует пакет команд
func init() {
	cobra.OnInitialize(initConfig)

	// Здесь определяются флаги конфигурации
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+AppName+".yaml)")

	// Cobra also supports local flags, which will only run when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig функция-инициализатор, которая вызовется перед вызовом Execute() команды
func initConfig() {
	if cfgFile != "" {
		// Используем конифигурационный файл из флага
		viper.SetConfigFile(cfgFile)
	} else {
		// Список директорий, в которых будет производиться поиск конфигурации
		viper.AddConfigPath("$HOME/." + AppName) // в домашней директории как .AppName.yaml
		viper.AddConfigPath("./configs")         // в директории configs относительно текущей
		//viper.AddConfigPath(".")  // в текущей директории
	}

	// Добавляем реплейсер точек на подчеркивание для имен переменных окружения
	// т.к. с точками в env работать не очень удобно :)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // Считывает переменных окружения, для всех ключей в конфиге

	// Если найден файл конфигурации используем его настройки
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
