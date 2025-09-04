import socket
import logging
import signal
import os
from common.utils import Bet, store_bets, load_bets, has_won

CLIENT_COUNT = int(os.getenv("CLIENT_COUNT"))

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._clients_ready = {}

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
            len_bytes = self.__read_all(client_sock, 4)
            len_bytes = int.from_bytes(len_bytes, "big")

            # Client ready for draw
            if len_bytes == 0:
                id_agency_bytes = self.__read_all(client_sock, 4)
                id_agency = int.from_bytes(id_agency_bytes, "big")
                self.__handle_client_ready(client_sock, id_agency)
                return

            msg = self.__read_all(client_sock, len_bytes).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')

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

    def __handle_client_ready(self, client_sock, agency_id):
        # Add this client
        client_address = client_sock.getpeername()
        self._clients_ready[agency_id] = [client_sock, client_address]
        logging.info(f'action: new_client_ready | result: success | cant: {len(self._clients_ready)} | clients: {self._clients_ready.keys()}')

        if len(self._clients_ready) == CLIENT_COUNT:
            logging.info('action: sorteo | result: success')
            self.__handle_draw()

    def __handle_draw(self):
        bets = load_bets()
        winning_bets = {}
        winning_bets_count = 0
        for bet in bets:
            if has_won(bet):
                winning_bets_count += 1
                if bet.agency in winning_bets:
                    winning_bets[bet.agency].append(bet.document)
                else:
                    winning_bets[bet.agency] = [bet.document]
        logging.info(f'action: sorteo_ganadores | result: success | cant_ganadores: {winning_bets_count}')
        
        for client_key in self._clients_ready.keys():
            if client_key in winning_bets:
                msg = "|".join(winning_bets[client_key])
            else:
                msg = ""
            self.__write_all(self._clients_ready[client_key][0], (msg + "\n").encode('utf-8'))
            self._clients_ready[client_key][0].close()

    def __read_all(self, client_sock, len_bytes):
        buffer = bytearray()
        while len(buffer) < len_bytes:
            partial_read = client_sock.recv(len_bytes - len(buffer))
            buffer.extend(partial_read)
        return bytes(buffer)
    
    def __write_all(self, client_sock, buffer):
        bytes_written = 0
        while bytes_written < len(buffer):
            n = client_sock.send(buffer[bytes_written:])
            bytes_written += n
        return bytes_written