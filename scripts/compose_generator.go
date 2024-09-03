package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	nombre_de_archivo := flag.String("nombre", "docker-compose-dev.yaml", "El nombre del archivo en el que se escribirá el docker compose")
	cantidad_de_clientes := flag.Int("clientes", 1, "La cantidad de clientes que se conectarán al servidor")

	flag.Parse()

	compose := `name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini

`

	for i := 1; i <= *cantidad_de_clientes; i++ {
		nombre_cliente := fmt.Sprintf("client%d", i)
		compose += fmt.Sprintf(`  %s:
    container_name: %s
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=%d
      - AGENCIA=%d
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/dataset/agency-%d.csv:/app/agency.csv
    depends_on:
      - server

`, nombre_cliente, nombre_cliente, i, i, i)
	}

	compose += `networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
`

	err := os.WriteFile(*nombre_de_archivo, []byte(compose), 0644)
	if err != nil {
		fmt.Printf("Error al escribir el archivo: %v\n", err)
		return
	}
}
