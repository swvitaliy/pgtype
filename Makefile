test:
	PGX_TEST_DATABASE='postgresql://postgres:12345678@localhost:5432/pgtype_test?sslmode=disable' go test
