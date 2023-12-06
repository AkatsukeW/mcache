package tcp

import (
	"bufio"
	"cache/pkg/cache"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	cache.Cache
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) Listen() {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	for {
		accept, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go s.process(accept)
	}
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		op, err := r.ReadByte()
		if err != nil && err != io.EOF {
			log.Default().Print("close connection:", err.Error())
			return
		}

		if op == 'S' {
			err = s.set(conn, r)
		} else if op == 'G' {
			err = s.get(conn, r)
		} else if op == 'D' {
			err = s.del(conn, r)
		} else {
			log.Default().Print("invalid operate:", op)
			return
		}

		if err != nil {
			log.Default().Print("err:", err.Error())
			return
		}
	}
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	key, value, err := s.readKeyAndValue(r)
	if err != nil {
		return err
	}

	return s.sendResponse(nil, s.Set(key, value), conn)
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}

	val, err := s.Get(key)
	return s.sendResponse(val, err, conn)

}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	return s.sendResponse(nil, s.Del(key), conn)

}

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	kLen, err := s.readLen(r)
	if err != nil {
		return "", err
	}

	key := make([]byte, kLen)
	_, err = io.ReadFull(r, key)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

func (s *Server) readLen(r *bufio.Reader) (int, error) {
	tmp, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}

	length, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	keyLen, err := s.readLen(r)
	if err != nil {
		return "", nil, err
	}

	valueLen, err := s.readLen(r)
	if err != nil {
		return "", nil, err
	}

	key := make([]byte, keyLen)
	_, err = io.ReadFull(r, key)
	if err != nil {
		return "", nil, err
	}

	value := make([]byte, valueLen)
	_, err = io.ReadFull(r, value)
	if err != nil {
		return "", nil, err
	}

	return string(key), value, nil
}

func (s *Server) sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errStr := err.Error()
		tmp := fmt.Sprintf("-%d", len(errStr)) + errStr

		_, err = conn.Write([]byte(tmp))
		return err
	}

	valueLen := fmt.Sprintf("%d", len(value))
	_, err = conn.Write(append([]byte(valueLen), value...))
	return err
}
