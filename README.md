## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

## Solucion

Para agregar la lógica de enviar lotes o **batches de apuestas**, tuvimos que realizar pequeñas modificaciones a nuestro protocolo, lo que demuestra la robustez del mismo para escalarlo. Seguimos enviando los 4 bytes iniciales para indicar la longitud del mensaje, pero ahora en nuestro mensaje vamos a separar cada apuesta con un `\n`:
```
len"id_agency#name1#surname1#dni1#birth_date1#bet1\nid_agency#name2#surname2#dni2#birth_date2#bet2\n..."
```
Por lo tanto, desde el lado del cliente, la lógica no cambia mucho. Ahora en vez de enviar una apuesta por mensaje, enviamos el máximo que nos permita `MaxBetsPerBatch`. Este valor decidí definirlo en **135**, aunque es configurable en el `config.yaml` del servidor. Dado que la apuesta más larga que pude encontrar en los archivos de prueba era de 59 bytes, en el peor de los casos podría enviar 135 apuestas de esa longitud, ya que 59*135 nos da **7965** que está justo por debajo de los 8kb que la cátedra nos sugiere no superar.

Nuestro servidor, por otro lado, tuvo que ser modificado levemente para poder leer cada apuesta individualmente de los batches que recibe:
```
msg = self.__read_all(client_sock, len_bytes).rstrip().decode('utf-8')
bets_msg = msg.split("\n")
for i, bet in enumerate(bets_msg):
    bet_data = bet.split("#")
    bet = Bet(*bet_data)
    bets_msg[i] = bet
store_bets(bets_msg)
logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets_msg)}')
```
Nótese como nuestro protocolo splitea las apuestas del batch por `\n`, y cada apuesta individual por `#`. Simplemente hubo que agregar un **for** para procesar cada apuesta del batch.
Vale destacar que seguimos utilizando nuestras funciones `read_all()` y `write_all()` para evitar short read y short write.

Para la lectura de archivos tuvimos que **modificar nuestro script de Python** para que los inyecte en los containers correspondientes de cada cliente, y se persistan por fuera de la imagen. Para eso modificamos nuestro volumes de cada cliente:
```
"volumes": [
    "./client/config.yaml:/config.yaml",
    f"./.data/agency-{i}.csv:/agency-{i}.csv",
]
```
Y para cargar estos archivos en memoria, en nuestro archivo `bet_msg_file.go` definimos la estructura **BetLoader**, que tiene la capacidad de leer y cargar en memoria solamente el batch de apuestas que le toca en cada iteración. Es decir, **no carga todo el archivo en memoria**, va leyendo lo que necesita.
Esto es posible gracias a las funciones **LoadBets()**, que crea un lector CSV y mantiene un "cursor" para saber cuanto leyó del archivo, y **LoadBetBatch()** que solamente carga en memoria el batch actual:
```
func LoadBetBatch() ([]byte, uint32, bool) {
    var buffer bytes.Buffer
    eof := false
    lineCount := 0

    for lineCount < br.max_bets_per_batch {
        record, err := br.reader.Read()
        // ...chequeo EOF...
        // ...validaciones...
        bet := Bet{ ... }
        buffer.Write(bet.to_bytes())
        lineCount++
    }
    return buffer.Bytes(), uint32(buffer.Len()), ...
```