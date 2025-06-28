package main

import (
    "fmt"
    "io"
    "log"
    "net"

    "github.com/ergagnon/ginder/protos"

    "google.golang.org/grpc"
)

// RawTextServerImpl implements the RawText service.
type RawTextServerImpl struct {
    protos.UnimplementedRawTextServer
}

// Extract implements the Extract method of the RawText service.
func (s *RawTextServerImpl) Extract(stream protos.RawText_ExtractServer) error {
    fmt.Println("New stream started")

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            // No more requests, return success
            fmt.Println("Stream ended")
            return nil
        }
        if err != nil {
            fmt.Printf("Error receiving request: %v\n", err)
            return err
        }

        // Process the received content byte by byte
        content := req.GetContent()
        contentType := "text/plain" // You can set this based on your logic

        for _, b := range content {
            // Create a reply for each byte
            reply := &protos.RawTextReply{
                Type:    contentType,
                Content: []byte{b},
            }

            // Send the reply back to the client
            if err := stream.Send(reply); err != nil {
                fmt.Printf("Error sending reply: %v\n", err)
                return err
            }
        }
    }
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    protos.RegisterRawTextServer(s, &RawTextServerImpl{})

    fmt.Println("Server is running on port 50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
