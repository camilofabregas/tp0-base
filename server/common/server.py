import socket
import logging
import signal
from common.utils import Bet, store_bets

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        # Set graceful shutdown flag
        self._shutdown = False
        signal.signal(signal.SIGTERM, self.__handle_sigterm)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while not self._shutdown:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError as e:
                if not self._shutdown:
                    logging.error(f'action: accept_connection | result: fail | error: {e}')

        logging.info("action: server_shutdown | result: success")

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            len = self.__read_all(client_sock, 4)
            len = int.from_bytes(len, "big")

            msg = self.__read_all(client_sock, len).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')

            #bet_msg = msg.split('#')
            #bet = Bet(*bet_msg)
            #store_bets([bet])
            #logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')
            try:
                bets_msg = msg.split("\n")
                for i, bet in enumerate(bets_msg):
                    bet_data = bet.split("#")
                    bet = Bet(*bet_data)
                    bets_msg[i] = bet
                store_bets(bets_msg)
                logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets_msg)}')

                response = "ACK BATCH\n".encode('utf-8')
                self.__write_all(client_sock, response)
                logging.info('action: bet_acknowledged | result: success')
            
            except:
                logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets_msg)}')
                response = "ERR BATCH\n".encode('utf-8')
                self.__write_all(client_sock, response)
                logging.info('action: error_sent | result: success')

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def __handle_sigterm(self, signum, frame):
        logging.info(f'action: server_shutdown | result: in_progress | signal: {signum}')
        self._shutdown = True
        self._server_socket.close()
        logging.info("action: close_server_socket | result: success")

    def __read_all(self, client_sock, length_bytes):
        buffer = bytearray()
        while len(buffer) < length_bytes:
            partial_read = client_sock.recv(length_bytes - len(buffer))
            buffer.extend(partial_read)
        return bytes(buffer)
    
    def __write_all(self, client_sock, buffer):
        bytes_written = 0
        while bytes_written < len(buffer):
            n = client_sock.send(buffer[bytes_written:])
            bytes_written += n
        return bytes_written