## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

## Solucion

La solución propuesta levanta un contenedor temporal de `busybox` en la misma red de Docker que el servidor (tp0_testing_net). Esto es debido a que incluye `nc` y por lo tanto no requerimos de netcat por separado, además de que permite una conexión interna en Docker para no exponer los puertos del servidor. Por medio de un pipe y con `nc`, enviamos un mensaje al servidor (utilizando su direccion y puerto del archivo de configuracion) y esperamos que la respuesta sea la misma. Esto es validado en la parte final del script.

Ejecución:
```
./validar-echo-server.sh
```