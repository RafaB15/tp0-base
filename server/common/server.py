import socket
import logging
import signal
from common.utils import *
from common.worker import Worker
from common.coordinator import Coordinator
import multiprocessing

BETS_MESSAGE = 1
WINNERS_REQUEST_MESSAGE = 2

class Server:
    def __init__(self, port, listen_backlog, expected_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._received_sigterm = False
        self._expected_clients = expected_clients
        self._processHandlers = []

        signal.signal(signal.SIGTERM, self.shutdown)
        
    def run(self):
        """
        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # Hacemos una barrera para sincronizar a los manejadores de clientes con el encargado de hacer el sorteo
        barrier = multiprocessing.Barrier(self._expected_clients + 1)
        
        # Hacemos un lock con el que sincronizaremos las llamadas a la función de store bets
        lock = multiprocessing.Lock()
        
        # Hacemos una queue con la que mandaremos información de los resultados del sorteo a los diferentes procesos
        queue = multiprocessing.Queue()

        coordinator_process = multiprocessing.Process(target=coordinate_clients, args=(self._expected_clients, barrier, queue))
        coordinator_process.start()

        current_clients = 0
        while (not self._received_sigterm) and (current_clients < self._expected_clients):
            try:
                client_sock = self.__accept_new_connection()
                current_clients += 1
                p = multiprocessing.Process(target=handle_client, args=(client_sock, barrier, lock, queue))
                p.start()
                self._processHandlers.append(p)
            except OSError as e:
                logging.error(f"action: run_server | result: fail | error: {e}")
                return
        
        logging.info("action: sorteo | result: success")
        
        coordinator_process.join()
        for p in self._processHandlers:
            p.join()
    
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
    
    def shutdown(self, signum, frame):
        """
        Función que se llama al recibir la señal SIGTERM
        Se cierra el socket para de
        """
        logging.info('action: shutdown server | result: in_progress')
        self._received_sigterm = True
        self._server_socket.close()
        logging.info('action: shutdown server | result: success')

def handle_client(client_sock, barrier, lock, queue):
    worker = Worker()
    worker.handle_client(client_sock, barrier, lock, queue)
    
def coordinate_clients(expected_clients, barrier, queue):
    coordinator = Coordinator()
    coordinator.send_winners(expected_clients, barrier, queue)