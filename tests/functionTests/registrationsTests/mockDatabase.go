package functionTests

/*
var (
	mockClient FirestoreClient
	mockCtx    context.Context
)*/

// Refactor everything here.

/*
type (
	// FirestoreClient defines the methods used by Firestore client.
	FirestoreClient interface {
		Set(ctx context.Context, docRef *firestore.DocumentRef, data interface{},
			opts ...firestore.SetOption) (*firestore.WriteResult, error)
		Collection(string) *firestore.CollectionRef
		Where(path string, op string, value interface{}) firestore.Query
		Limit(n int) firestore.Query
		//Documents(ctx context.Context) *firestore.DocumentIterator
		//GetAll() ([]*firestore.DocumentSnapshot, error)
		Documents(ctx context.Context) ([]*firestore.DocumentSnapshot, error)
		// Add other methods used by Firestore client if necessary
	}

	// MockFirestoreClient is a mock implementation of FirestoreClient for testing.
	MockFirestoreClient struct {
		//SetFunc func(ctx context.Context, docRef *firestore.DocumentRef, data interface{},
		//	opts ...firestore.SetOption) (*firestore.WriteResult, error)
		//CollectionFunc func(path string) *firestore.CollectionRef
		//WhereFunc      func(path string, op string, value interface{}) firestore.Query
		//LimitFunc      func(n int) firestore.Query
		//DocumentsFunc  func(ctx context.Context) *firestore.DocumentIterator
		mock.Mock
		// Add other mock functions if necessary
	}
)*/

/*
var client = MockFirestoreClient{}

func NewMockFirestoreClient() *MockFirestoreClient {
	return &MockFirestoreClient{
		SetFunc:        client.SetFunc,
		CollectionFunc: client.CollectionFunc,
		WhereFunc:      client.WhereFunc,
		LimitFunc:      client.LimitFunc,
		DocumentsFunc:  client.DocumentsFunc,
	}
}*/
/*
func (m *MockFirestoreClient) Where(path string, op string, value interface{}) firestore.Query {
	args := m.Called(path, op, value)
	return args.Get(0).(firestore.Query)
}

func (m *MockFirestoreClient) Limit(n int) firestore.Query {
	args := m.Called(n)
	return args.Get(0).(firestore.Query)
}

func (m *MockFirestoreClient) Collection(path string) *firestore.CollectionRef {
	args := m.Called(path)
	return args.Get(0).(*firestore.CollectionRef)
}

func (m *MockFirestoreClient) Documents(ctx context.Context) *firestore.DocumentIterator {
	args := m.Called(ctx)
	return args.Get(0).(*firestore.DocumentIterator)
}

func (m *MockFirestoreClient) Documents(ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*firestore.DocumentSnapshot), args.Error(1)
}

func (m *MockFirestoreClient) GetAll() ([]*firestore.DocumentSnapshot, error) {
	args := m.Called()
	return args.Get(0).([]*firestore.DocumentSnapshot), args.Error(1)
}

// Set implements the Set method for the mock Firestore client.
func (m *MockFirestoreClient) Set(ctx context.Context,
	docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	args := m.Called(ctx, docRef, data, opts)
	return args.Get(0).(*firestore.WriteResult), args.Error(1)
}*/
