## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).

## Solucion

El objetivo de este ejercicio es implementar la logica de envio y recepcion de apuestas entre cliente y servidor. Para ello, tuvimos que diseñar nuestro **propio protocolo**:
```
len"id_agency#name#surname#dni#birth_date#bet\n"
```
Donde `len` es un Big Endian de 4 bytes que nos indica el largo del mensaje (para saber hasta donde leer), seguido por los componentes de la Bet separados por #, hasta el `\n` que indica el final del mensaje. El pasaje a bytes del mensaje con este protocolo lo podemos observar en la función `to_bytes()` del cliente (**bet_msg.go**).

Entonces, para la comunicación entre cliente y servidor, primero se leen los primeros 4 bytes del mensaje para obtener la longitud del mismo, y luego el mensaje como tal. Enviar la longitud del mensaje fue necesario para poder implementar las funciones write_all y read_all, que voy a explicar a continuación y que tienen como objetivo **evitar los short read y short write**.

Para **evitar short write en el cliente**, simplemente repetimos la operación Write con un for (puede ocurrir que no se envie el mensaje completo en un solo Write), hasta garantizar que enviemos la totalidad del mensaje.
```
func write_all(buffer []byte) {
    bytes_sent := 0
    for bytes_sent < len(buffer) {
        n, err := c.conn.Write(buffer[bytes_sent:])
        bytes_sent += n
    }
```

Del lado del servidor tenemos nuestra equivalente write_all en Python para evitar los short write, pero también necesitamos una función read_all para **evitar los short read**:
```
def __read_all(length_bytes):
        while len(buffer) < length_bytes:
            partial_read = client_sock.recv(length_bytes - len(buffer))
            buffer.extend(partial_read)
        return bytes(buffer)
```
Para asegurar que `client_sock.recv` lea la totalidad del mensaje, utilizamos el largo del mensaje recibido al comienzo del mismo, y ciclamos con while hasta recibir el mensaje completo. Por eso, en cada lectura parcial, leemos `length_bytes - len(buffer)`, es decir **el total menos lo que 'ya tengo'**.