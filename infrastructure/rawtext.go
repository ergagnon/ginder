package infrastructure

import (
	"context"
	"io"
	"log"

	"github.com/ergagnon/ginder/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RawTextService interface {
    Extract(reader io.Reader) io.Reader
    io.Closer
}

func NewRawTextService() RawTextService {
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    conn, err := grpc.Dial("localhost:50051", opts...)

    if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := protos.NewRawTextClient(conn)

    return &rawTextService{client: client, conn: conn}
}

type rawTextService struct {
    client protos.RawTextClient
    conn *grpc.ClientConn
}

func (me *rawTextService) Extract(reader io.Reader) io.Reader {
    ctx := context.Background()
	stream, err := me.client.Extract(ctx)

    if err != nil {
		log.Fatalf("client.RawTextExtract failed: %v", err)
	}

    pr, pw := io.Pipe()

    go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
                pw.Close()
				return
			}
			if err != nil {
				log.Fatalf("client.RawTextExtract failed: %v", err)
			}

            pw.Write(in.Content)
			log.Println("Extract content from: ", in.Type)
		}
	}()

    buf := make([]byte, 1024)
    go func ()  {
        for {
            n, err := reader.Read(buf)

            if err == io.EOF {
                stream.CloseSend()
                break
            }

            if err := stream.Send(&protos.FileRequest{Content:  buf[:n]}); err != nil {
                log.Fatalf("client.Extract: stream.Send() failed: %v", err)
            }
        }
    }()

    return pr;
}

func (me *rawTextService) Close() error {
    return me.conn.Close()
}
