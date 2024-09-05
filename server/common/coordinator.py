from common.utils import *

class Coordinator:
    def send_winners(self, expected_clients, barrier, queue):
        barrier.wait()
        winners = {}
        for i in range(1, expected_clients + 1):
            winners[i] = []
        for bet in load_bets():
            if has_won(bet):
                winners[bet.agency].append(bet)
        for i in range(expected_clients):
            queue.put(winners)