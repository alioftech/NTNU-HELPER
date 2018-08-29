from threading import Thread
from threading import Lock

mtx = Lock()

x = 0

def add():
	global x
	for i in range(0, 1000000):
		mtx.acquire()
		x += 1
		mtx.release()

def sub():
	global x
	for i in range(0, 1000000):
		mtx.acquire()
		x -= 1
		mtx.release()

def startUp():
	add_thread = Thread(target = add)
	sub_thread = Thread(target = sub)

	add_thread.start()
	sub_thread.start()

	for i in range(0, 50):
		mtx.acquire()
		print("Value: " + str(x))
		mtx.release()

	add_thread.join()
	sub_thread.join()

	print("Done, final value: " + str(x))

startUp()