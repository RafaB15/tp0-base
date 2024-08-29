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

# Ejercicio 3

En el ejercicio 3 se escribió un script que levanta un container basándose en la imagen **busybox**, la cual es una imágne liviana que incluye *netcat*. 

Se utiliza la flag --rm para que se elimine el container una vez terminada la operación. Como network le pasamos el nombre de la network levantada por docker-compose para nuestro cliente y servidor. Obtenemos el nombre haciendo `docker network ls`.

El comando `sh -c "echo '$mensaje' | nc server 12345" 2>/dev/null` tiene muchas partes:
- ``sh -c``: Es para ejecutar un comando en la shell.
- ``"echo '$mensaje' | nc server 12345"``: Manda el mensaje al ip `server` (nos podemos referir así a la ip del servidor por estar dentro de la red) y puerto 12345 (definido en el archivo de configuración).
- ``2>/dev/null``: Redirijimos stderr para que en caso de que la red no se encuentre activa, se imprima solamente "action: test_echo_server | result: fail" en vez del error que nos lanza docker, para que los test automatizados no se vean perjudicados.

Finalmente, si el último comando corrió bien ($? -eq 0) y la respuesta es igual al mensaje mandado, entonces imprimimos "action: test_echo_server | result: success". En caso contrario, imprimimos "action: test_echo_server | result: fail".