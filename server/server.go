package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

/*
*Para que algo sea considerado como servidor, tendrá que tener algo llamado Config
*Que retorne una configuración (El struct Config)
* Config: tiene las caracteristicas del servidor.
* El puerto en el que se va ejecutar, la clave secreta para generar tokens y la conexion a base de datos.
! server: esta interface implementa el modelo de datos o estructura de config.
*/
type Server interface {
	Config() *Config
}

/*
* Broker que se encargara de manejar estos servidores
* broker va tener una propiedad config que va para tener un archivo de configuración
* con las propiedades definidas en el struct de config
* también un ruteador que va definir las rutas que va tener la API
* broker: Nos ayuda a tener varias instancias de servidor corriendo.
* Esta estructura a su vez tiene la estructura Config y el metodo Config, para ser de tipo Server.
 */

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return broker, nil
}

/*
* Metodo a broker que permita levantarse o ejecutarse
 */

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	log.Println("Starting server on port", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
