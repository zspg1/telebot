package telebot

type Pool struct {
	size  int
	deep  int
	stop  chan int
	queue []chan Update
}

func NewPool(size int, deep int) *Pool {
	if size <= 0 {
		size = 10
	}
	if deep <= 0 {
		deep = 100
	}

	t := make([]chan Update, size)
	for i := 0; i < size; i++ {
		t[i] = make(chan Update, deep)
	}

	return &Pool{
		size:  size,
		queue: t,
	}
}

func (b *Bot) startPool() {
	for i := 0; i < b.pool.size; i++ {
		go b.worker(b.pool.queue[i])
	}

}
func (b *Bot) worker(c chan Update) {
	for {
		select {
		case u := <-c:
			b.ProcessUpdate(u)
		}
	}
}

func (b *Bot) Submit(update Update) {
	if update.Message == nil || update.Message.Chat == nil {
		return
	}
	id := update.Message.Chat.ID
	if id < 0 {
		id = -id
	}
	idx := int(id)
	c := b.pool.queue[idx%b.pool.size]
	c <- update
}
