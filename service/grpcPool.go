package service

// import (
// 	"errors"
// 	"sync"

// 	"google.golang.org/grpc"
// )

// // type GrpcConnectionPool struct {
// // 	addr      string
// // 	maxConns  int
// // 	idleConns chan net.Conn
// // 	mu        sync.Mutex
// // 	closed    bool
// // }

// // TCPConnectionPool TCP 连接池结构体
// type TCPConnectionPool struct {
// 	// addr      string
// 	maxConns  int
// 	idleConns chan *grpc.ClientConn
// 	mu        sync.Mutex
// 	// closed    bool
// }

// // newTCPConnectionPool 创建一个新的 TCP 连接池
// func newTCPConnectionPool(maxConns int) *TCPConnectionPool {
// 	return &TCPConnectionPool{
// 		// addr:      addr,
// 		maxConns:  maxConns,
// 		idleConns: make(chan *grpc.ClientConn, maxConns),
// 	}
// }

// // get 获取一个连接
// func (p *TCPConnectionPool) get() (*grpc.ClientConn, error) {
// 	select {
// 	case conn := <-p.idleConns:
// 		// 检查连接是否有效
// 		// if err := p.checkConn(conn); err != nil {
// 		// 	conn.Close()
// 		// 	return p.createNewConn(ctx)
// 		// }
// 		return conn, nil
// 	default:
// 		return nil, errors.New("connection pool is empty")
// 	}
// }

// // put 释放连接回连接池
// func (p *TCPConnectionPool) put(conn *grpc.ClientConn) error {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// 	// if p.closed {
// 	// 	conn.Close()
// 	// 	return errors.New("connection pool is closed")
// 	// }
// 	p.idleConns <- conn
// 	return nil
// 	// if err := p.checkConn(conn); err != nil {
// 	// 	conn.Close()
// 	// 	return err
// 	// }

// 	// select {
// 	// case p.idleConns <- conn:
// 	// 	return nil
// 	// default:
// 	// 	// 连接池已满，关闭连接
// 	// 	conn.Close()
// 	// 	return errors.New("connection pool is full")
// 	// }
// }

// // len长度
// func (p *TCPConnectionPool) Len() int {
// 	return len(p.idleConns)
// }

// // Close 关闭连接池
// func (p *TCPConnectionPool) Close() error {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// 	if p.closed {
// 		return nil
// 	}

// 	p.closed = true
// 	close(p.idleConns)

// 	for conn := range p.idleConns {
// 		conn.Close()
// 	}

// 	return nil
// }

// // createNewConn 创建新连接
// func (p *TCPConnectionPool) createNewConn() (*grpc.ClientConn, error) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// 	if len(p.idleConns) >= p.maxConns {
// 		return nil, errors.New("connection pool is full")
// 	}

// 	// dialer := &net.Dialer{
// 	// 	Timeout:   5 * time.Second,
// 	// 	KeepAlive: 30 * time.Second,
// 	// }
// 	kacp := keepalive.ClientParameters{
// 		Time:                10 * time.Second, // 每 10 秒发送一次心跳
// 		Timeout:             5 * time.Second,  // 心跳超时时间 5 秒
// 		PermitWithoutStream: true,             // 允许在没有活跃流的情况下发送心跳
// 	}

// 	// return dialer.DialContext(ctx, "tcp", p.addr)
// 	conn, err := grpc.NewClient("localhost:50051",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithKeepaliveParams(kacp),
// 	)
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	return conn, nil
// }

// checkConn 检查连接是否有效
// func (p *TCPConnectionPool) checkConn(conn net.Conn) error {
// 	if conn == nil {
// 		return errors.New("connection is nil")
// 	}

// 	// 发送一个空字节检查连接是否正常
// 	if _, err := conn.Write([]byte{0}); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func grpcConn(addr string) {
// 	kacp := keepalive.ClientParameters{
// 		Time:                10 * time.Second, // 每 10 秒发送一次心跳
// 		Timeout:             5 * time.Second,  // 心跳超时时间 5 秒
// 		PermitWithoutStream: true,             // 允许在没有活跃流的情况下发送心跳
// 	}

// 	conn, err := grpc.NewClient(addr,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithKeepaliveParams(kacp),
// 		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
// 	)
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	// defer conn.Close()
// 	grpc_pool.put(conn)
// 	fmt.Println("getConn: ", conn)

// 	// c := pb.NewMatchBattleCalV1Client(conn)
// 	// log.Println(c)
// 	// c.AddGift(first_ctx, &battlecalv1pb.AddGiftToGroupReq{
// 	// 	GroupId: "1",
// 	// })
// }

// grpc send
// func grpcSend(msg *pb.AddGiftToGroupReq, count int) error {
// 	count++
// 	var endErr error
// 	conn, err := grpc_pool.get()
// 	if err != nil {
// 		if err1 := oneGetGrpcDomain(); err1 != nil {
// 			ziLog.Error("grpc send error: grpc domain error", debug)
// 			return errors.New("grpc send error: grpc domain error")
// 		}
// 		if count > 3 {
// 			ziLog.Error("grpc send error: grpc count greater 3", debug)
// 			return errors.New("grpc send error: grpc count greater 3")
// 		}
// 		return grpcSend(msg, count)
// 	}
// 	defer func() {
// 		if endErr != nil {
// 			conn.Close()
// 			return
// 		}
// 		log.Println("回收")
// 		grpc_pool.put(conn)
// 	}()

// 	c := pb.NewMatchBattleCalV1Client(conn)
// 	_, endErr = c.AddGift(first_ctx, msg)
// 	if endErr != nil {
// 		oneGetGrpcDomain()
// 		if grpc_pool.Len() == 0 {
// 			ziLog.Error("grpc send endErr: grpc pool len equire 0", debug)
// 			return errors.New("grpc send endErr: grpc pool len equire 0")
// 		}
// 		if count > 3 {
// 			ziLog.Error("grpc send endErr: grpc count greater 3", debug)
// 			return errors.New("grpc send endErr: grpc count greater 3")
// 		}
// 		if err := grpcSend(msg, count); err != nil {
// 			ziLog.Error(fmt.Sprintf("grpc send endErr: %v", err), debug)
// 		}
// 	}
// 	return nil
// }
