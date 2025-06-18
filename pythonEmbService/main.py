from sentence_transformers import SentenceTransformer
import emb_pb2, emb_pb2_grpc
from concurrent import futures
import grpc

class EmbedderService(emb_pb2_grpc.EmbedderServicer):
    def __init__(self):
        print("Loading embedding model...")
        self.model = SentenceTransformer('BAAI/bge-small-en-v1.5')
        print("Model loaded.")

    def Embed(self, request, context):
        vector = self.model.encode(request.text)
        return emb_pb2.EmbeddingResponse(embedding=vector.tolist())

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    emb_pb2_grpc.add_EmbedderServicer_to_server(EmbedderService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC server running on port 50051...")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
