package service

type ChannelWriter struct {
	channel chan string
}

func (c *ChannelWriter) Write(p []byte) (n int, err error) {
	c.channel <- string(p)
	return len(p), nil
}
