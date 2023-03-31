package main

import (
	"fmt"
	"path/filepath"

	"github.com/bluele/slack"
)

func sendMessage(message string, channel string) {
	api := slack.New(configuration.SlackTockenAPI)
	err := api.ChatPostMessage(channel, message, &slack.ChatPostMessageOpt{
		/*Username : configuration.SlackUserName,
		IconUrl :configuration.SlackAvatar,*/
		AsUser: true,
	})
	if err != nil {
		panic(err)
	}
}

func sendFile(message string, fileName string, fileType string, channel string) {

	//Inicializamos la api de slack
	api := slack.New(configuration.SlackTockenAPI)

	//Configuramos y enviamos el archivo
	info, err := api.FilesUpload(&slack.FilesUploadOpt{
		//Ruta del archivo
		Filepath: fileName,
		//Tipo de archivo
		Filetype: fileType,
		//Nombre del Archivo
		Filename: filepath.Base(fileName),
		//Titulo del archivo
		Title: message,
		//Canales a donde se envia
		Channels: []string{channel},
		//Mensaje opcional
		Content: message,
	})

	//Si pasa algo para el programa y muestra el error
	if err != nil {
		panic(err)
	}

	//Indicamos que se envio el archivo satisfactoriamente
	fmt.Println(fmt.Sprintf("Completed file upload with the ID: '%s'.", info.ID))
}
