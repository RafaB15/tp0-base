from common.utils import *
import logging

BETS_MESSAGE = 1
WINNERS_REQUEST_MESSAGE = 2

class Worker:
    def handle_client(self, client_sock, barrier, lock, queue):
        while True:
            try:
                message_type = self.__get_message_type(client_sock)
                if message_type == BETS_MESSAGE:
                    self.__handle_bets_message(client_sock, lock)
                else:
                    self.__handle_winners_request_message(client_sock, barrier, queue)
                    break
            except OSError as e:
                    logging.error(f"action: run_worker | result: fail | error: {e}")
                    return
            
    def __get_message_type(self, client_sock):
        return read_exact(client_sock, 1)[0]
    
    def __handle_bets_message(self, client_sock, lock):
        try:
            bets = deserialize_bets(client_sock)
            addr = client_sock.getpeername()
            
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')
            self.__store_bets_secure(bets, lock)
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
            
            self.__send_confirmation(client_sock, len(bets))
        except DeserializationError as e:
            logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
            raise e
        
    def __store_bets_secure(self, bets, lock):
        with lock:
            store_bets(bets)
            
    def __send_confirmation(self, socket, amount_read):
        confirmation = struct.pack('>B', amount_read)
        write_exact(socket, confirmation)
        
    def __handle_winners_request_message(self, client_sock, barrier, queue):
        """
        Read a byte from the socket and handle the winners request message
        """

        try:
            agency = read_exact(client_sock, 1)[0]
            barrier.wait()
            winner_bets = queue.get()[agency]
            self.__send_winners(client_sock, winner_bets)
        except OSError as e:
            client_sock.close()
            raise e
    
    def __send_winners(self, client_socket, winner_bets):
        # Primero escribbimos la cantidad de documentos ganadores como un u16
        length = len(winner_bets)
        data = struct.pack('>H', length)
        write_exact(client_socket, data)

        # Escribimos cada documento en utf8
        for bet in winner_bets:
            document = bet.document.encode('utf8')
            write_exact(client_socket, document)