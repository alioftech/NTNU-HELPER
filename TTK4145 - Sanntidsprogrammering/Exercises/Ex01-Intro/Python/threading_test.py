from threading import Thread

x = 0

def add():
	global x
	for i in range(0, 1000000):
		x += 1

def sub():
	global x
	for i in range(0, 1000000):
		x -= 1

def startUp():
	add_thread = Thread(target = add)
	sub_thread = Thread(target = sub)

	add_thread.start()
	sub_thread.start()

	for i in range(0, 50):
		print("Value: " + str(x))

	add_thread.join()
	sub_thread.join()

	print("Done, final value: " + str(x))

startUp()