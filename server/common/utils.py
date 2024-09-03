import csv
import datetime
import time
import struct

""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574

TOTAL_LENGTH_BYTES = 2
AGENCIA_OFFSET = 1
DOCUMENTO_LENGTH = 8
NACIMIENTO_LENGTH = 10
NUMERO_LENGTH = 2

""" A lottery bet registry. """
class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: int):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)

""" Checks whether a bet won the prize or not. """
def has_won(bet: Bet) -> bool:
    return bet.number == LOTTERY_WINNER_NUMBER

"""
Persist the information of each bet in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def store_bets(bets: list[Bet]) -> None:
    with open(STORAGE_FILEPATH, 'a+') as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow([bet.agency, bet.first_name, bet.last_name,
                             bet.document, bet.birthdate, bet.number])

"""
Loads the information all the bets in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def load_bets() -> list[Bet]:
    with open(STORAGE_FILEPATH, 'r') as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])


"""
Function to write to a socket and ensure the written amount is as expected.
Made to avoid short writes.
"""
def write_exact(socket, data):
    sent_bytes = 0
    while sent_bytes < len(data):
        sent_bytes += socket.send(data[sent_bytes:])

"""
Function to read from a socket and ensure the read amount is as expected.
Made to avoid short reads.
"""
def read_exact(socket, length):
    data = bytearray()
    while len(data) < length:
        packet = socket.recv(length - len(data))
        if not packet:
            return None
        data.extend(packet)
    return data

"""
Function to deserialize a bet from a socket connection.
Returns the deserialized bet in the form of a Bet object.
"""
def deserialize_bet(socket):
    # Obtenemos el primer campo, con la totalidad del mensaje
    total_length_bytes = read_exact(socket, TOTAL_LENGTH_BYTES)
    total_length = struct.unpack('>H', total_length_bytes)[0]
    
    # Leemos el mensaje entero del buffer
    data = read_exact(socket, total_length)
    
    # Leemos la agencia
    agencia = data[0]
    offset = AGENCIA_OFFSET
    
    # Obtenemos la longitud del nombre y la usamos para deserializarlo
    nombre_len = data[offset]
    offset += 1
    nombre = data[offset:offset + nombre_len].decode('utf-8')
    offset += nombre_len

    # Obtenemos la longitud del apellido y la usamos para deserializarlo
    apellido_len = data[offset]
    offset += 1
    apellido = data[offset:offset + apellido_len].decode('utf-8')
    offset += apellido_len
        
    # Obtenemos el DNI, el cual es un u32.
    documento = data[offset:offset + DOCUMENTO_LENGTH].decode('utf-8')
    offset += DOCUMENTO_LENGTH
        
    # Obtenemos la fecha de nacimiento en formato de string
    nacimiento = data[offset:offset + NACIMIENTO_LENGTH].decode('utf-8')
    offset += NACIMIENTO_LENGTH
        
    # Obtenemos el numero, el cual es un u16
    numero = struct.unpack('>H', data[offset:offset + NUMERO_LENGTH])[0]
    offset += NUMERO_LENGTH
        
    return Bet(agencia, nombre, apellido, documento, nacimiento, numero)

def deserialize_bets(socket):
    amount_of_bets = read_exact(socket, 1)[0]
    bets = []
    for i in range(amount_of_bets):
        bet = deserialize_bet(socket)
        bets.append(bet)
    return bets