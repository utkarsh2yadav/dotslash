package service

import (
	"bufio"
	"context"
	. "dotslash/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"unicode/utf8"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gofiber/websocket/v2"
)

func (l *Language) Handler(c *websocket.Conn) {

	defer c.Close()

	var (
		msg []byte
		err error
	)

	errorResponse := []byte("Server Error")

	if _, msg, err = c.ReadMessage(); err != nil {
		log.Println(err)
		return

	}

	body := new(WsBody)
	if err := json.Unmarshal(msg, body); err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, errorResponse)
		return
	}

	directory, err := os.MkdirTemp("", "")
	if err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, errorResponse)
		return
	}
	defer os.RemoveAll(directory)

	if err := os.WriteFile(fmt.Sprintf("%v/%v.%v", directory, l.Filename, l.getExtension()), []byte(body.Code), 0664); err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, errorResponse)
		return
	}

	var dockerClient *client.Client
	dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, errorResponse)
		return
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var resp container.ContainerCreateCreatedBody
	resp, err = dockerClient.ContainerCreate(timeoutContext, &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          true,
		Image:        l.Name,
		WorkingDir:   "/work",
		Cmd:          l.getCommand(body),
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: directory,
				Target: "/work",
			},
		},
	}, nil, nil, "")

	if err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, errorResponse)
		return
	}

	defer func(dockerClient *client.Client) {
		if err := dockerClient.ContainerRemove(timeoutContext, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			log.Println(err)
		}
	}(dockerClient)

	if err = dockerClient.ContainerStart(timeoutContext, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, errorResponse)
		return
	}

	statusCh, errorCh := dockerClient.ContainerWait(timeoutContext, resp.ID, container.WaitConditionNotRunning)

	var containerResp types.HijackedResponse
	containerResp, err = dockerClient.ContainerAttach(timeoutContext, resp.ID, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		log.Println(err)
		c.WriteMessage(websocket.TextMessage, errorResponse)
		return
	}
	defer containerResp.Close()

	bufin := bufio.NewReader(containerResp.Reader)
	inout := make(chan []byte)
	output := make(chan []byte)
	errChan := make(chan error)

	// Write to docker container
	go func(w io.WriteCloser) {
		for {
			data, ok := <-inout
			if !ok {
				w.Close()
				return
			}
			w.Write(append(data, '\n'))
		}
	}(containerResp.Conn)

	// Receive from docker container
	go func() {
		for {
			buffer := make([]byte, 4096, 4096)
			c, err := bufin.Read(buffer)
			if err != nil {
				errChan <- err
				break
			}
			if c > 0 {
				output <- buffer[:c]
			}
			if c == 0 {
				output <- []byte{' '}
			}
			if err != nil {
				break
			}
		}
	}()

	// Connect STDOUT to websocket
	go func() {
		for {
			data, ok := <-output
			if !ok {
				break
			}
			stringData := string(data[:])
			if !utf8.ValidString(stringData) {
				v := make([]rune, 0, len(stringData))
				for i, r := range stringData {
					if r == utf8.RuneError {
						_, size := utf8.DecodeRuneInString(stringData[i:])
						if size == 1 {
							continue
						}
					}
					v = append(v, r)
				}
				stringData = string(v)
			}
			if err := c.WriteMessage(websocket.TextMessage, []byte(stringData)); err != nil {
				errChan <- err
				break
			}
		}
	}()

	// Websocket to docker
	go func(c *websocket.Conn) {
		for {
			if _, msg, err := c.ReadMessage(); err != nil {
				log.Println(err)
				break
			} else {
				inout <- msg
			}
		}
	}(c)

loop:
	for {
		select {
		case err := <-errorCh:
			log.Println(err)
			close(inout)
			close(output)
			break loop
		case <-statusCh:
			close(inout)
			close(output)
			break loop
		case err := <-errChan:
			close(inout)
			close(output)
			log.Println(err)
			if err != io.EOF {
				log.Println(err)
			}
			break loop
		}
	}
}
