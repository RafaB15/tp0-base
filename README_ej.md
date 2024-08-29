# Ejercicio 1

En el ejercicio 1 se hizo un archivo generar-compose.sh, al cual se le tiene que pasar por parámetro el nombre dle archivo que se va a escribir y la cantidad de clientes que se requiere

```bash
./generar-compose.sh docker-compose-dev.yml 5
```

En este archivo se llama a un script de go, en el cual se concatenan strings formateados para generar a los clientes.

Se cambia la variable de entorno PYTHONUNBUFFERED para que sea igual a la cantidad de clientes, de manera que los clientes que no estén siendo atendidos se queden en el buffer y sus mensajes no se pierdan.

A cada cliente se le pone el nombre **clientn** y se le cambia el id a **n**.

# Ejercicio 2

En el ejercicio 2, con el objetivo de que cambios en los archivos de configuración del cliente y del servidor no desembocaran en la recontrucción de la imágen, se montaron volúmenes para que la información de los archivos *config.ini* y *config.yaml* pueda ser accedida por el container directamente. 

También se agregaron estos archivos en el dockerignore, para que no fuera copiados dentro del container al momento de crear la imágen con COPY, sino que fueran montados en el docker-compose-dev.yaml.