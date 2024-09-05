import socket
import logging
import signal
from common.utils import *
import struct

BETS_MESSAGE = 1
WINNERS_REQUEST_MESSAGE = 2

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._received_sigterm = False
        self._expected_clients = listen_backlog
        self._waiting_clients = {}

        signal.signal(signal.SIGTERM, self.shutdown)
        
    def run(self):
        """
        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while (not self._received_sigterm) and (len(self._waiting_clients) < self._expected_clients):
            try:
                client_sock = self.__accept_new_connection()
                message_type = self.__get_message_type(client_sock)
                if message_type == BETS_MESSAGE:
                    self.__handle_bets_message(client_sock)
                else:
                    self.__handle_winners_request_message(client_sock)
            except OSError as e:
                logging.error(f"action: run_server | result: fail | error: {e}")
                return
        
        logging.info("action: sorteo | result: success")
        
        try:
            self.__send_winners()
        except OSError as e:
            logging.error(f"action: run_server | result: fail | error: {e}")
            return
        
    def __get_message_type(self, client_sock):
        return read_exact(client_sock, 1)[0]


    def __handle_bets_message(self, client_sock):
        """:
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bets = deserialize_bets(client_sock)
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')
            store_bets(bets)
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
            self.__sendConfirmation(client_sock, len(bets))
        except DeserializationError as e:
            logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
            raise e
        finally:
            client_sock.close()

    def __sendConfirmation(self, socket, amount_read):
        confirmation = struct.pack('>B', amount_read)
        write_exact(socket, confirmation)
    
    def __handle_winners_request_message(self, client_sock):
        """
        Read a byte from the socket and handle the winners request message
        """

        try:
            agencia = read_exact(client_sock, 1)[0]
            self._waiting_clients[agencia] = client_sock
            return agencia
        except OSError as e:
            client_sock.close()
            raise e
    
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

    def __send_winners(self):
        winners = {}
        for i in range(1, self._expected_clients + 1):
            winners[i] = []
        for bet in load_bets():
            if has_won(bet):
                winners[bet.agency].append(bet)
        for agency, socket in self._waiting_clients.items():
            logging.debug(f'action: send winners | result: in progress | agency: {agency}')
            self.__send_winners_to_agency(socket, winners[agency])
            logging.debug(f'action: send winners | result: success | agency: {agency}')
           

    def __send_winners_to_agency(self, agency_socket, bets):
        # Primero escribbimos la cantidad de documentos ganadores como un u16
        length = len(bets)
        agency_socket.send(struct.pack('>H', length))

        # Escribimos cada documento en utf8
        for bet in bets:
            document = bet.document.encode('utf8')
            write_exact(agency_socket, document)
    
    def shutdown(self, signum, frame):
        """
        Función que se llama al recibir la señal SIGTERM
        Se cierra el socket para de
        """
        logging.info('action: shutdown server | result: in_progress')
        self._received_sigterm = True
        self._server_socket.close()
        logging.info('action: shutdown server | result: success')
